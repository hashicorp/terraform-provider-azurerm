// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_devops_pool

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2025-09-20/pools"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/manageddevopspools/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-.]{1,42}[a-zA-Z0-9-]$`),
					"`name` can only include alphanumeric characters, periods (.) and hyphens (-). It must also start with alphanumeric characters and cannot end with periods (.) and length between 3 and 44.",
				),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"azure_devops_organization": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"organization": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"url": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.All(
										validation.IsURLWithHTTPS,
										validation.StringMatch(
											regexp.MustCompile(`[a-zA-Z0-9]$`),
											"url must end with a letter or number",
										),
									),
								},

								// There's an issue with API that if parallelism is omitted, it's always set to `0` instead of being computed dynamically.
								// To workaround this, mark it as Required which is also consistent with portal behaviour.
								// Relevant GH issue: https://github.com/Azure/azure-rest-api-specs/issues/40986
								"parallelism": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 10000),
								},

								"projects": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.ProjectName,
									},
								},
							},
						},
					},

					"permission": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"kind": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForAzureDevOpsPermissionType(), false),
								},

								"administrator_account": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									ForceNew: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"groups": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.IsEmailAddress,
												},
												AtLeastOneOf: []string{"azure_devops_organization.0.permission.0.administrator_account.0.groups", "azure_devops_organization.0.permission.0.administrator_account.0.users"},
											},

											"users": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												ForceNew: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.IsEmailAddress,
												},
												AtLeastOneOf: []string{"azure_devops_organization.0.permission.0.administrator_account.0.groups", "azure_devops_organization.0.permission.0.administrator_account.0.users"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},

		"dev_center_project_id": commonschema.ResourceIDReferenceRequired(&projects.ProjectId{}),

		"maximum_concurrency": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 10000),
		},

		"virtual_machine_scale_set_fabric": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"image": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"aliases": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},

								"buffer": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  "*",
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^(?:\*|[0-9][0-9]?|100)$`),
										`Buffer must be "*" or value between 0 and 100.`,
									),
								},

								"id": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: azure.ValidateResourceID,
								},

								"well_known_image_name": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"sku_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"os_disk_storage_account_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(pools.OsDiskStorageAccountTypeStandard),
						ValidateFunc: validation.StringInSlice(pools.PossibleValuesForOsDiskStorageAccountType(), false),
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"security": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"interactive_logon_enabled": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},

								"key_vault_management": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"key_export_enabled": {
												Type:     pluginsdk.TypeBool,
												Default:  false,
												Optional: true,
											},

											"key_vault_certificate_ids": {
												Type:     pluginsdk.TypeList,
												Required: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: keyvault.ValidateNestedItemID(keyvault.VersionTypeVersionless, keyvault.NestedItemTypeSecret),
												},
											},

											"certificate_store_location": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"certificate_store_name": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCertificateStoreNameOption(), false),
											},
										},
									},
								},
							},
						},
					},

					"storage": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"caching": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(pools.CachingTypeReadOnly),
										string(pools.CachingTypeReadWrite),
									}, false),
								},

								"disk_size_in_gb": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 32767),
								},

								"drive_letter": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ValidateFunc: validation.StringMatch(
										regexp.MustCompile(`^[FGHIJKLMNOPQRSTUVWXYZfghijklmnopqrstuvwxyz]$`),
										"drive_letter must be a single letter and cannot be A, C, D, or E (case insensitive)",
									),
								},

								"storage_account_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      pools.StorageAccountTypeStandardLRS,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForStorageAccountType(), false),
								},
							},
						},
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityOptional(),

		"stateful_agent": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"grace_period_time_span": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "00:00:00",
						ValidateFunc: validate.AgentLifetime,
					},

					"maximum_agent_lifetime": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      "7.00:00:00",
						ValidateFunc: validate.AgentLifetime,
					},

					"manual_resource_prediction": manualResourcePredictionSchema("stateful_agent.0"),

					"automatic_resource_prediction": automaticResourcePredictionSchema("stateful_agent.0"),
				},
			},
			ExactlyOneOf: []string{"stateful_agent", "stateless_agent"},
		},

		"stateless_agent": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"manual_resource_prediction": manualResourcePredictionSchema("stateless_agent.0"),

					"automatic_resource_prediction": automaticResourcePredictionSchema("stateless_agent.0"),
				},
			},
			ExactlyOneOf: []string{"stateful_agent", "stateless_agent"},
		},

		"work_folder": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedDevOpsPoolModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("DecodeDiff: %+v", err)
			}

			if err := validateVirtualMachineScaleSetFabricImages(metadata, model.VirtualMachineScaleSetFabric); err != nil {
				return err
			}

			maxConcurrency := model.MaximumConcurrency

			for _, org := range model.AzureDevOpsOrganization {
				for _, perm := range org.Permission {
					if perm.Kind != string(pools.AzureDevOpsPermissionTypeSpecificAccounts) {
						if len(perm.AdministratorAccounts) > 0 {
							return fmt.Errorf("`administrator_account` must not be set when `permission` kind is `%s`", perm.Kind)
						}
					}
				}

				var parallelismSum int64
				for _, o := range org.Organizations {
					if o.Parallelism > 0 {
						parallelismSum += o.Parallelism
					}
				}
				if parallelismSum != maxConcurrency {
					return fmt.Errorf("the sum of `parallelism` across all organizations (%d) must equal `maximum_concurrency` (%d)", parallelismSum, maxConcurrency)
				}
			}

			for _, stateful := range model.StatefulAgent {
				for _, manualPredictions := range stateful.ManualResourcePrediction {
					if err := validateManualAgentCounts("stateful_agent", manualPredictions, maxConcurrency); err != nil {
						return err
					}
				}
			}

			for _, stateless := range model.StatelessAgent {
				for _, manualPredictions := range stateless.ManualResourcePrediction {
					if err := validateManualAgentCounts("stateless_agent", manualPredictions, maxConcurrency); err != nil {
						return err
					}
				}
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}
