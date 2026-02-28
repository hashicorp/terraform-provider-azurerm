package compute

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// This is to serve input validation for the customizer block.
var (
	fieldsOfFileCustomizer           = []string{"file_source_uri", "file_sha256_checksum", "file_destination_path"}
	fieldsOfPowerShellCustomizer     = []string{"powershell_commands", "powershell_run_as_system", "powershell_run_elevated", "powershell_script_uri", "powershell_sha256_checksum", "powershell_valid_exit_codes"}
	fieldsOfShellCustomizer          = []string{"shell_commands", "shell_script_uri", "shell_sha256_checksum"}
	fieldsOfWindowsRestartCustomizer = []string{"windows_restart_check_command", "windows_restart_command", "windows_restart_timeout"}
	fieldsOfWindowsUpdateCustomizer  = []string{"windows_update_filters", "windows_update_search_criteria", "windows_update_limit"}
)

func resourceImageBuilderTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmImageBuilderTemplateCreate,
		Read:   resourceArmImageBuilderTemplateRead,
		Update: resourceArmImageBuilderTemplateUpdate,
		Delete: resourceArmImageBuilderTemplateDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualmachineimagetemplate.ParseImageTemplateID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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

						// If not specify this property but only "file_source_uri" to the service, the service will calculate the sha256 of the file and return it.
						// So set this property as Computed for possible future usage. So forth to other similar properties in the `customizer` block.
						"file_sha256_checksum": {
							Type:         schema.TypeString,
							Optional:     true,
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
							Type:         schema.TypeString,
							Optional:     true,
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
							Type:         schema.TypeString,
							Computed:     true,
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
				ValidateFunc: validate.DiskSizeGB,
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
				ValidateFunc: validate.SharedImageVersionID,
				ExactlyOneOf: []string{"source_managed_image_id", "source_platform_image", "source_shared_image_version_id"},
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceArmImageBuilderTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Compute.VirtualMachineImageTemplateClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	id := virtualmachineimagetemplate.NewImageTemplateID(subscriptionId, resourceGroup, name)
	resp, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for the presence of existing Image Builder Template %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_image_builder_template", id.ID())
	}

	location := d.Get("location").(string)

	distribution, err := expandBasicImageTemplateDistributor(d.Get("distributions").([]interface{}), subscriptionId)
	if err != nil {
		return err
	}

	customizer, err := expandBasicImageTemplateCustomizer(d)
	if err != nil {
		return err
	}
	t := d.Get("tags").(map[string]interface{})

	identityExpanded, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := virtualmachineimagetemplate.ImageTemplate{
		Location: location,
		Identity: *identityExpanded,
		Properties: &virtualmachineimagetemplate.ImageTemplateProperties{
			VMProfile: &virtualmachineimagetemplate.ImageTemplateVMProfile{
				VMSize:       pointer.To(d.Get("size").(string)),
				OsDiskSizeGB: pointer.To(int64(d.Get("disk_size_gb").(int))),
			},

			Source:                expandBasicImageTemplateSource(d.Get("source_managed_image_id").(string), d.Get("source_platform_image").([]interface{}), d.Get("source_shared_image_version_id").(string)),
			Distribute:            *distribution,
			Customize:             customizer,
			BuildTimeoutInMinutes: pointer.To(int64(d.Get("build_timeout_minutes").(int))),
		},

		Tags: tags.Expand(t),
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		parameters.Properties.VMProfile.VnetConfig = &virtualmachineimagetemplate.VirtualNetworkConfig{
			SubnetId: pointer.To(v.(string)),
		}
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating image builder template %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceArmImageBuilderTemplateRead(d, meta)
}

func resourceArmImageBuilderTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineImageTemplateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachineimagetemplate.ParseImageTemplateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
	}

	d.Set("name", id.ImageTemplateName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		identityFlattened, err := identity.FlattenUserAssignedMap(&model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", identityFlattened); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if imageTemplateProperties := resp.Model.Properties; imageTemplateProperties != nil {
			managedImageId, platformImage, sharedImageVersionId := flattenBasicImageTemplateSource(imageTemplateProperties.Source)

			// only one among managedImageId / platformImage / sharedImageVersionId would be returned.
			if managedImageId != "" {
				if err := d.Set("source_managed_image_id", managedImageId); err != nil {
					return fmt.Errorf("setting image template managed image source: %+v", err)
				}
			}

			if len(platformImage) > 0 {
				if err := d.Set("source_platform_image", platformImage); err != nil {
					return fmt.Errorf("setting image template platform image source: %+v", err)
				}
			}

			if sharedImageVersionId != "" {
				if err := d.Set("source_shared_image_version_id", sharedImageVersionId); err != nil {
					return fmt.Errorf("setting image template shared image version source: %+v", err)
				}
			}

			flattenedDistributions, err := flattenBasicImageTemplateDistributor(&imageTemplateProperties.Distribute)
			if err != nil {
				return err
			}

			if err := d.Set("distributions", flattenedDistributions); err != nil {
				return fmt.Errorf("setting image template distribution: %+v", err)
			}

			if err := d.Set("customizer", flattenBasicImageTemplateCustomizer(imageTemplateProperties.Customize)); err != nil {
				return fmt.Errorf("setting `customizer`: %+v", err)
			}

			if err := d.Set("build_timeout_minutes", imageTemplateProperties.BuildTimeoutInMinutes); err != nil {
				return fmt.Errorf("setting `build timeout minutes`: %+v", err)
			}

			if vmProfile := resp.Model.Properties.VMProfile; vmProfile != nil {
				d.Set("size", vmProfile.VMSize)

				d.Set("disk_size_gb", vmProfile.OsDiskSizeGB)

				if vnetConfig := vmProfile.VnetConfig; vnetConfig != nil {
					d.Set("subnet_id", vnetConfig.SubnetId)
				}
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

	}

	return nil
}

func resourceArmImageBuilderTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineImageTemplateClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachineimagetemplate.ParseImageTemplateID(d.Id())
	if err != nil {
		return err
	}

	parameters := virtualmachineimagetemplate.ImageTemplateUpdateParameters{}

	if d.HasChange("identity") {
		identityExpanded, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		parameters.Identity = identityExpanded
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	err = client.UpdateThenPoll(ctx, *id, parameters)
	if err != nil {
		return fmt.Errorf("updating image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
	}

	return resourceArmImageBuilderTemplateRead(d, meta)
}

func resourceArmImageBuilderTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VirtualMachineImageTemplateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachineimagetemplate.ParseImageTemplateID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
	}

	return nil
}

func expandBasicImageTemplateSource(managedImageId string, platformImage []interface{}, sharedImageId string) virtualmachineimagetemplate.ImageTemplateSource {
	if len(managedImageId) != 0 {
		return &virtualmachineimagetemplate.ImageTemplateManagedImageSource{
			ImageId: managedImageId,
		}
	}

	if len(platformImage) != 0 && platformImage[0] != nil {
		v := platformImage[0].(map[string]interface{})
		result := &virtualmachineimagetemplate.ImageTemplatePlatformImageSource{
			Publisher: pointer.To(v["publisher"].(string)),
			Offer:     pointer.To(v["offer"].(string)),
			Sku:       pointer.To(v["sku"].(string)),
			Version:   pointer.To(v["version"].(string)),
		}

		planRaw := v["plan"].([]interface{})
		if len(planRaw) > 0 {
			planTemp := planRaw[0].(map[string]interface{})
			result.PlanInfo = &virtualmachineimagetemplate.PlatformImagePurchasePlan{
				PlanName:      planTemp["name"].(string),
				PlanProduct:   planTemp["product"].(string),
				PlanPublisher: planTemp["publisher"].(string),
			}
		}

		return result
	}

	if len(sharedImageId) != 0 {
		return &virtualmachineimagetemplate.ImageTemplateSharedImageVersionSource{
			ImageVersionId: sharedImageId,
		}
	}

	// since there is ExactlyOneOf on the source schema, below nil won't be reached
	return nil
}

func flattenBasicImageTemplateSource(input virtualmachineimagetemplate.ImageTemplateSource) (string, []interface{}, string) {
	if input != nil {
		switch source := input.(type) {
		case virtualmachineimagetemplate.ImageTemplateManagedImageSource:
			return flattenSourceManagedImage(&source), nil, ""
		case virtualmachineimagetemplate.ImageTemplatePlatformImageSource:
			return "", flattenSourcePlatformImage(&source), ""
		case virtualmachineimagetemplate.ImageTemplateSharedImageVersionSource:
			return "", nil, flattenSourceSharedImageVersion(&source)
		}
	}

	return "", nil, ""
}

