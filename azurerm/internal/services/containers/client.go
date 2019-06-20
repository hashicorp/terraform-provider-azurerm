package containers

import (
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-02-01/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	KubernetesClustersClient containerservice.ManagedClustersClient
	GroupsClient             containerinstance.ContainerGroupsClient
	RegistriesClient         containerregistry.RegistriesClient
	ReplicationsClient       containerregistry.ReplicationsClient
	ServicesClient           containerservice.ContainerServicesClient
}

func BuildClient(endpoint string, authorizer autorest.Authorizer, o *common.ClientOptions) *Client {
	c := Client{}

	c.RegistriesClient = containerregistry.NewRegistriesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RegistriesClient.Client, authorizer)

	c.ReplicationsClient = containerregistry.NewReplicationsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ReplicationsClient.Client, authorizer)

	c.GroupsClient = containerinstance.NewContainerGroupsClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GroupsClient.Client, authorizer)

	// ACS
	c.ServicesClient = containerservice.NewContainerServicesClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServicesClient.Client, authorizer)

	// AKS
	c.KubernetesClustersClient = containerservice.NewManagedClustersClientWithBaseURI(endpoint, o.SubscriptionId)
	o.ConfigureClient(&c.KubernetesClustersClient.Client, authorizer)

	return &c
}
