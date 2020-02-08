# resender
[![Build Status](https://travis-ci.org/freundallein/resender.svg?branch=master)](https://travis-ci.org/freundallein/resender)
[![Coverage Status](https://coveralls.io/repos/github/freundallein/resender/badge.svg?branch=master)](https://coveralls.io/github/freundallein/resender?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/freundallein/resender)](https://goreportcard.com/report/github.com/freundallein/resender)

Receive json packages via http, validate data and produce it to:  
ElasticSearch  
Kafka  
other HTTP server  


## Configuration
Application supports configuration via environment variables:
```
PORT=8000
MACHINE_ID=1                             (used for unique id generator)
EXTERNAL_URL=http://192.168.52.138:8080/
ELASTIC_URL=http://elastic:9200/
KAFKA_URL=kafka:19092
```
## Installation
### With docker  
```
$> docker pull freundallein/resender
```
### With source
```
$> git clone git@github.com:freundallein/resender.git
$> cd resender
$> make build
```

## Usage
Docker-compose

```
version: "3.5"

networks:
  network:
    name: network
    driver: bridge

volumes:
  kafdata:
    driver: local
  esdata:
    driver: local

services:
  loadbalancer:
    image: freundallein/go-lb:latest
    container_name: loadbalancer
    restart: always
    environment: 
      - PORT=8000
      - ADDRS=http://resender-one:8000,http://resender-two:8000
      - STALE_TIMEOUT=60
    networks: 
      - network
    ports:
      - 8000:8000

  resender-one:
    image: freundallein/resender:latest
    container_name: resender-one
    restart: always
    environment: 
      - PORT=8000
      - MACHINE_ID=1
      - EXTERNAL_URL=http://192.168.52.138:8080/
      - ELASTIC_URL=http://elastic:9200/
      - KAFKA_URL=kafka:19092
    depends_on: 
      - kafka
      - elastic
    networks: 
      - network
    expose:
      - 8000

  resender-two:
    image: freundallein/resender:latest
    container_name: resender-two
    restart: always
    environment: 
      - PORT=8000
      - MACHINE_ID=2
      - EXTERNAL_URL=http://192.168.52.138:8080/
      - ELASTIC_URL=http://elastic:9200/
      - KAFKA_URL=kafka:19092
    depends_on: 
      - kafka
      - elastic
    networks: 
      - network
    expose:
      - 8000

  elastic:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.5.2
    container_name: elastic
    environment:
      - discovery.type=single-node
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - esdata:/usr/share/elasticsearch/data
    ports:
      - 9200:9200
    networks:
      - network

  zookeeper:
    image: zookeeper:3.4.9
    container_name: zookeeper
    ports:
      - 2181:2181
    environment:
        ZOO_MY_ID: 1
        ZOO_PORT: 2181
        ZOO_SERVERS: server.1=zookeeper:2888:3888
    networks:
      - network
    
  kafka:
    image: confluentinc/cp-kafka:5.3.1
    container_name: kafka
    environment:
      KAFKA_ADVERTISED_LISTENERS: LISTENER_DOCKER_INTERNAL://kafka:19092,LISTENER_DOCKER_EXTERNAL://${DOCKER_HOST_IP:-127.0.0.1}:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: LISTENER_DOCKER_INTERNAL:PLAINTEXT,LISTENER_DOCKER_EXTERNAL:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: LISTENER_DOCKER_INTERNAL
      KAFKA_ZOOKEEPER_CONNECT: "zookeeper:2181"
      KAFKA_BROKER_ID: 1
      KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    networks:
      - network
    ports:
      - 9092:9092
    volumes:
      - kafdata:/var/lib/kafka/data
    depends_on:
      - zookeeper
```

## Request example
Server consumes only POST requests
```
POST http://0.0.0.0:8000 HTTP/1.1
content-Type: application/json

{
"ap_id" : "A8-F9-4B-B6-87-FF",
"probe_requests" : [
{
"mac" : "88-1D-FC-DF-6F-C1",
"timestamp" : "1579782767"
},
{
"mac" : "F8-59-71-PK-95-36",
"bssid" : "04-BF-6D-04-09-8C",
"timestamp" : "1579782767",
"ssid" : "SKOLTECH"
},
{
"mac" : "F8-59-71-PK-95-BB",
"timestamp" : "1579782767"
}
]
}
```
Service will valdiate incoming request, fill default values  and resend it to next services (http/Elastic/Kafka) is they are available.
## Metrics
Default prometheus metrics are available on `/metrics`  
