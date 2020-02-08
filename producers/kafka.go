package producers

import (
	"context"
	"log"
	"time"

	"github.com/freundallein/resender/data"
	"github.com/segmentio/kafka-go"
)

// KafkaProducer - send packages to Apache Kafka
type KafkaProducer struct {
	Name  string
	Topic string
	conn  *kafka.Conn
}

//NewKafka - constructor
func NewKafka(kafkaUrl, entityName string) (*KafkaProducer, error) {
	log.Println("[kafka] trying to connect", kafkaUrl)
	var conn *kafka.Conn
	for i := 1; i <= attempts; i++ {
		kfk, err := kafka.DialLeader(context.Background(), "tcp", kafkaUrl, entityName, 0)
		if err != nil {
			log.Println("[kafka] attempting", i, "of", attempts, err)
			time.Sleep(connTimeout)
			continue
		}
		conn = kfk
		break
	}
	if conn == nil {
		return nil, ErrNoConnection
	}
	log.Println("[kafka] connecton succeed to", kafkaUrl)
	return &KafkaProducer{
		Name:  "[kafka]",
		Topic: entityName,
		conn:  conn,
	}, nil
}

// GetName - getter
func (p *KafkaProducer) GetName() string {
	return p.Name
}

// Validate - can valdiate pacakages, if you need to.
func (p *KafkaProducer) Validate(pkg data.Package) error {
	return nil
}

// Produce - package indexing
func (p *KafkaProducer) Produce(uid string, data []byte) {
	log.Println(p.Name, uid, "posting to", p.Topic)
	message := kafka.Message{
		Key:   []byte(uid),
		Value: data,
	}
	_, err := p.conn.WriteMessages(message)
	if err != nil {
		log.Println(p.Name, err)
		return
	}
	log.Println(p.Name, uid, "succeeded")
}
