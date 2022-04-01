package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

type Measurement struct {
	Location string
	Temp     *float64
	Pressure *float64
	Humidity *float64
}

type Provider interface {
	Measure() Measurement
}

func saveMeasurement(conn *pgx.Conn, m *Measurement) error {
	_, err := conn.Exec(context.Background(), "insert into inside_weather(time, location, temp, pressure, humidity)  values(NOW(), $1, $2, $3, $4)", m.Location, m.Temp, m.Pressure, m.Humidity)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:5432"
	}
	url := fmt.Sprintf("postgres://postgres:password@%s/sensors", dbHost)

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	sensorName := os.Getenv("SENSOR_NAME")
	if sensorName == "" {
		fmt.Println("SENSOR_NAME is not set")
		os.Exit(2)
	}

	i2cBus := os.Getenv("I2C_BUS")

	gpio := os.Getenv("GPIO")

	if gpio == "" {
		gpio = "GPIO19"
	}

	fmt.Fprintf(os.Stderr, "dbHost=%s, sensorName=%s, i2cBus=%s, gpio=%s\n", dbHost, sensorName, i2cBus, gpio)

	var provider Provider
	if len(os.Args) > 1 && os.Args[1] == "mcp9808" {
		fmt.Println("sensor is mcp9808")
		provider = &Mcp9808{i2cBus, sensorName}
	} else if len(os.Args) > 1 && os.Args[1] == "dht22" {
		fmt.Println("sensor is dht22")
		provider = &Dht22{gpio, i2cBus, sensorName}
	} else if len(os.Args) > 1 && os.Args[1] == "bme280" || len(os.Args) == 1 {
		fmt.Println("sensor is bme280")
		provider = &Bme280{i2cBus, sensorName}
	} else {
		fmt.Println("usage: rpi-temp [mcp9808|dht22|bme280]")
		os.Exit(2)
	}

	// TODO: provider.Measure should be a channel and we read from that channel here
	for {
		m := provider.Measure()

		if err = saveMeasurement(conn, &m); err != nil {
			log.Fatal(err)
			return
		}

		time.Sleep(5 * time.Second)
	}
}
