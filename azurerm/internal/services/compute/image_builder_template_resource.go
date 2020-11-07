package compute

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"

	"github.com/Azure/azure-sdk-for-go/services/virtualmachineimagebuilder/mgmt/2020-02-01/virtualmachineimagebuilder"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	msiValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/validate"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// This is to serve input validation for the customizer block.
var fieldsOfFileCustomizer = []string{"file_source_uri", "file_sha256_checksum", "file_destination_path"}
var fieldsOfPowerShellCustomizer = []string{"powershell_commands", "powershell_run_elevated", "powershell_script_uri", "powershell_sha256_checksum", "powershell_valid_exit_codes"}
var fieldsOfShellCustomizer = []string{"shell_commands", "shell_script_uri", "shell_sha256_checksum"}
var fieldsOfWindowsRestartCustomizer = []string{"windows_restart_check_command", "windows_restart_command", "windows_restart_timeout"}
var fieldsOfWindowsUpdateCustomizer = []string{"windows_update_filters", "windows_update_search_criteria", "windows_update_limit"}

func resourceArmImageBuilderTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmImageBuilderTemplateCreate,
		Read:   resourceArmImageBuilderTemplateRead,
		Update: resourceArmImageBuilderTemplateUpdate,
		Delete: resourceArmImageBuilderTemplateDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ImageBuilderTemplateID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

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
								string(virtualmachineimagebuilder.UserAssigned),
							}, false),
						},
						"identity_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: msiValidate.UserAssignedIdentityId,
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
				ValidateFunc: validateDiskSizeGB,
			},

			// The array of this block is order insensitive. But there is a bug preventing using TypeSet for this block having a sub field `location` using StateFunc:
			// https://github.com/hashicorp/terraform-plugin-sdk/issues/160.
			// In detail, the symptom is specifying user friendly region say "West US 2" in the sub field `location` generated two blocks even if only specifying one block in .tf.
			// Using normalized region say "westus2" does not have the symptom.
			// Given this bug would break the core logic, using TypeList as a workaround.
			// This justification also applies to the "distribution_shared_image" block below.
			"distribution_managed_image": {
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

						"resource_group_name": azure.SchemaResourceGroupName(),

						"location": azure.SchemaLocation(),

						"run_output_name": distributionRunOutputNameSchema(),

						"tags": tags.ForceNewSchema(),
					},
				},
				AtLeastOneOf: []string{"distribution_managed_image", "distribution_shared_image", "distribution_vhd"},
			},

			"distribution_shared_image": {
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
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualmachineimagebuilder.StandardLRS),
								string(virtualmachineimagebuilder.StandardZRS),
							}, false),
						},

						"tags": tags.ForceNewSchema(),
					},
				},
				AtLeastOneOf: []string{"distribution_managed_image", "distribution_shared_image", "distribution_vhd"},
			},

			"distribution_vhd": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_output_name": distributionRunOutputNameSchema(),

						"tags": tags.ForceNewSchema(),
					},
				},
				AtLeastOneOf: []string{"distribution_managed_image", "distribution_shared_image", "distribution_vhd"},
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
				ValidateFunc: validate.ImageID,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmImageBuilderTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Compute.VMImageBuilderTemplateClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for the presence of existing Image Builder Template %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_image_builder_template", *existing.ID)
	}

	location := d.Get("location").(string)

	distribution, err := expandBasicImageTemplateDistributor(d.Get("distribution_managed_image").([]interface{}),
		d.Get("distribution_shared_image").([]interface{}),
		d.Get("distribution_vhd").([]interface{}),
		subscriptionId)
	if err != nil {
		return err
	}

	customizer, err := expandBasicImageTemplateCustomizer(d)
	if err != nil {
		return err
	}

	parameters := virtualmachineimagebuilder.ImageTemplate{
		Location: utils.String(location),
		Identity: expandImageTemplateIdentity(d.Get("identity").([]interface{})),
		ImageTemplateProperties: &virtualmachineimagebuilder.ImageTemplateProperties{
			VMProfile: &virtualmachineimagebuilder.ImageTemplateVMProfile{
				VMSize:       utils.String(d.Get("size").(string)),
				OsDiskSizeGB: utils.Int32(int32(d.Get("disk_size_gb").(int))),
			},

			Source:                expandBasicImageTemplateSource(d.Get("source_managed_image_id").(string), d.Get("source_platform_image").([]interface{}), d.Get("source_shared_image_version_id").(string)),
			Distribute:            distribution,
			Customize:             customizer,
			BuildTimeoutInMinutes: utils.Int32(int32(d.Get("build_timeout_minutes").(int))),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		parameters.ImageTemplateProperties.VMProfile.VnetConfig = &virtualmachineimagebuilder.VirtualNetworkConfig{
			SubnetID: utils.String(v.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, parameters, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("creating image builder template %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for image builder template %q (Resource Group %q) to finish provisioning: %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving image builder template %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for image builder template %q (Resource Group %q)", name, resourceGroup)
	}

	id, err := parse.ImageBuilderTemplateID(*resp.ID)
	if err != nil {
		return err
	}

	d.SetId(id.ID(subscriptionId))

	return resourceArmImageBuilderTemplateRead(d, meta)
}

func resourceArmImageBuilderTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMImageBuilderTemplateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ImageBuilderTemplateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving image builder template %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenImageBuilderTemplateIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if imageTemplateProperties := resp.ImageTemplateProperties; imageTemplateProperties != nil {
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

		distributionManagedImages, distributionSharedImages, distributionVhds, err := flattenBasicImageTemplateDistributor(imageTemplateProperties.Distribute)
		if err != nil {
			return err
		}

		if err := d.Set("distribution_managed_image", distributionManagedImages); err != nil {
			return fmt.Errorf("setting image template managed image distribution: %+v", err)
		}

		if err := d.Set("distribution_shared_image", distributionSharedImages); err != nil {
			return fmt.Errorf("setting image template shared image distribution: %+v", err)
		}

		if err := d.Set("distribution_vhd", distributionVhds); err != nil {
			return fmt.Errorf("setting image template vhd distribution: %+v", err)
		}

		if err := d.Set("customizer", flattenBasicImageTemplateCustomizer(imageTemplateProperties.Customize)); err != nil {
			return fmt.Errorf("setting `customizer`: %+v", err)
		}

		if err := d.Set("build_timeout_minutes", imageTemplateProperties.BuildTimeoutInMinutes); err != nil {
			return fmt.Errorf("setting `build timeout minutes`: %+v", err)
		}
	}

	if vmProfile := resp.VMProfile; vmProfile != nil {
		d.Set("size", vmProfile.VMSize)

		d.Set("disk_size_gb", vmProfile.OsDiskSizeGB)

		if vnetConfig := vmProfile.VnetConfig; vnetConfig != nil {
			d.Set("subnet_id", vnetConfig.SubnetID)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmImageBuilderTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMImageBuilderTemplateClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ImageBuilderTemplateID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving image builder template %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	parameters := virtualmachineimagebuilder.ImageTemplateUpdateParameters{}

	if d.HasChange("identity") {
		parameters.Identity = expandImageTemplateIdentity(d.Get("identity").([]interface{}))
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, parameters, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("updating image builder template %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of image builder template %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmImageBuilderTemplateRead(d, meta)
}

func resourceArmImageBuilderTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.VMImageBuilderTemplateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ImageBuilderTemplateID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting image builder template %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of image builder template %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandImageTemplateIdentity(input []interface{}) *virtualmachineimagebuilder.ImageTemplateIdentity {
	if len(input) == 0 {
		return nil
	}

	identity := input[0].(map[string]interface{})

	identityIds := make(map[string]*virtualmachineimagebuilder.ImageTemplateIdentityUserAssignedIdentitiesValue)
	for _, v := range identity["identity_ids"].(*schema.Set).List() {
		identityIds[v.(string)] = &virtualmachineimagebuilder.ImageTemplateIdentityUserAssignedIdentitiesValue{}
	}

	identityType := virtualmachineimagebuilder.ResourceIdentityType(identity["type"].(string))

	return &virtualmachineimagebuilder.ImageTemplateIdentity{
		Type:                   identityType,
		UserAssignedIdentities: identityIds,
	}
}

func flattenImageBuilderTemplateIdentity(input *virtualmachineimagebuilder.ImageTemplateIdentity) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	identityIds := make([]string, 0)
	if input.UserAssignedIdentities != nil {
		for k := range input.UserAssignedIdentities {
			identityIds = append(identityIds, k)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
		},
	}
}

func expandBasicImageTemplateSource(managedImageId string, platformImage []interface{}, sharedImageId string) virtualmachineimagebuilder.BasicImageTemplateSource {
	if len(managedImageId) != 0 {
		return &virtualmachineimagebuilder.ImageTemplateManagedImageSource{
			ImageID: utils.String(managedImageId),
		}
	}

	if len(platformImage) != 0 && platformImage[0] != nil {
		v := platformImage[0].(map[string]interface{})
		result := &virtualmachineimagebuilder.ImageTemplatePlatformImageSource{
			Publisher: utils.String(v["publisher"].(string)),
			Offer:     utils.String(v["offer"].(string)),
			Sku:       utils.String(v["sku"].(string)),
			Version:   utils.String(v["version"].(string)),
		}

		planRaw := v["plan"].([]interface{})
		if len(planRaw) > 0 {
			planTemp := planRaw[0].(map[string]interface{})
			result.PlanInfo = &virtualmachineimagebuilder.PlatformImagePurchasePlan{
				PlanName:      utils.String(planTemp["name"].(string)),
				PlanProduct:   utils.String(planTemp["product"].(string)),
				PlanPublisher: utils.String(planTemp["publisher"].(string)),
			}
		}

		return result
	}

	if len(sharedImageId) != 0 {
		return &virtualmachineimagebuilder.ImageTemplateSharedImageVersionSource{
			ImageVersionID: utils.String(sharedImageId),
		}
	}

	// since there is ExactlyOneOf on the source schema, below nil won't be reached
	return nil
}

func flattenBasicImageTemplateSource(input virtualmachineimagebuilder.BasicImageTemplateSource) (string, []interface{}, string) {
	if input != nil {
		switch source := input.(type) {
		case virtualmachineimagebuilder.ImageTemplateManagedImageSource:
			return flattenSourceManagedImage(&source), nil, ""
		case virtualmachineimagebuilder.ImageTemplatePlatformImageSource:
			return "", flattenSourcePlatformImage(&source), ""
		case virtualmachineimagebuilder.ImageTemplateSharedImageVersionSource:
			return "", nil, flattenSourceSharedImageVersion(&source)
		}
	}

	return "", nil, ""
}

func flattenSourceManagedImage(input *virtualmachineimagebuilder.ImageTemplateManagedImageSource) string {
	if input != nil && input.ImageID != nil {
		return *input.ImageID
	}

	return ""
}

func flattenSourcePlatformImage(input *virtualmachineimagebuilder.ImageTemplatePlatformImageSource) []interface{} {
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

func flattenImageBuilderTemplateSourcePlatformImagePlan(input *virtualmachineimagebuilder.PlatformImagePurchasePlan) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.PlanName != nil {
		name = *input.PlanName
	}

	product := ""
	if input.PlanProduct != nil {
		product = *input.PlanProduct
	}

	publisher := ""
	if input.PlanPublisher != nil {
		publisher = *input.PlanPublisher
	}

	return []interface{}{
		map[string]interface{}{
			"name":      name,
			"product":   product,
			"publisher": publisher,
		},
	}
}

func flattenSourceSharedImageVersion(input *virtualmachineimagebuilder.ImageTemplateSharedImageVersionSource) string {
	if input != nil && input.ImageVersionID != nil {
		return *input.ImageVersionID
	}

	return ""
}

func expandBasicImageTemplateDistributor(distributionManagedImages []interface{}, distributionSharedImages []interface{}, distributionVhds []interface{}, subscriptionId string) (*[]virtualmachineimagebuilder.BasicImageTemplateDistributor, error) {
	results := make([]virtualmachineimagebuilder.BasicImageTemplateDistributor, 0)
	runOutputNameSet := make(map[string]bool)

	if len(distributionManagedImages) > 0 {
		for _, v := range distributionManagedImages {
			if v != nil {
				distributionManagedImage := v.(map[string]interface{})
				resourceGroupName := distributionManagedImage["resource_group_name"].(string)
				runOutputName := distributionManagedImage["run_output_name"].(string)

				_, existing := runOutputNameSet[runOutputName]
				if existing {
					return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
				} else {
					runOutputNameSet[runOutputName] = true

					results = append(results, virtualmachineimagebuilder.ImageTemplateManagedImageDistributor{
						ImageID:       utils.String("/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroupName + "/providers/Microsoft.Compute/images/" + distributionManagedImage["name"].(string)),
						Location:      utils.String(distributionManagedImage["location"].(string)),
						RunOutputName: utils.String(runOutputName),
						ArtifactTags:  tags.Expand(distributionManagedImage["tags"].(map[string]interface{})),
					})
				}
			}
		}
	}

	if len(distributionVhds) > 0 {
		for _, v := range distributionVhds {
			if v != nil {
				distributionVhd := v.(map[string]interface{})
				runOutputName := distributionVhd["run_output_name"].(string)

				_, existing := runOutputNameSet[runOutputName]
				if existing {
					return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
				} else {
					runOutputNameSet[runOutputName] = true
					results = append(results, virtualmachineimagebuilder.ImageTemplateVhdDistributor{
						RunOutputName: utils.String(distributionVhd["run_output_name"].(string)),
						ArtifactTags:  tags.Expand(distributionVhd["tags"].(map[string]interface{})),
					})
				}
			}
		}
	}

	if len(distributionSharedImages) > 0 {
		for _, v := range distributionSharedImages {
			if v != nil {
				distributionSharedImage := v.(map[string]interface{})
				runOutputName := distributionSharedImage["run_output_name"].(string)

				_, existing := runOutputNameSet[runOutputName]
				if existing {
					return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
				} else {
					runOutputNameSet[runOutputName] = true
					results = append(results, virtualmachineimagebuilder.ImageTemplateSharedImageDistributor{
						GalleryImageID:     utils.String(distributionSharedImage["id"].(string)),
						ReplicationRegions: expandImageTemplateSharedImageDistributorReplicaRegions(distributionSharedImage["replica_regions"].([]interface{})),
						RunOutputName:      utils.String(distributionSharedImage["run_output_name"].(string)),
						ExcludeFromLatest:  utils.Bool(distributionSharedImage["exclude_from_latest"].(bool)),
						StorageAccountType: virtualmachineimagebuilder.SharedImageStorageAccountType(distributionSharedImage["storage_account_type"].(string)),
						ArtifactTags:       tags.Expand(distributionSharedImage["tags"].(map[string]interface{})),
					})
				}
			}
		}
	}

	return &results, nil
}

func flattenBasicImageTemplateDistributor(input *[]virtualmachineimagebuilder.BasicImageTemplateDistributor) ([]interface{}, []interface{}, []interface{}, error) {
	distributionManagedImages := make([]interface{}, 0)
	distributionVhds := make([]interface{}, 0)
	distributionSharedImages := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			switch distribute := v.(type) {
			case virtualmachineimagebuilder.ImageTemplateManagedImageDistributor:
				flattenedManagedImg, err := flattenDistributionManagedImage(&distribute)
				if err != nil {
					return nil, nil, nil, fmt.Errorf("setting image template managed image source: %+v", err)
				}

				distributionManagedImages = append(distributionManagedImages, flattenedManagedImg)
			case virtualmachineimagebuilder.ImageTemplateSharedImageDistributor:
				distributionSharedImages = append(distributionSharedImages, flattenDistributionSharedImage(&distribute))
			case virtualmachineimagebuilder.ImageTemplateVhdDistributor:
				distributionVhds = append(distributionVhds, flattenDistributionVhd(&distribute))
			}
		}
	}

	return distributionManagedImages, distributionSharedImages, distributionVhds, nil
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

func flattenDistributionManagedImage(input *virtualmachineimagebuilder.ImageTemplateManagedImageDistributor) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	if input == nil {
		return result, nil
	}

	imageName := ""
	resourceGroupName := ""
	if input.ImageID != nil {
		imageNameReturned, resourceGroupNameReturned, err := imageBuilderTemplateManagedImageNameAndResourceGroupName(*input.ImageID)
		if err != nil {
			return result, err
		}

		imageName = imageNameReturned
		resourceGroupName = resourceGroupNameReturned
	}

	result["name"] = imageName
	result["resource_group_name"] = resourceGroupName

	if input.Location != nil {
		result["location"] = *input.Location
	}

	if input.RunOutputName != nil {
		result["run_output_name"] = *input.RunOutputName
	}

	if input.ArtifactTags != nil {
		result["tags"] = tags.Flatten(input.ArtifactTags)
	}

	return result, nil
}

func flattenDistributionSharedImage(input *virtualmachineimagebuilder.ImageTemplateSharedImageDistributor) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	if input.GalleryImageID != nil {
		results["id"] = *input.GalleryImageID
	}

	if input.ReplicationRegions != nil {
		results["replica_regions"] = flattenImageTemplateSharedImageDistributorReplicaRegions(input.ReplicationRegions)
	}

	if input.RunOutputName != nil {
		results["run_output_name"] = *input.RunOutputName
	}

	if input.ExcludeFromLatest != nil {
		results["exclude_from_latest"] = *input.ExcludeFromLatest
	}

	results["storage_account_type"] = string(input.StorageAccountType)

	if input.ArtifactTags != nil {
		results["tags"] = tags.Flatten(input.ArtifactTags)
	}

	return results
}

func flattenDistributionVhd(input *virtualmachineimagebuilder.ImageTemplateVhdDistributor) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}

	if input.RunOutputName != nil {
		results["run_output_name"] = *input.RunOutputName
	}

	if input.ArtifactTags != nil {
		results["tags"] = tags.Flatten(input.ArtifactTags)
	}

	return results
}

