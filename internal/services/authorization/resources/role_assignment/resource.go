// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_assignment

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// TODO: this wants splitting into virtual resources with Virtual IDs

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmRoleAssignmentCreate,
		Read:   resourceArmRoleAssignmentRead,
		Update: resourceArmRoleAssignmentUpdate,
		Delete: resourceArmRoleAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.RoleAssignmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: roleAssignmentSchema(),
	}
}
