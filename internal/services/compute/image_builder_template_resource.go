package compute

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/imagebuilder/2024-02-01/virtualmachineimagetemplate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = ImageBuilderTemplateResource{}
	_ sdk.ResourceWithUpdate = ImageBuilderTemplateResource{}
)

type ImageBuilderTemplateResource struct{}

func (ImageBuilderTemplateResource) ModelObject() interface{} {
	return &ImageBuilderTemplateResourceModel{}
}

func (ImageBuilderTemplateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualmachineimagetemplate.ValidateImageTemplateID
}

func (ImageBuilderTemplateResource) ResourceType() string {
	return "azurerm_image_builder_template"
}

type ImageBuilderTemplateResourceModel struct {
	Name                       string                                    `tfschema:"name"`
	ResourceGroupName          string                                    `tfschema:"resource_group_name"`
	Location                   string                                    `tfschema:"location"`
	Identity                   []identity.ModelUserAssigned              `tfschema:"identity"`
	BuildTimeoutMinutes        int64                                     `tfschema:"build_timeout_minutes"`
	Customizer                 []ImageBuilderTemplateCustomizer          `tfschema:"customizer"`
	DiskSizeGb                 int64                                     `tfschema:"disk_size_gb"`
	Distributions              []ImageBuilderTemplateDistributions       `tfschema:"distributions"`
	Size                       string                                    `tfschema:"size"`
	SourceManagedImageId       string                                    `tfschema:"source_managed_image_id"`
	SourcePlatformImage        []ImageBuilderTemplateSourcePlatformImage `tfschema:"source_platform_image"`
	SourceSharedImageVersionId string                                    `tfschema:"source_shared_image_version_id"`
	SubnetId                   string                                    `tfschema:"subnet_id"`
	Tags                       map[string]string                         `tfschema:"tags"`
}

type ImageBuilderTemplateCustomizer struct {
	Type                        string   `tfschema:"type"`
	FileDestinationPath         string   `tfschema:"file_destination_path"`
	FileSha256Checksum          string   `tfschema:"file_sha256_checksum"`
	FileSourceUri               string   `tfschema:"file_source_uri"`
	Name                        string   `tfschema:"name"`
	PowershellCommands          []string `tfschema:"powershell_commands"`
	PowershellRunAsSystem       bool     `tfschema:"powershell_run_as_system"`
	PowershellRunElevated       bool     `tfschema:"powershell_run_elevated"`
	PowershellScriptUri         string   `tfschema:"powershell_script_uri"`
	PowershellSha256Checksum    string   `tfschema:"powershell_sha256_checksum"`
	PowershellValidExitCodes    []int64  `tfschema:"powershell_valid_exit_codes"`
	ShellCommands               []string `tfschema:"shell_commands"`
	ShellScriptUri              string   `tfschema:"shell_script_uri"`
	ShellSha256Checksum         string   `tfschema:"shell_sha256_checksum"`
	WindowsRestartCheckCommand  string   `tfschema:"windows_restart_check_command"`
	WindowsRestartCommand       string   `tfschema:"windows_restart_command"`
	WindowsRestartTimeout       string   `tfschema:"windows_restart_timeout"`
	WindowsUpdateFilters        []string `tfschema:"windows_update_filters"`
	WindowsUpdateSearchCriteria string   `tfschema:"windows_update_search_criteria"`
	WindowsUpdateLimit          int64    `tfschema:"windows_update_limit"`
}

type ImageBuilderTemplateDistributions struct {
	ManagedImage []ImageBuilderTemplateDistributionsManagedImage `tfschema:"managed_image"`
	SharedImage  []ImageBuilderTemplateDistributionsSharedImage  `tfschema:"shared_image"`
	Vhd          []ImageBuilderTemplateDistributionsVhd          `tfschema:"vhd"`
}

