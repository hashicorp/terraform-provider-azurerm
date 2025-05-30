package helper

type DataFlowType string

const (
	DataFlowTypeFlowlet           = "Flowlet"
	DataFlowTypeMappingDataFlow   = "MappingDataFlow"
	DataFlowTypeWranglingDataFlow = "WranglingDataFlow"
)
