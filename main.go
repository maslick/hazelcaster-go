package main

import (
	"fmt"
	"log"
	"time"
)

var i = 0

func main() {
	hz := newHzClient()
	createNewReading(hz)
	for {
		<-time.After(time.Duration(10) * time.Second)
		createNewReading(hz)
	}
}

func createNewReading(hazelcaster *Hazelcaster) {
	reading := Reading{
		Name:      fmt.Sprintf("reading%d", i),
		Timestamp: time.Now().Unix(),
	}
	log.Println("Pushing reading to Hazelcast:", reading)
	err := hazelcaster.persist(reading)
	if err != nil {
		log.Println(err)
	}
	i++
}
