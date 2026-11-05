// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_deployment

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:cognitive_account_id,resource_group_name:cognitive_account_id,account_name:cognitive_account_id" -test-expect-non-empty

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &deployments.DeploymentId{}
}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_cognitive_deployment"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return deployments.ValidateDeploymentID
}
