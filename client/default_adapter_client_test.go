package client

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net"
	"testing"
	"time"
)

func NewClientAndStart() AdapterClient {
	var options = NewAdapterOptions("brfcmmayllmwdon9", "test", net.ParseIP("127.0.0.1"), 10011, "a5efaa0de7600f1fb009adca17782fd7537a17d26cbfd0ffb7bbedfcd1eca8b1", 3)
	var client = NewDefaultAdapterClient(options)
	var _, err = client.Start().BlockingGet()
	if err != nil {
		panic(err)
	}
	return client
}

func TestStart(t *testing.T) {
	NewClientAndStart()
}

func TestPropReport(t *testing.T) {
	var adapter = NewClientAndStart()
	var _, _ = adapter.DeviceOf("gw_01").WithCreateIfNotExist(false).Op().ReportProps().WithValue("voltage", 220).WithValue("current", 0.25).NeedReply().Execute().BlockingGet()
	var _, _ = adapter.GwDeviceOf("gw_01").WithCreateIfNotExist(false).Op().ReportProps().WithValue("voltage", 220).WithValue("current", 0.25).NeedReply().Execute().BlockingGet()
}
func TestEventReport(t *testing.T) {
	var adapter = NewClientAndStart()
	var val = make(map[string]any)
	val["voltage"] = 250
	adapter.DeviceOf("gw_02").WithCreateIfNotExist(true).WithNameIfCreate("自动创建指定设备名称").Op().ReportEvents().NeedReply().WithValue("high_voltage_alarm", val).Execute().BlockingGet()
}

func TestService(t *testing.T) {
	var adapter = NewClientAndStart()
	adapter.Configurator().HandleServiceWith("read_config", func(device Device, params map[string]any, client AdapterClient) (bool, map[string]any, error) {
		var result = make(map[string]any)
		result["endpoint"] = "https://www.baidu.com"
		result["test_config"] = viper.GetString("application.name")
		return true, result, nil
	})
	time.Sleep(time.Hour * 1)
}

func TestPropReportSub(t *testing.T) {
	var adapter = NewClientAndStart()
	var reply, err = adapter.GwDeviceOf("gw_01").SubDeviceOf("thermometer", "sub_01").WithCreateIfNotExist(true).Op().ReportProps().NeedReply().WithValue("temperature", 23.4).Execute().BlockingGet()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}
}

func TestAskThingModel(t *testing.T) {
	var adapter = NewClientAndStart()
	var reply, err = adapter.AdapterOperation().AskThingModelOperation().Execute().BlockingGet()
	if err == nil {
		var props = reply.GetProps()
		for _, val := range props {
			log.Println(val.ToString())
		}
	}
}
