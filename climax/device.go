package climax

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DeviceType uint8

const (
	Undefined          DeviceType = 0
	Smoke_Detector     DeviceType = 11
	Temperature_Sensor DeviceType = 20
	Power_Switch       DeviceType = 24
	Power_Switch_Meter DeviceType = 48
	Room_Sensor        DeviceType = 54
	Hue_Sensor         DeviceType = 74
)

type DeviceInterface interface {
	Identify() string
	DeviceType() DeviceType
}

type Device struct {
	Id      string
	Zone    uint8
	Area    uint8
	Type    DeviceType
	Name    string
	Status  string
	Cond    uint8 `json:"cond_ok,string,omitempty"`
	Battery uint8 `json:"battery_ok,string,omitempty"`
	Rssi    string
}

type TemperatureSensor struct {
	Device
	Temperature float32
}

func (t Device) DeviceType() DeviceType {
	return t.Type
}

func (t Device) Identify() string {
	return strings.Replace(t.Id, ":", "", 1)
}

type PowerSwitchMeter struct {
	Device
	OnOff  bool
	Power  float64
	Energy float64
}

func parseTemperature(status string) float32 {
	var temperature float32
	fmt.Sscanf(status, "%f °C", &temperature)
	return temperature
}

func parsePowerSwitchMeterStatus(status string) (bool, float64, float64) {
	parts := strings.Split(status, ",")
	if len(parts) < 3 {
		return false, 0.0, 0.0
	}

	onOff := strings.TrimSpace(parts[0]) == "On"

	powerStr := strings.TrimSpace(strings.TrimSuffix(parts[1], "W"))
	power, err := strconv.ParseFloat(powerStr, 64)
	if err != nil {
		power = 0.0
	}

	energyStr := strings.TrimSpace(strings.TrimSuffix(parts[2], "kWh"))
	energy, err := strconv.ParseFloat(energyStr, 64)
	if err != nil {
		energy = 0.0
	}

	return onOff, power, energy
}

type KeyType string

const (
	ActivePowerKey KeyType = "Active Power"
	EnergyKey      KeyType = "Energy"
	TemperatureKey KeyType = "Temperature"
)

type DeviceHistory struct {
	Id       string `json:"device_id"`
	DateTime CustomTime
	Area     string
	Zone     string
	Name     string
	Key      KeyType
	Value    string
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	loc, _ := time.LoadLocation("Europe/Stockholm")
	date, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(b), loc)
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (dv DeviceHistory) ToDevice(t DeviceType) Device {
	zone, _ := strconv.ParseInt(dv.Zone, 10, 8)
	area, _ := strconv.ParseInt(dv.Zone, 10, 8)
	device := Device{Id: dv.Id, Zone: uint8(zone), Area: uint8(area), Type: t, Name: dv.Name, Status: string(dv.Key) + " " + dv.Value}

	return device
}
