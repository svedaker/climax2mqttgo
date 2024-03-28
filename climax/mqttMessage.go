package climax

import (
	"encoding/json"
	"fmt"
)

// type Config struct {
// 	StateTopic        string `json:"state_topic"`
// 	DeviceClass       string `json:"device_class"`
// 	Name              string `json:"name"`
// 	UnitOfMeasurement string `json:"unit_of_measurement"`
// 	ValueTemplate     string `json:"value_template"`
// }

func (ts TemperatureSensor) MqttDiscoveryMessageTemperature() (string, []byte, error) {
	topic := fmt.Sprintf("homeassistant/sensor/%s/temperature/config", ts.Id)
	payload := map[string]interface{}{
		"unique_id":           fmt.Sprintf("%s_temperature", ts.Id),
		"state_topic":         fmt.Sprintf("climax2mqtt/sensors/%s/state", ts.Id),
		"name":                fmt.Sprintf("%s Temperature", ts.Name),
		"device_class":        "temperature",
		"unit_of_measurement": "°C",
		"value_template":      "{{ value_json.temperature }}",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return topic, nil, err
	}
	return topic, jsonData, nil
}

func (psm PowerSwitchMeter) MqttDiscoveryMessagePower() (string, []byte, error) {
	topic := fmt.Sprintf("homeassistant/sensor/%s/power/config", psm.Id)
	payload := map[string]interface{}{
		"unique_id":           fmt.Sprintf("%s_power", psm.Id),
		"state_topic":         fmt.Sprintf("climax2mqtt/sensors/%s/state", psm.Id),
		"name":                fmt.Sprintf("%s Power", psm.Name),
		"device_class":        "power",
		"unit_of_measurement": "W",
		"value_template":      "{{ value_json.power }}",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return topic, nil, err
	}
	return topic, jsonData, nil
}

func (psm PowerSwitchMeter) MqttDiscoveryMessageOnOff() (string, []byte, error) {
	topic := fmt.Sprintf("homeassistant/binary_sensor/%s/onoff/config", psm.Id)
	payload := map[string]interface{}{
		"unique_id":      fmt.Sprintf("%s_onoff", psm.Id),
		"state_topic":    fmt.Sprintf("climax2mqtt/sensors/%s/state", psm.Id),
		"name":           fmt.Sprintf("%s On/Off", psm.Name),
		"device_class":   "power", // Optional
		"payload_on":     "ON",
		"payload_off":    "OFF",
		"value_template": "{{ value_json.on_off }}",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return "", nil, err
	}
	return topic, jsonData, nil
}

func (psm PowerSwitchMeter) MqttDiscoveryMessageEnergy() (string, []byte, error) {
	topic := fmt.Sprintf("homeassistant/sensor/%s/energy/config", psm.Id)
	payload := map[string]interface{}{
		"unique_id":           fmt.Sprintf("%s_energy", psm.Id),
		"state_topic":         fmt.Sprintf("climax2mqtt/sensors/%s/state", psm.Id),
		"name":                fmt.Sprintf("%s Energy Usage", psm.Name),
		"device_class":        "energy",
		"unit_of_measurement": "kWh",
		"value_template":      "{{ value_json.energy }}",
		"icon":                "mdi:counter",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return "", nil, err
	}
	return topic, jsonData, nil
}

func (ts TemperatureSensor) MqttUpdateValueMessage() (string, []byte, error) {
	// Topic for publishing temperature updates
	topic := fmt.Sprintf("climax2mqtt/sensors/%s/state", ts.Id)

	// Payload structure reflecting the current state/value
	// Adjust this structure to match your actual sensor data format and Home Assistant configuration
	payload := map[string]interface{}{
		"temperature": ts.Temperature,
	}

	// Serialize the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return topic, nil, fmt.Errorf("error serializing temperature update to JSON: %w", err)
	}

	return topic, jsonData, nil
}
