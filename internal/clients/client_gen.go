package clients

// NOTE: this file is generated - manual changes will be overwritten.

import (
	containerservice_v2022_09_02_preview "github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2022-09-02-preview"
	loadtestservice_v2021_12_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview"
	managedidentity_v2022_01_31_preview "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2022-01-31-preview"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	containers "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	loadtestservice "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice/client"
	managedidentity "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedidentity/client"
)

type autoClient struct {
	ContainerService *containerservice_v2022_09_02_preview.Client
	LoadTestService  *loadtestservice_v2021_12_01_preview.Client
	ManagedIdentity  *managedidentity_v2022_01_31_preview.Client
}

func buildAutoClients(client *autoClient, o *common.ClientOptions) error {
	client.ContainerService = containers.NewClient(o)
	client.LoadTestService = loadtestservice.NewClient(o)
	client.ManagedIdentity = managedidentity.NewClient(o)
	return nil
}
