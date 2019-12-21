package main

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-go-client/config"
	"github.com/hazelcast/hazelcast-go-client/config/property"
	"github.com/hazelcast/hazelcast-go-client/core/logger"
	"log"
	"os"
	"sort"
)

type Hazelcaster struct {
	client hazelcast.Client
}

const collectionName = "hazelcaster"

func newHzClient(clearOnStartup ...bool) *Hazelcaster {
	hzAddress := getEnv("HZ_SERVER_ADDR", "192.168.99.100:5701")
	hzUsername := getEnv("HZ_USERNAME", "dev")
	hzPassword := getEnv("HZ_PASSWORD", "dev-pass")

	cfg := hazelcast.NewConfig()
	cfg.NetworkConfig().SSLConfig().SetEnabled(false)
	cfg.NetworkConfig().AddAddress(hzAddress)
	cfg.GroupConfig().SetName(hzUsername)
	cfg.GroupConfig().SetPassword(hzPassword)

	discoveryCfg := config.NewCloudConfig()
	discoveryCfg.SetEnabled(false)
	cfg.NetworkConfig().SetCloudConfig(discoveryCfg)

	cfg.SetProperty(property.LoggingLevel.Name(), logger.DebugLevel)
	cfg.SetProperty(property.StatisticsEnabled.Name(), "true")
	cfg.SetProperty(property.StatisticsPeriodSeconds.Name(), "1")

	cfg.SerializationConfig().AddPortableFactory(portableFactorID, &ReadingPortableFactory{})
	hazelcastClient, err := hazelcast.NewClientWithConfig(cfg)
	if err != nil {
		fmt.Println(err)
	}

	if len(clearOnStartup) == 0 || len(clearOnStartup) == 1 && clearOnStartup[0] {
		l, err := hazelcastClient.GetList(collectionName)
		if err == nil {
			err := l.Clear()
			log.Println("Clearing list... success:", err == nil)
		}
	}
	return &Hazelcaster{client: hazelcastClient}
}

func (hz *Hazelcaster) persist(reading Reading) error {
	readingsList, err := hz.client.GetList(collectionName)
	if err != nil {
		return err
	}

	ok, err := readingsList.Add(reading)
	if err != nil {
		return err
	}
	log.Println("Saved to Hazelcast cloud:", ok)
	return nil
}

func (hz *Hazelcaster) fetch() ([]Reading, error) {
	list, err := hz.client.GetList(collectionName)
	if err != nil {
		return nil, err
	}

	size, _ := list.Size()
	log.Println("List size:", size)
	slice, err := list.ToSlice()
	if err != nil {
		return nil, err
	}

	var result []Reading
	for _, r := range slice {
		result = append(result, r.(Reading))
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp < result[j].Timestamp
	})

	return result, nil
}

func (hz *Hazelcaster) shutdown() {
	hz.client.Shutdown()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
