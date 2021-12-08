package application

type ServiceTypeHealthPolicy struct {
	MaxPercentUnhealthyPartitionsPerService int64 `json:"maxPercentUnhealthyPartitionsPerService"`
	MaxPercentUnhealthyReplicasPerPartition int64 `json:"maxPercentUnhealthyReplicasPerPartition"`
	MaxPercentUnhealthyServices             int64 `json:"maxPercentUnhealthyServices"`
}
