// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2021-12-01/backup" // nolint: staticcheck
	vmwaremachines "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/machines"
	vmwarerunasaccounts "github.com/hashicorp/go-azure-sdk/resource-manager/migrate/2020-01-01/runasaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotectableitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/backupresourcevaultconfigs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/protectionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicesbackup/2023-02-01/resourceguardproxy"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworkmappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotecteditems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationprotectioncontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryservicesproviders"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationvaultsetting"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/azuresdkhacks"
)

type Client struct {
	ProtectableItemsClient *backupprotectableitems.BackupProtectableItemsClient
	ProtectedItemsClient   *protecteditems.ProtectedItemsClient
	// the Swagger lack of LRO mark, so we are using track-1 sdk to get the LRO client. tracked on https://github.com/Azure/azure-rest-api-specs/issues/22758
	ProtectedItemOperationResultsClient       *backup.ProtectedItemOperationResultsClient
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
	FabricClient                              *replicationfabrics.ReplicationFabricsClient
	ProtectionContainerClient                 *replicationprotectioncontainers.ReplicationProtectionContainersClient
	ReplicationPoliciesClient                 *replicationpolicies.ReplicationPoliciesClient
	ContainerMappingClient                    *replicationprotectioncontainermappings.ReplicationProtectionContainerMappingsClient
	NetworkMappingClient                      *replicationnetworkmappings.ReplicationNetworkMappingsClient
	ReplicationProtectedItemsClient           *replicationprotecteditems.ReplicationProtectedItemsClient
	ReplicationRecoveryPlansClient            *replicationrecoveryplans.ReplicationRecoveryPlansClient
	ReplicationNetworksClient                 *replicationnetworks.ReplicationNetworksClient
	ResourceGuardProxyClient                  *resourceguardproxy.ResourceGuardProxyClient
	VMWareMachinesClient                      *vmwaremachines.MachinesClient
	VMWareRunAsAccountsClient                 *vmwarerunasaccounts.RunAsAccountsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	vaultConfigsClient := backupresourcevaultconfigs.NewBackupResourceVaultConfigsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultConfigsClient.Client, o.ResourceManagerAuthorizer)

	vaultSettingsClient, err := replicationvaultsetting.NewReplicationVaultSettingClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationVaultSettings client: %+v", err)
	}
	o.Configure(vaultSettingsClient.Client, o.Authorizers.ResourceManager)

	vaultsClient, err := vaults.NewVaultsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Vaults client: %+v", err)
	}
	o.Configure(vaultsClient.Client, o.Authorizers.ResourceManager)

	vaultCertificatesClient := azuresdkhacks.NewVaultCertificatesClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&vaultCertificatesClient.Client, o.ResourceManagerAuthorizer)

	protectableItemsClient := backupprotectableitems.NewBackupProtectableItemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&protectableItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemsClient := protecteditems.NewProtectedItemsClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&protectedItemsClient.Client, o.ResourceManagerAuthorizer)

	protectedItemOperationResultClient := backup.NewProtectedItemOperationResultsClientWithBaseURI(o.ResourceManagerEndpoint, o.SubscriptionId)
	o.ConfigureClient(&protectedItemOperationResultClient.Client, o.ResourceManagerAuthorizer)

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

	fabricClient, err := replicationfabrics.NewReplicationFabricsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationFabrics client: %+v", err)
	}
	o.Configure(fabricClient.Client, o.Authorizers.ResourceManager)

	protectionContainerClient, err := replicationprotectioncontainers.NewReplicationProtectionContainersClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationProtectionContainers client: %+v", err)
	}
	o.Configure(protectionContainerClient.Client, o.Authorizers.ResourceManager)

	replicationPoliciesClient, err := replicationpolicies.NewReplicationPoliciesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationPolicies client: %+v", err)
	}
	o.Configure(replicationPoliciesClient.Client, o.Authorizers.ResourceManager)

	containerMappingClient, err := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationProtectionContainerMappings client: %+v", err)
	}
	o.Configure(containerMappingClient.Client, o.Authorizers.ResourceManager)

	networkMappingClient, err := replicationnetworkmappings.NewReplicationNetworkMappingsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationNetworks client: %+v", err)
	}
	o.Configure(networkMappingClient.Client, o.Authorizers.ResourceManager)

	replicationMigrationItemsClient, err := replicationprotecteditems.NewReplicationProtectedItemsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationNetworks client: %+v", err)
	}
	o.Configure(replicationMigrationItemsClient.Client, o.Authorizers.ResourceManager)

	replicationRecoveryPlanClient, err := replicationrecoveryplans.NewReplicationRecoveryPlansClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationNetworks client: %+v", err)
	}
	o.Configure(replicationRecoveryPlanClient.Client, o.Authorizers.ResourceManager)

	replicationNetworksClient, err := replicationnetworks.NewReplicationNetworksClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building ReplicationNetworks client: %+v", err)
	}
	o.Configure(replicationNetworksClient.Client, o.Authorizers.ResourceManager)

	resourceGuardProxyClient := resourceguardproxy.NewResourceGuardProxyClientWithBaseURI(o.ResourceManagerEndpoint)
	o.ConfigureClient(&resourceGuardProxyClient.Client, o.ResourceManagerAuthorizer)

	vmwareMachinesClient, err := vmwaremachines.NewMachinesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VMWare Machine client: %+v", err)
	}
	o.Configure(vmwareMachinesClient.Client, o.Authorizers.ResourceManager)

	vmwareRunAsAccountsClient, err := vmwarerunasaccounts.NewRunAsAccountsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building VMWare Run As Accounts client: %+v", err)
	}
	o.Configure(vmwareRunAsAccountsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		ProtectableItemsClient:                    &protectableItemsClient,
		ProtectedItemsClient:                      &protectedItemsClient,
		ProtectedItemsGroupClient:                 &protectedItemsGroupClient,
		ProtectionPoliciesClient:                  &protectionPoliciesClient,
		ProtectionContainerOperationResultsClient: &backupProtectionContainerOperationResultsClient,
		BackupProtectionContainersClient:          &backupProtectionContainersClient,
		ProtectedItemOperationResultsClient:       &protectedItemOperationResultClient,
		BackupOperationStatusesClient:             &backupOperationStatusesClient,
		BackupOperationResultsClient:              &backupOperationResultClient,
		VaultsClient:                              vaultsClient,
		VaultsConfigsClient:                       &vaultConfigsClient,
		VaultCertificatesClient:                   &vaultCertificatesClient,
		VaultsSettingsClient:                      vaultSettingsClient,
		FabricClient:                              fabricClient,
		ProtectionContainerClient:                 protectionContainerClient,
		ReplicationPoliciesClient:                 replicationPoliciesClient,
		ContainerMappingClient:                    containerMappingClient,
		NetworkMappingClient:                      networkMappingClient,
		ReplicationProtectedItemsClient:           replicationMigrationItemsClient,
		ReplicationRecoveryPlansClient:            replicationRecoveryPlanClient,
		ReplicationNetworksClient:                 replicationNetworksClient,
		ResourceGuardProxyClient:                  &resourceGuardProxyClient,
		VMWareMachinesClient:                      vmwareMachinesClient,
		VMWareRunAsAccountsClient:                 vmwareRunAsAccountsClient,
	}, nil
}
