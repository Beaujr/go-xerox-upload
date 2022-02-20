package client

import (
	"encoding/json"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"time"
)

var (
	f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}
	mqttEnabled = flag.Bool("mqtt", false, "Publish to MQTT")
	mqttServer  = flag.String("mqttServer", "tcp://broker.emqx.io:1883", "MQTT broker")
	topic       = flag.String("topic", "go-xerox-upload/ocr", "Topic to publish to")
	clientId    = flag.String("clientId", "go-xerox-upload", "Client Id for MQTT")
)

type OCRMessage struct {
	FileId   string `json:"fileId"`
	ParentId string `json:"parentId"`
	Name     string `json:"name"`
}

func (google *GoogleClient) submitToPubSub(fileID *string, parentId *string, name *string) error {
	opts := mqtt.NewClientOptions().AddBroker(*mqttServer).SetClientID(*clientId)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe(*topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	b, err := json.Marshal(&OCRMessage{Name: *name, ParentId: *parentId, FileId: *fileID})
	if err != nil {
		return err
	}
	log.Println("publishing:" + string(b))
	token := c.Publish(*topic, 0, false, string(b))
	token.Wait()
	if token := c.Unsubscribe(*topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	c.Disconnect(250)
	return nil
}
