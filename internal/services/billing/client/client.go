package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2019-10-01-preview/enrollmentaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2020-05-01/billingaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	BillingAccountsClient           *billingaccounts.BillingAccountsClient
	BillingEnrollmentAccountsClient *enrollmentaccounts.EnrollmentAccountsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	billingAccountsClient, err := billingaccounts.NewBillingAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BillingAccounts client: %+v", err)
	}
	o.Configure(billingAccountsClient.Client, o.Authorizers.ResourceManager)

	billingEnrollmentAccountsClient, err := enrollmentaccounts.NewEnrollmentAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building BillingEnrollmentAccounts client: %+v", err)
	}
	o.Configure(billingEnrollmentAccountsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		BillingAccountsClient:           billingAccountsClient,
		BillingEnrollmentAccountsClient: billingEnrollmentAccountsClient,
	}, nil
}
