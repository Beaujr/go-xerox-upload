package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	xclient "github.com/beaujr/go-xerox-upload/client"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	qos        = flag.Int("qos", 0, "The QoS to subscribe to messages at")
	username   = flag.String("username", "", "A username to authenticate to the MQTT server")
	password   = flag.String("password", "", "Password to match username")
	mqttServer = flag.String("subServer", "tcp://broker.emqx.io:1883", "MQTT broker")
	topic      = flag.String("subTopic", "go-xerox-upload/ocr", "Topic to publish to")
)

func onMessageReceived(client MQTT.Client, message MQTT.Message) {
	fmt.Printf("Received message on topic: %s\nMessage: %s\n", message.Topic(), message.Payload())
	var rrq xclient.OCRMessage
	x, err := xclient.NewGoogleClient()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := json.Unmarshal(message.Payload(), &rrq); err != nil {
		panic(err)
	}
	_, err = x.OCRFile(rrq.FileId, rrq.ParentId, rrq.Name)
	if err != nil {
		log.Printf("Error occurred for file %s, %s\n", rrq.Name, err.Error())
	}
	log.Println("Submitted")
}

func main() {
	flag.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	hostname, _ := os.Hostname()
	clientid := flag.String("clientid", hostname+strconv.Itoa(time.Now().Second()), "A clientid for the connection")
	flag.Parse()

	connOpts := MQTT.NewClientOptions().AddBroker(*mqttServer).SetClientID(*clientid).SetCleanSession(true)
	if *username != "" {
		connOpts.SetUsername(*username)
		if *password != "" {
			connOpts.SetPassword(*password)
		}
	}
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	connOpts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(*topic, byte(*qos), onMessageReceived); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to %s\n", *mqttServer)
	}
	<-c
}
