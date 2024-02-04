package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2020-05-01/billingaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BillingAccountsClient *billingaccounts.BillingAccountsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	billingaccountsClient, err := billingaccounts.NewBillingAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BillingAccounts client: %+v", err)
	}
	o.Configure(billingaccountsClient.Client, o.Authorizers.ResourceManager)
	
	return &Client{
		BillingAccountsClient: billingaccountsClient,
	}, nil
}