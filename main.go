package main

import (
	"climax/climax"
	"climax/mqttService"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Mqtt   mqttService.MqttConfig `yaml:"mqtt"`
	Climax climax.ClimaxConfig    `yaml:"climax"`
}

func main() {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Println("Error loading config:", err)
		return
	}

	log.Printf("MQTT Config: %+v\n", cfg.Mqtt)
	log.Printf("Climax Config: %+v\n", cfg.Climax)

	//server(&cfg)
	// cfg := climax.ClimaxConfig{BaseUrl: "http://192.168.1.187/"}
	// devices, _ := cfg.GetDevices()
	// //fmt.Printf("%+v\n", devices)

	// repository := climax.NewMemoryDeviceRepository()
	// for _, device := range devices {
	// 	repository.AddOrUpdate(device)
	// }
	// fmt.Printf("%+v\n", repository)

	// // deviceHistory, _ := cfg.GetDeviceHistory()
	// // fmt.Printf("%+v\n", deviceHistory)

	// //mqttCfg := mqttService.MqttConfig{BaseUrl: "broker.emqx.io", Port: 1883}
	// mqttCfg := mqttService.MqttConfig{BaseUrl: "192.168.1.142", Port: 1883, Username: "mqtt", Password: "mqtt"}
	// mqtt := mqttService.Connect(mqttCfg)
	// mqttService.Subscribe(mqtt, "zigbee2mqtt/#", nil)
	// mqttService.Publish(mqtt, "zigbee2mqtt/#", "Test message")

	// time.Sleep(time.Second)
	// mqtt.Disconnect(250)
}

func server(config *Config) {
	mqttCfg := mqttService.MqttConfig{BaseUrl: "192.168.1.142", Port: 1883, Username: "mqtt", Password: "mqtt"}
	mqttClient := mqttService.Connect(mqttCfg)
	repo := climax.NewMemoryDeviceRepository() // Initialize your device repository

	cfg := climax.ClimaxConfig{BaseUrl: "http://192.168.1.187/"}

	// Periodically fetch devices and publish updates
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		devices, err := cfg.GetDevices()
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
		publishIfNoError(dev.MqttDiscoveryMessageOnOff(), device, mqttClient)
		publishIfNoError(dev.MqttDiscoveryMessagePower(), device, mqttClient)
	default:
		log.Printf("Unsupported device type for device Id %s", dev.Identify())
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