func flattenSourceManagedImage(input *virtualmachineimagetemplate.ImageTemplateManagedImageSource) string {
	if input != nil && input.ImageId != "" {
		return input.ImageId
	}

	return ""
}

func flattenSourcePlatformImage(input *virtualmachineimagetemplate.ImageTemplatePlatformImageSource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}

	offer := ""
	if input.Offer != nil {
		offer = *input.Offer
	}

	sku := ""
	if input.Sku != nil {
		sku = *input.Sku
	}

	version := ""
	if input.Version != nil {
		version = *input.Version
	}

	return []interface{}{
		map[string]interface{}{
			"publisher": publisher,
			"offer":     offer,
			"sku":       sku,
			"version":   version,
			"plan":      flattenImageBuilderTemplateSourcePlatformImagePlan(input.PlanInfo),
		},
	}
}

func flattenImageBuilderTemplateSourcePlatformImagePlan(input *virtualmachineimagetemplate.PlatformImagePurchasePlan) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.PlanName != "" {
		name = input.PlanName
	}

	product := ""
	if input.PlanProduct != "" {
		product = input.PlanProduct
	}

	publisher := ""
	if input.PlanPublisher != "" {
		publisher = input.PlanPublisher
	}

	return []interface{}{
		map[string]interface{}{
			"name":      name,
			"product":   product,
			"publisher": publisher,
		},
	}
}

func flattenSourceSharedImageVersion(input *virtualmachineimagetemplate.ImageTemplateSharedImageVersionSource) string {
	if input != nil && input.ImageVersionId != "" {
		return input.ImageVersionId
	}

	return ""
}

func expandBasicImageTemplateDistributor(distributions []interface{}, subscriptionId string) (*[]virtualmachineimagetemplate.ImageTemplateDistributor, error) {
	results := make([]virtualmachineimagetemplate.ImageTemplateDistributor, 0)
	runOutputNameSet := make(map[string]bool)

	if len(distributions) > 0 {
		for _, v := range distributions {
			if v != nil {
				distributionItems := v.(map[string]interface{})
				managedImages := distributionItems["managed_image"].([]interface{})
				sharedImages := distributionItems["shared_image"].([]interface{})
				vhds := distributionItems["vhd"].([]interface{})

				if len(managedImages) > 0 {
					for _, managedImageRaw := range managedImages {
						if managedImageRaw != nil {
							managedImage := managedImageRaw.(map[string]interface{})
							resourceGroupName := managedImage["resource_group_name"].(string)
							runOutputName := managedImage["run_output_name"].(string)

							_, existing := runOutputNameSet[runOutputName]
							if existing {
								return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
							} else {
								runOutputNameSet[runOutputName] = true

								results = append(results, virtualmachineimagetemplate.ImageTemplateManagedImageDistributor{
									ImageId:       "/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroupName + "/providers/Microsoft.Compute/images/" + managedImage["name"].(string),
									Location:      managedImage["location"].(string),
									RunOutputName: runOutputName,
									ArtifactTags:  tags.Expand(managedImage["tags"].(map[string]interface{})),
								})
							}
						}
					}
				}

				if len(sharedImages) > 0 {
					for _, sharedImageRaw := range sharedImages {
						if sharedImageRaw != nil {
							sharedImage := sharedImageRaw.(map[string]interface{})
							runOutputName := sharedImage["run_output_name"].(string)

							_, existing := runOutputNameSet[runOutputName]
							if existing {
								return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
							} else {
								runOutputNameSet[runOutputName] = true
								results = append(results, virtualmachineimagetemplate.ImageTemplateSharedImageDistributor{
									GalleryImageId:     sharedImage["id"].(string),
									ReplicationRegions: expandImageTemplateSharedImageDistributorReplicaRegions(sharedImage["replica_regions"].([]interface{})),
									RunOutputName:      sharedImage["run_output_name"].(string),
									ExcludeFromLatest:  pointer.To(sharedImage["exclude_from_latest"].(bool)),
									StorageAccountType: pointer.To(virtualmachineimagetemplate.SharedImageStorageAccountType(sharedImage["storage_account_type"].(string))),
									ArtifactTags:       tags.Expand(sharedImage["tags"].(map[string]interface{})),
									Versioning:         expandImageTemplateSharedImageDistributorVersioning(sharedImage["versioning"].([]interface{})),
								})
							}
						}
					}
				}

				if len(vhds) > 0 {
					for _, vhdRaw := range vhds {
						if vhdRaw != nil {
							vhd := vhdRaw.(map[string]interface{})
							runOutputName := vhd["run_output_name"].(string)

							_, existing := runOutputNameSet[runOutputName]
							if existing {
								return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
							} else {
								runOutputNameSet[runOutputName] = true
								results = append(results, virtualmachineimagetemplate.ImageTemplateVhdDistributor{
									RunOutputName: vhd["run_output_name"].(string),
									ArtifactTags:  tags.Expand(vhd["tags"].(map[string]interface{})),
								})
							}
						}
					}
				}
			} else {
				return &results, fmt.Errorf("at least one of `managed_image`, `shared_image` and `vhd` is required to specify in `distributions`")
			}
		}
	}

	return &results, nil
}

