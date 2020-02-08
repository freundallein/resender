package producers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/freundallein/resender/data"
)

// ElasticProducer - send pacakges to ElasticSearch
type ElasticProducer struct {
	Name  string
	Index string
	cli   *elasticsearch.Client
}

//NewElastic - constructor
func NewElastic(elasticUrl, entityName string) (*ElasticProducer, error) {
	log.Println("[elastic] trying to connect", elasticUrl)
	cfg := elasticsearch.Config{
		Addresses: []string{elasticUrl},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	err = checkElastic(es)
	if err != nil {
		return nil, err
	}
	log.Println("[elastic] connecton succeed to", elasticUrl)
	prdc := &ElasticProducer{
		Name:  "[elastic]",
		Index: entityName,
		cli:   es,
	}
	err = prdc.createIndex()
	if err != nil {
		return nil, err
	}
	return prdc, nil
}

// GetName - getter
func (p *ElasticProducer) GetName() string {
	return p.Name
}

// Validate - can valdiate pacakages, if you need to.
func (p *ElasticProducer) Validate(pkg data.Package) error {
	return nil
}

// Produce - package indexing
func (p *ElasticProducer) Produce(uid string, data []byte) {
	log.Println(p.Name, uid, "sending to", p.Index)
	req := esapi.IndexRequest{
		Index:      p.Index,
		DocumentID: uid,
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), p.cli)
	if err != nil {
		log.Println(p.Name, "Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("%s %s Error indexing document ID=%s", p.Name, res.Status(), uid)
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Printf("%s Error parsing the response body: %s", p.Name, err)
	} else {
		log.Printf("%s %s succeeded; version=%d", p.Name, r["_id"], int(r["_version"].(float64)))
	}
}

func (p *ElasticProducer) createIndex() error {
	log.Println(p.Name, "checking index existance:", p.Index)
	res, err := p.cli.Indices.Create(p.Index)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return err
	}
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		} else {
			if e["error"].(map[string]interface{})["type"] != "resource_already_exists_exception" {
				return errors.New(e["error"].(map[string]interface{})["reason"].(string))
			}
		}
	}
	return nil
}

func checkElastic(cli *elasticsearch.Client) error {
	for i := 1; i <= attempts; i++ {
		_, err := cli.Info()
		if err != nil {
			log.Println("[elastic] attempting", i, "of", attempts, err)
			time.Sleep(connTimeout)
			if i == attempts {
				return ErrNoConnection
			}
			continue
		}
		break
	}
	return nil
}
