package main

import (
	"fmt"
	"log"

	"github.com/MichaelS11/go-dht"
)

type Dht22 struct {
	gpio     string
	bus      string
	location string
}

func (m *Dht22) Measure() Measurement {
	err := dht.HostInit()
	if err != nil {
		log.Fatal(err)
	}

	dht, err := dht.NewDHT(m.gpio, dht.Celsius, "")
	if err != nil {
		log.Fatal(err)
	}

	humidity, temperature, err := dht.ReadRetry(11)
	if err != nil {
		log.Fatal(err)
	}

	var _null *float64

	fmt.Printf("humidity=%v%%, temp=%v*C\n", humidity, temperature)

	return Measurement{m.location, &temperature, _null, &humidity}
}
