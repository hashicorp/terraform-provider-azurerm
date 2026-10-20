// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package extended_location_custom_location

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2021-10-01/connectedclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2024-11-01/extensions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z\d.\-_]*[A-Za-z\d]$`),
				"supported alphanumeric characters and periods, underscores, hyphens. Name should end with an alphanumeric character.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"namespace": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]{0,252}$"),
				"namespace must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number.",
			),
		},

		"host_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: connectedclusters.ValidateConnectedClusterID,
		},

		"cluster_extension_ids": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem:     commonschema.ResourceIDReferenceElem(&extensions.ScopedExtensionId{}),
		},

		"host_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(customlocations.PossibleValuesForHostType(), false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authentication": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsBase64,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
