// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package serviceconnector

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -properties "name" -compare-values "container_app_name:container_app_id,resource_group_name:container_app_id,subscription_id:container_app_id" -test-name storageBlob

var (
	_ sdk.ResourceWithIdentity = ContainerAppConnectorResource{}
	_ sdk.ResourceWithUpdate   = ContainerAppConnectorResource{}
)

type ContainerAppConnectorResource struct{}

type ContainerAppConnectorResourceModel struct {
	Name             string             `tfschema:"name"`
	ContainerAppId   string             `tfschema:"container_app_id"`
	TargetResourceId string             `tfschema:"target_resource_id"`
	ClientType       string             `tfschema:"client_type"`
	AuthInfo         []AuthInfoModel    `tfschema:"authentication"`
	VnetSolution     string             `tfschema:"vnet_solution"`
	SecretStore      []SecretStoreModel `tfschema:"secret_store"`
	Scope            string             `tfschema:"scope"`
}

// Expand the SDK's scope segment so identity includes the Container App, not just the connection name.
type containerAppConnectionIdentity struct {
	servicelinker.ScopedLinkerId
}

func (r ContainerAppConnectorResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"container_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: containerapps.ValidateContainerAppID,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"client_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(servicelinker.ClientTypeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(servicelinker.ClientTypeNone),
				string(servicelinker.ClientTypeDotnet),
				string(servicelinker.ClientTypeJava),
				string(servicelinker.ClientTypePython),
				string(servicelinker.ClientTypeGo),
				string(servicelinker.ClientTypePhp),
				string(servicelinker.ClientTypeRuby),
				string(servicelinker.ClientTypeDjango),
				string(servicelinker.ClientTypeNodejs),
				string(servicelinker.ClientTypeSpringBoot),
			}, false),
		},

		"secret_store": secretStoreSchema(),

		"vnet_solution": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(servicelinker.PossibleValuesForVNetSolutionType(), false),
		},

		"authentication": authInfoSchema(),
	}
}

func (r ContainerAppConnectorResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ContainerAppConnectorResource) ModelObject() any {
	return &ContainerAppConnectorResourceModel{}
}

func (r ContainerAppConnectorResource) ResourceType() string {
	return "azurerm_container_app_connection"
}

func (r ContainerAppConnectorResource) Identity() resourceids.ResourceId {
	return &containerAppConnectionIdentity{}
}

func (r ContainerAppConnectorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ContainerAppConnectorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.ServiceConnector.ServiceLinkerClient

			id := servicelinker.NewScopedLinkerID(model.ContainerAppId, model.Name)
			existing, err := client.LinkerGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			authInfo, err := expandServiceConnectorAuthInfoForCreate(model.AuthInfo)
			if err != nil {
				return fmt.Errorf("expanding `authentication`: %+v", err)
			}

			serviceConnectorProperties := servicelinker.LinkerProperties{
				AuthInfo: authInfo,
				Scope:    pointer.To(model.Scope),
			}

			if storageAccountId, err := commonids.ParseStorageAccountID(model.TargetResourceId); err == nil {
				targetResourceId := fmt.Sprintf("%s/blobServices/default", storageAccountId.ID())
				serviceConnectorProperties.TargetService = servicelinker.AzureResource{
					Id:   &targetResourceId,
					Type: servicelinker.TargetServiceTypeAzureResource,
				}
			} else {
				serviceConnectorProperties.TargetService = servicelinker.AzureResource{
					Id: &model.TargetResourceId,
				}
			}

			if model.SecretStore != nil {
				serviceConnectorProperties.SecretStore = expandSecretStore(model.SecretStore)
			}

			if model.ClientType != "" {
				serviceConnectorProperties.ClientType = pointer.ToEnum[servicelinker.ClientType](model.ClientType)
			}

			if model.VnetSolution != "" {
				vNetSolution := servicelinker.VNetSolution{
					Type: pointer.ToEnum[servicelinker.VNetSolutionType](model.VnetSolution),
				}
				serviceConnectorProperties.VNetSolution = &vNetSolution
			}

			props := servicelinker.LinkerResource{
				Id:         pointer.To(id.ID()),
				Name:       pointer.To(model.Name),
				Properties: serviceConnectorProperties,
			}

			if err := client.LinkerCreateOrUpdateThenPoll(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &containerAppConnectionIdentity{ScopedLinkerId: id})
		},
	}
}

