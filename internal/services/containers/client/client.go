package client

import (
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2021-10-01/containerinstance"
	containerregistry_v2019_06_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2019-06-01-preview"
	containerregistry_v2021_08_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2019-08-01/containerservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/agentpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/maintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-02-02-preview/managedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AgentPoolsClient                            *agentpools.AgentPoolsClient
	ContainerInstanceClient                     *containerinstance.ContainerInstanceClient
	ContainerRegistryClient_v2021_08_01_preview *containerregistry_v2021_08_01_preview.Client
	// v2019_06_01_preview is needed for container registry agent pools and tasks
	ContainerRegistryClient_v2019_06_01_preview *containerregistry_v2019_06_01_preview.Client
	KubernetesClustersClient                    *managedclusters.ManagedClustersClient
	KubernetesExtensionsClient                  *extensions.ExtensionsClient
	MaintenanceConfigurationsClient             *maintenanceconfigurations.MaintenanceConfigurationsClient
	ServicesClient                              *containerservices.ContainerServicesClient
	Environment                                 azure.Environment
}

func NewContainersClient(o *common.ClientOptions) (*Client, error) {
	containerInstanceClient := containerinstance.NewContainerInstanceClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&containerInstanceClient.Client, o.ResourceManagerAuthorizer)

	containerRegistryClient_v2019_06_01_preview, err := containerregistry_v2019_06_01_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}

	containerRegistryClient_v2021_08_01_preview, err := containerregistry_v2021_08_01_preview.NewClientWithBaseURI(o.Environment.ResourceManager, func(c *resourcemanager.Client) {
		o.Configure(c, o.Authorizers.ResourceManager)
	})
	if err != nil {
		return nil, err
	}

	// AKS
	kubernetesClustersClient := managedclusters.NewManagedClustersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&kubernetesClustersClient.Client, o.ResourceManagerAuthorizer)

	kubernetesExtensionsClient := extensions.NewExtensionsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&kubernetesExtensionsClient.Client, o.ResourceManagerAuthorizer)

	agentPoolsClient := agentpools.NewAgentPoolsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&agentPoolsClient.Client, o.ResourceManagerAuthorizer)

	maintenanceConfigurationsClient := maintenanceconfigurations.NewMaintenanceConfigurationsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&maintenanceConfigurationsClient.Client, o.ResourceManagerAuthorizer)

	servicesClient := containerservices.NewContainerServicesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&servicesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AgentPoolsClient:                            &agentPoolsClient,
		ContainerInstanceClient:                     &containerInstanceClient,
		ContainerRegistryClient_v2021_08_01_preview: containerRegistryClient_v2021_08_01_preview,
		ContainerRegistryClient_v2019_06_01_preview: containerRegistryClient_v2019_06_01_preview,
		KubernetesClustersClient:                    &kubernetesClustersClient,
		KubernetesExtensionsClient:                  &kubernetesExtensionsClient,
		MaintenanceConfigurationsClient:             &maintenanceConfigurationsClient,
		ServicesClient:                              &servicesClient,
		Environment:                                 o.AzureEnvironment,
	}, nil
}
