// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package serviceconnector

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2022-05-01/links"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicelinker/2024-04-01/servicelinker"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = KubernetesClusterConnectorResource{}

type KubernetesClusterConnectorResource struct{}

type KubernetesClusterConnectorResourceModel struct {
	Name                string             `tfschema:"name"`
	KubernetesClusterId string             `tfschema:"kubernetes_cluster_id"`
	TargetResourceId    string             `tfschema:"target_resource_id"`
	ClientType          string             `tfschema:"client_type"`
	AuthInfo            []AuthInfoModel    `tfschema:"authentication"`
	VnetSolution        string             `tfschema:"vnet_solution"`
	SecretStore         []SecretStoreModel `tfschema:"secret_store"`
}

func (r KubernetesClusterConnectorResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"kubernetes_cluster_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateKubernetesClusterID,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(servicelinker.VNetSolutionTypeServiceEndpoint),
				string(servicelinker.VNetSolutionTypePrivateLink),
			}, false),
		},

		"authentication": authInfoSchema(),
	}
}

func (r KubernetesClusterConnectorResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesClusterConnectorResource) ModelObject() interface{} {
	return &KubernetesClusterConnectorResourceModel{}
}

func (r KubernetesClusterConnectorResource) ResourceType() string {
	return "azurerm_kubernetes_cluster_connection"
}

func (r KubernetesClusterConnectorResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model KubernetesClusterConnectorResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.ServiceConnector.ServiceLinkerClient

			id := servicelinker.NewScopedLinkerID(model.KubernetesClusterId, model.Name)
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
			}

			if storageAccountId, err := commonids.ParseStorageAccountID(model.TargetResourceId); err == nil {
				targetResourceId := fmt.Sprintf("%s/blobServices/default", storageAccountId.ID())
				serviceConnectorProperties.TargetService = servicelinker.AzureResource{
					Id: &targetResourceId,
				}
			} else {
				serviceConnectorProperties.TargetService = servicelinker.AzureResource{
					Id: &model.TargetResourceId,
				}
			}

			if model.SecretStore != nil {
				secretStore := expandSecretStore(model.SecretStore)
				serviceConnectorProperties.SecretStore = secretStore
			}

			if model.ClientType != "" {
				clientType := servicelinker.ClientType(model.ClientType)
				serviceConnectorProperties.ClientType = &clientType
			}

			if model.VnetSolution != "" {
				vNetSolutionType := servicelinker.VNetSolutionType(model.VnetSolution)
				vNetSolution := servicelinker.VNetSolution{
					Type: &vNetSolutionType,
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
			return nil
		},
	}
}

func (r KubernetesClusterConnectorResource) Read() sdk.ResourceFunc {
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

			pwd := metadata.ResourceData.Get("authentication.0.secret").(string)

			if model := resp.Model; model != nil {
				props := model.Properties
				if props.AuthInfo == nil || props.TargetService == nil {
					return nil
				}

				state := KubernetesClusterConnectorResourceModel{
					Name:                id.LinkerName,
					KubernetesClusterId: id.ResourceUri,
					TargetResourceId:    flattenTargetService(props.TargetService),
					AuthInfo:            flattenServiceConnectorAuthInfo(props.AuthInfo, pwd),
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

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r KubernetesClusterConnectorResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceConnector.LinksClient
			id, err := links.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.LinkerDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesClusterConnectorResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ServiceConnector.LinksClient
			id, err := links.ParseScopedLinkerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state KubernetesClusterConnectorResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			d := metadata.ResourceData
			linkerProps := links.LinkerProperties{}

			if d.HasChange("client_type") {
				clientType := links.ClientType(state.ClientType)
				linkerProps.ClientType = &clientType
			}

			if d.HasChange("vnet_solution") {
				vnetSolutionType := links.VNetSolutionType(state.VnetSolution)
				vnetSolution := links.VNetSolution{
					Type: &vnetSolutionType,
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

func (r KubernetesClusterConnectorResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return servicelinker.ValidateScopedLinkerID
}