func flattenBasicImageTemplateDistributor(input *[]virtualmachineimagetemplate.ImageTemplateDistributor) ([]interface{}, error) {
	results := make([]interface{}, 0)

	distributionManagedImages := make([]interface{}, 0)
	distributionSharedImages := make([]interface{}, 0)
	distributionVhds := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			switch distribute := v.(type) {
			case virtualmachineimagetemplate.ImageTemplateManagedImageDistributor:
				flattenedManagedImg, err := flattenDistributionManagedImage(&distribute)
				if err != nil {
					return nil, fmt.Errorf("setting image template managed image source: %+v", err)
				}

				distributionManagedImages = append(distributionManagedImages, flattenedManagedImg)
			case virtualmachineimagetemplate.ImageTemplateSharedImageDistributor:
				distributionSharedImages = append(distributionSharedImages, flattenDistributionSharedImage(&distribute))
			case virtualmachineimagetemplate.ImageTemplateVhdDistributor:
				distributionVhds = append(distributionVhds, flattenDistributionVhd(&distribute))
			}
		}
	}

	output := map[string]interface{}{
		"managed_image": distributionManagedImages,
		"shared_image":  distributionSharedImages,
		"vhd":           distributionVhds,
	}

	results = append(results, output)
	return results, nil
}

func expandImageTemplateSharedImageDistributorReplicaRegions(input []interface{}) *[]string {
	if input == nil {
		return nil
	}

	result := make([]string, 0)
	for _, v := range input {
		result = append(result, v.(map[string]interface{})["name"].(string))
	}

	return &result
}

func expandImageTemplateSharedImageDistributorVersioning(input []interface{}) virtualmachineimagetemplate.DistributeVersioner {
	var result virtualmachineimagetemplate.DistributeVersioner
	if len(input) > 0 {
		config := input[0].(map[string]interface{})

		if v := config["scheme"].(string); v == "Latest" {
			result = virtualmachineimagetemplate.DistributeVersionerLatest{
				Scheme: config["scheme"].(string),
				Major:  pointer.To(config["major"].(int64)),
			}
		}

		if v := config["scheme"].(string); v == "Source" {
			result = virtualmachineimagetemplate.DistributeVersionerSource{
				Scheme: config["scheme"].(string),
			}
		}
	}
	return result
}

func flattenImageTemplateSharedImageDistributorVersioning(input virtualmachineimagetemplate.DistributeVersioner) []interface{} {
	results := make([]interface{}, 0)
	output := make(map[string]interface{})

	switch versioner := input.(type) {
	case virtualmachineimagetemplate.DistributeVersionerLatest:
		output["scheme"] = versioner.Scheme
		output["major"] = versioner.Major

	case virtualmachineimagetemplate.DistributeVersionerSource:
		output["scheme"] = versioner.Scheme
	}

	results = append(results, output)

	return results
}

func flattenImageTemplateSharedImageDistributorReplicaRegions(input *[]string) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			output := make(map[string]interface{})
			output["name"] = v
			results = append(results, output)
		}
	}

	return results
}

