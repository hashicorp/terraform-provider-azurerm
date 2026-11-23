// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_instance_elastic_san_volume_group

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/backupinstanceresources"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_instance_elastic_san_volume_group -service-package-name dataprotection -properties "name" -compare-values "subscription_id:vault_id,resource_group_name:vault_id,backup_vault_name:vault_id"

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &backupinstanceresources.BackupInstanceId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_data_protection_backup_instance_elastic_san_volume_group"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return backupinstanceresources.ValidateBackupInstanceID
}

func pollElasticSanVolumeGroupBackupInstance(ctx context.Context, client *backupinstanceresources.BackupInstanceResourcesClient, id backupinstanceresources.BackupInstanceId, pendingState backupinstanceresources.CurrentProtectionState) error {
	pollerType := custompollers.NewDataProtectionBackupInstance20260601Poller(client, id, backupinstanceresources.CurrentProtectionStateProtectionConfigured, []backupinstanceresources.CurrentProtectionState{pendingState})
	poller := pollers.NewPoller(pollerType, time.Minute, pollers.DefaultNumberOfDroppedConnectionsToAllow)
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}
	return nil
}
