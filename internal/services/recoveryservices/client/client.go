package client

import (
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-10-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupprotectableitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupresourcestorageconfigsnoncrr"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/backupresourcevaultconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/protectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2021-12-01/protectionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationnetworkmappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryservicesproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationvaultsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/azuresdkhacks"
)

type Client struct {
	ProtectableItemsClient                    *backupprotectableitems.BackupProtectableItemsClient
	ProtectedItemsClient                      *protecteditems.ProtectedItemsClient
	ProtectedItemsGroupClient                 *backupprotecteditems.BackupProtectedItemsClient
	ProtectionPoliciesClient                  *protectionpolicies.ProtectionPoliciesClient
	ProtectionContainerOperationResultsClient *backup.ProtectionContainerOperationResultsClient
	BackupProtectionContainersClient          *protectioncontainers.ProtectionContainersClient
	BackupOperationStatusesClient             *backup.OperationStatusesClient
	BackupOperationResultsClient              *backup.OperationResultsClient
	VaultsClient                              *vaults.VaultsClient
	VaultsConfigsClient                       *backupresourcevaultconfigs.BackupResourceVaultConfigsClient
	VaultCertificatesClient                   *azuresdkhacks.VaultCertificatesClient
	VaultReplicationProvider                  *replicationrecoveryservicesproviders.ReplicationRecoveryServicesProvidersClient
	VaultsSettingsClient                      *replicationvaultsetting.ReplicationVaultSettingClient
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

	vaultSettingsClient := replicationvaultsetting.NewReplicationVaultSettingClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultSettingsClient.Client, o.ResourceManagerAuthorizer)

	storageConfigsClient := backupresourcestorageconfigsnoncrr.NewBackupResourceStorageConfigsNonCRRClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&storageConfigsClient.Client, o.ResourceManagerAuthorizer)

	vaultsClient := vaults.NewVaultsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultsClient.Client, o.ResourceManagerAuthorizer)

	vaultCertificatesClient := azuresdkhacks.NewVaultCertificatesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultCertificatesClient.Client, o.ResourceManagerAuthorizer)

	protectableItemsClient := backupprotectableitems.NewBackupProtectableItemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&protectableItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsClient := protecteditems.NewProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&protectedItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsGroupClient := backupprotecteditems.NewBackupProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&protectedItemsGroupClient.Client, o.ResourceManagerAuthorizer)

	protectionPoliciesClient := protectionpolicies.NewProtectionPoliciesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&protectionPoliciesClient.Client, o.ResourceManagerAuthorizer)

	backupProtectionContainersClient := protectioncontainers.NewProtectionContainersClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&backupProtectionContainersClient.Client, o.ResourceManagerAuthorizer)

	backupOperationStatusesClient := backup.NewOperationStatusesClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupOperationStatusesClient.Client, o.ResourceManagerAuthorizer)

	backupOperationResultClient := backup.NewOperationResultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&backupOperationResultClient.Client, o.ResourceManagerAuthorizer)

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
		BackupOperationResultsClient:              &backupOperationResultClient,
		VaultsClient:                              &vaultsClient,
		VaultsConfigsClient:                       &vaultConfigsClient,
		VaultCertificatesClient:                   &vaultCertificatesClient,
		VaultsSettingsClient:                      &vaultSettingsClient,
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
