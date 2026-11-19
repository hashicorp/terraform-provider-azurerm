// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_pipeline

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_monitor_pipeline"
}

func (r Resource) Identity() resourceids.ResourceId {
	return &pipelinegroups.PipelineGroupId{}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return pipelinegroups.ValidatePipelineGroupID
}
