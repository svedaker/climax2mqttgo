package mqttService

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	BaseUrl  string `yaml:"base_url" env:"MQTT_BASEURL"`
	Port     int    `yaml:"port" env:"MQTT_PORT" env-default:"1883"`
	Username string `yaml:"username" env:"MQTT_USERNAME" env-default:"mqtt"`
	Password string `yaml:"password" env:"MQTT_PASSWORD" env-default:"mqtt"`
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("mqtt: Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("mqtt: Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("mqtt: Connect lost: %v", err)
}

func Connect(cfg MqttConfig) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.BaseUrl, cfg.Port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func Subscribe(client mqtt.Client, topic string, callback mqtt.MessageHandler) {
	token := client.Subscribe(topic, 1, callback)
	token.Wait()
	log.Printf("mqtt: Subscribed to topic %s", topic)
}

func Publish(client mqtt.Client, topic string, message interface{}) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
}