// Passing d as input rather its business subset because this function needs d.GetOK() to verify users' input to the `customizer` block is valid.
func expandBasicImageTemplateCustomizer(d *schema.ResourceData) (*[]virtualmachineimagebuilder.BasicImageTemplateCustomizer, error) {
	input := d.Get("customizer").([]interface{})

	if len(input) == 0 {
		return nil, nil
	}

	results := make([]virtualmachineimagebuilder.BasicImageTemplateCustomizer, 0)

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

			results = append(results, virtualmachineimagebuilder.ImageTemplateFileCustomizer{
				Name:           utils.String(customizer["name"].(string)),
				SourceURI:      utils.String(customizer["file_source_uri"].(string)),
				Sha256Checksum: utils.String(customizer["file_sha256_checksum"].(string)),
				Destination:    utils.String(customizer["file_destination_path"].(string)),
			})

			i++
		case "PowerShell":
			if err := validateImageTemplateCustomizerInputForPowerShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagebuilder.ImageTemplatePowerShellCustomizer{
				Name:           utils.String(customizer["name"].(string)),
				ScriptURI:      utils.String(customizer["powershell_script_uri"].(string)),
				Sha256Checksum: utils.String(customizer["powershell_sha256_checksum"].(string)),
				Inline:         utils.ExpandStringSlice(customizer["powershell_commands"].([]interface{})),
				RunElevated:    utils.Bool(customizer["powershell_run_elevated"].(bool)),
				ValidExitCodes: utils.ExpandInt32Slice(customizer["powershell_valid_exit_codes"].([]interface{})),
			})

			i++
		case "Shell":
			if err := validateImageTemplateCustomizerInputForShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagebuilder.ImageTemplateShellCustomizer{
				Name:           utils.String(customizer["name"].(string)),
				ScriptURI:      utils.String(customizer["shell_script_uri"].(string)),
				Sha256Checksum: utils.String(customizer["shell_sha256_checksum"].(string)),
				Inline:         utils.ExpandStringSlice(customizer["shell_commands"].([]interface{})),
			})

			i++
		case "WindowsRestart":
			if err := validateImageTemplateCustomizerInputForWindowsRestartType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagebuilder.ImageTemplateRestartCustomizer{
				Name:                utils.String(customizer["name"].(string)),
				RestartCommand:      utils.String(customizer["windows_restart_command"].(string)),
				RestartCheckCommand: utils.String(customizer["windows_restart_check_command"].(string)),
				RestartTimeout:      utils.String(customizer["windows_restart_timeout"].(string)),
			})

			i++
		case "WindowsUpdate":
			if err := validateImageTemplateCustomizerInputForWindowsUpdateType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagebuilder.ImageTemplateWindowsUpdateCustomizer{
				Name:           utils.String(customizer["name"].(string)),
				SearchCriteria: utils.String(customizer["windows_update_search_criteria"].(string)),
				Filters:        utils.ExpandStringSlice(customizer["windows_update_filters"].([]interface{})),
				UpdateLimit:    utils.Int32(int32(customizer["windows_update_limit"].(int))),
			})

			i++
		}
	}

	return &results, nil
}

