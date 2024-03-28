package climax

import (
	"encoding/json"
	"fmt"
)

type MqttMessage struct {
	Topic   string
	Message []byte
	Err     error
}

func (ts TemperatureSensor) MqttDiscoveryMessageTemperature() MqttMessage {
	id := ts.Identify()
	topic := fmt.Sprintf("homeassistant/sensor/%s/temperature/config", id)
	payload := map[string]interface{}{
		"unique_id":           fmt.Sprintf("%s_temperature", id),
		"state_topic":         fmt.Sprintf("climax2mqtt/sensors/%s/state", id),
		"name":                fmt.Sprintf("%s Temperature", ts.Name),
		"device_class":        "temperature",
		"unit_of_measurement": "°C",
		"value_template":      "{{ value_json.temperature }}",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return MqttMessage{topic, nil, err}
	}
	return MqttMessage{topic, jsonData, nil}
}

func (psm PowerSwitchMeter) MqttDiscoveryMessagePower() MqttMessage {
	id := psm.Identify()
	topic := fmt.Sprintf("homeassistant/sensor/%s/power/config", id)
	payload := map[string]interface{}{
		"unique_id":           fmt.Sprintf("%s_power", id),
		"state_topic":         fmt.Sprintf("climax2mqtt/sensors/%s/state", id),
		"name":                fmt.Sprintf("%s Power", psm.Name),
		"device_class":        "power",
		"unit_of_measurement": "W",
		"value_template":      "{{ value_json.power }}",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return MqttMessage{topic, nil, err}
	}
	return MqttMessage{topic, jsonData, nil}
}

func (psm PowerSwitchMeter) MqttDiscoveryMessageOnOff() MqttMessage {
	id := psm.Identify()
	topic := fmt.Sprintf("homeassistant/binary_sensor/%s/onoff/config", id)
	payload := map[string]interface{}{
		"unique_id":      fmt.Sprintf("%s_onoff", id),
		"state_topic":    fmt.Sprintf("climax2mqtt/sensors/%s/state", id),
		"name":           fmt.Sprintf("%s On/Off", psm.Name),
		"device_class":   "power", // Optional
		"payload_on":     "ON",
		"payload_off":    "OFF",
		"value_template": "{{ value_json.on_off }}",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return MqttMessage{"", nil, err}
	}
	return MqttMessage{topic, jsonData, nil}
}

func (psm PowerSwitchMeter) MqttDiscoveryMessageEnergy() MqttMessage {
	id := psm.Identify()
	topic := fmt.Sprintf("homeassistant/sensor/%s/energy/config", id)
	payload := map[string]interface{}{
		"unique_id":           fmt.Sprintf("%s_energy", id),
		"state_topic":         fmt.Sprintf("climax2mqtt/sensors/%s/state", id),
		"name":                fmt.Sprintf("%s Energy Usage", psm.Name),
		"device_class":        "energy",
		"unit_of_measurement": "kWh",
		"value_template":      "{{ value_json.energy }}",
		"icon":                "mdi:counter",
	}

	jsonData, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return MqttMessage{"", nil, err}
	}
	return MqttMessage{topic, jsonData, nil}
}

func (ts TemperatureSensor) MqttUpdateValueMessage() MqttMessage {
	id := ts.Identify()
	topic := fmt.Sprintf("climax2mqtt/sensors/%s/state", id)

	// Payload structure reflecting the current state/value
	payload := map[string]interface{}{
		"temperature": ts.Temperature,
	}

	// Serialize the payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return MqttMessage{topic, nil, fmt.Errorf("error serializing temperature update to JSON: %w", err)}
	}

	return MqttMessage{topic, jsonData, nil}
}

func (psm PowerSwitchMeter) MqttUpdateValueMessage() MqttMessage {
	id := psm.Identify()
	topic := fmt.Sprintf("climax2mqtt/sensors/%s/state", id)

	onOffState := "OFF"
	if psm.OnOff {
		onOffState = "ON"
	}
	payload := map[string]interface{}{
		"on_off": onOffState,
		"power":  psm.Power,
		"energy": psm.Energy,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return MqttMessage{"", nil, fmt.Errorf("error serializing power switch meter state update to JSON: %w", err)}
	}

	return MqttMessage{topic, jsonData, nil}
}