type ImageBuilderTemplateDistributionsManagedImage struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	RunOutputName     string            `tfschema:"run_output_name"`
	Tags              map[string]string `tfschema:"tags"`
}
type ImageBuilderTemplateDistributionsSharedImage struct {
	Id                 string                                                       `tfschema:"id"`
	ReplicaRegions     []ImageBuilderTemplateDistributionsSharedImageReplicaRegions `tfschema:"replica_regions"`
	RunOutputName      string                                                       `tfschema:"run_output_name"`
	ExcludeFromLatest  bool                                                         `tfschema:"exclude_from_latest"`
	StorageAccountType string                                                       `tfschema:"storage_account_type"`
	Tags               map[string]string                                            `tfschema:"tags"`
	Versioning         []ImageBuilderTemplateDistributionsSharedImageVersioning     `tfschema:"versioning"`
}

type ImageBuilderTemplateDistributionsVhd struct {
	RunOutputName string            `tfschema:"run_output_name"`
	Tags          map[string]string `tfschema:"tags"`
}

type ImageBuilderTemplateDistributionsSharedImageReplicaRegions struct {
	Name string `tfschema:"name"`
}
type ImageBuilderTemplateDistributionsSharedImageVersioning struct {
	Scheme string `tfschema:"scheme"`
	Major  int64  `tfschema:"major"`
}

type ImageBuilderTemplateSourcePlatformImage struct {
	Publisher string                                        `tfschema:"publisher"`
	Offer     string                                        `tfschema:"offer"`
	Sku       string                                        `tfschema:"sku"`
	Version   string                                        `tfschema:"version"`
	Plan      []ImageBuilderTemplateSourcePlatformImagePlan `tfschema:"plan"`
}

type ImageBuilderTemplateSourcePlatformImagePlan struct {
	Name      string `tfschema:"name"`
	Product   string `tfschema:"product"`
	Publisher string `tfschema:"publisher"`
}

// This is to serve input validation for the customizer block.
var (
	fieldsOfFileCustomizer           = []string{"file_source_uri", "file_sha256_checksum", "file_destination_path"}
	fieldsOfPowerShellCustomizer     = []string{"powershell_commands", "powershell_run_as_system", "powershell_run_elevated", "powershell_script_uri", "powershell_sha256_checksum", "powershell_valid_exit_codes"}
	fieldsOfShellCustomizer          = []string{"shell_commands", "shell_script_uri", "shell_sha256_checksum"}
	fieldsOfWindowsRestartCustomizer = []string{"windows_restart_check_command", "windows_restart_command", "windows_restart_timeout"}
	fieldsOfWindowsUpdateCustomizer  = []string{"windows_update_filters", "windows_update_search_criteria", "windows_update_limit"}
)

func (ImageBuilderTemplateResource) Arguments() map[string]*pluginsdk.Schema {
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
	}
}