func flattenDistributionManagedImage(input *virtualmachineimagetemplate.ImageTemplateManagedImageDistributor) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if input == nil {
		return result, nil
	}

	imageName := ""
	resourceGroupName := ""
	if input.ImageId != "" {
		imageNameReturned, resourceGroupNameReturned, err := imageBuilderTemplateManagedImageNameAndResourceGroupName(input.ImageId)
		if err != nil {
			return result, err
		}

		imageName = imageNameReturned
		resourceGroupName = resourceGroupNameReturned
	}

	result["name"] = imageName
	result["resource_group_name"] = resourceGroupName

	if input.Location != "" {
		result["location"] = input.Location
	}

	if input.RunOutputName != "" {
		result["run_output_name"] = input.RunOutputName
	}

	if input.ArtifactTags != nil {
		result["tags"] = tags.Flatten(input.ArtifactTags)
	}

	return result, nil
}

func flattenDistributionSharedImage(input *virtualmachineimagetemplate.ImageTemplateSharedImageDistributor) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	if input.GalleryImageId != "" {
		results["id"] = input.GalleryImageId
	}

	if input.ReplicationRegions != nil {
		results["replica_regions"] = flattenImageTemplateSharedImageDistributorReplicaRegions(input.ReplicationRegions)
	}

	if input.RunOutputName != "" {
		results["run_output_name"] = input.RunOutputName
	}

	if input.ExcludeFromLatest != nil {
		results["exclude_from_latest"] = *input.ExcludeFromLatest
	}

	if input.StorageAccountType != nil {
		results["storage_account_type"] = string(*input.StorageAccountType)
	}

	if input.ArtifactTags != nil {
		results["tags"] = tags.Flatten(input.ArtifactTags)
	}

	results["versioning"] = flattenImageTemplateSharedImageDistributorVersioning(input.Versioning)
	return results
}

func flattenDistributionVhd(input *virtualmachineimagetemplate.ImageTemplateVhdDistributor) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	if input.RunOutputName != "" {
		results["run_output_name"] = input.RunOutputName
	}

	if input.ArtifactTags != nil {
		results["tags"] = tags.Flatten(input.ArtifactTags)
	}

	return results
}

// Passing d as input rather its business subset because this function needs d.GetOK() to verify users' input to the `customizer` block is valid.
func expandBasicImageTemplateCustomizer(d *schema.ResourceData) (*[]virtualmachineimagetemplate.ImageTemplateCustomizer, error) {
	input := d.Get("customizer").([]interface{})

	if len(input) == 0 {
		return nil, nil
	}

	results := make([]virtualmachineimagetemplate.ImageTemplateCustomizer, 0)

	// This index serves identifying the invalid input customizer block from the customizer block list.
	i := 0

	for _, customizerRaw := range input {
		customizer := customizerRaw.(map[string]interface{})
		t := customizer["type"].(string)

		switch t {
		case "File":
			if err := validateImageTemplateCustomizerInputForFileType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateFileCustomizer{
				Name:           pointer.To(customizer["name"].(string)),
				SourceUri:      pointer.To(customizer["file_source_uri"].(string)),
				Sha256Checksum: pointer.To(customizer["file_sha256_checksum"].(string)),
				Destination:    pointer.To(customizer["file_destination_path"].(string)),
			})

			i++
		case "PowerShell":
			if err := validateImageTemplateCustomizerInputForPowerShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer{
				Name:           pointer.To(customizer["name"].(string)),
				ScriptUri:      pointer.To(customizer["powershell_script_uri"].(string)),
				Sha256Checksum: pointer.To(customizer["powershell_sha256_checksum"].(string)),
				Inline:         utils.ExpandStringSlice(customizer["powershell_commands"].([]interface{})),
				RunAsSystem:    pointer.To(customizer["powershell_run_as_system"].(bool)),
				RunElevated:    pointer.To(customizer["powershell_run_elevated"].(bool)),
				ValidExitCodes: utils.ExpandInt64Slice(customizer["powershell_valid_exit_codes"].([]interface{})),
			})

			i++
		case "Shell":
			if err := validateImageTemplateCustomizerInputForShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateShellCustomizer{
				Name:           pointer.To(customizer["name"].(string)),
				ScriptUri:      pointer.To(customizer["shell_script_uri"].(string)),
				Sha256Checksum: pointer.To(customizer["shell_sha256_checksum"].(string)),
				Inline:         utils.ExpandStringSlice(customizer["shell_commands"].([]interface{})),
			})

			i++
		case "WindowsRestart":
			if err := validateImageTemplateCustomizerInputForWindowsRestartType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateRestartCustomizer{
				Name:                pointer.To(customizer["name"].(string)),
				RestartCommand:      pointer.To(customizer["windows_restart_command"].(string)),
				RestartCheckCommand: pointer.To(customizer["windows_restart_check_command"].(string)),
				RestartTimeout:      pointer.To(customizer["windows_restart_timeout"].(string)),
			})

			i++
		case "WindowsUpdate":
			if err := validateImageTemplateCustomizerInputForWindowsUpdateType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer{
				Name:           pointer.To(customizer["name"].(string)),
				SearchCriteria: pointer.To(customizer["windows_update_search_criteria"].(string)),
				Filters:        utils.ExpandStringSlice(customizer["windows_update_filters"].([]interface{})),
				UpdateLimit:    pointer.To(int64(customizer["windows_update_limit"].(int))),
			})

			i++
		}
	}

	return &results, nil
}

