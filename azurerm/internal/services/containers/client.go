package containers

import (
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2018-09-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-06-01/containerservice"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	KubernetesClustersClient containerservice.ManagedClustersClient
	GroupsClient             containerinstance.ContainerGroupsClient
	RegistriesClient         containerregistry.RegistriesClient
	ReplicationsClient       containerregistry.ReplicationsClient
	ServicesClient           containerservice.ContainerServicesClient
}

func BuildClient(o *common.ClientOptions) *Client {
	c := Client{}

	c.RegistriesClient = containerregistry.NewRegistriesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.RegistriesClient.Client, o.ResourceManagerAuthorizer)

	c.ReplicationsClient = containerregistry.NewReplicationsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ReplicationsClient.Client, o.ResourceManagerAuthorizer)

	c.GroupsClient = containerinstance.NewContainerGroupsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.GroupsClient.Client, o.ResourceManagerAuthorizer)

	// ACS
	c.ServicesClient = containerservice.NewContainerServicesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.ServicesClient.Client, o.ResourceManagerAuthorizer)

	// AKS
	c.KubernetesClustersClient = containerservice.NewManagedClustersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&c.KubernetesClustersClient.Client, o.ResourceManagerAuthorizer)

	return &c
}
