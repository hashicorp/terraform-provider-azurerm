// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_customized_accelerator

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (s Resource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_accelerator_id": commonschema.ResourceIDReferenceRequiredForceNew(&appplatform.ApplicationAcceleratorId{}, "Spring Cloud Accelerator"),

		"git_repository": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"basic_auth": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						MaxItems:      1,
						ConflictsWith: []string{"git_repository.0.ssh_auth"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"username": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"password": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"ssh_auth": {
						Type:          pluginsdk.TypeList,
						Optional:      true,
						ForceNew:      true,
						MaxItems:      1,
						ConflictsWith: []string{"git_repository.0.basic_auth"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"private_key": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"host_key": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"host_key_algorithm": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"branch": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
					},

					"ca_certificate_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.AsGeneratedID(appplatform.ParseCertificateIDInsensitively),
					},

					"commit": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
					},

					"git_tag": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						ExactlyOneOf: []string{"git_repository.0.branch", "git_repository.0.commit", "git_repository.0.git_tag"},
					},

					"interval_in_seconds": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(10),
					},

					"path": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"accelerator_tags": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"accelerator_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      appplatform.CustomizedAcceleratorTypeAccelerator,
			ValidateFunc: validation.StringInSlice(appplatform.PossibleValuesForCustomizedAcceleratorType(), false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"icon_url": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (s Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
