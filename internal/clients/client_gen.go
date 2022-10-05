package clients

// NOTE: this file is generated - manual changes will be overwritten.

import (
	managedidentity_2018_11_30 "github.com/hashicorp/go-azure-sdk/resource-manager/managedidentity/2018-11-30"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	managedidentity "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedidentity/client"
)

type autoClient struct {
	ManagedIdentity *managedidentity_2018_11_30.Client
}

func buildAutoClients(client *autoClient, o *common.ClientOptions) error {
	client.ManagedIdentity = managedidentity.NewClient(o)
	return nil
}
