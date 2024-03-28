package climax

type DeviceRepository struct {
	devices map[string]DeviceInterface
}

func NewMemoryDeviceRepository() DeviceRepository {
	return DeviceRepository{
		devices: make(map[string]DeviceInterface),
	}
}

func (d *DeviceRepository) IsNewDevice(deviceID string) bool {
	_, exists := d.devices[deviceID]
	return !exists
}

func (d DeviceRepository) GetDevice(id string) (DeviceInterface, bool) {
	lastValue, exist := d.devices[id]
	return lastValue, exist
}

func (d DeviceRepository) AddOrUpdate(device DeviceInterface) bool {
	deviceId := device.Identify()
	lastValue, exist := d.devices[deviceId]

	if !exist || lastValue.State() != device.State() {
		d.devices[deviceId] = device
		return true
	}
	return false
}
