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

func connect(clientId string, config models.Config) mqtt.Client {
	opts := createClientOptions(clientId, config)
	client := mqtt.NewClient(opts)
	token := client.Connect()

	for !token.WaitTimeout(3 * time.Second) {
	}

	if err := token.Error(); err != nil {
		log.Fatal(err)
	}

	return client
}

func createClientOptions(clientId string, config models.Config) *mqtt.ClientOptions {
	fmt.Printf("config: %+v\n", config)

	broker := fmt.Sprintf("%s://%s:%d", config.Protocol, config.Host, config.Port)
	fmt.Printf("broker: %s\n", broker)

	uriString := fmt.Sprintf("tcp://%s:%s@%s:%v", config.User.Username, config.User.Password, config.Host, config.Port)
	uri, err := url.Parse(uriString)
	if err != nil {
		log.Fatal(err)
	}

	password, _ := uri.User.Password()
	fmt.Printf("uriString: %s\n", uriString)

	fmt.Printf("config.User.Username: %s\n", config.User.Username)
	fmt.Printf("uri.User.Username(): %s\n", uri.User.Username())
	fmt.Printf("config.User.Password: %s\n", config.User.Password)
	fmt.Printf("config.User.Password(): %s\n", config.User.Password())
	fmt.Printf("password: %s\n", password)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetUsername(uri.User.Username())
	opts.SetPassword(password)
	opts.SetClientID(clientId)

	return opts
}

func main() {
	config := loadConfigs()
	client := connect("pub", config)

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
