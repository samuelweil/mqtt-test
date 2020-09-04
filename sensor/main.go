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

type sensorType struct {
	Name  string
	Min   int
	Max   int
	Units string
}

var sensorTypes []sensorType = []sensorType{
	{
		Name:  "temperature",
		Min:   0,
		Max:   100,
		Units: "F",
	},
	{
		Name:  "humidity",
		Min:   0,
		Max:   100,
		Units: "%",
	},
	{
		Name:  "pressure",
		Min:   0,
		Max:   30,
		Units: "psi",
	},
}

type Sensor struct {
	mqttClient mqtt.Client
	ID         string
	typ        sensorType
	zoneID     string
}

func NewSensor(mqttBroker string) *Sensor {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttBroker)
	mqttClient := connect(opts)

	sensorID, ok := os.LookupEnv("SENSOR_ID")
	if !ok {
		sensorID = uuid.New().String()
	}

	sensTypeInt := rand.Intn(len(sensorTypes))

	return &Sensor{
		mqttClient: mqttClient,
		ID:         sensorID,
		typ:        sensorTypes[sensTypeInt],
		zoneID:     uuid.New().String(),
	}
}

func (s *Sensor) topic() string {
	return fmt.Sprintf("sensors/%s/%s/%s", s.ID, s.zoneID, s.typ.Name)
}

func (s *Sensor) randValue() float32 {
	return float32(rand.Intn(s.typ.Max-s.typ.Min) + s.typ.Min)
}

type SensorPayload struct {
	Timestamp string  `json:"timeStamp"`
	Value     float32 `json:"value"`
	Units     string  `json:"units"`
}

func (s *Sensor) Poll() {

	payload, err := json.Marshal(SensorPayload{
		Value:     s.randValue(),
		Units:     s.typ.Units,
		Timestamp: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		log.Printf("Error creating message payload %s", err)
		return
	}

	token := s.mqttClient.Publish(s.topic(), byte(0), false, payload)
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
