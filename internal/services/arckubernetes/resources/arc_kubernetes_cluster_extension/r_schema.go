// Copyright IBM Corp. 2014, 2026
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
			Description: "Specifies the name which should be used for this Arc Kubernetes Cluster Extension.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]{0,252}$"),
				"name must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number.",
			),
		},

		"cluster_id": {
			Description:  "Specifies the Cluster ID.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: connectedclusters.ValidateConnectedClusterID,
		},

		"extension_type": {
			Description:  "Specifies the type of extension. For more information, please refer to [Available Extensions for Arc-enabled Kubernetes clusters](https://learn.microsoft.com/azure/azure-arc/kubernetes/extensions-release).",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"configuration_protected_settings": {
			Description: "Configuration settings that are sensitive, as name-value pairs for configuring this extension.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"configuration_settings": {
			Description: "Configuration settings, as name-value pairs for configuring this extension.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"release_train": {
			Description:  "The release train used by this extension.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"release_namespace": {
			Description:   "Namespace where the extension release must be placed for a cluster scoped extension. If this namespace does not exist, it will be created.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ForceNew:      true,
			ConflictsWith: []string{"target_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"target_namespace": {
			Description:   "Namespace where the extension will be created for a namespace scoped extension. If this namespace does not exist, it will be created.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ForceNew:      true,
			ConflictsWith: []string{"release_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"version": {
			Description:  "User-specified version that the extension should pin to. If it is not set, Azure will use the latest version and auto upgrade it.",
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
			Description: "The current version of the extension.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
