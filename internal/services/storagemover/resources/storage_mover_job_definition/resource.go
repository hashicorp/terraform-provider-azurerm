// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_mover_job_definition

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2025-07-01/jobdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name storage_mover_job_definition -service-package-name storagemover -properties "name" -compare-values "subscription_id:storage_mover_project_id,resource_group_name:storage_mover_project_id,storage_mover_name:storage_mover_project_id,project_name:storage_mover_project_id"

type Resource struct{}

var (
	_ sdk.ResourceWithIdentity = Resource{}
	_ sdk.ResourceWithUpdate   = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &jobdefinitions.JobDefinitionId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_storage_mover_job_definition"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return jobdefinitions.ValidateJobDefinitionID
}
