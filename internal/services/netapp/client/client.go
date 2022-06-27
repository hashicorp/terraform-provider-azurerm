package client

import (
	"github.com/Azure/azure-sdk-for-go/services/netapp/mgmt/2021-10-01/netapp"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient          *netapp.AccountsClient
	PoolClient             *netapp.PoolsClient
	VolumeClient           *netapp.VolumesClient
	SnapshotClient         *netapp.SnapshotsClient
	SnapshotPoliciesClient *netapp.SnapshotPoliciesClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := netapp.NewAccountsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	poolClient := netapp.NewPoolsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&poolClient.Client, o.ResourceManagerAuthorizer)

	volumeClient := netapp.NewVolumesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&volumeClient.Client, o.ResourceManagerAuthorizer)

	snapshotClient := netapp.NewSnapshotsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&snapshotClient.Client, o.ResourceManagerAuthorizer)

	snapshotPoliciesClient := netapp.NewSnapshotPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&snapshotPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:          &accountClient,
		PoolClient:             &poolClient,
		VolumeClient:           &volumeClient,
		SnapshotClient:         &snapshotClient,
		SnapshotPoliciesClient: &snapshotPoliciesClient,
	}
}
