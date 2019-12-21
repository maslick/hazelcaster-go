package main

import (
	"log"
	"testing"
	"time"
)

func Test_SimpleHazelcastClient(t *testing.T) {
	hz := newHzClient()
	for {
		<-time.After(time.Duration(3) * time.Second)
		readings, err := hz.fetch()
		if err != nil {
			log.Println(err)
		}

		for _, reading := range readings {
			log.Println(reading)
		}
	}
}
