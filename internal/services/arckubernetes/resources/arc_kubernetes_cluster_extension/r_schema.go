// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_cluster_extension

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]{0,252}$"),
				"name must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number.",
			),
		},

		"cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: connectedclusters.ValidateConnectedClusterID,
		},

		"extension_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"configuration_protected_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"configuration_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"release_train": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"release_namespace": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ForceNew:      true,
			ConflictsWith: []string{"target_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"target_namespace": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ForceNew:      true,
			ConflictsWith: []string{"release_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"current_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
