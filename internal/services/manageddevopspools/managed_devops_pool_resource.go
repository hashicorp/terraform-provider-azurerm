package manageddevopspools

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devopsinfrastructure/2024-10-19/pools"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = ManagedDevOpsPoolResource{}

type ManagedDevOpsPoolResource struct{}

func (ManagedDevOpsPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"agent_profile": AgentProfileSchema(),
		"dev_center_project_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"fabric_profile": FabricProfileSchema(),
		"identity":       commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"location":       commonschema.Location(),
		"maximum_concurrency": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"organization_profile": OrganizationProfileSchema(),
		"tags":                 commonschema.Tags(),
	}
}

func (ManagedDevOpsPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (ManagedDevOpsPoolResource) ModelObject() interface{} {
	return &ManagedDevOpsPoolModel{}
}

func (ManagedDevOpsPoolResource) ResourceType() string {
	return "azurerm_managed_devops_pool"
}

func (r ManagedDevOpsPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			var config ManagedDevOpsPoolModel
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
			if err := r.mapResourceModelToPool(config, &payload); err != nil {
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
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

			var config ManagedDevOpsPoolModel
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

			if err := r.mapResourceModelToPool(config, properties); err != nil {
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
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

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

			state := ManagedDevOpsPoolModel{
				Name: id.PoolName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					// if there are properties to set into state do that here
					state.DevCenterProjectResourceId = props.DevCenterProjectResourceId
					state.MaximumConcurrency = props.MaximumConcurrency
					// Add other props
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (ManagedDevOpsPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ManagedDevOpsPools.PoolsClient

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

func (r ManagedDevOpsPoolResource) mapResourceModelToPool(input ManagedDevOpsPoolModel, output *pools.Pool) error {
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

	if err := r.mapAgentProfileModelToPoolProperties(input.AgentProfile, output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "AgentProfile", "Properties", err)
	}

	if err := r.mapFabricProfileSchemaToPoolProperties(input.FabricProfile, output.Properties); err != nil {
		return fmt.Errorf("mapping schema model to sdk model: %+v", err)
	}

	// Add Organization Profile

	return nil
}

func (r ManagedDevOpsPoolResource) mapAgentProfileModelToPoolProperties(input AgentProfileModel, output *pools.PoolProperties) error {
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

func (r ManagedDevOpsPoolResource) mapFabricProfileSchemaToPoolProperties(input FabricProfileModel, output *pools.PoolProperties) error {
	if input.Kind == "Vmss" {
		vmssFabricProfile := pools.VMSSFabricProfile{
			Images:         mapImageModelToPoolImages(input.Images),
			NetworkProfile: mapNetworkProfileModelToNetworkProfile(input.NetworkProfile),
			OsProfile:      mapOsProfileModelToOsProfile(input.OsProfile),
			Sku:            mapDevOpsAzureSkuModelToDevOpsAzureSku(input.Sku),
			StorageProfile: mapStorageProfileModelToDevOpsAzureStorageProfile(input.StorageProfile),
			Kind:           input.Kind,
		}
		output.FabricProfile = vmssFabricProfile
	}

	return fmt.Errorf("Invalid Fabric Profile Kind Provided: %s", input.Kind)
}

func mapImageModelToPoolImages(input []ImageModel) []pools.PoolImage {
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

func mapNetworkProfileModelToNetworkProfile(input *NetworkProfileModel) *pools.NetworkProfile {
	if input == nil {
		return nil
	}
	output := &pools.NetworkProfile{
		SubnetId: input.SubnetId,
	}

	return output
}

func mapOsProfileModelToOsProfile(input *OsProfileModel) *pools.OsProfile {
	if input == nil {
		return nil
	}
	loginType := pools.LogonType(input.LogonType)
	output := &pools.OsProfile{
		LogonType:                 &loginType,
		SecretsManagementSettings: mapSecretsManagementSettingsModelToSecretsManagementSettings(input.SecretsManagementSettings),
	}

	return output
}

func mapDevOpsAzureSkuModelToDevOpsAzureSku(input DevOpsAzureSkuModel) pools.DevOpsAzureSku {
	output := pools.DevOpsAzureSku{
		Name: input.Name,
	}

	return output
}

func mapStorageProfileModelToDevOpsAzureStorageProfile(input *StorageProfileModel) *pools.StorageProfile {
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

func mapSecretsManagementSettingsModelToSecretsManagementSettings(input *SecretsManagementSettingsModel) *pools.SecretsManagementSettings {
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
