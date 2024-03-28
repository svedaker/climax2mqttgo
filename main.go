package main

import (
	"climax/climax"
	"fmt"
)

func main() {
	cfg := climax.ClimaxConfig{BaseUrl: "http://192.168.1.187/"}
	devices, _ := cfg.GetDevices()
	//fmt.Printf("%+v\n", devices)

	repository := climax.NewMemoryDeviceRepository()
	for _, device := range devices {
		repository.AddOrUpdate(device)
	}
	fmt.Printf("%+v\n", repository)

	// deviceHistory, _ := cfg.GetDeviceHistory()
	// fmt.Printf("%+v\n", deviceHistory)

	//mqttCfg := mqttService.MqttConfig{BaseUrl: "broker.emqx.io", Port: 1883}
	// mqttCfg := mqttService.MqttConfig{BaseUrl: "192.168.1.142", Port: 1883, Username: "mqtt", Password: "mqtt"}
	// mqtt := mqttService.Connect(mqttCfg)
	// mqttService.Subscribe(mqtt, "zigbee2mqtt/#", nil)
	// //mqttService.Publish(mqtt, "zigbee2mqtt/#", "Test message")

	// time.Sleep(time.Second)
	// mqtt.Disconnect(250)
}
