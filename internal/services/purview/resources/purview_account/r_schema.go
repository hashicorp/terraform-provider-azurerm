// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package purview_account

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func purviewAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][-a-zA-Z0-9]{1,61}[a-zA-Z0-9]$`),
				"The Purview account name must be between 3 and 63 characters long, it can contain only letters, numbers and hyphens, and the first and last characters must be a letter or number.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": resourcePurviewAccountIdentitySchema(),

		"managed_event_hub_enabled": {
			Type:     pluginsdk.TypeBool,
			Default:  true,
			Optional: true,
		},

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: resourcegroups.ValidateName,
		},

		"public_network_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"managed_resources": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_group_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"storage_account_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"event_hub_namespace_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"catalog_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"guardian_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"scan_endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"atlas_kafka_endpoint_primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"atlas_kafka_endpoint_secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"aws_external_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}

// resourcePurviewAccountIdentitySchema will take the common `SystemAssignedUserAssignedIdentityRequired` schema and make
// small changes to match this purview resource schemas with what is supported by API.
// The API inconsistency has been reported: https://github.com/Azure/azure-rest-api-specs/issues/22257
func resourcePurviewAccountIdentitySchema() *schema.Schema {
	customSchema := commonschema.SystemAssignedUserAssignedIdentityRequired()

	// remove UserAssigned from type validation, it is not supported in API
	customSchema.Elem.(*schema.Resource).Schema["type"].ValidateFunc = validation.StringInSlice([]string{
		string(identity.TypeSystemAssigned),
		string(identity.TypeSystemAssignedUserAssigned),
	}, false)

	// API is returning lower-case IDs which cause false diff report if mixed case was used by client.
	// Change the Set hash function to hash lower case always, which will suppress this invalid diff.
	customSchema.Elem.(*schema.Resource).Schema["identity_ids"].Set = resourcePurviewAccountIdentityIdsSetHash

	return customSchema
}
