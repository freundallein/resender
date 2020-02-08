package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/freundallein/resender/httpserv"
	"github.com/freundallein/resender/producers"
	"github.com/freundallein/resender/uidgen"
)

const (
	timeFormat = "02.01.2006 15:04:05"

	portKey     = "PORT"
	machineKey  = "MACHINE_ID"
	defaultPort = "8000"

	externalUrlKey      = "EXTERNAL_URL"
	defaultExteranalUrl = "http://0.0.0.0:8000/"

	elasticUrlKey = "ELASTIC_URL"
	kafkaUrlKey   = "KAFKA_URL"
)

type logWriter struct{}

// Write - custom logger formatting
func (writer logWriter) Write(bytes []byte) (int, error) {
	msg := fmt.Sprintf("%s | [resender] %s", time.Now().UTC().Format(timeFormat), string(bytes))
	return fmt.Print(msg)
}

func getEnv(key string, fallback string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}
	return fallback, nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	port, err := getEnv(portKey, defaultPort)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	mID, err := getEnv(machineKey, "1")
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	machineID, err := strconv.Atoi(mID)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	externalUrl, err := getEnv(externalUrlKey, defaultExteranalUrl)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	elasticUrl, err := getEnv(elasticUrlKey, "")
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	kafkaUrl, err := getEnv(kafkaUrlKey, "")
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	prds := producers.Gather(externalUrl, elasticUrl, kafkaUrl)

	gen := uidgen.New(uint8(machineID))
	httpOptions := &httpserv.Options{
		Port:      port,
		Producers: prds,
		Gen:       gen,
	}
	srv, err := httpserv.New(httpOptions)
	if err != nil {
		log.Fatalf("[config] %s\n", err.Error())
	}
	if err := srv.Run(); err != nil {
		log.Fatalf("[server] %s\n", err.Error())
	}
}
