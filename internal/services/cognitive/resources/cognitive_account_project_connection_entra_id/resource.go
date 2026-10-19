// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_entra_id

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name cognitive_account_project_connection_entra_id -properties "name" -compare-values "subscription_id:cognitive_account_project_id,resource_group_name:cognitive_account_project_id,account_name:cognitive_account_project_id,project_name:cognitive_account_project_id" -test-name "basic" -test-expect-non-empty

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
)

type Resource struct{}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return cognitiveAccountProjectConnectionImporter(projectconnectionresource.ConnectionAuthTypeAAD, r.ResourceType())
}

func (r Resource) Identity() resourceids.ResourceId {
	return new(projectconnectionresource.ProjectConnectionId)
}

func (r Resource) ResourceType() string {
	return "azurerm_cognitive_account_project_connection_entra_id"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projectconnectionresource.ValidateProjectConnectionID
}