func (ImageBuilderTemplateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ImageBuilderTemplateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineImageTemplateClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ImageBuilderTemplateResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := virtualmachineimagetemplate.NewImageTemplateID(subscriptionId, model.ResourceGroupName, model.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for the presence of existing Image Builder Template %q (Resource Group %q): %s", model.Name, model.ResourceGroupName, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			distribution, err := expandBasicImageTemplateDistributor(model.Distributions, subscriptionId)
			if err != nil {
				return err
			}

			customizer, err := expandBasicImageTemplateCustomizer(metadata.ResourceData, model.Customizer)
			if err != nil {
				return err
			}

			identityExpanded, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			parameters := virtualmachineimagetemplate.ImageTemplate{
				Location: model.Location,
				Identity: *identityExpanded,
				Properties: &virtualmachineimagetemplate.ImageTemplateProperties{
					VMProfile: &virtualmachineimagetemplate.ImageTemplateVMProfile{
						VMSize:       &model.Size,
						OsDiskSizeGB: &model.DiskSizeGb,
					},

					Source:                expandBasicImageTemplateSource(model.SourceManagedImageId, model.SourcePlatformImage, model.SourceSharedImageVersionId),
					Distribute:            *distribution,
					Customize:             customizer,
					BuildTimeoutInMinutes: &model.BuildTimeoutMinutes,
				},

				Tags: &model.Tags,
			}

			if model.SubnetId != "" {
				parameters.Properties.VMProfile.VnetConfig = &virtualmachineimagetemplate.VirtualNetworkConfig{
					SubnetId: &model.SubnetId,
				}
			}

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating image builder template %q (Resource Group %q): %+v", model.Name, model.ResourceGroupName, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ImageBuilderTemplateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineImageTemplateClient

			state := ImageBuilderTemplateResourceModel{}

			id, err := virtualmachineimagetemplate.ParseImageTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
			}

			state.Name = id.ImageTemplateName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)

				identityFlattened, err := identity.FlattenUserAssignedMapToModel(&model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = *identityFlattened

				if imageTemplateProperties := resp.Model.Properties; imageTemplateProperties != nil {
					managedImageId, platformImage, sharedImageVersionId := flattenBasicImageTemplateSource(imageTemplateProperties.Source)

					// only one among managedImageId / platformImage / sharedImageVersionId would be returned.
					if managedImageId != "" {
						state.SourceManagedImageId = managedImageId
					}
					if len(platformImage) > 0 {
						state.SourcePlatformImage = platformImage
					}

					if sharedImageVersionId != "" {
						state.SourceSharedImageVersionId = sharedImageVersionId
					}

					flattenedDistributions, err := flattenBasicImageTemplateDistributor(&imageTemplateProperties.Distribute)
					if err != nil {
						return err
					}

					state.Distributions = flattenedDistributions
					state.Customizer = flattenBasicImageTemplateCustomizer(imageTemplateProperties.Customize)
					state.BuildTimeoutMinutes = *imageTemplateProperties.BuildTimeoutInMinutes

					if vmProfile := resp.Model.Properties.VMProfile; vmProfile != nil {
						state.Size = *vmProfile.VMSize
						state.DiskSizeGb = *vmProfile.OsDiskSizeGB

						if vnetConfig := vmProfile.VnetConfig; vnetConfig != nil {
							state.SubnetId = *vnetConfig.SubnetId
						}
					}
				}
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ImageBuilderTemplateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineImageTemplateClient

			var model ImageBuilderTemplateResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			id, err := virtualmachineimagetemplate.ParseImageTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			parameters := virtualmachineimagetemplate.ImageTemplateUpdateParameters{}

			if metadata.ResourceData.HasChange("identity") {
				identityExpanded, err := identity.ExpandUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}

				parameters.Identity = identityExpanded
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = &model.Tags
			}

			err = client.UpdateThenPoll(ctx, *id, parameters)
			if err != nil {
				return fmt.Errorf("updating image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}

func (r ImageBuilderTemplateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineImageTemplateClient

			id, err := virtualmachineimagetemplate.ParseImageTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting image builder template %q (Resource Group %q): %+v", id.ImageTemplateName, id.ResourceGroupName, err)
			}

			return nil
		},
	}
}