func flattenBasicImageTemplateCustomizer(input *[]virtualmachineimagebuilder.BasicImageTemplateCustomizer) []interface{} {
	customizerList := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
			switch customizer := v.(type) {
			case virtualmachineimagebuilder.ImageTemplateFileCustomizer:
				customizerList = append(customizerList, flattenCustomizerFile(&customizer))
			case virtualmachineimagebuilder.ImageTemplatePowerShellCustomizer:
				customizerList = append(customizerList, flattenCustomizerPowerShell(&customizer))
			case virtualmachineimagebuilder.ImageTemplateShellCustomizer:
				customizerList = append(customizerList, flattenCustomizerShell(&customizer))
			case virtualmachineimagebuilder.ImageTemplateRestartCustomizer:
				customizerList = append(customizerList, flattenCustomizerWindowsRestart(&customizer))
			case virtualmachineimagebuilder.ImageTemplateWindowsUpdateCustomizer:
				customizerList = append(customizerList, flattenCustomizerWindowsUpdate(&customizer))
			}
		}

		return customizerList
	}

	return customizerList
}

func flattenCustomizerFile(input *virtualmachineimagebuilder.ImageTemplateFileCustomizer) map[string]interface{} {
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

	if input.SourceURI != nil {
		result["file_source_uri"] = *input.SourceURI
	}

	if input.Sha256Checksum != nil {
		result["file_sha256_checksum"] = *input.Sha256Checksum
	}

	if input.Destination != nil {
		result["file_destination_path"] = *input.Destination
	}

	return result
}

