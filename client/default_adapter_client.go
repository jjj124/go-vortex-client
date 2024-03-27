package client

import (
	"bytes"
	"github.com/Vortex-ECO/Vortex-SDk-GO/msg"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	futures "github.com/jjj124/go-future"
	"github.com/jjj124/go-metrics"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type defaultAdapterClient struct {
	adapterOptions      AdapterOptions
	mqttClients         []mqtt.Client
	receivedMsgHandlers []ReceivedMsgHandler
	futures             *sync.Map
	adapterCache        AdapterCaches
	component           DefaultAdapterComponent
	vp                  *viper.Viper
}

func (d *defaultAdapterClient) Viper() *viper.Viper {
	return d.vp
}

func (d *defaultAdapterClient) Configurator() Configurator {
	return &configurator{
		component: d.component,
	}
}

func (d *defaultAdapterClient) Components() AdapterComponents {
	return &d.component
}

func (d *defaultAdapterClient) GwDeviceOf(hardwareId string) GatewayDeviceSelector {
	return NewGatewayOf(d, d.Pid(), hardwareId)
}
func (d *defaultAdapterClient) Shutdown() {

}
func (d *defaultAdapterClient) DeviceOf(hardwareId string) DeviceSelector {
	return MakeDeviceOf(d, d.Pid(), hardwareId)
}
func (d *defaultAdapterClient) Caches() AdapterCaches {
	return d.adapterCache
}
func (d *defaultAdapterClient) AdapterOperation() AdapterOperations {
	return NewAdapterOperations(d)
}
func (d *defaultAdapterClient) Delivery(m *msg.DeliveryMsg) futures.Future[msg.ReceivedMsg] {
	var mm = *m
	var mqttClient = m.MqttClient()
	if mqttClient == nil {
		var index = rand.Intn(len(d.mqttClients))
		mqttClient = d.mqttClients[index]
	}
	var ret = futures.NewFuture[msg.ReceivedMsg]()
	d.futures.Store(mm.MsgId(), ret)
	var bytes, err = mm.Marshal()
	if err != nil {

	} else {
		log.Println("send --> ", mm.ToString())
		go func() {
			var token = mqttClient.Publish("up/"+d.Pid(), 1, false, bytes)
			if token.Wait() {
				metrics.GetOrRegisterCounter("adapter.msg.send", d.component.MetricsRegistry()).Inc(1)
				d.component.RecentSendMsg().Push(m)
			} else {
				ret.CompleteExceptionally(errors.New("send msg fail!"))
			}
		}()
	}
	return ret
}
func (d *defaultAdapterClient) Pid() string {
	return d.adapterOptions.Pid()
}
func (d *defaultAdapterClient) ClientId() string {
	return d.adapterOptions.ClientId()
}
func (d *defaultAdapterClient) Start() futures.Future[string] {
	var waitGroup = sync.WaitGroup{}
	var ok = "ok"
	waitGroup.Add(d.adapterOptions.ConnectNum())
	for i := 0; i < d.adapterOptions.ConnectNum(); i++ {
		go d.connect(i, &waitGroup)
	}
	var connectFuture = futures.NewFuture[string]()
	go func() {
		waitGroup.Wait()
		connectFuture.Complete(&ok)
	}()
	var ret = futures.NewFuture[string]()
	connectFuture.WhenComplete(func(s *string, err error) {
		if err != nil {
			ret.CompleteExceptionally(err)
			return
		}
		d.AdapterOperation().AskConfigOperation().Execute().WhenComplete(func(a *AskConfigReply, err error) {
			if err != nil {
				ret.CompleteExceptionally(err)
			} else {
				var content = a.GetContent()
				if content != "" {
					var reader = bytes.NewReader([]byte(content))
					d.Viper().SetConfigType(a.GetFormat())
					var err = d.Viper().ReadConfig(reader)
					if err != nil {
						log.Fatalln(err)
					} else {
						log.Println("sync config success !")
						var keys = d.Viper().AllKeys()
						for _, key := range keys {
							log.Println(key, "=", d.Viper().Get(key))
						}
					}
				}
				ret.Complete(&ok)
			}
		})
	})
	return ret
}
func (d *defaultAdapterClient) connect(i int, waitGroup *sync.WaitGroup) {
	var url = "tcp://" + d.adapterOptions.Ip().String() + ":" + strconv.Itoa(d.adapterOptions.Port())
	var opts = mqtt.NewClientOptions().AddBroker(url)
	opts.Username = d.adapterOptions.Pid()
	opts.Password = d.adapterOptions.Secret()
	opts.AutoReconnect = true
	opts.CleanSession = true
	opts.ClientID = d.ClientId()
	opts.ConnectTimeout = 3 * time.Second
	opts.ConnectRetry = false
	opts.ConnectRetryInterval = 3 * time.Second
	opts.OnReconnecting = func(client mqtt.Client, options *mqtt.ClientOptions) {
		log.Println("connection lost ,reconnecting...")
	}
	var client = mqtt.NewClient(opts)
	var token = client.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Println("connect adapter gateway timeout , reconnect in 3000ms!")
		time.Sleep(3 * time.Second)
		go d.connect(i, waitGroup)
	} else {
		log.Printf("connect adapter gateway success [con_id=%d]!", i)
		token = client.Subscribe("down/"+d.Pid(), 1, newMsgHandler(d))
		token.Wait()
		d.mqttClients[i] = client
		waitGroup.Done()
	}
}
func (d *defaultAdapterClient) Options() AdapterOptions {
	return d.adapterOptions
}
func NewDefaultAdapterClient(adapterOptions AdapterOptions) AdapterClient {

	var receivedMsgHandlers = make([]ReceivedMsgHandler, 2)
	receivedMsgHandlers[0] = NewServiceInvokeMsgHandler()
	receivedMsgHandlers[1] = NewAskDeviceConnectStateHandler()

	var ret = &defaultAdapterClient{adapterOptions: adapterOptions,
		mqttClients:         make([]mqtt.Client, adapterOptions.ConnectNum()),
		receivedMsgHandlers: receivedMsgHandlers,
		futures:             &sync.Map{},
		adapterCache:        NewCaches(),
		component:           NewAdapterComponent(),
		vp:                  viper.New(),
	}

	return ret
}
func newMsgHandler(d *defaultAdapterClient) mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		var receivedMsg, err = msg.NewReceivedMsg(message.Payload(), client)
		if err != nil {

		} else {
			log.Println("recv <-- ", receivedMsg.ToString())
			go func() {
				for _, handler := range d.receivedMsgHandlers {
					handler(receivedMsg, d)
				}
			}()
			metrics.GetOrRegisterCounter("adapter.msg.recv", d.component.MetricsRegistry()).Inc(1)
			d.component.RecentReceivedMsg().Push(receivedMsg)
			var f, b = d.futures.LoadAndDelete(receivedMsg.MsgId())
			if b {
				var v, suc = f.(futures.Future[msg.ReceivedMsg])
				if suc {
					var err = receivedMsg.Error()
					if err != nil {
						var errMsg = (err["msg"]).(string)
						v.CompleteExceptionally(errors.New(errMsg))
						return
					} else {
						v.Complete(receivedMsg)
					}
				}
			}
		}
	}
}
