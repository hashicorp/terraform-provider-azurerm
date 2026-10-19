// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ssh_public_key

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sshPublicKeySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[-a-zA-Z0-9(_).]{1,128}$`),
				"Public SSH Key name must be 1 - 128 characters long, can contain letters, numbers, underscores, dots and hyphens (but the first and last character must be a letter or number).",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"public_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.SSHKey,
		},

		"tags": commonschema.Tags(),
	}
}
