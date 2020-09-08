package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	fmt.Println("Starting client")

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://broker:1883")

	client := connect(opts)

	token := client.Subscribe("sensors/#", byte(0), handleMessage)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
	}
}

func handleMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message received: %s: %s\n", msg.Topic(), msg.Payload())
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
