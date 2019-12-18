package client

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/sqlvirtualmachine/mgmt/2017-03-01-preview/sqlvirtualmachine"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/common"
)

type Client struct {
	SQLVirtualMachinesClient                                 *sqlvirtualmachine.SQLVirtualMachinesClient
}

func NewClient(o *common.ClientOptions) *Client {
	SQLVirtualMachinesClient := sqlvirtualmachine.NewSQLVirtualMachinesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&SQLVirtualMachinesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		SQLVirtualMachinesClient: &SQLVirtualMachinesClient,
	}
}
