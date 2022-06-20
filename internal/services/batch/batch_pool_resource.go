package batch

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2022-01-01/batch"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceBatchPool() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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
			_, err := parse.PoolID(id)
			return err
		}),
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PoolName,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},
			"vm_size": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
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
			//TODO not able to determine support application licenses
			"application_licenses": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"application_packages": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.ApplicationID,
						},
						"version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
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
			"deployment_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cloud_service_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"os_family": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"os_version": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      utils.String("*"),
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"virtual_machine_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
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
													Default:  string(batch.CachingTypeReadOnly),
													ValidateFunc: validation.StringInSlice([]string{
														string(batch.CachingTypeNone),
														string(batch.CachingTypeReadOnly),
														string(batch.CachingTypeReadWrite),
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
													Default:  batch.StorageAccountTypeStandardLRS,
													ValidateFunc: validation.StringInSlice([]string{
														string(batch.StorageAccountTypeStandardLRS),
														string(batch.StorageAccountTypePremiumLRS),
													}, false),
												},
											},
										},
									},
									"disk_encryption_configuration": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"disk_encryption_target": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(batch.DiskEncryptionTargetTemporaryDisk),
														string(batch.DiskEncryptionTargetOsDisk),
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
												"settings": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsJSON,
												},
												"protected_settings": {
													Type:      pluginsdk.TypeString,
													Optional:  true,
													Sensitive: true,
												},
												"provision_after_extensions": {
													Type:     pluginsdk.TypeList,
													Optional: true,
													Elem: &schema.Schema{
														Type:         pluginsdk.TypeString,
														Optional:     true,
														ValidateFunc: validation.StringIsNotEmpty,
													},
												},
											},
										},
									},
									"image_reference": {
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
									"license_type": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"node_agent_sku_id": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ForceNew: true,
									},
									"node_placement_configuration": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"policy": {
													Type:     pluginsdk.TypeString,
													Optional: true,
													Default:  string(batch.NodePlacementPolicyTypeRegional),
													ValidateFunc: validation.StringInSlice([]string{
														string(batch.NodePlacementPolicyTypeZonal),
														string(batch.NodePlacementPolicyTypeRegional),
													}, false),
												},
											},
										},
									},
									"os_disk": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"ephemeral_os_disk_settings": {
													Type:     pluginsdk.TypeList,
													Required: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"placement": {
																Type:     pluginsdk.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice(
																	[]string{
																		string(batch.DiffDiskPlacementCacheDisk),
																	}, false),
															},
														},
													},
												},
											},
										},
									},
									"windows_configuration": {
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
							},
						},
					},
				},
			},
			"display_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"fixed_scale": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"node_deallocation_option": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(batch.ComputeNodeDeallocationOptionRequeue),
							ValidateFunc: validation.StringInSlice([]string{
								string(batch.ComputeNodeDeallocationOptionRequeue),
								string(batch.ComputeNodeDeallocationOptionRetainedData),
								string(batch.ComputeNodeDeallocationOptionTaskCompletion),
								string(batch.ComputeNodeDeallocationOptionTerminate),
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
			"inter_node_communication": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(batch.InterNodeCommunicationStateEnabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(batch.InterNodeCommunicationStateEnabled),
					string(batch.InterNodeCommunicationStateDisabled),
				}, false),
			},
			"metadata": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"mount_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"azure_blob_file_system_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
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
										ExactlyOneOf: []string{"account_key", "sas_key", "identity_reference"},
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"sas_key": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ExactlyOneOf: []string{"account_key", "sas_key", "identity_reference"},
										Sensitive:    true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"identity_reference": {
										Type:         pluginsdk.TypeList,
										Optional:     true,
										ExactlyOneOf: []string{"account_key", "sas_key", "identity_reference"},
										Elem:         identityReference(),
									},
									"blobfuse_options": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"azure_file_share_configuration": {
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
						"cifs_mount_configuration": {
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
						"nfs_mount_configuration": {
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
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						//TODO test whether this is forcenew
						"dynamic_vnet_assignment_scope": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  string(batch.DynamicVNetAssignmentScopeNone),
							ValidateFunc: validation.StringInSlice([]string{
								string(batch.DynamicVNetAssignmentScopeNone),
								string(batch.DynamicVNetAssignmentScopeJob),
							}, false),
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
								string(batch.IPAddressProvisioningTypeBatchManaged),
								string(batch.IPAddressProvisioningTypeUserManaged),
								string(batch.IPAddressProvisioningTypeNoPublicIPAddresses),
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
											string(batch.InboundEndpointProtocolTCP),
											string(batch.InboundEndpointProtocolUDP),
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
														string(batch.NetworkSecurityGroupRuleAccessAllow),
														string(batch.NetworkSecurityGroupRuleAccessDeny),
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
													ForceNew: true,
													Elem: &schema.Schema{
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
						},
					},
				},
			},
			"start_task": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: startTaskSchema(),
				},
			},
			"task_scheduling_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"node_fill_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(batch.ComputeNodeFillTypeSpread),
								string(batch.ComputeNodeFillTypePack),
							}, false),
						},
					},
				},
			},
			"max_tasks_per_node": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"stop_pending_resize_operation": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
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
								string(batch.ElevationLevelNonAdmin),
								string(batch.ElevationLevelAdmin),
							}, false),
						},
						"linux_user_configuration": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"uid": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										RequiredWith: []string{"gid"},
									},
									"gid": {
										Type:         pluginsdk.TypeInt,
										Optional:     true,
										RequiredWith: []string{"uid"},
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
											string(batch.LoginModeBatch),
											string(batch.LoginModeInteractive),
										}, false),
									},
								},
							},
						},
					},
				},
			},
			"identity": commonschema.UserAssignedIdentityOptional(),
		},
	}
}

func resourceBatchPoolCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_batch_pool", id.ID())
		}
	}

	parameters := batch.Pool{
		PoolProperties: &batch.PoolProperties{
			VMSize:                 utils.String(d.Get("vm_size").(string)),
			DisplayName:            utils.String(d.Get("display_name").(string)),
			InterNodeCommunication: batch.InterNodeCommunicationState(d.Get("inter_node_communication").(string)),
			TaskSlotsPerNode:       utils.Int32(int32(d.Get("max_tasks_per_node").(int))),
		},
	}

	applicationLicences, err := ExpandBatchPoolApplicationLicenses(d)
	if err != nil {
		return fmt.Errorf(`expanding "application_licenses": %v`, err)
	}
	parameters.PoolProperties.ApplicationLicenses = applicationLicences

	applicationPackages, err := ExpendBatchPoolApplicationPackages(d)
	if err != nil {
		return fmt.Errorf(`expanding "application_packages": %v`, err)
	}
	parameters.PoolProperties.ApplicationPackages = applicationPackages

	taskSchedulingPolicy, err := ExpandBatchPoolTaskSchedulingPolicy(d)
	if err != nil {
		return fmt.Errorf(`expanding "task_scheduling_policy": %v`, err)
	}
	parameters.PoolProperties.TaskSchedulingPolicy = taskSchedulingPolicy

	userAccounts, err := ExpandBatchPoolUserAccounts(d)
	if err != nil {
		return fmt.Errorf(`expanding "user_accounts": %v`, err)
	}
	parameters.PoolProperties.UserAccounts = userAccounts

	identity, err := expandBatchPoolIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}
	parameters.Identity = identity

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("expanding scale settings: %+v", err)
	}
	parameters.PoolProperties.ScaleSettings = scaleSettings

	deploymentConfiguration, err := ExpandBatchPoolDeploymentConfiguration(d)
	if err != nil {
		return fmt.Errorf(`expanding "deployment_configuration": %v`, err)
	}
	parameters.PoolProperties.DeploymentConfiguration = deploymentConfiguration

	mountConfiguration, err := ExpandBatchPoolMountConfigurations(d)
	if err != nil {
		return fmt.Errorf(`expanding "mount_configuration": %v`, err)
	}
	parameters.PoolProperties.MountConfiguration = mountConfiguration

	//nodeAgentSkuID := d.Get("node_agent_sku_id").(string)
	//
	//storageImageReferenceSet := d.Get("storage_image_reference").([]interface{})
	//imageReference, err := ExpandBatchPoolImageReference(storageImageReferenceSet)
	//if err != nil {
	//	return fmt.Errorf("creating %s: %+v", id, err)
	//}

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

		parameters.PoolProperties.StartTask = startTask
	}

	//containerConfiguration, err := ExpandBatchPoolContainerConfiguration(d.Get("container_configuration").([]interface{}))
	//if err != nil {
	//	return fmt.Errorf("creating %s: %+v", id, err)
	//}

	//parameters.PoolProperties.DeploymentConfiguration = &batch.DeploymentConfiguration{
	//	VirtualMachineConfiguration: &batch.VirtualMachineConfiguration{
	//		NodeAgentSkuID:         &nodeAgentSkuID,
	//		ImageReference:         imageReference,
	//		ContainerConfiguration: containerConfiguration,
	//	},
	//}

	certificates := d.Get("certificate").([]interface{})
	certificateReferences, err := ExpandBatchPoolCertificateReferences(certificates)
	if err != nil {
		return fmt.Errorf("expanding `certificate`: %+v", err)
	}
	parameters.PoolProperties.Certificates = certificateReferences

	if err := validateBatchPoolCrossFieldRules(&parameters); err != nil {
		return err
	}

	metaDataRaw := d.Get("metadata").(map[string]interface{})
	parameters.PoolProperties.Metadata = ExpandBatchMetaData(metaDataRaw)

	networkConfiguration := d.Get("network_configuration").([]interface{})
	parameters.PoolProperties.NetworkConfiguration, err = ExpandBatchPoolNetworkConfiguration(networkConfiguration)
	if err != nil {
		return fmt.Errorf("expanding `network_configuration`: %+v", err)
	}

	_, err = client.Create(ctx, id.ResourceGroup, id.BatchAccountName, id.Name, parameters, "", "")
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	// if the pool is not Steady after the create operation, wait for it to be Steady
	if props := read.PoolProperties; props != nil && props.AllocationState != batch.AllocationStateSteady {
		if err = waitForBatchPoolPendingResizeOperation(ctx, client, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
			return fmt.Errorf("waiting for %s", id)
		}
	}

	return resourceBatchPoolRead(d, meta)
}

func resourceBatchPoolUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.PoolProperties.AllocationState != batch.AllocationStateSteady {
		log.Printf("[INFO] there is a pending resize operation on this pool...")
		stopPendingResizeOperation := d.Get("stop_pending_resize_operation").(bool)
		if !stopPendingResizeOperation {
			return fmt.Errorf("updating %s because of pending resize operation. Set flag `stop_pending_resize_operation` to true to force update", *id)
		}

		log.Printf("[INFO] stopping the pending resize operation on this pool...")
		if _, err = client.StopResize(ctx, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
			return fmt.Errorf("stopping resize operation for %s: %+v", *id, err)
		}

		// waiting for the pool to be in steady state
		if err = waitForBatchPoolPendingResizeOperation(ctx, client, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
			return fmt.Errorf("waiting for %s", *id)
		}
	}

	parameters := batch.Pool{
		PoolProperties: &batch.PoolProperties{},
	}

	identity, err := expandBatchPoolIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}
	parameters.Identity = identity

	applicationLicences, err := ExpandBatchPoolApplicationLicenses(d)
	if err != nil {
		return fmt.Errorf(`expanding "application_licenses": %v`, err)
	}
	parameters.PoolProperties.ApplicationLicenses = applicationLicences

	applicationPackages, err := ExpendBatchPoolApplicationPackages(d)
	if err != nil {
		return fmt.Errorf(`expanding "application_packages": %v`, err)
	}
	parameters.PoolProperties.ApplicationPackages = applicationPackages

	taskSchedulingPolicy, err := ExpandBatchPoolTaskSchedulingPolicy(d)
	if err != nil {
		return fmt.Errorf(`expanding "task_scheduling_policy": %v`, err)
	}
	parameters.PoolProperties.TaskSchedulingPolicy = taskSchedulingPolicy

	userAccounts, err := ExpandBatchPoolUserAccounts(d)
	if err != nil {
		return fmt.Errorf(`expanding "user_accounts": %v`, err)
	}
	parameters.PoolProperties.UserAccounts = userAccounts

	deploymentConfiguration, err := ExpandBatchPoolDeploymentConfiguration(d)
	if err != nil {
		return fmt.Errorf(`expanding "deployment_configuration": %v`, err)
	}
	parameters.PoolProperties.DeploymentConfiguration = deploymentConfiguration

	mountConfiguration, err := ExpandBatchPoolMountConfigurations(d)
	if err != nil {
		return fmt.Errorf(`expanding "mount_configuration": %v`, err)
	}
	parameters.PoolProperties.MountConfiguration = mountConfiguration

	networkConfiguration := d.Get("network_configuration").([]interface{})
	parameters.PoolProperties.NetworkConfiguration, err = ExpandBatchPoolNetworkConfiguration(networkConfiguration)
	if err != nil {
		return fmt.Errorf("expanding `network_configuration`: %+v", err)
	}

	scaleSettings, err := expandBatchPoolScaleSettings(d)
	if err != nil {
		return fmt.Errorf("expanding scale settings: %+v", err)
	}
	parameters.PoolProperties.ScaleSettings = scaleSettings

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

		parameters.PoolProperties.StartTask = startTask
	}

	certificates := d.Get("certificate").([]interface{})
	certificateReferences, err := ExpandBatchPoolCertificateReferences(certificates)
	if err != nil {
		return fmt.Errorf("expanding `certificate`: %+v", err)
	}
	parameters.PoolProperties.Certificates = certificateReferences

	if err := validateBatchPoolCrossFieldRules(&parameters); err != nil {
		return err
	}

	if d.HasChange("metadata") {
		log.Printf("[DEBUG] Updating the MetaData for %s", *id)
		metaDataRaw := d.Get("metadata").(map[string]interface{})
		parameters.PoolProperties.Metadata = ExpandBatchMetaData(metaDataRaw)
	}

	result, err := client.Update(ctx, id.ResourceGroup, id.BatchAccountName, id.Name, parameters, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// if the pool is not Steady after the update, wait for it to be Steady
	if props := result.PoolProperties; props != nil && props.AllocationState != batch.AllocationStateSteady {
		if err := waitForBatchPoolPendingResizeOperation(ctx, client, id.ResourceGroup, id.BatchAccountName, id.Name); err != nil {
			return fmt.Errorf("waiting for %s", *id)
		}
	}

	return resourceBatchPoolRead(d, meta)
}

func resourceBatchPoolRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", *id)
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroup)

	identity, err := flattenBatchPoolIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.PoolProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("vm_size", props.VMSize)
		d.Set("max_tasks_per_node", props.TaskSlotsPerNode)
		d.Set("inter_node_communication", string(props.InterNodeCommunication))

		if props.ApplicationLicenses != nil {
			applicationLicenses := make([]interface{}, len(*props.ApplicationLicenses))
			for _, license := range *props.ApplicationLicenses {
				applicationLicenses = append(applicationLicenses, license)
			}
			d.Set("application_licenses", applicationLicenses)
		}

		if props.ApplicationPackages != nil {
			applicationPackages := make([]interface{}, len(*props.ApplicationPackages))
			for _, pkg := range *props.ApplicationPackages {
				appPkg := make(map[string]interface{}, 1)
				appPkg["id"] = *pkg.ID
				if pkg.Version != nil {
					appPkg["version"] = *pkg.Version
				}
				applicationPackages = append(applicationPackages, appPkg)
			}
			d.Set("application_packages", applicationPackages)
		}

		if props.TaskSchedulingPolicy != nil {
			taskSchedulingPolicy := make([]interface{}, 1)
			nodeFillType := make(map[string]interface{}, 1)
			nodeFillType["node_fill_type"] = string(props.TaskSchedulingPolicy.NodeFillType)
			taskSchedulingPolicy = append(taskSchedulingPolicy, nodeFillType)
			d.Set("task_scheduling_policy", taskSchedulingPolicy)
		}

		if props.UserAccounts != nil {
			userAccounts := make([]interface{}, len(*props.UserAccounts))
			for _, userAccount := range *props.UserAccounts {
				userAccounts = append(userAccounts, flattenBatchPoolUserAccount(&userAccount))
			}
			d.Set("user_accounts", userAccounts)
		}

		if props.MountConfiguration != nil {
			mountConfigs := make([]interface{}, len(*props.MountConfiguration))
			for _, mountConfig := range *props.MountConfiguration {
				mountConfigs = append(mountConfigs, flattenBatchPoolMountConfig(&mountConfig))
			}
			d.Set("mount_configuration", mountConfigs)
		}

		if props.DeploymentConfiguration != nil {
			deploymentConfigMap := make(map[string]interface{}, 0)
			if props.DeploymentConfiguration.CloudServiceConfiguration != nil {
				deploymentConfigMap["cloud_service_configuration"] = flattenBatchPoolCloudServiceConfiguration(props.DeploymentConfiguration.CloudServiceConfiguration)

			}
			if props.DeploymentConfiguration.VirtualMachineConfiguration != nil {
				deploymentConfigMap["cloud_service_configuration"] = flattenBatchPoolVirtualMachineConfiguration(props.DeploymentConfiguration.VirtualMachineConfiguration)
			}

			d.Set("deployment_configuration", deploymentConfigMap)
		}

		if scaleSettings := props.ScaleSettings; scaleSettings != nil {
			if err := d.Set("auto_scale", flattenBatchPoolAutoScaleSettings(scaleSettings.AutoScale)); err != nil {
				return fmt.Errorf("flattening `auto_scale`: %+v", err)
			}
			if err := d.Set("fixed_scale", flattenBatchPoolFixedScaleSettings(scaleSettings.FixedScale)); err != nil {
				return fmt.Errorf("flattening `fixed_scale `: %+v", err)
			}
		}

		//if props.DeploymentConfiguration != nil &&
		//	props.DeploymentConfiguration.VirtualMachineConfiguration != nil &&
		//	props.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference != nil {
		//	imageReference := props.DeploymentConfiguration.VirtualMachineConfiguration.ImageReference
		//
		//	d.Set("storage_image_reference", flattenBatchPoolImageReference(imageReference))
		//	d.Set("node_agent_sku_id", props.DeploymentConfiguration.VirtualMachineConfiguration.NodeAgentSkuID)
		//}
		//
		//if dcfg := props.DeploymentConfiguration; dcfg != nil {
		//	if vmcfg := dcfg.VirtualMachineConfiguration; vmcfg != nil {
		//		d.Set("container_configuration", flattenBatchPoolContainerConfiguration(d, vmcfg.ContainerConfiguration))
		//	}
		//}

		if err := d.Set("certificate", flattenBatchPoolCertificateReferences(props.Certificates)); err != nil {
			return fmt.Errorf("flattening `certificate`: %+v", err)
		}

		d.Set("start_task", flattenBatchPoolStartTask(props.StartTask))
		d.Set("metadata", FlattenBatchMetaData(props.Metadata))

		if err := d.Set("network_configuration", flattenBatchPoolNetworkConfiguration(props.NetworkConfiguration)); err != nil {
			return fmt.Errorf("setting `network_configuration`: %v", err)
		}
	}

	return nil
}

func resourceBatchPoolDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.PoolClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PoolID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.BatchAccountName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}
	return nil
}

func waitForBatchPoolPendingResizeOperation(ctx context.Context, client *batch.PoolClient, resourceGroup string, accountName string, poolName string) error {
	// waiting for the pool to be in steady state
	log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped...")
	isSteady := false
	for !isSteady {
		resp, err := client.Get(ctx, resourceGroup, accountName, poolName)
		if err != nil {
			return fmt.Errorf("retrieving the Batch pool %q (Resource Group %q): %+v", poolName, resourceGroup, err)
		}

		isSteady = resp.PoolProperties.AllocationState == batch.AllocationStateSteady
		time.Sleep(time.Second * 30)
		log.Printf("[INFO] waiting for the pending resize operation on this pool to be stopped... New try in 30 seconds...")
	}

	return nil
}

// validateUserIdentity validates that the user identity for a start task has been well specified
// it should have a auto_user block or a user_name defined, but not both at the same time.
func validateUserIdentity(userIdentity *batch.UserIdentity) error {
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

func validateBatchPoolCrossFieldRules(pool *batch.Pool) error {
	// Perform validation across multiple fields as per https://docs.microsoft.com/en-us/rest/api/batchmanagement/pool/create#resourcefile

	if pool.StartTask != nil {
		startTask := *pool.StartTask
		if startTask.ResourceFiles != nil {
			for _, referenceFile := range *startTask.ResourceFiles {
				// Must specify exactly one of AutoStorageContainerName, StorageContainerUrl or HttpUrl
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

func identityReference() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"identity_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},
	}
}

func containerRegistry() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"username": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"password"},
			AtLeastOneOf: []string{"username", "identity_reference"},
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"username"},
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"registry_server": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      utils.String("docker.io"),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"identity_reference": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			ForceNew:     true,
			AtLeastOneOf: []string{"username", "identity_reference"},
			Elem:         identityReference(),
		},
	}
}

func startTaskSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"command_line": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"container_run_options": {
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
							string(batch.ContainerWorkingDirectoryTaskWorkingDirectory),
							string(batch.ContainerWorkingDirectoryContainerImageDefault),
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
									Default:  string(batch.ElevationLevelNonAdmin),
									ValidateFunc: validation.StringInSlice([]string{
										string(batch.ElevationLevelNonAdmin),
										string(batch.ElevationLevelAdmin),
									}, false),
								},
								"scope": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									Default:  string(batch.AutoUserScopeTask),
									ValidateFunc: validation.StringInSlice([]string{
										string(batch.AutoUserScopeTask),
										string(batch.AutoUserScopePool),
									}, false),
								},
							},
						},
						AtLeastOneOf: []string{"start_task.0.user_identity.0.user_name", "start_task.0.user_identity.0.auto_user"},
					},
				},
			},
		},
		//lintignore:XS003
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
					"identity_reference": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     identityReference(),
					},
					"storage_container_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}
