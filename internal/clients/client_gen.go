package clients

// NOTE: this file is generated - manual changes will be overwritten.

import (
	loadtests_v2021_12_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview"
  managedidentity_2018_11_30 "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	loadtest "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice/client"
)

type autoClient struct {
	LoadTestService *loadtests_v2021_12_01_preview.Client
	ManagedIdentity *managedidentity_2018_11_30.Client
}

func buildAutoClients(client *autoClient, o *common.ClientOptions) error {
	client.LoadTestService = loadtest.NewClient(o)
	client.ManagedIdentity = managedidentity.NewClient(o)
	return nil
}
