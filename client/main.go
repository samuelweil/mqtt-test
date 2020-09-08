package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Client struct{}

func (c *Client) Subscribe(topic string) chan map[string]interface{} {
	result := make(chan map[string]interface{})
	return result
}

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

type JsonObj map[string]interface{}

func (obj JsonObj) String() string {

	fields := make([]string, 0)

	for k, v := range obj {

		var str string

		switch v.(type) {
		case int:
			str = fmt.Sprintf("%d", v)
		case string:
			str = v.(string)
		default:
			str = fmt.Sprintf("%v", v)
		}

		fields = append(fields, fmt.Sprintf("%s: %s", k, str))
	}

	var result strings.Builder
	result.WriteString("{\n  ")
	result.WriteString(strings.Join(fields, ",\n  "))
	result.WriteString("\n}")
	return result.String()
}

func handleMessage(client mqtt.Client, msg mqtt.Message) {

	var payload JsonObj

	e := json.Unmarshal(msg.Payload(), &payload)
	if e == nil {
		fmt.Println(payload)
	} else {
		fmt.Println(e)
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
