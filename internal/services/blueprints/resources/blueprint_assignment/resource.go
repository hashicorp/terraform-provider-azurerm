// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package blueprint_assignment

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/assignment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBlueprintAssignmentCreateUpdate,
		Update: resourceBlueprintAssignmentCreateUpdate,
		Read:   resourceBlueprintAssignmentRead,
		Delete: resourceBlueprintAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := assignment.ParseScopedBlueprintAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: blueprintAssignmentSchema(),
	}
}
