package trafficmanager

import "github.com/Azure/azure-sdk-for-go/services/trafficmanager/mgmt/2018-04-01/trafficmanager"

type Client struct {
	GeographialHierarchiesClient trafficmanager.GeographicHierarchiesClient
	ProfilesClient               trafficmanager.ProfilesClient
	EndpointsClient              trafficmanager.EndpointsClient
}
