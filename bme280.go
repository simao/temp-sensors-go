package main

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/bmxx80"
	"periph.io/x/host/v3"
)

type Bme280 struct {
	bus      string
	location string
}

func (b *Bme280) Measure() Measurement {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open(b.bus)
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	// Open a handle to a bme280/bmp280 connected on the I²C bus using default
	// settings:
	dev, err := bmxx80.NewI2C(bus, 0x76, &bmxx80.DefaultOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Halt()

	// Read temperature from the sensor:
	var env physic.Env
	if err = dev.Sense(&env); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%8s %10s %9s\n", env.Temperature, env.Pressure, env.Humidity)

	t := env.Temperature.Celsius()
	p := float64(env.Pressure) / 1000000000.0
	h := float64(env.Humidity) / 100000.0

	return Measurement{
		b.location,
		&t,
		&p,
		&h,
	}
}