func expandBasicImageTemplateSource(managedImageId string, platformImage []ImageBuilderTemplateSourcePlatformImage, sharedImageId string) virtualmachineimagetemplate.ImageTemplateSource {
	if len(managedImageId) != 0 {
		return &virtualmachineimagetemplate.ImageTemplateManagedImageSource{
			ImageId: managedImageId,
		}
	}

	if len(platformImage) != 0 {
		v := platformImage[0]
		result := &virtualmachineimagetemplate.ImageTemplatePlatformImageSource{
			Publisher: &v.Publisher,
			Offer:     &v.Offer,
			Sku:       &v.Sku,
			Version:   &v.Version,
		}

		if len(v.Plan) > 0 {
			planTemp := v.Plan[0]
			result.PlanInfo = &virtualmachineimagetemplate.PlatformImagePurchasePlan{
				PlanName:      planTemp.Name,
				PlanProduct:   planTemp.Product,
				PlanPublisher: planTemp.Publisher,
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

func flattenBasicImageTemplateSource(input virtualmachineimagetemplate.ImageTemplateSource) (string, []ImageBuilderTemplateSourcePlatformImage, string) {
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

func flattenSourcePlatformImage(input *virtualmachineimagetemplate.ImageTemplatePlatformImageSource) []ImageBuilderTemplateSourcePlatformImage {
	if input == nil {
		return nil
	}

	result := make([]ImageBuilderTemplateSourcePlatformImage, 0)

	result = append(result, ImageBuilderTemplateSourcePlatformImage{
		Plan: flattenImageBuilderTemplateSourcePlatformImagePlan(input.PlanInfo),
	})

	if input.Publisher != nil {
		result[0].Publisher = *input.Publisher
	}

	if input.Offer != nil {
		result[0].Offer = *input.Offer
	}

	if input.Sku != nil {
		result[0].Sku = *input.Sku
	}

	if input.Version != nil {
		result[0].Version = *input.Version
	}
	return result
}

func flattenImageBuilderTemplateSourcePlatformImagePlan(input *virtualmachineimagetemplate.PlatformImagePurchasePlan) []ImageBuilderTemplateSourcePlatformImagePlan {
	if input == nil {
		return nil
	}

	result := make([]ImageBuilderTemplateSourcePlatformImagePlan, 0)

	result = append(result, ImageBuilderTemplateSourcePlatformImagePlan{
		Name:      input.PlanName,
		Product:   input.PlanProduct,
		Publisher: input.PlanPublisher,
	})
	return result
}

func flattenSourceSharedImageVersion(input *virtualmachineimagetemplate.ImageTemplateSharedImageVersionSource) string {
	if input != nil && input.ImageVersionId != "" {
		return input.ImageVersionId
	}

	return ""
}

func expandBasicImageTemplateDistributor(distributions []ImageBuilderTemplateDistributions, subscriptionId string) (*[]virtualmachineimagetemplate.ImageTemplateDistributor, error) {
	results := make([]virtualmachineimagetemplate.ImageTemplateDistributor, 0)
	runOutputNameSet := make(map[string]bool)

	if len(distributions) > 0 {
		for _, v := range distributions {
			managedImages := v.ManagedImage
			sharedImages := v.SharedImage
			vhds := v.Vhd

			if len(managedImages) == 0 && len(sharedImages) == 0 && len(vhds) == 0 {
				return &results, fmt.Errorf("at least one of `managed_image`, `shared_image` and `vhd` is required to specify in `distributions`")
			}

			if len(managedImages) > 0 {
				for _, managedImage := range managedImages {
					resourceGroupName := managedImage.ResourceGroupName
					runOutputName := managedImage.RunOutputName

					_, existing := runOutputNameSet[runOutputName]
					if existing {
						return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
					} else {
						runOutputNameSet[runOutputName] = true

						results = append(results, virtualmachineimagetemplate.ImageTemplateManagedImageDistributor{
							ImageId:       "/subscriptions/" + subscriptionId + "/resourceGroups/" + resourceGroupName + "/providers/Microsoft.Compute/images/" + managedImage.Name,
							Location:      managedImage.Location,
							RunOutputName: runOutputName,
							ArtifactTags:  &managedImage.Tags,
						})
					}
				}
			}

			if len(sharedImages) > 0 {
				for _, sharedImage := range sharedImages {
					// sharedImage := sharedImageRaw.(map[string]interface{})
					runOutputName := sharedImage.RunOutputName

					_, existing := runOutputNameSet[runOutputName]
					if existing {
						return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
					} else {
						runOutputNameSet[runOutputName] = true
						results = append(results, virtualmachineimagetemplate.ImageTemplateSharedImageDistributor{
							GalleryImageId:     sharedImage.Id,
							ReplicationRegions: expandImageTemplateSharedImageDistributorReplicaRegions(sharedImage.ReplicaRegions),
							RunOutputName:      sharedImage.RunOutputName,
							ExcludeFromLatest:  &sharedImage.ExcludeFromLatest,
							StorageAccountType: pointer.To(virtualmachineimagetemplate.SharedImageStorageAccountType(sharedImage.StorageAccountType)),
							ArtifactTags:       &sharedImage.Tags,
							Versioning:         expandImageTemplateSharedImageDistributorVersioning(sharedImage.Versioning),
						})
					}
				}
			}

			if len(vhds) > 0 {
				for _, vhd := range vhds {
					// vhd := vhdRaw.(map[string]interface{})
					runOutputName := vhd.RunOutputName

					_, existing := runOutputNameSet[runOutputName]
					if existing {
						return &results, fmt.Errorf("`run_output_name` must be unique among all distribution destinations. %q already exists", runOutputName)
					} else {
						runOutputNameSet[runOutputName] = true
						results = append(results, virtualmachineimagetemplate.ImageTemplateVhdDistributor{
							RunOutputName: vhd.RunOutputName,
							ArtifactTags:  &vhd.Tags,
						})
					}
				}
			}
		}
	}

	return &results, nil
}

func flattenBasicImageTemplateDistributor(input *[]virtualmachineimagetemplate.ImageTemplateDistributor) ([]ImageBuilderTemplateDistributions, error) {
	results := make([]ImageBuilderTemplateDistributions, 0)

	distributionManagedImages := make([]ImageBuilderTemplateDistributionsManagedImage, 0)
	distributionSharedImages := make([]ImageBuilderTemplateDistributionsSharedImage, 0)
	distributionVhds := make([]ImageBuilderTemplateDistributionsVhd, 0)

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

	output := ImageBuilderTemplateDistributions{
		ManagedImage: distributionManagedImages,
		SharedImage:  distributionSharedImages,
		Vhd:          distributionVhds,
	}

	results = append(results, output)
	return results, nil
}

func expandImageTemplateSharedImageDistributorReplicaRegions(input []ImageBuilderTemplateDistributionsSharedImageReplicaRegions) *[]string {
	if input == nil {
		return nil
	}

	result := make([]string, 0)
	for _, v := range input {
		result = append(result, v.Name)
	}

	return &result
}

func expandImageTemplateSharedImageDistributorVersioning(input []ImageBuilderTemplateDistributionsSharedImageVersioning) virtualmachineimagetemplate.DistributeVersioner {
	var result virtualmachineimagetemplate.DistributeVersioner
	if len(input) > 0 {
		config := input[0]

		if v := config.Scheme; v == "Latest" {
			result = virtualmachineimagetemplate.DistributeVersionerLatest{
				Scheme: config.Scheme,
				Major:  &config.Major,
			}
		}

		if v := config.Scheme; v == "Source" {
			result = virtualmachineimagetemplate.DistributeVersionerSource{
				Scheme: config.Scheme,
			}
		}
	}
	return result
}

func flattenImageTemplateSharedImageDistributorVersioning(input virtualmachineimagetemplate.DistributeVersioner) []ImageBuilderTemplateDistributionsSharedImageVersioning {
	results := make([]ImageBuilderTemplateDistributionsSharedImageVersioning, 0)
	output := ImageBuilderTemplateDistributionsSharedImageVersioning{}

	switch versioner := input.(type) {
	case virtualmachineimagetemplate.DistributeVersionerLatest:
		output.Scheme = versioner.Scheme
		output.Major = *versioner.Major

	case virtualmachineimagetemplate.DistributeVersionerSource:
		output.Scheme = versioner.Scheme
	}

	results = append(results, output)

	return results
}

func flattenImageTemplateSharedImageDistributorReplicaRegions(input *[]string) []ImageBuilderTemplateDistributionsSharedImageReplicaRegions {
	results := make([]ImageBuilderTemplateDistributionsSharedImageReplicaRegions, 0)

	if input != nil {
		for _, v := range *input {
			results = append(results, ImageBuilderTemplateDistributionsSharedImageReplicaRegions{
				Name: v,
			})
		}
	}

	return results
}

func flattenDistributionManagedImage(input *virtualmachineimagetemplate.ImageTemplateManagedImageDistributor) (ImageBuilderTemplateDistributionsManagedImage, error) {
	result := ImageBuilderTemplateDistributionsManagedImage{}
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

	result.Name = imageName
	result.ResourceGroupName = resourceGroupName
	result.Location = input.Location
	result.RunOutputName = input.RunOutputName
	result.Tags = *input.ArtifactTags

	return result, nil
}

func flattenDistributionSharedImage(input *virtualmachineimagetemplate.ImageTemplateSharedImageDistributor) ImageBuilderTemplateDistributionsSharedImage {
	results := ImageBuilderTemplateDistributionsSharedImage{}
	if input == nil {
		return results
	}

	results.Id = input.GalleryImageId

	if input.ReplicationRegions != nil {
		results.ReplicaRegions = flattenImageTemplateSharedImageDistributorReplicaRegions(input.ReplicationRegions)
	}

	results.RunOutputName = input.RunOutputName

	if input.ExcludeFromLatest != nil {
		results.ExcludeFromLatest = *input.ExcludeFromLatest
	}

	if input.StorageAccountType != nil {
		results.StorageAccountType = string(*input.StorageAccountType)
	}

	if input.ArtifactTags != nil {
		results.Tags = *input.ArtifactTags
	}

	results.Versioning = flattenImageTemplateSharedImageDistributorVersioning(input.Versioning)
	return results
}

func flattenDistributionVhd(input *virtualmachineimagetemplate.ImageTemplateVhdDistributor) ImageBuilderTemplateDistributionsVhd {
	results := ImageBuilderTemplateDistributionsVhd{}
	if input == nil {
		return results
	}

	if input.RunOutputName != "" {
		results.RunOutputName = input.RunOutputName
	}

	if input.ArtifactTags != nil {
		results.Tags = *input.ArtifactTags
	}

	return results
}

// Passing d as input rather its business subset because this function needs d.GetOK() to verify users' input to the `customizer` block is valid.
func expandBasicImageTemplateCustomizer(d *schema.ResourceData, input []ImageBuilderTemplateCustomizer) (*[]virtualmachineimagetemplate.ImageTemplateCustomizer, error) {
	if len(input) == 0 {
		return nil, nil
	}

	results := make([]virtualmachineimagetemplate.ImageTemplateCustomizer, 0)

	// This index serves identifying the invalid input customizer block from the customizer block list.
	i := 0

	for _, customizer := range input {
		t := customizer.Type

		switch t {
		case "File":
			if err := validateImageTemplateCustomizerInputForFileType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateFileCustomizer{
				Name:           pointer.To(customizer.Name),
				SourceUri:      pointer.To(customizer.FileSourceUri),
				Sha256Checksum: pointer.To(customizer.FileSha256Checksum),
				Destination:    pointer.To(customizer.FileDestinationPath),
			})

			i++
		case "PowerShell":
			if err := validateImageTemplateCustomizerInputForPowerShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer{
				Name:           pointer.To(customizer.Name),
				ScriptUri:      pointer.To(customizer.PowershellScriptUri),
				Sha256Checksum: pointer.To(customizer.PowershellSha256Checksum),
				Inline:         pointer.To(customizer.PowershellCommands),
				RunAsSystem:    pointer.To(customizer.PowershellRunAsSystem),
				RunElevated:    pointer.To(customizer.PowershellRunElevated),
				ValidExitCodes: pointer.To(customizer.PowershellValidExitCodes),
			})

			i++
		case "Shell":
			if err := validateImageTemplateCustomizerInputForShellType(d, customizer, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateShellCustomizer{
				Name:           pointer.To(customizer.Name),
				ScriptUri:      pointer.To(customizer.ShellScriptUri),
				Sha256Checksum: pointer.To(customizer.ShellSha256Checksum),
				Inline:         pointer.To(customizer.ShellCommands),
			})

			i++
		case "WindowsRestart":
			if err := validateImageTemplateCustomizerInputForWindowsRestartType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateRestartCustomizer{
				Name:                pointer.To(customizer.Name),
				RestartCommand:      pointer.To(customizer.WindowsRestartCommand),
				RestartCheckCommand: pointer.To(customizer.WindowsRestartCheckCommand),
				RestartTimeout:      pointer.To(customizer.WindowsRestartTimeout),
			})

			i++
		case "WindowsUpdate":
			if err := validateImageTemplateCustomizerInputForWindowsUpdateType(d, i); err != nil {
				return nil, err
			}

			results = append(results, virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer{
				Name:           pointer.To(customizer.Name),
				SearchCriteria: pointer.To(customizer.WindowsUpdateSearchCriteria),
				Filters:        pointer.To(customizer.WindowsUpdateFilters),
				UpdateLimit:    pointer.To(customizer.WindowsUpdateLimit),
			})

			i++
		}
	}

	return &results, nil
}

func flattenBasicImageTemplateCustomizer(input *[]virtualmachineimagetemplate.ImageTemplateCustomizer) []ImageBuilderTemplateCustomizer {
	customizerList := make([]ImageBuilderTemplateCustomizer, 0)

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

func flattenCustomizerFile(input *virtualmachineimagetemplate.ImageTemplateFileCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}

	if input == nil {
		return result
	}

	result.Type = "File"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.SourceUri != nil {
		result.FileSourceUri = *input.SourceUri
	}

	if input.Sha256Checksum != nil {
		result.FileSha256Checksum = *input.Sha256Checksum
	}

	if input.Destination != nil {
		result.FileDestinationPath = *input.Destination
	}

	return result
}

func flattenCustomizerShell(input *virtualmachineimagetemplate.ImageTemplateShellCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "Shell"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.ScriptUri != nil {
		result.ShellScriptUri = *input.ScriptUri
	}

	if input.Sha256Checksum != nil {
		result.ShellSha256Checksum = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result.ShellCommands = *input.Inline
	}

	return result
}

func flattenCustomizerPowerShell(input *virtualmachineimagetemplate.ImageTemplatePowerShellCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "PowerShell"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.ScriptUri != nil {
		result.PowershellScriptUri = *input.ScriptUri
	}

	if input.Sha256Checksum != nil {
		result.PowershellSha256Checksum = *input.Sha256Checksum
	}

	if input.Inline != nil {
		result.PowershellCommands = *input.Inline
	}

	if input.RunAsSystem != nil {
		result.PowershellRunAsSystem = pointer.From(input.RunAsSystem)
	}

	if input.RunElevated != nil {
		result.PowershellRunElevated = *input.RunElevated
	}

	if input.ValidExitCodes != nil {
		result.PowershellValidExitCodes = *input.ValidExitCodes
	}

	return result
}

func flattenCustomizerWindowsRestart(input *virtualmachineimagetemplate.ImageTemplateRestartCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "WindowsRestart"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.RestartCommand != nil {
		result.WindowsRestartCommand = *input.RestartCommand
	}

	if input.RestartCheckCommand != nil {
		result.WindowsRestartCheckCommand = *input.RestartCheckCommand
	}

	if input.RestartTimeout != nil {
		result.WindowsRestartTimeout = *input.RestartTimeout
	}

	return result
}

func flattenCustomizerWindowsUpdate(input *virtualmachineimagetemplate.ImageTemplateWindowsUpdateCustomizer) ImageBuilderTemplateCustomizer {
	result := ImageBuilderTemplateCustomizer{}
	if input == nil {
		return result
	}

	result.Type = "WindowsUpdate"

	if input.Name != nil {
		result.Name = *input.Name
	}

	if input.SearchCriteria != nil {
		result.WindowsUpdateSearchCriteria = *input.SearchCriteria
	}

	if input.Filters != nil {
		result.WindowsUpdateFilters = *input.Filters
	}

	if input.UpdateLimit != nil {
		result.WindowsUpdateLimit = *input.UpdateLimit
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
func validateImageTemplateCustomizerInputForFileType(d *schema.ResourceData, customizer ImageBuilderTemplateCustomizer, index int) error {
	sourceUri := customizer.FileSourceUri
	if sourceUri == "" {
		return fmt.Errorf("`file_source_uri` must be specified if the customizer type is File")
	}

	destinationPath := customizer.FileDestinationPath
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

func validateImageTemplateCustomizerInputForPowerShellType(d *schema.ResourceData, customizer ImageBuilderTemplateCustomizer, index int) error {
	scriptUri := customizer.PowershellScriptUri
	commandsRaw := customizer.PowershellCommands

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

func validateImageTemplateCustomizerInputForShellType(d *schema.ResourceData, customizer ImageBuilderTemplateCustomizer, index int) error {
	scriptUri := customizer.ShellScriptUri
	commandsRaw := customizer.ShellCommands

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
	var imageName string
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