func flattenBasicImageTemplateCustomizer(input *[]virtualmachineimagetemplate.ImageTemplateCustomizer) []interface{} {
	customizerList := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			switch customizer := v.(type) {
			case virtualmachineimagetemplate.ImageTemplateFileCustomizer:
				customizerList = append(customizerList, flattenCustomizerFile(&customizer))
			case virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer:
				customizerList = append(customizerList, flattenCustomizerPowerShell(&customizer))
			case virtualmachineimagetemplate.ImageTemplateShellCustomizer:
				customizerList = append(customizerList, flattenCustomizerShell(&customizer))
			case virtualmachineimagetemplate.ImageTemplateRestartCustomizer:
				customizerList = append(customizerList, flattenCustomizerWindowsRestart(&customizer))
			case virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer:
				customizerList = append(customizerList, flattenCustomizerWindowsUpdate(&customizer))
			}
		}

		return customizerList
	}

	return customizerList
}

func flattenCustomizerFile(input *virtualmachineimagetemplate.ImageTemplateFileCustomizer) map[string]interface{} {
	if input == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["type"] = "File"

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	result["name"] = name

	if input.SourceUri != nil {
		result["file_source_uri"] = *input.SourceUri
	}

	if input.Sha256Checksum != nil {
		result["file_sha256_checksum"] = *input.Sha256Checksum
	}

	if input.Destination != nil {
		result["file_destination_path"] = *input.Destination
	}

	return result
}

func flattenCustomizerShell(input *virtualmachineimagetemplate.ImageTemplateShellCustomizer) map[string]interface{} {
	if input == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["type"] = "Shell"

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	result["name"] = name

	if input.ScriptUri != nil {
		result["shell_script_uri"] = *input.ScriptUri
	}

	if input.Sha256Checksum != nil {
		result["shell_sha256_checksum"] = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result["shell_commands"] = *input.Inline
	}

	return result
}

func flattenCustomizerPowerShell(input *virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer) map[string]interface{} {
	if input == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["type"] = "PowerShell"

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	result["name"] = name

	if input.ScriptUri != nil {
		result["powershell_script_uri"] = *input.ScriptUri
	}

	if input.Sha256Checksum != nil {
		result["powershell_sha256_checksum"] = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result["powershell_commands"] = *input.Inline
	}

	if input.RunAsSystem != nil {
		result["powershell_run_as_system"] = *input.RunAsSystem
	}

	if input.RunElevated != nil {
		result["powershell_run_elevated"] = *input.RunElevated
	}

	if input.ValidExitCodes != nil {
		result["powershell_valid_exit_codes"] = *input.ValidExitCodes
	}

	return result
}

