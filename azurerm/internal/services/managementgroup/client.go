package managementgroup

import "github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2018-03-01-preview/managementgroups"

type Client struct {
	GroupsClient       managementgroups.Client
	SubscriptionClient managementgroups.SubscriptionsClient
}
