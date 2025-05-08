// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBatchPool() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceBatchPoolCreate,
		Read:   resourceBatchPoolRead,
		Update: resourceBatchPoolUpdate,
		Delete: resourceBatchPoolDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := pool.ParsePoolID(id)
			return err
		}),
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PoolName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},
			"display_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vm_size": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"max_tasks_per_node": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"fixed_scale": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// Property `node_deallocation_method` is set to be a writeOnly property by service team
						// It can only perform on PUT operation and is not able to perform GET operation
						// Here we treat `node_deallocation_method` the same as a secret value.
						// Issue link: https://github.com/Azure/azure-rest-api-specs/issues/20948
						"node_deallocation_method": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.ComputeNodeDeallocationOptionRequeue),
								string(pool.ComputeNodeDeallocationOptionRetainedData),
								string(pool.ComputeNodeDeallocationOptionTaskCompletion),
								string(pool.ComputeNodeDeallocationOptionTerminate),
							}, false),
						},
						"target_dedicated_nodes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(0, 2000),
						},
						"target_low_priority_nodes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
						"resize_timeout": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "PT15M",
						},
					},
				},
			},
			"auto_scale": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"evaluation_interval": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "PT15M",
						},
						"formula": {
							Type:     pluginsdk.TypeString,
							Required: true,
							DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
								return strings.TrimSpace(old) == strings.TrimSpace(new)
							},
						},
					},
				},
			},
			"container_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
						},
						"container_image_names": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
						},
						"container_registries": {
							Type:       pluginsdk.TypeList,
							Optional:   true,
							ForceNew:   true,
							ConfigMode: pluginsdk.SchemaConfigModeAttr,
							Elem: &pluginsdk.Resource{
								Schema: containerRegistry(),
							},
							AtLeastOneOf: []string{"container_configuration.0.type", "container_configuration.0.container_image_names", "container_configuration.0.container_registries"},
						},
					},
				},
			},
			"storage_image_reference": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
							AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
						},

						"publisher": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
						},

						"offer": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
						},

						"sku": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ForceNew:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.StringIsNotEmpty,
							AtLeastOneOf:     []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
						},

						"version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{"storage_image_reference.0.id", "storage_image_reference.0.publisher", "storage_image_reference.0.offer", "storage_image_reference.0.sku", "storage_image_reference.0.version"},
						},
					},
				},
			},
			"node_agent_sku_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stop_pending_resize_operation": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
			"certificate": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
							// The ID returned for the certificate in the batch account and the certificate applied to the pool
							// are not consistent in their casing which causes issues when referencing IDs across resources
							// (as Terraform still sees differences to apply due to the casing)
							// Handling by ignoring casing for now. Raised as an issue: https://github.com/Azure/azure-rest-api-specs/issues/5574
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"store_location": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CurrentUser",
								"LocalMachine",
							}, false),
						},
						"store_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"visibility": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"StartTask",
									"Task",
									"RemoteUser",
								}, false),
							},
						},
					},
				},
			},

			"identity": commonschema.UserAssignedIdentityOptional(),

			"start_task": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: startTaskSchema(),
				},
			},
			"metadata": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"mount": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"azure_blob_file_system": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"account_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"container_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"relative_mount_path": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"account_key": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"sas_key": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"identity_id": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: commonids.ValidateUserAssignedIdentityID,
									},
									"blobfuse_options": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"azure_file_share": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"account_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"azure_file_url": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsURLWithHTTPS,
									},
									"account_key": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"relative_mount_path": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"mount_options": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"cifs_mount": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"user_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"source": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"relative_mount_path": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"mount_options": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
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
						"nfs_mount": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"source": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"relative_mount_path": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"mount_options": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
			"network_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dynamic_vnet_assignment_scope": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ForceNew:         true,
							Default:          string(pool.DynamicVNetAssignmentScopeNone),
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.DynamicVNetAssignmentScopeNone),
								string(pool.DynamicVNetAssignmentScopeJob),
							}, false),
						},
						"accelerated_networking_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"public_ips": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},
						"public_address_provisioning_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.IPAddressProvisioningTypeBatchManaged),
								string(pool.IPAddressProvisioningTypeUserManaged),
								string(pool.IPAddressProvisioningTypeNoPublicIPAddresses),
							}, false),
						},
						"endpoint_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"protocol": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(pool.InboundEndpointProtocolTCP),
											string(pool.InboundEndpointProtocolUDP),
										}, false),
									},
									"backend_port": {
										Type:     pluginsdk.TypeInt,
										Required: true,
										ForceNew: true,
										ValidateFunc: validation.All(
											validation.IsPortNumber,
											validation.IntNotInSlice([]int{29876, 29877}),
										),
									},
									"frontend_port_range": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validate.FrontendPortRange,
									},
									"network_security_group_rules": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"priority": {
													Type:         pluginsdk.TypeInt,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.IntAtLeast(150),
												},
												"access": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(pool.NetworkSecurityGroupRuleAccessAllow),
														string(pool.NetworkSecurityGroupRuleAccessDeny),
													}, false),
												},
												"source_address_prefix": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ForceNew:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"source_port_ranges": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													Computed: true,
													ForceNew: true,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeString,
														Default:      "*",
														ValidateFunc: validation.StringIsNotEmpty,
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
			},
			"data_disks": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"lun": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(0, 63),
						},
						"caching": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(pool.CachingTypeReadOnly),
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.CachingTypeNone),
								string(pool.CachingTypeReadOnly),
								string(pool.CachingTypeReadWrite),
							}, false),
						},
						"disk_size_gb": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},
						"storage_account_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  pool.StorageAccountTypeStandardLRS,
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.StorageAccountTypeStandardLRS),
								string(pool.StorageAccountTypePremiumLRS),
							}, false),
						},
					},
				},
			},
			"disk_encryption": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disk_encryption_target": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.DiskEncryptionTargetTemporaryDisk),
								string(pool.DiskEncryptionTargetOsDisk),
							}, false),
						},
					},
				},
			},
			"extensions": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"publisher": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type_handler_version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"auto_upgrade_minor_version": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"automatic_upgrade_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
						"settings_json": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"protected_settings": { // todo 4.0 - should this actually be a map of key value pairs?
							Type:      pluginsdk.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"provision_after_extensions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
			"node_placement": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"policy": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(pool.NodePlacementPolicyTypeRegional),
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.NodePlacementPolicyTypeZonal),
								string(pool.NodePlacementPolicyTypeRegional),
							}, false),
						},
					},
				},
			},
			"license_type": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"os_disk_placement": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(pool.DiffDiskPlacementCacheDisk),
					}, false),
			},
			"inter_node_communication": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(pool.InterNodeCommunicationStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(pool.InterNodeCommunicationStateEnabled),
					string(pool.InterNodeCommunicationStateDisabled),
				}, false),
			},

			"security_profile": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"host_encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							ForceNew: true,
							Optional: true,
						},
						"security_type": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(pool.PossibleValuesForSecurityTypes(), false),
						},
						"secure_boot_enabled": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							ForceNew:     true,
							RequiredWith: []string{"security_profile.0.security_type"},
						},
						"vtpm_enabled": {
							Type:         pluginsdk.TypeBool,
							Optional:     true,
							ForceNew:     true,
							RequiredWith: []string{"security_profile.0.security_type"},
						},
					},
				},
			},

			"target_node_communication_mode": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(pool.PossibleValuesForNodeCommunicationMode(), false),
			},

			"task_scheduling_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"node_fill_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.ComputeNodeFillTypeSpread),
								string(pool.ComputeNodeFillTypePack),
							}, false),
						},
					},
				},
			},
			"user_accounts": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
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
						"elevation_level": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(pool.ElevationLevelNonAdmin),
								string(pool.ElevationLevelAdmin),
							}, false),
						},
						"linux_user_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"uid": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
									},
									"gid": {
										Type:     pluginsdk.TypeInt,
										Optional: true,
									},
									"ssh_private_key": {
										Type:      pluginsdk.TypeString,
										Optional:  true,
										Sensitive: true,
									},
								},
							},
						},
						"windows_user_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"login_mode": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(pool.LoginModeBatch),
											string(pool.LoginModeInteractive),
										}, false),
									},
								},
							},
						},
					},
				},
			},
			"windows": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enable_automatic_updates": {
							Type:     pluginsdk.TypeBool,
							Default:  true,
							Optional: true,
						},
					},
				},
			},
		},
	}

	return resource
}

func resourceBatchPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := pool.NewPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_batch_pool", id.ID())
		}
	}

	parameters := pool.Pool{
		Properties: &pool.PoolProperties{
			VMSize:                 utils.String(d.Get("vm_size").(string)),
			DisplayName:            utils.String(d.Get("display_name").(string)),
			InterNodeCommunication: pointer.To(pool.InterNodeCommunicationState(d.Get("inter_node_communication").(string))),
			TaskSlotsPerNode:       utils.Int64(int64(d.Get("max_tasks_per_node").(int))),
		},
	}

	userAccounts, err := ExpandBatchPoolUserAccounts(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "user_accounts": %v`, err)
	}
	parameters.Properties.UserAccounts = userAccounts

	taskSchedulingPolicy, err := ExpandBatchPoolTaskSchedulingPolicy(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "task_scheduling_policy": %v`, err)
	}
	parameters.Properties.TaskSchedulingPolicy = taskSchedulingPolicy

	identityResult, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}
	parameters.Identity = identityResult

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("expanding scale settings: %+v", err)
	}

	parameters.Properties.ScaleSettings = scaleSettings

	if startTaskValue, startTaskOk := d.GetOk("start_task"); startTaskOk {
		startTaskList := startTaskValue.([]interface{})
		startTask, startTaskErr := ExpandBatchPoolStartTask(startTaskList)

		if startTaskErr != nil {
			return fmt.Errorf("creating %s: %+v", id, startTaskErr)
		}

		// start task should have a user identity defined
		userIdentity := startTask.UserIdentity
		if userIdentityError := validateUserIdentity(userIdentity); userIdentityError != nil {
			return fmt.Errorf("creating %s: %+v", id, userIdentityError)
		}

		parameters.Properties.StartTask = startTask
	}

	if vmDeploymentConfiguration, deploymentErr := expandBatchPoolVirtualMachineConfig(d); deploymentErr == nil {
		parameters.Properties.DeploymentConfiguration = &pool.DeploymentConfiguration{
			VirtualMachineConfiguration: vmDeploymentConfiguration,
		}
	} else {
		return deploymentErr
	}

	if v, ok := d.GetOk("certificate"); ok {
		certificateReferences, err := ExpandBatchPoolCertificateReferences(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `certificate`: %+v", err)
		}
		parameters.Properties.Certificates = certificateReferences
	}

	if err := validateBatchPoolCrossFieldRules(parameters.Properties); err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	parameters.Properties.Metadata = ExpandBatchMetaData(metaDataRaw)

	mountConfiguration, err := ExpandBatchPoolMountConfigurations(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "mount": %v`, err)
	}
	parameters.Properties.MountConfiguration = mountConfiguration

	networkConfiguration := d.Get("network_configuration").([]interface{})
	parameters.Properties.NetworkConfiguration, err = ExpandBatchPoolNetworkConfiguration(networkConfiguration)
	if err != nil {
		return fmt.Errorf("expanding `network_configuration`: %+v", err)
	}

	if v, ok := d.GetOk("target_node_communication_mode"); ok {
		parameters.Properties.TargetNodeCommunicationMode = pointer.To(pool.NodeCommunicationMode(v.(string)))
	}

	_, err = client.Create(ctx, id, parameters, pool.CreateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	read, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// if the pool is not Steady after the create operation, wait for it to be Steady
	if model := read.Model; model != nil {
		if props := model.Properties; props != nil && props.AllocationState != nil && *props.AllocationState != pool.AllocationStateSteady {
			if err = waitForBatchPoolPendingResizeOperation(ctx, client, id); err != nil {
				return fmt.Errorf("waiting for %s", id)
			}
		}
	}

	return resourceBatchPoolRead(d, meta)
}

func resourceBatchPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pool.ParsePoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.AllocationState != nil && *props.AllocationState != pool.AllocationStateSteady {
			log.Printf("[INFO] there is a pending resize operation on this pool...")
			stopPendingResizeOperation := d.Get("stop_pending_resize_operation").(bool)
			if !stopPendingResizeOperation {
				return fmt.Errorf("updating %s because of pending resize operation. Set flag `stop_pending_resize_operation` to true to force update", *id)
			}

			log.Printf("[INFO] stopping the pending resize operation on this pool...")
			if _, err = client.StopResize(ctx, *id); err != nil {
				return fmt.Errorf("stopping resize operation for %s: %+v", *id, err)
			}

			// waiting for the pool to be in steady state
			if err = waitForBatchPoolPendingResizeOperation(ctx, client, *id); err != nil {
				return fmt.Errorf("waiting for %s", *id)
			}
		}
	}

	parameters := pool.Pool{
		Properties: &pool.PoolProperties{},
	}

	identity, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}
	parameters.Identity = identity

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("expanding scale settings: %+v", err)
	}

	parameters.Properties.ScaleSettings = scaleSettings

	taskSchedulingPolicy, err := ExpandBatchPoolTaskSchedulingPolicy(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "task_scheduling_policy": %v`, err)
	}
	parameters.Properties.TaskSchedulingPolicy = taskSchedulingPolicy

	userAccounts, err := ExpandBatchPoolUserAccounts(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "user_accounts": %v`, err)
	}
	parameters.Properties.UserAccounts = userAccounts

	if startTaskValue, startTaskOk := d.GetOk("start_task"); startTaskOk {
		startTaskList := startTaskValue.([]interface{})
		startTask, startTaskErr := ExpandBatchPoolStartTask(startTaskList)

		if startTaskErr != nil {
			return fmt.Errorf("updating %s: %+v", *id, startTaskErr)
		}

		// start task should have a user identity defined
		userIdentity := startTask.UserIdentity
		if userIdentityError := validateUserIdentity(userIdentity); userIdentityError != nil {
			return fmt.Errorf("creating %s: %+v", *id, userIdentityError)
		}

		parameters.Properties.StartTask = startTask
	}
	certificates := d.Get("certificate").([]interface{})
	certificateReferences, err := ExpandBatchPoolCertificateReferences(certificates)
	if err != nil {
		return fmt.Errorf("expanding `certificate`: %+v", err)
	}
	parameters.Properties.Certificates = certificateReferences

	if err := validateBatchPoolCrossFieldRules(parameters.Properties); err != nil {
		return err
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for %s", *id)
		metaDataRaw := d.Get("metadata").(map[string]interface{})

		parameters.Properties.Metadata = ExpandBatchMetaData(metaDataRaw)
	}

	mountConfiguration, err := ExpandBatchPoolMountConfigurations(d)
	if err != nil {
		log.Printf(`[DEBUG] expanding "mount": %v`, err)
	}
	parameters.Properties.MountConfiguration = mountConfiguration

	if d.HasChange("target_node_communication_mode") {
		parameters.Properties.TargetNodeCommunicationMode = pointer.To(pool.NodeCommunicationMode(d.Get("target_node_communication_mode").(string)))
	}

	result, err := client.Update(ctx, *id, parameters, pool.UpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// if the pool is not Steady after the update, wait for it to be Steady
	if model := result.Model; model != nil {
		if props := model.Properties; props != nil && props.AllocationState != nil && *props.AllocationState != pool.AllocationStateSteady {
			if err := waitForBatchPoolPendingResizeOperation(ctx, client, *id); err != nil {
				return fmt.Errorf("waiting for %s", *id)
			}
		}
	}

	return resourceBatchPoolRead(d, meta)
}

func resourceBatchPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pool.ParsePoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PoolName)
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		identityResult, err := identity.FlattenUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityResult); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("display_name", props.DisplayName)
			d.Set("vm_size", props.VMSize)
			d.Set("inter_node_communication", string(pointer.From(props.InterNodeCommunication)))

			if scaleSettings := props.ScaleSettings; scaleSettings != nil {
				if err := d.Set("auto_scale", flattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
					return fmt.Errorf("flattening `auto_scale`: %+v", err)
				}
				if err := d.Set("fixed_scale", flattenBatchPoolFixedScaleSettings(d, scaleSettings.FixedScale)); err != nil {
					return fmt.Errorf("flattening `fixed_scale `: %+v", err)
				}
			}

			if props.TaskSchedulingPolicy != nil && props.TaskSchedulingPolicy.NodeFillType != "" {
				taskSchedulingPolicy := make([]interface{}, 0)
				nodeFillType := make(map[string]interface{})
				nodeFillType["node_fill_type"] = string(props.TaskSchedulingPolicy.NodeFillType)
				taskSchedulingPolicy = append(taskSchedulingPolicy, nodeFillType)
				d.Set("task_scheduling_policy", taskSchedulingPolicy)
			}

			if props.UserAccounts != nil {
				userAccounts := make([]interface{}, 0)
				for _, userAccount := range *props.UserAccounts {
					userAccounts = append(userAccounts, flattenBatchPoolUserAccount(d, &userAccount))
				}
				d.Set("user_accounts", userAccounts)
			}

			d.Set("max_tasks_per_node", props.TaskSlotsPerNode)

			if props.DeploymentConfiguration != nil {
				if props.DeploymentConfiguration.VirtualMachineConfiguration != nil {
					config := props.DeploymentConfiguration.VirtualMachineConfiguration
					if config.ContainerConfiguration != nil {
						d.Set("container_configuration", flattenBatchPoolContainerConfiguration(d, config.ContainerConfiguration))
					}
					if config.DataDisks != nil {
						dataDisks := make([]interface{}, 0)
						for _, item := range *config.DataDisks {
							dataDisk := make(map[string]interface{})
							dataDisk["lun"] = item.Lun
							dataDisk["disk_size_gb"] = item.DiskSizeGB

							caching := ""
							if item.Caching != nil {
								caching = string(*item.Caching)
							}
							dataDisk["caching"] = caching

							storageAccountType := ""
							if item.StorageAccountType != nil {
								storageAccountType = string(*item.StorageAccountType)
							}
							dataDisk["storage_account_type"] = storageAccountType

							dataDisks = append(dataDisks, dataDisk)
						}
						d.Set("data_disks", dataDisks)
					}
					if config.DiskEncryptionConfiguration != nil {
						diskEncryptionConfiguration := make([]interface{}, 0)
						if config.DiskEncryptionConfiguration.Targets != nil {
							for _, item := range *config.DiskEncryptionConfiguration.Targets {
								target := make(map[string]interface{})
								target["disk_encryption_target"] = string(item)
								diskEncryptionConfiguration = append(diskEncryptionConfiguration, target)
							}
						}
						d.Set("disk_encryption", diskEncryptionConfiguration)
					}
					if config.Extensions != nil {
						extensions := make([]interface{}, 0)
						n := len(*config.Extensions)
						for _, item := range *config.Extensions {
							extension := make(map[string]interface{})
							extension["name"] = item.Name
							extension["publisher"] = item.Publisher
							extension["type"] = item.Type
							if item.TypeHandlerVersion != nil {
								extension["type_handler_version"] = *item.TypeHandlerVersion
							}
							if item.AutoUpgradeMinorVersion != nil {
								extension["auto_upgrade_minor_version"] = *item.AutoUpgradeMinorVersion
							}
							if item.EnableAutomaticUpgrade != nil {
								extension["automatic_upgrade_enabled"] = *item.EnableAutomaticUpgrade
							}
							if item.Settings != nil {
								settingValue, err := json.Marshal((*item.Settings).(map[string]interface{}))
								if err != nil {
									return fmt.Errorf("flattening `settings_json`: %+v", err)
								}
								extension["settings_json"] = string(settingValue)
							}

							for i := 0; i < n; i++ {
								if v, ok := d.GetOk(fmt.Sprintf("extensions.%d.name", i)); ok && v == item.Name {
									extension["protected_settings"] = d.Get(fmt.Sprintf("extensions.%d.protected_settings", i))
									break
								}
							}

							if item.ProvisionAfterExtensions != nil {
								extension["provision_after_extensions"] = *item.ProvisionAfterExtensions
							}
							extensions = append(extensions, extension)
						}
						d.Set("extensions", extensions)
					}

					d.Set("storage_image_reference", flattenBatchPoolImageReference(&config.ImageReference))
					d.Set("license_type", config.LicenseType)
					d.Set("node_agent_sku_id", config.NodeAgentSkuId)

					if config.NodePlacementConfiguration != nil {
						nodePlacementConfiguration := make([]interface{}, 0)
						nodePlacementConfig := make(map[string]interface{})
						nodePlacementConfig["policy"] = string(*config.NodePlacementConfiguration.Policy)
						nodePlacementConfiguration = append(nodePlacementConfiguration, nodePlacementConfig)
						d.Set("node_placement", nodePlacementConfiguration)
					}

					osDiskPlacement := ""
					if config.OsDisk != nil && config.OsDisk.EphemeralOSDiskSettings != nil && config.OsDisk.EphemeralOSDiskSettings.Placement != nil {
						osDiskPlacement = string(*config.OsDisk.EphemeralOSDiskSettings.Placement)
					}
					d.Set("os_disk_placement", osDiskPlacement)

					if config.SecurityProfile != nil {
						d.Set("security_profile", flattenBatchPoolSecurityProfile(config.SecurityProfile))
					}

					if config.WindowsConfiguration != nil {
						windowsConfig := []interface{}{
							map[string]interface{}{
								"enable_automatic_updates": *config.WindowsConfiguration.EnableAutomaticUpdates,
							},
						}
						d.Set("windows", windowsConfig)
					}
				}
			}

			if err := d.Set("certificate", flattenBatchPoolCertificateReferences(props.Certificates)); err != nil {
				return fmt.Errorf("flattening `certificate`: %+v", err)
			}

			d.Set("start_task", flattenBatchPoolStartTask(d, props.StartTask))
			d.Set("metadata", FlattenBatchMetaData(props.Metadata))

			if props.MountConfiguration != nil {
				mountConfigs := make([]interface{}, 0)
				for _, mountConfig := range *props.MountConfiguration {
					mountConfigs = append(mountConfigs, flattenBatchPoolMountConfig(d, &mountConfig))
				}
				d.Set("mount", mountConfigs)
			}

			targetNodeCommunicationMode := ""
			if props.TargetNodeCommunicationMode != nil {
				targetNodeCommunicationMode = string(*props.TargetNodeCommunicationMode)
			}
			d.Set("target_node_communication_mode", targetNodeCommunicationMode)

			if err := d.Set("network_configuration", flattenBatchPoolNetworkConfiguration(props.NetworkConfiguration)); err != nil {
				return fmt.Errorf("setting `network_configuration`: %v", err)
			}
		}
	}

	return nil
}

func resourceBatchPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pool.ParsePoolID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandBatchPoolScaleSettings(d *pluginsdk.ResourceData) (*pool.ScaleSettings, error) {
	scaleSettings := &pool.ScaleSettings{}

	autoScaleValue, autoScaleOk := d.GetOk("auto_scale")
	fixedScaleValue, fixedScaleOk := d.GetOk("fixed_scale")

	if !autoScaleOk && !fixedScaleOk {
		return nil, fmt.Errorf("auto_scale block or fixed_scale block need to be specified")
	}

	if autoScaleOk && fixedScaleOk {
		return nil, fmt.Errorf("auto_scale and fixed_scale blocks cannot be specified at the same time")
	}

	if autoScaleOk {
		autoScale := autoScaleValue.([]interface{})
		if len(autoScale) == 0 {
			return nil, fmt.Errorf("when scale mode is Auto, auto_scale block is required")
		}

		autoScaleSettings := autoScale[0].(map[string]interface{})

		autoScaleEvaluationInterval := autoScaleSettings["evaluation_interval"].(string)
		autoScaleFormula := autoScaleSettings["formula"].(string)

		scaleSettings.AutoScale = &pool.AutoScaleSettings{
			EvaluationInterval: &autoScaleEvaluationInterval,
			Formula:            autoScaleFormula,
		}
	} else if fixedScaleOk {
		fixedScale := fixedScaleValue.([]interface{})
		if len(fixedScale) == 0 {
			return nil, fmt.Errorf("when scale mode is Fixed, fixed_scale block is required")
		}

		fixedScaleSettings := fixedScale[0].(map[string]interface{})
		nodeDeallocationOption := pool.ComputeNodeDeallocationOption(fixedScaleSettings["node_deallocation_method"].(string))
		targetDedicatedNodes := int32(fixedScaleSettings["target_dedicated_nodes"].(int))
		targetLowPriorityNodes := int32(fixedScaleSettings["target_low_priority_nodes"].(int))
		resizeTimeout := fixedScaleSettings["resize_timeout"].(string)

		scaleSettings.FixedScale = &pool.FixedScaleSettings{
			NodeDeallocationOption: &nodeDeallocationOption,
			ResizeTimeout:          &resizeTimeout,
			TargetDedicatedNodes:   utils.Int64(int64(targetDedicatedNodes)),
			TargetLowPriorityNodes: utils.Int64(int64(targetLowPriorityNodes)),
		}
	}

	return scaleSettings, nil
}

func waitForBatchPoolPendingResizeOperation(ctx context.Context, client *pool.PoolClient, id pool.PoolId) error {
	// waiting for the pool to be in steady state
	log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped...")
	isSteady := false
	for !isSteady {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.AllocationState != nil {
			isSteady = *resp.Model.Properties.AllocationState == pool.AllocationStateSteady
			if isSteady {
				break
			}
		}
		time.Sleep(time.Second * 30)
		log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped... New try in 30 seconds...")
	}
	return nil
}

// validateUserIdentity validates that the user identity for a start task has been well specified
// it should have a auto_user block or a user_name defined, but not both at the same time.
func validateUserIdentity(userIdentity *pool.UserIdentity) error {
	if userIdentity == nil {
		return errors.New("user_identity block needs to be specified")
	}

	if userIdentity.AutoUser == nil && userIdentity.UserName == nil {
		return errors.New("auto_user or user_name needs to be specified in the user_identity block")
	}

	if userIdentity.AutoUser != nil && userIdentity.UserName != nil && *userIdentity.UserName != "" {
		return errors.New("auto_user and user_name cannot be specified in the user_identity at the same time")
	}

	return nil
}

func validateBatchPoolCrossFieldRules(pool *pool.PoolProperties) error {
	// Perform validation across multiple fields as per https://docs.microsoft.com/en-us/rest/api/batchmanagement/pool/create#resourcefile

	if pool.StartTask != nil {
		startTask := *pool.StartTask
		if startTask.ResourceFiles != nil {
			for _, referenceFile := range *startTask.ResourceFiles {
				// Must specify exactly one of AutoStorageContainerName, StorageContainerURL or HttpUrl
				sourceCount := 0
				if referenceFile.AutoStorageContainerName != nil {
					sourceCount++
				}
				if referenceFile.StorageContainerURL != nil {
					sourceCount++
				}
				if referenceFile.HTTPURL != nil {
					sourceCount++
				}
				if sourceCount != 1 {
					return fmt.Errorf("exactly one of auto_storage_container_name, storage_container_url and http_url must be specified")
				}

				if referenceFile.BlobPrefix != nil {
					if referenceFile.AutoStorageContainerName == nil && referenceFile.StorageContainerURL == nil {
						return fmt.Errorf("auto_storage_container_name or storage_container_url must be specified when using blob_prefix")
					}
				}

				if referenceFile.HTTPURL != nil {
					if referenceFile.FilePath == nil {
						return fmt.Errorf("file_path must be specified when using http_url")
					}
				}
			}
		}
	}

	return nil
}

func startTaskSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"command_line": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"run_options": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"image_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"registry": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: containerRegistry(),
						},
					},
					"working_directory": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(pool.ContainerWorkingDirectoryTaskWorkingDirectory),
							string(pool.ContainerWorkingDirectoryContainerImageDefault),
						}, false),
					},
				},
			},
		},

		"task_retry_maximum": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"wait_for_success": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"common_environment_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"user_identity": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"user_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"start_task.0.user_identity.0.user_name", "start_task.0.user_identity.0.auto_user"},
					},
					"auto_user": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"elevation_level": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(pool.ElevationLevelNonAdmin),
									ValidateFunc: validation.StringInSlice([]string{
										string(pool.ElevationLevelNonAdmin),
										string(pool.ElevationLevelAdmin),
									}, false),
								},
								"scope": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(pool.AutoUserScopeTask),
									ValidateFunc: validation.StringInSlice([]string{
										string(pool.AutoUserScopeTask),
										string(pool.AutoUserScopePool),
									}, false),
								},
							},
						},
						AtLeastOneOf: []string{"start_task.0.user_identity.0.user_name", "start_task.0.user_identity.0.auto_user"},
					},
				},
			},
		},
		// lintignore:XS003
		"resource_file": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"auto_storage_container_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"blob_prefix": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"file_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"file_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"http_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"storage_container_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"user_assigned_identity_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func containerRegistry() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"registry_server": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user_assigned_identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
			Description:  "The User Assigned Identity to use for Container Registry access.",
		},
		"user_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
