package common

import (
	"github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2019-08-01/documentdb"
)

func GetThroughputFromResult(throughputResponse documentdb.ThroughputSettingsGetResults) *int32 {
	props := throughputResponse.ThroughputSettingsGetProperties
	if props == nil {
		return nil
	}

	res := props.Resource
	if res == nil {
		return nil
	}

	return res.Throughput
}
