package resourceconnector

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ResourceConnectorApplianceResource{}

var _ sdk.ResourceWithCustomizeDiff = ResourceConnectorApplianceResource{}

type ResourceConnectorApplianceResource struct{}

type ApplianceModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Location          string                         `tfschema:"location"`
	Distro            appliances.Distro              `tfschema:"distro"`
	Identity          []identity.ModelSystemAssigned `tfschema:"identity"`
	Provider          appliances.Provider            `tfschema:"infrastructure_provider"`
	PublicKey         string                         `tfschema:"public_key"`
	Tags              map[string]string              `tfschema:"tags"`
	Version           string                         `tfschema:"version"`
}

func (r ResourceConnectorApplianceResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 260),
				validation.StringMatch(regexp.MustCompile(`[^+#%&'?/,%\\]+$`), "any of '+', '#', '%', '&', ''', '?', '/', ',', '%', '&', '\\', are not allowed"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"distro": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(appliances.DistroAKSEdge),
			}, false),
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"infrastructure_provider": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(appliances.ProviderHCI),
				string(appliances.ProviderSCVMM),
				string(appliances.ProviderVMWare),
			}, false),
		},

		"public_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),

		"version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"latest",
			}, false),
		},
	}
}

func (r ResourceConnectorApplianceResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ResourceConnectorApplianceResource) ModelObject() interface{} {
	return &ResourceConnectorApplianceResource{}
}

func (r ResourceConnectorApplianceResource) ResourceType() string {
	return "azurerm_resource_connector_appliance"
}

func (r ResourceConnectorApplianceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApplianceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.ResourceConnector.AppliancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := appliances.NewApplianceID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			if model.PublicKey != "" {
				return fmt.Errorf("the public key can not be set when creating %s. it could be set after a deploy as an update operation", id)
			}

			if model.Version != "" {
				return fmt.Errorf("the version cannot be set when creating %s. it could be set from upgrade call", id)
			}

			parameters := appliances.Appliance{
				Location: model.Location,
				Properties: &appliances.ApplianceProperties{
					Distro: pointer.To(model.Distro),
					InfrastructureConfig: &appliances.AppliancePropertiesInfrastructureConfig{
						Provider: pointer.To(model.Provider),
					},
				},
				Tags: pointer.To(model.Tags),
			}

			identity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding SystemAssigned Identity: %+v", err)
			}

			parameters.Identity = identity

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ResourceConnectorApplianceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ResourceConnector.AppliancesClient

			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("Decoding state for %s", *id)
			var model ApplianceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := resp.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			identity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding SystemAssigned Identity: %+v", err)
			}

			parameters.Identity = identity

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if parameters.Properties == nil {
				parameters.Properties = &appliances.ApplianceProperties{}
			}

			if metadata.ResourceData.HasChange("public_key") {
				parameters.Properties.PublicKey = pointer.To(model.PublicKey)
			}

			if metadata.ResourceData.HasChange("version") {
				parameters.Properties.Version = pointer.To(model.Version)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ResourceConnectorApplianceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.ResourceConnector.AppliancesClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %q: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			state := ApplianceModel{
				Name:              pointer.From(resp.Model.Name),
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.NormalizeNilable(pointer.To(resp.Model.Location)),
				Identity:          identity.FlattenSystemAssignedToModel(resp.Model.Identity),
			}

			state.Tags = pointer.From(resp.Model.Tags)
			if v := resp.Model.Properties; v != nil {
				state.Distro = pointer.From(v.Distro)
				state.PublicKey = pointer.From(v.PublicKey)
				state.Version = pointer.From(v.Version)
				if p := v.InfrastructureConfig; p != nil {
					state.Provider = pointer.From(p.Provider)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r ResourceConnectorApplianceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ResourceConnector.AppliancesClient
			id, err := appliances.ParseApplianceID(metadata.ResourceData.Id())
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

func (r ResourceConnectorApplianceResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			if rd.HasChange("public_key") {
				old, _ := rd.GetChange("public_key")
				if old.(string) != "" {
					return fmt.Errorf("the public_key can not be updated once it is set")
				}
			}
			return nil
		},
	}
}

func (r ResourceConnectorApplianceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return appliances.ValidateApplianceID
}
