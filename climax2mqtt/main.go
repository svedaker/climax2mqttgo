package main

import (
	"climax/climax"
	"climax/mqttService"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Mqtt   mqttService.MqttConfig `json:"mqtt"`
	Climax climax.ClimaxConfig    `json:"climax"`
}

func main() {
	var cfg Config
	optionsPath := "data/options.json"

	file, err := os.Open(optionsPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&cfg.Mqtt)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	file.Seek(0, 0)
	err = json.NewDecoder(file).Decode(&cfg.Climax)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	log.Printf("MQTT Config: %+v\n", cfg.Mqtt)
	log.Printf("Climax Config: %+v\n", cfg.Climax)

	server(&cfg)
}

func server(config *Config) {
	mqttClient := mqttService.Connect(config.Mqtt)
	repo := climax.NewMemoryDeviceRepository()

	// Periodically fetch devices and publish updates
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		devices, err := config.Climax.GetDevices()
		if err != nil {
			log.Printf("Error fetching devices: %v", err)
			continue
		}
		for _, device := range devices {
			deviceId := device.Identify()

			if repo.IsNewDevice(deviceId) {
				publishDiscoveryMessages(device, mqttClient)
			}
			if repo.AddOrUpdate(device) {
				publishUpdateValueMessage(device, mqttClient)
			}
		}
	}
}

func publishDiscoveryMessages(device climax.DeviceInterface, mqttClient mqtt.Client) {
	switch dev := device.(type) {
	case climax.TemperatureSensor:
		publishIfNoError(dev.MqttDiscoveryMessageTemperature(), device, mqttClient)
	case climax.PowerSwitchMeter:
		publishIfNoError(dev.MqttDiscoveryMessageEnergy(), device, mqttClient)
		publishIfNoError(dev.MqttDiscoveryMessageSwitch(), device, mqttClient)
		publishIfNoError(dev.MqttDiscoveryMessagePower(), device, mqttClient)
	default:
		log.Printf("Unsupported device type for device Id %s %s", dev.Identify(), dev)
	}
}

func publishUpdateValueMessage(device climax.DeviceInterface, mqttClient mqtt.Client) {
	switch dev := device.(type) {
	case climax.TemperatureSensor:
		publishIfNoError(dev.MqttUpdateValueMessage(), device, mqttClient)
	case climax.PowerSwitchMeter:
		publishIfNoError(dev.MqttUpdateValueMessage(), device, mqttClient)
	default:
		log.Printf("Unsupported device type for device ID %s", dev.Identify())

	}
}

func publishIfNoError(mqttMessage climax.MqttMessage, device climax.DeviceInterface, mqttClient mqtt.Client) {
	if mqttMessage.Err != nil {
		log.Printf("Error generating discovery message for device %s: %v", device.Identify(), mqttMessage.Err)
		return
	}
	log.Printf("Publis to topic %s with message %s", mqttMessage.Topic, string(mqttMessage.Message))
	mqttService.Publish(mqttClient, mqttMessage.Topic, mqttMessage.Message)
}
