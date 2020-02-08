package producers

import (
	"errors"
	"log"
	"time"

	"github.com/freundallein/resender/data"
)

const (
	entityName = "packages"

	attempts    = 5
	connTimeout = 5 * time.Second
)

var (
	// ErrNoConnection - you should check connection to producer
	ErrNoConnection = errors.New("no connection")
)

// Producer - common interface for package producers
type Producer interface {
	GetName() string
	Validate(data.Package) error
	Produce(string, []byte)
}

// Gather - collect all avalilable producers
func Gather(externalUrl, elasticUrl, kafkaUrl string) []Producer {
	prds := []Producer{
		NewHttp(externalUrl),
	}

	kafka, err := NewKafka(kafkaUrl, entityName)
	if err != nil {
		log.Println("[kafka]", err)
	}
	if kafka != nil {
		prds = append(prds, kafka)
	}
	elastic, err := NewElastic(elasticUrl, entityName)
	if err != nil {
		log.Println("[elastic]", err)
	}
	if elastic != nil {
		prds = append(prds, elastic)
	}
	return prds
}
