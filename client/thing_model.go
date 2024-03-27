package client

type ThingModel struct {
	Name     string
	DataType string
	Symbol   string
}

func (a *ThingModel) ToString() string {
	return "[name=" + a.Name + ",data_type=" + a.DataType + ",symbol=" + a.Symbol + "]"
}
