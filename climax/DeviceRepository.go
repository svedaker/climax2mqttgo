package climax

type DeviceRepository struct {
	devices map[string]DeviceInterface
}

func NewMemoryDeviceRepository() DeviceRepository {
	return DeviceRepository{
		devices: make(map[string]DeviceInterface),
	}
}

func (d DeviceRepository) GetDevice(id string) DeviceInterface {
	return d.devices[id]
}

func (d DeviceRepository) AddOrUpdate(device DeviceInterface) {
	d.devices[device.Identify()] = device
}
