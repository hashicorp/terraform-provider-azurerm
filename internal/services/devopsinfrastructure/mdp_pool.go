package devopsinfrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-10-19/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	AgentProfileKindStateless = "Stateless"
	AgentProfileKindStateful  = "Stateful"
)

var _ sdk.Resource = PoolResource{}

type PoolResource struct{}

func (PoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"agent_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"grace_period_time_span": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"kind": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							AgentProfileKindStateless,
							AgentProfileKindStateful,
						}, false),
					},
					"max_agent_lifetime": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					"resource_predictions": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"resource_predictions_profile": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"kind": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForResourcePredictionsProfileType(), false),
								},
								"prediction_preference": {
									Type:         pluginsdk.TypeString,
									Computed:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForPredictionPreference(), false),
								},
							},
						},
					},
				},
			},
		},
		"dev_center_project_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"fabric_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"image": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"aliases": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"buffer": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"resource_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"well_known_image_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"kind": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string("Vmss"),
						}, false),
					},
					"network_profile": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subnet_id": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"os_profile": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"logon_type": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForLogonType(), false),
								},
								"secrets_management_settings": {
									Type: pluginsdk.TypeList,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"certificate_store_location": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"key_exportable": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},
											"observed_certificates": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeString,
												},
											},
										},
									},
								},
							},
						},
					},
					"sku": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
					"storage_profile": {
						Type: pluginsdk.TypeList,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"data_disks": {
									Type: pluginsdk.TypeList,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"caching": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(pools.PossibleValuesForCachingType(), false),
											},
											"disk_size": {
												Type:     pluginsdk.TypeInt,
												Required: true,
											},
											"drive_letter": {
												Type:     pluginsdk.TypeString,
												Optional: true,
											},
											"storage_account_type": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(pools.PossibleValuesForStorageAccountType(), false),
											},
										},
									},
								},
								"os_disk_storage_account_type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForOsDiskStorageAccountType(), false),
								},
							},
						},
					},
				},
			},
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"location": commonschema.Location(),
		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"organization_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kind": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string("AzureDevOps"),
							string("GitHub"),
						}, false),
					},
					"organizations": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"parallelism": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
								"projects": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"repositories": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"url": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsURLWithHTTPS,
								},
							},
						},
					},
					"permission_profile": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"groups": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"kind": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(pools.PossibleValuesForAzureDevOpsPermissionType(), false),
								},
								"users": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},
		"tags": commonschema.Tags(),
		"type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (PoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (PoolResource) ModelObject() interface{} {
	return &PoolResourceSchema{}
}

