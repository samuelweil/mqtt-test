package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func main() {

	seed := time.Now().UnixNano()
	rand.Seed(seed)

	sensor := NewSensor("tcp://broker:1883")
	fmt.Printf("Starting sensor %s\n", sensor.ID)

	for {
		go sensor.Poll()
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
	}
}

type Sensor struct {
	mqttClient mqtt.Client
	ID         string
	topic      string
}

func NewSensor(mqttBroker string) *Sensor {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	mqttClient := connect(opts)

	sensorID, ok := os.LookupEnv("SENSOR_ID")
	if !ok {
		sensorID = uuid.New().String()
	}

	topic := fmt.Sprintf("/sensors/temperature/%s", sensorID)

	return &Sensor{
		mqttClient: mqttClient,
		ID:         sensorID,
		topic:      topic,
	}
}

const (
	sensorMax = 100
	sensorMin = 0
)

type SensorPayload struct {
	Timestamp  string  `json:"timeStamp"`
	SensorID   string  `json:"sensorId"`
	SensorType string  `json:"sensorType"`
	Value      float32 `json:"value"`
	Units      string  `json:"units"`
}

func (s *Sensor) Poll() {
	value := float32(rand.Intn(sensorMax) + sensorMin)

	payload, err := json.Marshal(SensorPayload{
		SensorID:   s.ID,
		SensorType: "temperature",
		Value:      value,
		Units:      "fahrenheit",
		Timestamp:  time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("Error creating message payload %s", err)
		return
	}

	token := s.mqttClient.Publish(s.topic, byte(0), false, payload)
	if token.Wait() && token.Error() != nil {
		log.Printf("Error sending sensor reading: %s", token.Error())
	}
}

func connect(opts *mqtt.ClientOptions) mqtt.Client {
	maxTries := 5
	client := mqtt.NewClient(opts)

	var e error

	for i := 0; i < maxTries; i++ {
		token := client.Connect()
		token.Wait()
		e = token.Error()
		if e != nil {
			time.Sleep(2 * time.Second)
		} else {
			return client
		}

	}

	panic(e)
}
