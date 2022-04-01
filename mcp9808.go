package main

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/mcp9808"
	"periph.io/x/host/v3"
)

type Mcp9808 struct {
	bus      string
	location string
}

func (m *Mcp9808) Measure() Measurement {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	// Open a handle to the first available I²C bus:
	bus, err := i2creg.Open(m.bus)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Open a handle to a bme280/bmp280 connected on the I²C bus using default
	// settings:
	dev, err := mcp9808.New(bus, &mcp9808.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	// defer dev.Halt()

	// Read temperature from the sensor:
	temp, err := dev.SenseTemp()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%8s\n", temp)

	var _null *float64
	var t = temp.Celsius()
	return Measurement{m.location, &t, _null, _null}
}
