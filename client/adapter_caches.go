package client

import (
	"github.com/dgraph-io/ristretto"
	"sync"
)

type DeviceCache interface {
	Pid() string
	PutDevice(device Device)
	GetDeviceByDid(did string) Device
	GetDeviceByHardwareId(did string) Device
	InvalidateDeviceCache(did string)
}

type deviceCache struct {
	pid             string
	didAsKey        *ristretto.Cache
	hardwareIdAsKey *ristretto.Cache
}

func (d *deviceCache) Pid() string {
	return d.pid
}

func (d *deviceCache) PutDevice(device Device) {
	d.didAsKey.Set(device.Did(), device, 1)
	d.hardwareIdAsKey.Set(device.HardwareId(), device, 1)
}

func (d *deviceCache) GetDeviceByDid(did string) Device {
	var ret, found = d.didAsKey.Get(did)
	if found {
		return ret.(Device)
	}
	return nil
}

func (d *deviceCache) GetDeviceByHardwareId(did string) Device {
	var ret, found = d.hardwareIdAsKey.Get(did)
	if found {
		return ret.(Device)
	}
	return nil
}

func (d *deviceCache) InvalidateDeviceCache(did string) {
	var dev = d.GetDeviceByDid(did)
	if dev != nil {
		d.didAsKey.Del(dev.Did())
		d.hardwareIdAsKey.Del(dev.HardwareId())
	}
}

func NewDeviceCache(pid string) DeviceCache {
	ch1, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1024000,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	ch2, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1024000,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	var ret = deviceCache{pid, ch1, ch2}
	return &ret
}

type AdapterCaches interface {
	GetDeviceCacheByPid(pid string) DeviceCache
	CachePidNick(nick string, pid string)
	GetPidByNick(nick string) string
}
type adapterCaches struct {
	devCacheMap *sync.Map
	locker      sync.Locker
	pidNick     *sync.Map
}

func (d *adapterCaches) CachePidNick(nick string, pid string) {
	d.pidNick.Store(nick, pid)
}

func (d *adapterCaches) GetPidByNick(nick string) string {
	var val, ok = d.pidNick.Load(nick)
	if ok {
		return (val).(string)
	}
	return ""
}

func (d *adapterCaches) GetDeviceCacheByPid(pid string) DeviceCache {
	var v, b = d.devCacheMap.Load(pid)
	if b {
		return v.(DeviceCache)
	}
	d.locker.Lock()
	defer d.locker.Unlock()
	var v1, b1 = d.devCacheMap.Load(pid)
	if b1 {
		return v1.(DeviceCache)
	}
	var c = NewDeviceCache(pid)
	d.devCacheMap.Store(pid, c)
	return c
}

func NewCaches() AdapterCaches {
	return &adapterCaches{
		devCacheMap: &sync.Map{},
		locker:      &sync.Mutex{},
		pidNick:     &sync.Map{},
	}
}
