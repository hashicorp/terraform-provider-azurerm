package hdinsight

import "github.com/Azure/azure-sdk-for-go/services/preview/hdinsight/mgmt/2018-06-01-preview/hdinsight"

type Client struct {
	ApplicationsClient   hdinsight.ApplicationsClient
	ClustersClient       hdinsight.ClustersClient
	ConfigurationsClient hdinsight.ConfigurationsClient
}
