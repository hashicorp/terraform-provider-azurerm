// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database_sensitivity_label

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/sensitivitylabels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "column_name,schema_name,table_name" -compare-values "subscription_id:database_id,resource_group_name:database_id,server_name:database_id,database_name:database_id" -known-values "sensitivity_label_source:current"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

func (Resource) ResourceType() string {
	return "azurerm_mssql_database_sensitivity_label"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sensitivitylabels.ValidateSensitivityLabelSourceID
}

func (Resource) Identity() resourceids.ResourceId {
	return &sensitivitylabels.SensitivityLabelSourceId{}
}

func (Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := sensitivitylabels.ParseSensitivityLabelSourceID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		if id.SensitivityLabelSource != sensitivitylabels.SensitivityLabelSourceCurrent {
			return fmt.Errorf("only the `current` sensitivity label can be imported, got %q", id.SensitivityLabelSource)
		}
		return nil
	}
}