func flattenCustomizerWindowsRestart(input *virtualmachineimagetemplate.ImageTemplateRestartCustomizer) map[string]interface{} {
	if input == nil {
		return nil
	}

	result := make(map[string]interface{})
	result["type"] = "WindowsRestart"

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	result["name"] = name

	if input.RestartCommand != nil {
		result["windows_restart_command"] = *input.RestartCommand
	}

	if input.RestartCheckCommand != nil {
		result["windows_restart_check_command"] = *input.RestartCheckCommand
	}

	if input.RestartTimeout != nil {
		result["windows_restart_timeout"] = *input.RestartTimeout
	}

	return result
}

func flattenCustomizerWindowsUpdate(input *virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer) map[string]interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	result["type"] = "WindowsUpdate"

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	result["name"] = name

	if input.SearchCriteria != nil {
		result["windows_update_search_criteria"] = *input.SearchCriteria
	}

	if input.Filters != nil {
		result["windows_update_filters"] = *input.Filters
	}

	if input.UpdateLimit != nil {
		result["windows_update_limit"] = *input.UpdateLimit
	}

	return result
}

func distributionRunOutputNameSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		ValidateFunc: validation.StringMatch(
			regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
			"Run output name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
		),
	}
}

// In the `customizer` block, when a certain `type` say "File" is specified, do not allow users to specify fields that do not belong to `File`.
// Because once users specify those irrelevant fields, those values won't be sent to the backend service and they won't return from GET,
// thus next time there will be diff shown which will force creating a new resource. So at the very beginning forbid users from doing this.
func validateImageTemplateCustomizerInputForFileType(d *schema.ResourceData, customizer map[string]interface{}, index int) error {
	sourceUri := customizer["file_source_uri"].(string)
	if sourceUri == "" {
		return fmt.Errorf("`file_source_uri` must be specified if the customizer type is File")
	}

	destinationPath := customizer["file_destination_path"].(string)
	if destinationPath == "" {
		return fmt.Errorf("`destination_path` must be specified if the customizer type is File")
	}

	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "File"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForPowerShellType(d *schema.ResourceData, customizer map[string]interface{}, index int) error {
	scriptUri := customizer["powershell_script_uri"].(string)
	commandsRaw := customizer["powershell_commands"].([]interface{})

	if (scriptUri == "" && len(commandsRaw) == 0) ||
		(scriptUri != "" && len(commandsRaw) > 0) {
		return fmt.Errorf("exactly one of `powershell_script_uri` and `powershell_commands` must be specified if the customizer type is PowerShell")
	}

	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "PowerShell"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForShellType(d *schema.ResourceData, customizer map[string]interface{}, index int) error {
	scriptUri := customizer["shell_script_uri"].(string)
	commandsRaw := customizer["shell_commands"].([]interface{})

	if (scriptUri == "" && len(commandsRaw) == 0) ||
		(scriptUri != "" && len(commandsRaw) > 0) {
		return fmt.Errorf("exactly one of `shell_script_uri` and `shell_commands` must be specified if the customizer type is Shell")
	}

	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "Shell"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForWindowsRestartType(d *schema.ResourceData, index int) error {
	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsUpdateCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "WindowsRestart"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizerInputForWindowsUpdateType(d *schema.ResourceData, index int) error {
	excludeList := make([]string, 0)
	excludeList = append(excludeList, fieldsOfFileCustomizer...)
	excludeList = append(excludeList, fieldsOfPowerShellCustomizer...)
	excludeList = append(excludeList, fieldsOfShellCustomizer...)
	excludeList = append(excludeList, fieldsOfWindowsRestartCustomizer...)

	if err := validateImageTemplateCustomizer(d, excludeList, index, "WindowsUpdate"); err != nil {
		return err
	}

	return nil
}

func validateImageTemplateCustomizer(d *schema.ResourceData, exclude []string, index int, currentType string) error {
	if len(exclude) == 0 {
		return nil
	}

	for _, excludeField := range exclude {
		if _, ok := d.GetOk(fmt.Sprintf("customizer.%d.%s", index, excludeField)); ok {
			return fmt.Errorf("`%s` should not be set in the `customizer` block when the customizer type is `%s`", excludeField, currentType)
		}
	}

	return nil
}

func imageBuilderTemplateManagedImageNameAndResourceGroupName(input string) (string, string, error) {
	imageName := ""
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return "", "", err
	}

	if imageName, err = id.PopSegment("images"); err != nil {
		return "", "", err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return "", "", err
	}

	return imageName, id.ResourceGroup, nil
}
