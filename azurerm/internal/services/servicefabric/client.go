package servicefabric

import "github.com/Azure/azure-sdk-for-go/services/servicefabric/mgmt/2018-02-01/servicefabric"

type Client struct {
	ClustersClient servicefabric.ClustersClient
}
