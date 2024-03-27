package client

type Device interface {
	Pid() string
	Did() string
	HardwareId() string
}
type SubDevice interface {
	Device
	GwId() string
}
type defaultDevice struct {
	pid        string
	did        string
	hardwareId string
}
type defaultSubDevice struct {
	defaultDevice
	gwId string
}

func (d *defaultDevice) Pid() string {
	return d.pid
}
func (d *defaultDevice) Did() string {
	return d.did
}
func (d *defaultDevice) HardwareId() string {
	return d.hardwareId
}
func (d *defaultSubDevice) GwId() string {
	return d.gwId
}
func NewDevice(pid string, did string, hardwareId string) Device {
	return &defaultDevice{
		pid, did, hardwareId,
	}
}
func NewSubDevice(pid string, did string, hardwareId string, gwId string) SubDevice {
	return &defaultSubDevice{
		defaultDevice: defaultDevice{
			pid:        pid,
			did:        did,
			hardwareId: hardwareId,
		},
		gwId: gwId,
	}
}
