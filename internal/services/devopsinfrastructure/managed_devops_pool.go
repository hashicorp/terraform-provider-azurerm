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
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-04-04-preview/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	ManagedDevOpsPoolAgentProfileKindStateless = "Stateless"
	ManagedDevOpsPoolAgentProfileKindStateful  = "Stateful"
)

var _ sdk.Resource = ManagedDevOpsPoolResource{}

type ManagedDevOpsPoolResource struct{}

func (ManagedDevOpsPoolResource) Arguments() map[string]*pluginsdk.Schema {
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
							ManagedDevOpsPoolAgentProfileKindStateless,
							ManagedDevOpsPoolAgentProfileKindStateful,
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

func (ManagedDevOpsPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ManagedDevOpsPoolResource) ModelObject() interface{} {
	return &ManagedDevOpsPoolResourceSchema{}
}

type ManagedDevOpsPoolResourceSchema struct {
	AgentProfile               ManagedDevOpsPoolAgentProfileSchema        `tfschema:"agent_profile"`
	DevCenterProjectResourceId string                                     `tfschema:"dev_center_project_id"`
	FabricProfile              ManagedDevOpsPoolFabricProfileSchema       `tfschema:"fabric_profile"`
	Identity                   []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location                   string                                     `tfschema:"location"`
	MaximumConcurrency         int64                                      `tfschema:"maximum_concurrency"`
	Name                       string                                     `tfschema:"name"`
	OrganizationProfile        ManagedDevOpsPoolOrganizationProfileSchema `tfschema:"organization_profile"`
	Tags                       map[string]string                          `tfschema:"tags"`
	Type                       string                                     `tfschema:"type"`
}

type ManagedDevOpsPoolAgentProfileSchema struct {
	GracePeriodTimeSpan        *string                                        `tfschema:"grace_period_time_span"`
	Kind                       string                                         `tfschema:"kind"` // Stateless, Stateful
	MaxAgentLifetime           *string                                        `tfschema:"max_agent_lifetime"`
	PredictionPreference       *string                                        `tfschema:"prediction_preference"`
	ResourcePredictions        *interface{}                                   `tfschema:"resource_predictions"`
	ResourcePredictionsProfile ManagedDevOpsPoolAgentPredictionsProfileSchema `tfschema:"resource_predictions_profile"` // Automatic, Manual
}

type ManagedDevOpsPoolAgentPredictionsProfileSchema struct {
	Kind                 string  `tfschema:"kind"` // Stateless, Stateful
	PredictionPreference *string `tfschema:"prediction_preference"`
}

type ManagedDevOpsPoolFabricProfileSchema struct {
	Images         []ManagedDevOpsPoolImageSchema         `tfschema:"images"`
	Kind           string                                 `tfschema:"kind"`
	NetworkProfile *ManagedDevOpsPoolNetworkProfileSchema `tfschema:"network_profile"`
	OsProfile      *ManagedDevOpsPoolOSProfileSchema      `tfschema:"os_profile"`
	Sku            ManagedDevOpsPoolDevOpsAzureSkuSchema  `tfschema:"sku"`
	StorageProfile *ManagedDevOpsPoolStorageProfileSchema `tfschema:"storage_profile"`
}

type ManagedDevOpsPoolImageSchema struct {
	Aliases            *[]string `tfschema:"aliases"`
	Buffer             *string   `tfschema:"buffer"`
	ResourceId         *string   `tfschema:"resource_id"`
	WellKnownImageName *string   `tfschema:"well_known_image_name"`
}
type ManagedDevOpsPoolOSProfileSchema struct {
	LogonType                 string                                            `tfschema:"logon_type"` // PossibleValuesForLogonType()
	SecretsManagementSettings *ManagedDevOpsPoolSecretsManagementSettingsSchema `tfschema:"secrets_management_settings"`
}

type ManagedDevOpsPoolSecretsManagementSettingsSchema struct {
	CertificateStoreLocation *string  `tfschema:"certificate_store_location"`
	KeyExportable            bool     `tfschema:"key_exportable"`
	ObservedCertificates     []string `tfschema:"observed_certificates"`
}

type ManagedDevOpsPoolNetworkProfileSchema struct {
	SubnetId string `tfschema:"subnet_id"`
}

type ManagedDevOpsPoolDevOpsAzureSkuSchema struct {
	Name string `tfschema:"name"`
}

type ManagedDevOpsPoolStorageProfileSchema struct {
	DataDisks                *[]ManagedDevOpsPoolDataDiskSchema `tfschema:"data_disks"`
	OsDiskStorageAccountType string                             `tfschema:"os_disk_storage_account_type"`
}

type ManagedDevOpsPoolDataDiskSchema struct {
	Caching            string  `tfschema:"caching"`
	DiskSizeGiB        *int64  `tfschema:"disk_size"`
	DriveLetter        *string `tfschema:"drive_letter"`
	StorageAccountType string  `tfschema:"storage_account_type"`
}

type ManagedDevOpsPoolOrganizationProfileSchema struct {
	Organizations     []ManagedDevOpsPoolOrganizationSchema                `tfschema:"organizations"`
	PermissionProfile ManagedDevOpsPoolOrganizationPermissionProfileSchema `tfschema:"permission_profile"`
	Kind              string                                               `tfschema:"kind"`
}

type ManagedDevOpsPoolOrganizationSchema struct {
	Parallelism  int64     `tfschema:"parallelism"`
	Projects     []string  `tfschema:"projects"`
	Repositories *[]string `tfschema:"repositories"`
	Url          string    `tfschema:"url"`
}

type ManagedDevOpsPoolOrganizationPermissionProfileSchema struct {
	Groups []string `tfschema:"groups"`
	Kind   string   `tfschema:"kind"`
	Users  []string `tfschema:"users"`
}

func (ManagedDevOpsPoolResource) ResourceType() string {
	return "azurerm_managed_devops_pool"
}

func (r ManagedDevOpsPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevOpsInfrastructure.PoolsClient

			var config ManagedDevOpsPoolResourceSchema
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
			if err := r.mapManagedDevOpsPoolResourceSchemaToPool(config, &payload); err != nil {
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

func (r ManagedDevOpsPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DevOpsInfrastructure.PoolsClient

			var config ManagedDevOpsPoolResourceSchema
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

			if err := r.mapManagedDevOpsPoolResourceSchemaToPool(config, properties); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (ManagedDevOpsPoolResource) Read() sdk.ResourceFunc {
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
				resourceData := &ManagedDevOpsPoolResourceSchema{}
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

func (ManagedDevOpsPoolResource) Delete() sdk.ResourceFunc {
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

func (ManagedDevOpsPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return pools.ValidatePoolID
}

func (r ManagedDevOpsPoolResource) mapManagedDevOpsPoolResourceSchemaToPool(input ManagedDevOpsPoolResourceSchema, output *pools.Pool) error {
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

	if err := r.mapManagedDevOpsPoolAgentProfileSchemaToPoolProperties(input.AgentProfile, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "AgentProfile", "Properties", err)
	}

	if err := r.mapManagedDevOpsPoolFabricProfileSchemaToPoolProperties(input.FabricProfile, output.Properties); err != nil {
		return fmt.Errorf("mapping schema model to sdk model: %+v", err)
	}

	return nil
}

func (r ManagedDevOpsPoolResource) mapManagedDevOpsPoolAgentProfileSchemaToPoolProperties(input ManagedDevOpsPoolAgentProfileSchema, output *pools.PoolProperties) error {
	if input.Kind == ManagedDevOpsPoolAgentProfileKindStateful {
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

	if input.Kind == ManagedDevOpsPoolAgentProfileKindStateless {
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

func (r ManagedDevOpsPoolResource) mapManagedDevOpsPoolFabricProfileSchemaToPoolProperties(input ManagedDevOpsPoolFabricProfileSchema, output *pools.PoolProperties) error {
	if input.Kind == "Vmss" {
		vmssFabricProfile := pools.VMSSFabricProfile{
			Images:         expandManagedDevOpsPoolImageSchemaToPoolImages(input.Images),
			NetworkProfile: expandManagedDevOpsPoolNetworkProfileSchemaToNetworkProfile(input.NetworkProfile),
			OsProfile:      expandManagedDevOpsPoolOsProfileSchemaToOsProfile(input.OsProfile),
			Sku:            expandManagedDevOpsPoolDevOpsAzureSkuSchemaToDevOpsAzureSku(input.Sku),
			StorageProfile: expandManagedDevOpsPoolStorageProfileSchemaToDevOpsAzureStorageProfile(input.StorageProfile),
			Kind:           input.Kind,
		}
		output.FabricProfile = vmssFabricProfile
	}

	return fmt.Errorf("Invalid Fabric Profile Kind Provided: %s", input.Kind)
}

func expandManagedDevOpsPoolImageSchemaToPoolImages(input []ManagedDevOpsPoolImageSchema) []pools.PoolImage {
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

func expandManagedDevOpsPoolNetworkProfileSchemaToNetworkProfile(input *ManagedDevOpsPoolNetworkProfileSchema) *pools.NetworkProfile {
	if input == nil {
		return nil
	}
	output := &pools.NetworkProfile{
		SubnetId: input.SubnetId,
	}

	return output
}

func expandManagedDevOpsPoolOsProfileSchemaToOsProfile(input *ManagedDevOpsPoolOSProfileSchema) *pools.OsProfile {
	if input == nil {
		return nil
	}
	loginType := pools.LogonType(input.LogonType)
	output := &pools.OsProfile{
		LogonType:                 &loginType,
		SecretsManagementSettings: expandManagedDevOpsPoolSecretsManagementSettingsSchemaToSecretsManagementSettings(input.SecretsManagementSettings),
	}

	return output
}

func expandManagedDevOpsPoolDevOpsAzureSkuSchemaToDevOpsAzureSku(input ManagedDevOpsPoolDevOpsAzureSkuSchema) pools.DevOpsAzureSku {
	output := pools.DevOpsAzureSku{
		Name: input.Name,
	}

	return output
}

func expandManagedDevOpsPoolStorageProfileSchemaToDevOpsAzureStorageProfile(input *ManagedDevOpsPoolStorageProfileSchema) *pools.StorageProfile {
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

func expandManagedDevOpsPoolSecretsManagementSettingsSchemaToSecretsManagementSettings(input *ManagedDevOpsPoolSecretsManagementSettingsSchema) *pools.SecretsManagementSettings {
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
