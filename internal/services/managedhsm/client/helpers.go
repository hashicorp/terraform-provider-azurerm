package client

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-07-01/managedhsms"
)

// TODO: AddToCache and a Purge too

func (c *Client) BaseUriForManagedHSM(ctx context.Context, managedHsmId managedhsms.ManagedHSMId) (*string, error) {
	// TODO: implement me
	return nil, nil
}

func (c *Client) ManagedHSMIDFromBaseUrl(ctx context.Context, subscriptionId commonids.SubscriptionId, managedHsmBaseUrl string) (*managedhsms.ManagedHSMId, error) {
	// TODO: implement me
	return nil, nil
}
