package climax

import (
	"encoding/json"
	"testing"
)

func TestDeviceJson(t *testing.T) {
	rawJson := []byte(`{"area": 1, "zone": 6, "type": 20, "type_f": "Temperature Sensor", "name": "Kök", "cond": "", "cond_ok": "1", "battery": "", "battery_ok": "1", "tamper": "", "tamper_ok": "1", "bypass": "No", "temp_bypass": "0", "rssi": "Weak, 3", "status": "6.68 °C", "id": "ZB:445001", "su": 1}`)

	var device Device

	if err := json.Unmarshal(rawJson, &device); err != nil {
		t.Error(err)
	}

	t.Errorf("%+v", device)
}

func TestDeviceHistoryJson(t *testing.T) {
	rawJson := []byte(`{"datetime": "2024-03-26 08:33:46", "area": "1", "zone": "3", "name": "", "device_id": "ZB:b02a01", "key": "Active Power", "value": "0.0W" }`)

	var deviceHistory DeviceHistory

	if err := json.Unmarshal(rawJson, &deviceHistory); err != nil {
		t.Error(err)
	}

	t.Errorf("%+v", deviceHistory)
}