func flattenCustomizerShell(input *virtualmachineimagebuilder.ImageTemplateShellCustomizer) map[string]interface{} {
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

	if input.ScriptURI != nil {
		result["shell_script_uri"] = *input.ScriptURI
	}

	if input.Sha256Checksum != nil {
		result["shell_sha256_checksum"] = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result["shell_commands"] = *input.Inline
	}

	return result
}

func flattenCustomizerPowerShell(input *virtualmachineimagebuilder.ImageTemplatePowerShellCustomizer) map[string]interface{} {
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

	if input.ScriptURI != nil {
		result["powershell_script_uri"] = *input.ScriptURI
	}

	if input.Sha256Checksum != nil {
		result["powershell_sha256_checksum"] = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result["powershell_commands"] = *input.Inline
	}

	if input.RunElevated != nil {
		result["powershell_run_elevated"] = *input.RunElevated
	}

	if input.ValidExitCodes != nil {
		result["powershell_valid_exit_codes"] = *input.ValidExitCodes
	}

	return result
}

func flattenCustomizerWindowsRestart(input *virtualmachineimagebuilder.ImageTemplateRestartCustomizer) map[string]interface{} {
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

func flattenCustomizerWindowsUpdate(input *virtualmachineimagebuilder.ImageTemplateWindowsUpdateCustomizer) map[string]interface{} {
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
