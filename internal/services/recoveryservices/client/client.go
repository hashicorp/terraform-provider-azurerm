package client

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-10-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupresourcestorageconfigsnoncrr"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupresourcevaultconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworkmappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	ProtectableItemsClient                    *backup.ProtectableItemsClient
	ProtectedItemsClient                      *backup.ProtectedItemsClient
	ProtectedItemsGroupClient                 *backup.ProtectedItemsGroupClient
	ProtectionPoliciesClient                  *backup.ProtectionPoliciesClient
	ProtectionContainerOperationResultsClient *backup.ProtectionContainerOperationResultsClient
	BackupProtectionContainersClient          *backup.ProtectionContainersClient
	BackupOperationStatusesClient             *backup.OperationStatusesClient
	VaultsClient                              *vaults.VaultsClient
	VaultsConfigsClient                       *backupresourcevaultconfigs.BackupResourceVaultConfigsClient
	StorageConfigsClient                      *backupresourcestorageconfigsnoncrr.BackupResourceStorageConfigsNonCRRClient
	FabricClient                              *replicationfabrics.ReplicationFabricsClient
	ProtectionContainerClient                 *replicationprotectioncontainers.ReplicationProtectionContainersClient
	ReplicationPoliciesClient                 *replicationpolicies.ReplicationPoliciesClient
	ContainerMappingClient                    *replicationprotectioncontainermappings.ReplicationProtectionContainerMappingsClient
	NetworkMappingClient                      *replicationnetworkmappings.ReplicationNetworkMappingsClient
	ReplicationProtectedItemsClient           *replicationprotecteditems.ReplicationProtectedItemsClient
	ReplicationRecoveryPlansClient            *replicationrecoveryplans.ReplicationRecoveryPlansClient
}

func NewClient(o *common.ClientOptions) *Client {
	vaultConfigsClient := backupresourcevaultconfigs.NewBackupResourceVaultConfigsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultConfigsClient.Client, o.ResourceManagerAuthorizer)

	storageConfigsClient := backupresourcestorageconfigsnoncrr.NewBackupResourceStorageConfigsNonCRRClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&storageConfigsClient.Client, o.ResourceManagerAuthorizer)

	vaultsClient := vaults.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	protectableItemsClient := backup.NewProtectableItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectableItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsClient := backup.NewProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectedItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsGroupClient := backup.NewProtectedItemsGroupClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectedItemsGroupClient.Client, o.ResourceManagerAuthorizer)

	protectionPoliciesClient := backup.NewProtectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	backupProtectionContainersClient := backup.NewProtectionContainersClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupProtectionContainersClient.Client, o.ResourceManagerAuthorizer)

	backupOperationStatusesClient := backup.NewOperationStatusesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupOperationStatusesClient.Client, o.ResourceManagerAuthorizer)

	backupProtectionContainerOperationResultsClient := backup.NewProtectionContainerOperationResultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupProtectionContainerOperationResultsClient.Client, o.ResourceManagerAuthorizer)

	fabricClient := replicationfabrics.NewReplicationFabricsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&fabricClient.Client, o.ResourceManagerAuthorizer)

	protectionContainerClient := replicationprotectioncontainers.NewReplicationProtectionContainersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&protectionContainerClient.Client, o.ResourceManagerAuthorizer)

	replicationPoliciesClient := replicationpolicies.NewReplicationPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&replicationPoliciesClient.Client, o.ResourceManagerAuthorizer)

	containerMappingClient := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&containerMappingClient.Client, o.ResourceManagerAuthorizer)

	networkMappingClient := replicationnetworkmappings.NewReplicationNetworkMappingsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&networkMappingClient.Client, o.ResourceManagerAuthorizer)

	replicationMigrationItemsClient := replicationprotecteditems.NewReplicationProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&replicationMigrationItemsClient.Client, o.ResourceManagerAuthorizer)

	replicationRecoveryPlanClient := replicationrecoveryplans.NewReplicationRecoveryPlansClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&replicationRecoveryPlanClient.Client, o.ResourceManagerAuthorizer)

	return &Client{
		ProtectableItemsClient:                    &protectableItemsClient,
		ProtectedItemsClient:                      &protectedItemsClient,
		ProtectedItemsGroupClient:                 &protectedItemsGroupClient,
		ProtectionPoliciesClient:                  &protectionPoliciesClient,
		ProtectionContainerOperationResultsClient: &backupProtectionContainerOperationResultsClient,
		BackupProtectionContainersClient:          &backupProtectionContainersClient,
		BackupOperationStatusesClient:             &backupOperationStatusesClient,
		VaultsClient:                              &vaultsClient,
		VaultsConfigsClient:                       &vaultConfigsClient,
		StorageConfigsClient:                      &storageConfigsClient,
		FabricClient:                              &fabricClient,
		ProtectionContainerClient:                 &protectionContainerClient,
		ReplicationPoliciesClient:                 &replicationPoliciesClient,
		ContainerMappingClient:                    &containerMappingClient,
		NetworkMappingClient:                      &networkMappingClient,
		ReplicationProtectedItemsClient:           &replicationMigrationItemsClient,
		ReplicationRecoveryPlansClient:            &replicationRecoveryPlanClient,
	}
}