func (r ContainerAppConnectorResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceConnector.ServiceLinkerClient
			id, err := servicelinker.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.LinkerGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r ContainerAppConnectorResource) flatten(metadata sdk.ResourceMetaData, id *servicelinker.ScopedLinkerId, model *servicelinker.LinkerResource) error {
	if model == nil {
		return fmt.Errorf("retrieving %s: response model was nil", *id)
	}

	props := model.Properties
	if props.AuthInfo == nil || props.TargetService == nil {
		return fmt.Errorf("retrieving %s: authentication or target service was missing", *id)
	}

	containerAppId, err := containerapps.ParseContainerAppIDInsensitively(id.ResourceUri)
	if err != nil {
		return fmt.Errorf("parsing Container App ID: %+v", err)
	}

	pwd := metadata.ResourceData.Get("authentication.0.secret").(string)
	state := ContainerAppConnectorResourceModel{
		Name:             id.LinkerName,
		ContainerAppId:   containerAppId.ID(),
		TargetResourceId: flattenTargetService(props.TargetService),
		AuthInfo:         flattenServiceConnectorAuthInfo(props.AuthInfo, pwd),
		Scope:            pointer.From(props.Scope),
	}

	if props.ClientType != nil {
		state.ClientType = string(*props.ClientType)
	}

	if props.VNetSolution != nil && props.VNetSolution.Type != nil {
		state.VnetSolution = string(*props.VNetSolution.Type)
	}

	if props.SecretStore != nil {
		state.SecretStore = flattenSecretStore(*props.SecretStore)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &containerAppConnectionIdentity{ScopedLinkerId: *id}); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r ContainerAppConnectorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceConnector.LinksClient
			id, err := links.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.LinkerDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppConnectorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceConnector.LinksClient
			id, err := links.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ContainerAppConnectorResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			linkerProps := links.LinkerProperties{}
			d := metadata.ResourceData

			if d.HasChange("client_type") {
				linkerProps.ClientType = pointer.ToEnum[links.ClientType](state.ClientType)
			}

			if d.HasChange("vnet_solution") {
				vnetSolution := links.VNetSolution{
					Type: pointer.ToEnum[links.VNetSolutionType](state.VnetSolution),
				}
				linkerProps.VNetSolution = &vnetSolution
			}

			if d.HasChange("secret_store") {
				linkerProps.SecretStore = pointer.To(links.SecretStore{KeyVaultId: expandSecretStore(state.SecretStore).KeyVaultId})
			}

			if d.HasChange("authentication") {
				authInfo, err := expandServiceConnectorAuthInfoForUpdate(state.AuthInfo)
				if err != nil {
					return fmt.Errorf("expanding `authentication`: %+v", err)
				}

				linkerProps.AuthInfo = authInfo
			}

			props := links.LinkerPatch{
				Properties: &linkerProps,
			}

			if err := client.LinkerUpdateThenPoll(ctx, *id, props); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ContainerAppConnectorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return servicelinker.ValidateScopedLinkerID
}

func (id containerAppConnectionIdentity) Segments() []resourceids.Segment {
	segments := id.ScopedLinkerId.Segments()[1:]
	// The parser requires distinct names for the two static "providers" segments.
	segments[0].Name = "staticServiceLinkerProviders"
	return append(containerapps.ContainerAppId{}.Segments(), segments...)
}

func (id *containerAppConnectionIdentity) FromParseResult(input resourceids.ParseResult) error {
	containerAppId := containerapps.ContainerAppId{}
	if err := containerAppId.FromParseResult(input); err != nil {
		return err
	}

	name, ok := input.Parsed["linkerName"]
	if !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "linkerName", input)
	}

	id.ScopedLinkerId = servicelinker.NewScopedLinkerID(containerAppId.ID(), name)
	return nil
}
