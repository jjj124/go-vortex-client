package client

import (
	futures "github.com/jjj124/go-future"
)

type DeviceOperations interface {
	ReportProps() PropReportOperation
	ReportEvents() EventReportOperation
}

type deviceOperations struct {
	device  futures.Future[Device]
	adapter AdapterClient
}

func (d *deviceOperations) ReportEvents() EventReportOperation {
	return NewEventReportOperation(d.adapter, d.device)
}

func (d *deviceOperations) ReportProps() PropReportOperation {
	return NewPropReportOperation(d.adapter, d.device)
}

func NewDeviceOperations(adapter AdapterClient, device futures.Future[Device]) DeviceOperations {
	return &deviceOperations{
		device,
		adapter,
	}
}
