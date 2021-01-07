package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/peknur/ruuvitag"

	"gitlab.com/kirbo/addon-ruuvitag-mqtt/ruuvitag-mqtt/internal/models"
)

func loadConfigs() models.Config {
	log.Print("Loading configs")
	jsonFile, err := os.Open("/data/options.json")
	if err != nil {
		fmt.Print(err)
	}

	var config models.Config

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &config)

	return config
}

func connect(clientId string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()

	for !token.WaitTimeout(3 * time.Second) {
	}

	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	return client
}

func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
    host := fmt.Sprintf("tcp://%s", uri.Host)
    fmt.Printf("host: %s\n", host)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(host)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)

	return opts
}

func main() {
	config := loadConfigs()

	fmt.Printf("config: %+v\n", config)

	uriString := fmt.Sprintf("mqtt://%s:%s@%s:%v", config.User.Username, config.User.Password, config.Host, config.Port)
	fmt.Printf("uriString: %s\n", uriString)

	uri, err := url.Parse(uriString)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("uri: %+v\n", uri)

	client := connect("pub", uri)

	scanner, err := ruuvitag.OpenScanner(10)
	if err != nil {
		log.Fatal(err)
	}

	output := scanner.Start()
	for {
		data := <-output

		var (
			battery      = data.BatteryVoltage()
			humidity     = data.Humidity()
			pressure     = float32(data.Pressure()) / float32(100)
			temperature  = data.Temperature()
			acceleration = models.DeviceAcceleration{
				X: data.AccelerationX(),
				Y: data.AccelerationY(),
				Z: data.AccelerationZ(),
			}

			topic = "ruuvitag/" + data.DeviceID() + "/"

			topicA = topic + "acceleration"
			topicB = topic + "battery"
			topicH = topic + "humidity"
			topicP = topic + "pressure"
			topicT = topic + "t"
		)

		client.Publish(topicA, 0, true, acceleration)
		client.Publish(topicB, 0, true, battery)
		client.Publish(topicH, 0, true, humidity)
		client.Publish(topicP, 0, true, pressure)
		client.Publish(topicT, 0, true, temperature)
	}
}
