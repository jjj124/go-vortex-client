package client

type AdapterOperations interface {
	CreateDeviceOperation() CreateDeviceOperation
	AskPidByNickOperation() AskPidByNickOperation
	AskConfigOperation() AskProtocolConfigOperation
	AskThingModelOperation() AskThingModelOperation
}

type adapterOperations struct {
	adapter AdapterClient
}

func (d *adapterOperations) AskThingModelOperation() AskThingModelOperation {
	return NewAskThingModelOperation(d.adapter)
}

func (d *adapterOperations) AskConfigOperation() AskProtocolConfigOperation {
	return NewAskConfigOperation(d.adapter)
}

func (d *adapterOperations) AskPidByNickOperation() AskPidByNickOperation {
	return NewAskPidByNickOperation(d.adapter)
}

func (d *adapterOperations) CreateDeviceOperation() CreateDeviceOperation {
	return NewCreateDeviceOperation(d.adapter)
}
func NewAdapterOperations(adapter AdapterClient) AdapterOperations {
	var ret = &adapterOperations{
		adapter: adapter,
	}
	return ret
}
