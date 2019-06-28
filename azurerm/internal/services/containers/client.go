package containers

import (
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-02-01/containerservice"
	"github.com/Azure/go-autorest/autorest"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/ar"
)

type Client struct {
	KubernetesClustersClient containerservice.ManagedClustersClient
	GroupsClient             containerinstance.ContainerGroupsClient
	RegistriesClient         containerregistry.RegistriesClient
	ReplicationsClient       containerregistry.ReplicationsClient
	ServicesClient           containerservice.ContainerServicesClient
}

func BuildClient(endpoint, subscriptionId, partnerId string, auth autorest.Authorizer, skipProviderReg bool) *Client {
	c := Client{}

	c.RegistriesClient = containerregistry.NewRegistriesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.RegistriesClient.Client, auth, partnerId, skipProviderReg)

	c.ReplicationsClient = containerregistry.NewReplicationsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ReplicationsClient.Client, auth, partnerId, skipProviderReg)

	c.GroupsClient = containerinstance.NewContainerGroupsClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.GroupsClient.Client, auth, partnerId, skipProviderReg)

	// ACS
	c.ServicesClient = containerservice.NewContainerServicesClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.ServicesClient.Client, auth, partnerId, skipProviderReg)

	// AKS
	c.KubernetesClustersClient = containerservice.NewManagedClustersClientWithBaseURI(endpoint, subscriptionId)
	ar.ConfigureClient(&c.KubernetesClustersClient.Client, auth, partnerId, skipProviderReg)

	return &c
}
