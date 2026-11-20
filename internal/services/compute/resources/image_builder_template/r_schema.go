// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package image_builder_template

import (
	"regexp"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
				"Image template name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		// Though the 'None' type identity is declared in swagger, passing it to service triggers error: "Removing identity is not supported when creating or updating an image template.". So not expose it to users.
		"identity": {
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(identity.TypeUserAssigned),
						}, false),
					},
					"identity_ids": {
						Type:     schema.TypeSet,
						Required: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
				},
			},
		},

		"build_timeout_minutes": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      240,
			ValidateFunc: validation.IntBetween(0, 960),
		},

		// Fields in this block should not be assigned default values. Because fields mapped to Customizer Type A should not be specified when Customizer Type B is specified as the current customizer block type.
		"customizer": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							"File",
							"PowerShell",
							"Shell",
							"WindowsRestart",
							"WindowsUpdate",
						}, false),
					},

					"file_destination_path": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"file_sha256_checksum": {
						Type:     schema.TypeString,
						Optional: true,
						// Note: O+C because if not specify this property but only "file_source_uri" to the service, the service will calculate the sha256 of the file and return it.
						// So set this property as Computed for possible future usage. So forth to other similar properties in the `customizer` block.
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"file_source_uri": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"name": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"powershell_commands": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"powershell_run_as_system": {
						Type:     schema.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"powershell_run_elevated": {
						Type:     schema.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"powershell_script_uri": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"powershell_sha256_checksum": {
						Type:     schema.TypeString,
						Optional: true,
						// Note: O+C because if not specify this property but only "powershell_script_uri" to the service, the service will calculate the sha256 of the file and return it.
						// So set this property as Computed for possible future usage. So forth to other similar properties in the `customizer` block.
						Computed:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"powershell_valid_exit_codes": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type: schema.TypeInt,
						},
					},

					"shell_commands": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"shell_script_uri": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"shell_sha256_checksum": {
						Type:     schema.TypeString,
						Computed: true,
						// Note: O+C because if not specify this property but only "shell_script_uri" to the service, the service will calculate the sha256 of the file and return it.
						// So set this property as Computed for possible future usage. So forth to other similar properties in the `customizer` block.
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"windows_restart_check_command": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"windows_restart_command": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"windows_restart_timeout": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"windows_update_filters": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type:         schema.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"windows_update_search_criteria": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"windows_update_limit": {
						Type:     schema.TypeInt,
						Optional: true,
						ForceNew: true,
					},
				},
			},
		},

		"disk_size_gb": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 32767),
		},

		"distributions": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// The array of this block is order insensitive. But there is a bug preventing using TypeSet for this block having a sub field `location` using StateFunc:
					// https://github.com/hashicorp/terraform-plugin-sdk/issues/160.
					// In detail, the symptom is specifying user friendly region say "West US 2" in the sub field `location` generated two blocks even if only specifying one block in .tf.
					// Using normalized region say "westus2" does not have the symptom.
					// Given this bug would break the core logic, use TypeList as a workaround.
					// This justification also applies to the "distribution_shared_image" block below.
					"managed_image": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         schema.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"resource_group_name": commonschema.ResourceGroupName(),

								"location": commonschema.Location(),

								"run_output_name": distributionRunOutputNameSchema(),

								"tags": commonschema.TagsForceNew(),
							},
						},
					},

					"shared_image": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"id": {
									Type:         schema.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: azure.ValidateResourceID,
								},

								// The latest swagger (ver: 2020-02-14) defines this type as []string. However, in native SIG Image Version, not only region name but replica count and storage type are supported for each region.
								// So leave the type of this field as array of object here to serve future possible extensibility without bringing in breaking changes in user facing schema.
								"replica_regions": {
									Type:     schema.TypeList,
									Required: true,
									MinItems: 1,
									ForceNew: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"name": {
												Type:             schema.TypeString,
												Required:         true,
												ForceNew:         true,
												ValidateFunc:     location.EnhancedValidate,
												StateFunc:        location.StateFunc,
												DiffSuppressFunc: location.DiffSuppressFunc,
											},
										},
									},
								},

								"run_output_name": distributionRunOutputNameSchema(),

								"exclude_from_latest": {
									Type:     schema.TypeBool,
									Optional: true,
									ForceNew: true,
									Default:  false,
								},

								"storage_account_type": {
									Type:         schema.TypeString,
									Optional:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(virtualmachineimagetemplate.PossibleValuesForSharedImageStorageAccountType(), false),
								},

								"tags": commonschema.TagsForceNew(),

								"versioning": {
									Type:     schema.TypeList,
									Required: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"scheme": {
												Type:     schema.TypeString,
												Required: true,
												ForceNew: true,
												ValidateFunc: validation.StringInSlice([]string{
													"Source",
													"Latest",
												}, false),
											},
											"major": {
												Type:     schema.TypeInt,
												Optional: true,
												ForceNew: true,
											},
										},
									},
								},
							},
						},
					},

					"vhd": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"run_output_name": distributionRunOutputNameSchema(),

								"tags": commonschema.TagsForceNew(),
							},
						},
					},
				},
			},
		},

		"size": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "Standard_D1_v2",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"source_managed_image_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: images.ValidateImageID,
			ExactlyOneOf: []string{"source_managed_image_id", "source_platform_image", "source_shared_image_version_id"},
		},

		"source_platform_image": {
			Type:     schema.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"publisher": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"offer": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"sku": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					// If specify "Latest", the returned value is "latest ([a specific version])". I.e. in this situation the value returned by the service does not honor case sensitivity and adds more values.
					// A rest api bug filed: https://github.com/Azure/azure-rest-api-specs/issues/11313
					"version": {
						Type:         schema.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
							return strings.Contains(strings.ToLower(old), "latest (") && strings.EqualFold(new, "latest")
						},
					},

					"plan": {
						Type:     schema.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         schema.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"product": {
									Type:         schema.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"publisher": {
									Type:         schema.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
			ExactlyOneOf: []string{"source_managed_image_id", "source_platform_image", "source_shared_image_version_id"},
		},

		"source_shared_image_version_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: galleryimageversions.ValidateImageVersionID,
			ExactlyOneOf: []string{"source_managed_image_id", "source_platform_image", "source_shared_image_version_id"},
		},

		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validate.SubnetID,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
