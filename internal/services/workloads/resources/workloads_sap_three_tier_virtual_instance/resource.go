// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package workloads_sap_three_tier_virtual_instance

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2024-09-01/sapvirtualinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-params "10+(data.RandomInteger%90)" -test-expect-non-empty true

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (Resource) Identity() resourceids.ResourceId {
	return &sapvirtualinstances.SapVirtualInstanceId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_workloads_sap_three_tier_virtual_instance"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return sapvirtualinstances.ValidateSapVirtualInstanceID
}