type PoolResourceSchema struct {
	AgentProfile               AgentProfileSchema        				  `tfschema:"agent_profile"`
	DevCenterProjectResourceId string                                     `tfschema:"dev_center_project_id"`
	FabricProfile              FabricProfileSchema       				  `tfschema:"fabric_profile"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                   string                                     `tfschema:"location"`
	MaximumConcurrency         int64                                      `tfschema:"maximum_concurrency"`
	Name                       string                                     `tfschema:"name"`
	OrganizationProfile        OrganizationProfileSchema 				  `tfschema:"organization_profile"`
	Tags                       map[string]string                          `tfschema:"tags"`
	Type                       string                                     `tfschema:"type"`
}

type AgentProfileSchema struct {
	GracePeriodTimeSpan        *string                                        `tfschema:"grace_period_time_span"`
	Kind                       string                                         `tfschema:"kind"` // Stateless, Stateful
	MaxAgentLifetime           *string                                        `tfschema:"max_agent_lifetime"`
	PredictionPreference       *string                                        `tfschema:"prediction_preference"`
	ResourcePredictions        *interface{}                                   `tfschema:"resource_predictions"`
	ResourcePredictionsProfile AgentPredictionsProfileSchema 				  `tfschema:"resource_predictions_profile"` // Automatic, Manual
}

type AgentPredictionsProfileSchema struct {
	Kind                 string  `tfschema:"kind"` // Stateless, Stateful
	PredictionPreference *string `tfschema:"prediction_preference"`
}

type FabricProfileSchema struct {
	Images         []ImageSchema         `tfschema:"images"`
	Kind           string                                 `tfschema:"kind"`
	NetworkProfile *NetworkProfileSchema `tfschema:"network_profile"`
	OsProfile      *OSProfileSchema      `tfschema:"os_profile"`
	Sku            SkuSchema  `tfschema:"sku"`
	StorageProfile *StorageProfileSchema `tfschema:"storage_profile"`
}

type ImageSchema struct {
	Aliases            *[]string `tfschema:"aliases"`
	Buffer             *string   `tfschema:"buffer"`
	ResourceId         *string   `tfschema:"resource_id"`
	WellKnownImageName *string   `tfschema:"well_known_image_name"`
}
type OSProfileSchema struct {
	LogonType                 string                                            `tfschema:"logon_type"` // PossibleValuesForLogonType()
	SecretsManagementSettings *SecretsManagementSettingsSchema `tfschema:"secrets_management_settings"`
}

type SecretsManagementSettingsSchema struct {
	CertificateStoreLocation *string  `tfschema:"certificate_store_location"`
	KeyExportable            bool     `tfschema:"key_exportable"`
	ObservedCertificates     []string `tfschema:"observed_certificates"`
}

type NetworkProfileSchema struct {
	SubnetId string `tfschema:"subnet_id"`
}

type SkuSchema struct {
	Name string `tfschema:"name"`
}

type StorageProfileSchema struct {
	DataDisks                *[]DataDiskSchema `tfschema:"data_disks"`
	OsDiskStorageAccountType string                             `tfschema:"os_disk_storage_account_type"`
}

type DataDiskSchema struct {
	Caching            string  `tfschema:"caching"`
	DiskSizeGiB        *int64  `tfschema:"disk_size"`
	DriveLetter        *string `tfschema:"drive_letter"`
	StorageAccountType string  `tfschema:"storage_account_type"`
}

type OrganizationProfileSchema struct {
	Organizations     []OrganizationSchema                `tfschema:"organizations"`
	PermissionProfile OrganizationPermissionProfileSchema `tfschema:"permission_profile"`
	Kind              string                                               `tfschema:"kind"`
}

type OrganizationSchema struct {
	Parallelism  int64     `tfschema:"parallelism"`
	Projects     []string  `tfschema:"projects"`
	Repositories *[]string `tfschema:"repositories"`
	Url          string    `tfschema:"url"`
}

type OrganizationPermissionProfileSchema struct {
	Groups []string `tfschema:"groups"`
	Kind   string   `tfschema:"kind"`
	Users  []string `tfschema:"users"`
}

func (PoolResource) ResourceType() string {
	return "azurerm_mdp_pool"
}

func (r PoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevOpsInfrastructure.PoolsClient

			var config PoolResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			id := pools.NewPoolID(subscriptionId, config.Name, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload pools.Pool
			if err := r.mapPoolResourceSchemaToPool(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevOpsInfrastructure.PoolsClient

			var config PoolResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			properties := existing.Model

			if err := r.mapPoolResourceSchemaToPool(config, properties); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (PoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevOpsInfrastructure.PoolsClient

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				resourceData := &PoolResourceSchema{}
				metadata.Encode(resourceData)

				metadata.ResourceData.Set("location", model.Location)

				if props := model.Properties; props != nil {
					// if there are properties to set into state do that here
				}

				if err := tags.FlattenAndSet(metadata.ResourceData, model.Tags); err != nil {
					return fmt.Errorf("setting `tags`: %+v", err)
				}
			}
			return nil
		},
	}
}

func (PoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevOpsInfrastructure.PoolsClient

			id, err := pools.ParsePoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (PoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return pools.ValidatePoolID
}

func (r PoolResource) mapPoolResourceSchemaToPool(input PoolResourceSchema, output *pools.Pool) error {
	identity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(input.Identity)
	if err != nil {
		return fmt.Errorf("expanding SystemAndUserAssigned Identity: %+v", err)
	}

	output.Identity = identity
	output.Location = location.Normalize(input.Location)
	output.Name = &input.Name
	output.Tags = &input.Tags
	output.Type = &input.Type

	if output.Properties == nil {
		output.Properties = &pools.PoolProperties{}
	}

	if err := r.mapAgentProfileSchemaToPoolProperties(input.AgentProfile, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "AgentProfile", "Properties", err)
	}

	if err := r.mapFabricProfileSchemaToPoolProperties(input.FabricProfile, output.Properties); err != nil {
		return fmt.Errorf("mapping schema model to sdk model: %+v", err)
	}

	return nil
}

func (r PoolResource) mapAgentProfileSchemaToPoolProperties(input AgentProfileSchema, output *pools.PoolProperties) error {
	if input.Kind == AgentProfileKindStateful {
		stateful := &pools.Stateful{
			GracePeriodTimeSpan: input.GracePeriodTimeSpan,
			MaxAgentLifetime:    input.MaxAgentLifetime,
			ResourcePredictions: input.ResourcePredictions,
		}
		if input.ResourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeAutomatic) {
			automaticPredictionsProfile := &pools.AutomaticResourcePredictionsProfile{
				Kind:                 pools.ResourcePredictionsProfileTypeAutomatic,
				PredictionPreference: (*pools.PredictionPreference)(input.ResourcePredictionsProfile.PredictionPreference),
			}
			stateful.ResourcePredictionsProfile = automaticPredictionsProfile
		}
		if input.ResourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeManual) {
			manualPredictionsProfile := &pools.ManualResourcePredictionsProfile{
				Kind: pools.ResourcePredictionsProfileTypeAutomatic,
			}
			stateful.ResourcePredictionsProfile = manualPredictionsProfile
		}
		output.AgentProfile = stateful.AgentProfile()
	}

	if input.Kind == AgentProfileKindStateless {
		agentProfileStateless := &pools.StatelessAgentProfile{
			Kind:                input.Kind,
			ResourcePredictions: input.ResourcePredictions,
		}
		if input.ResourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeAutomatic) {
			automaticPredictionsProfile := &pools.AutomaticResourcePredictionsProfile{
				Kind:                 pools.ResourcePredictionsProfileTypeAutomatic,
				PredictionPreference: (*pools.PredictionPreference)(input.ResourcePredictionsProfile.PredictionPreference),
			}
			agentProfileStateless.ResourcePredictionsProfile = automaticPredictionsProfile
		}
		if input.ResourcePredictionsProfile.Kind == string(pools.ResourcePredictionsProfileTypeManual) {
			manualPredictionsProfile := &pools.ManualResourcePredictionsProfile{
				Kind: pools.ResourcePredictionsProfileTypeAutomatic,
			}
			agentProfileStateless.ResourcePredictionsProfile = manualPredictionsProfile
		}

		output.AgentProfile = agentProfileStateless.AgentProfile()
	}

	return fmt.Errorf("Invalid Agent Profile Kind Provided: %s", input.Kind)
}

func (r PoolResource) mapFabricProfileSchemaToPoolProperties(input FabricProfileSchema, output *pools.PoolProperties) error {
	if input.Kind == "Vmss" {
		vmssFabricProfile := pools.VMSSFabricProfile{
			Images:         expandImageSchemaToPoolImages(input.Images),
			NetworkProfile: expandNetworkProfileSchemaToNetworkProfile(input.NetworkProfile),
			OsProfile:      expandOsProfileSchemaToOsProfile(input.OsProfile),
			Sku:            expandSkuSchemaToSku(input.Sku),
			StorageProfile: expandStorageProfileSchemaToPoolAzureStorageProfile(input.StorageProfile),
			Kind:           input.Kind,
		}
		output.FabricProfile = vmssFabricProfile
	}

	return fmt.Errorf("Invalid Fabric Profile Kind Provided: %s", input.Kind)
}

func expandImageSchemaToPoolImages(input []ImageSchema) []pools.PoolImage {
	output := []pools.PoolImage{}

	for _, image := range input {
		imageOut := pools.PoolImage{
			Aliases:            image.Aliases,
			Buffer:             image.Buffer,
			ResourceId:         image.ResourceId,
			WellKnownImageName: image.WellKnownImageName,
		}
		output = append(output, imageOut)
	}
	return output
}

func expandNetworkProfileSchemaToNetworkProfile(input *NetworkProfileSchema) *pools.NetworkProfile {
	if input == nil {
		return nil
	}
	output := &pools.NetworkProfile{
		SubnetId: input.SubnetId,
	}

	return output
}

func expandOsProfileSchemaToOsProfile(input *OSProfileSchema) *pools.OsProfile {
	if input == nil {
		return nil
	}
	loginType := pools.LogonType(input.LogonType)
	output := &pools.OsProfile{
		LogonType:                 &loginType,
		SecretsManagementSettings: expandSecretsManagementSettingsSchemaToSecretsManagementSettings(input.SecretsManagementSettings),
	}

	return output
}

func expandSkuSchemaToSku(input SkuSchema) pools.DevOpsAzureSku {
	output := pools.DevOpsAzureSku{
		Name: input.Name,
	}

	return output
}

func expandStorageProfileSchemaToPoolAzureStorageProfile(input *StorageProfileSchema) *pools.StorageProfile {
	if input == nil {
		return nil
	}

	osDiskStorageAccountType := pools.OsDiskStorageAccountType(input.OsDiskStorageAccountType)
	output := &pools.StorageProfile{
		OsDiskStorageAccountType: &osDiskStorageAccountType,
	}

	if input.DataDisks != nil {
		dataDisksOut := []pools.DataDisk{}
		for _, disk := range *input.DataDisks {
			cachingType := pools.CachingType(disk.Caching)
			storageAccountType := pools.StorageAccountType(disk.StorageAccountType)
			diskOut := pools.DataDisk{
				Caching:            &cachingType,
				DiskSizeGiB:        disk.DiskSizeGiB,
				DriveLetter:        disk.DriveLetter,
				StorageAccountType: &storageAccountType,
			}
			dataDisksOut = append(dataDisksOut, diskOut)
		}
		output.DataDisks = &dataDisksOut
	}

	return output
}

func expandSecretsManagementSettingsSchemaToSecretsManagementSettings(input *SecretsManagementSettingsSchema) *pools.SecretsManagementSettings {
	if input == nil {
		return nil
	}

	output := &pools.SecretsManagementSettings{
		CertificateStoreLocation: input.CertificateStoreLocation,
		KeyExportable:            input.KeyExportable,
		ObservedCertificates:     input.ObservedCertificates,
	}

	return output
}