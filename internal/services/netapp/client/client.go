package client

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/capacitypools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/netappaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/snapshotpolicy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/snapshots"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2021-10-01/volumesreplication"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	AccountClient           *netappaccounts.NetAppAccountsClient
	PoolClient              *capacitypools.CapacityPoolsClient
	VolumeClient            *volumes.VolumesClient
	VolumeReplicationClient *volumesreplication.VolumesReplicationClient
	SnapshotClient          *snapshots.SnapshotsClient
	SnapshotPoliciesClient  *snapshotpolicy.SnapshotPolicyClient
}

func NewClient(o *common.ClientOptions) *Client {
	accountClient := netappaccounts.NewNetAppAccountsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&accountClient.Client, o.ResourceManagerAuthorizer)

	poolClient := capacitypools.NewCapacityPoolsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&poolClient.Client, o.ResourceManagerAuthorizer)

	volumeClient := volumes.NewVolumesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&volumeClient.Client, o.ResourceManagerAuthorizer)

	volumeReplicationClient := volumesreplication.NewVolumesReplicationClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&volumeReplicationClient.Client, o.ResourceManagerAuthorizer)

	snapshotClient := snapshots.NewSnapshotsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&snapshotClient.Client, o.ResourceManagerAuthorizer)

	snapshotPoliciesClient := snapshotpolicy.NewSnapshotPolicyClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&snapshotPoliciesClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		AccountClient:           &accountClient,
		PoolClient:              &poolClient,
		VolumeClient:            &volumeClient,
		VolumeReplicationClient: &volumeReplicationClient,
		SnapshotClient:          &snapshotClient,
		SnapshotPoliciesClient:  &snapshotPoliciesClient,
	}
}
