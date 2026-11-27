// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_protection_backup_policy_elastic_san_volume

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2026-06-01/basebackuppolicyresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const elasticSanVolumeDatasourceType = "Microsoft.ElasticSan/elasticSans/volumeGroups"

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name data_protection_backup_policy_elastic_san_volume -service-package-name dataprotection -properties "name" -compare-values "subscription_id:data_protection_backup_vault_id,resource_group_name:data_protection_backup_vault_id,backup_vault_name:data_protection_backup_vault_id"

type Resource struct{}

var (
	_ sdk.Resource                   = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &basebackuppolicyresources.BackupPolicyId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_data_protection_backup_policy_elastic_san_volume"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return basebackuppolicyresources.ValidateBackupPolicyID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := basebackuppolicyresources.ParseBackupPolicyIDInsensitively(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := metadata.Client.DataProtection.BackupPolicyClient20260601.BackupPoliciesGet(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: model was nil", *id)
		}
		if !backupPolicySupportsElasticSanVolumes(resp.Model.Properties) {
			return fmt.Errorf("%s is not a backup policy for `%s`", *id, elasticSanVolumeDatasourceType)
		}

		metadata.SetID(id)
		return nil
	}
}

func backupPolicySupportsElasticSanVolumes(input basebackuppolicyresources.BaseBackupPolicy) bool {
	policy, ok := input.(basebackuppolicyresources.BackupPolicy)
	if !ok {
		return false
	}
	for _, datasourceType := range policy.DatasourceTypes {
		if strings.EqualFold(datasourceType, elasticSanVolumeDatasourceType) {
			return true
		}
	}
	return false
}
