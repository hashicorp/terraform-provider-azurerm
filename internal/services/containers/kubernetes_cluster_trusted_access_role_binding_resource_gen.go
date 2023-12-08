package containers

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2023-03-02-preview/trustedaccess"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = KubernetesClusterTrustedAccessRoleBindingResource{}
var _ sdk.ResourceWithUpdate = KubernetesClusterTrustedAccessRoleBindingResource{}

type KubernetesClusterTrustedAccessRoleBindingResource struct{}

func (r KubernetesClusterTrustedAccessRoleBindingResource) ModelObject() interface{} {
	return &KubernetesClusterTrustedAccessRoleBindingResourceSchema{}
}

type KubernetesClusterTrustedAccessRoleBindingResourceSchema struct {
	KubernetesClusterId string   `tfschema:"kubernetes_cluster_id"`
	Name                string   `tfschema:"name"`
	Roles               []string `tfschema:"roles"`
	SourceResourceId    string   `tfschema:"source_resource_id"`
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return trustedaccess.ValidateTrustedAccessRoleBindingID
}
func (r KubernetesClusterTrustedAccessRoleBindingResource) ResourceType() string {
	return "azurerm_kubernetes_cluster_trusted_access_role_binding"
}
func (r KubernetesClusterTrustedAccessRoleBindingResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"kubernetes_cluster_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"roles": {
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Required: true,
			Type:     pluginsdk.TypeList,
		},
		"source_resource_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
func (r KubernetesClusterTrustedAccessRoleBindingResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (r KubernetesClusterTrustedAccessRoleBindingResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20230302Preview.TrustedAccess

			var config KubernetesClusterTrustedAccessRoleBindingResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId

			kubernetesClusterId, err := commonids.ParseKubernetesClusterID(config.KubernetesClusterId)
			if err != nil {
				return err
			}

			id := trustedaccess.NewTrustedAccessRoleBindingID(subscriptionId, kubernetesClusterId.ResourceGroupName, kubernetesClusterId.ManagedClusterName, config.Name)

			existing, err := client.RoleBindingsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var payload trustedaccess.TrustedAccessRoleBinding
			if err := r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.RoleBindingsCreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
func (r KubernetesClusterTrustedAccessRoleBindingResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20230302Preview.TrustedAccess
			schema := KubernetesClusterTrustedAccessRoleBindingResourceSchema{}

			id, err := trustedaccess.ParseTrustedAccessRoleBindingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			kubernetesClusterId := commonids.NewKubernetesClusterID(id.SubscriptionId, id.ResourceGroupName, id.ManagedClusterName)

			resp, err := client.RoleBindingsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.KubernetesClusterId = kubernetesClusterId.ID()
				schema.Name = id.TrustedAccessRoleBindingName
				if err := r.mapTrustedAccessRoleBindingToKubernetesClusterTrustedAccessRoleBindingResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
func (r KubernetesClusterTrustedAccessRoleBindingResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20230302Preview.TrustedAccess

			id, err := trustedaccess.ParseTrustedAccessRoleBindingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.RoleBindingsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
func (r KubernetesClusterTrustedAccessRoleBindingResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ContainerService.V20230302Preview.TrustedAccess

			id, err := trustedaccess.ParseTrustedAccessRoleBindingID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesClusterTrustedAccessRoleBindingResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.RoleBindingsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			if err := r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(config, &payload); err != nil {
				return fmt.Errorf("mapping schema model to sdk model: %+v", err)
			}

			if _, err := client.RoleBindingsCreateOrUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBindingProperties(input KubernetesClusterTrustedAccessRoleBindingResourceSchema, output *trustedaccess.TrustedAccessRoleBindingProperties) error {

	roles := make([]string, 0)
	for _, v := range input.Roles {
		roles = append(roles, v)
	}
	output.Roles = roles

	output.SourceResourceId = input.SourceResourceId
	return nil
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapTrustedAccessRoleBindingPropertiesToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input trustedaccess.TrustedAccessRoleBindingProperties, output *KubernetesClusterTrustedAccessRoleBindingResourceSchema) error {

	roles := make([]string, 0)
	for _, v := range input.Roles {
		roles = append(roles, v)
	}
	output.Roles = roles

	output.SourceResourceId = input.SourceResourceId
	return nil
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(input KubernetesClusterTrustedAccessRoleBindingResourceSchema, output *trustedaccess.TrustedAccessRoleBinding) error {

	if err := r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBindingProperties(input, &output.Properties); err != nil {
		return fmt.Errorf("mapping Schema to SDK Field %q / Model %q: %+v", "TrustedAccessRoleBindingProperties", "Properties", err)
	}

	return nil
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapTrustedAccessRoleBindingToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input trustedaccess.TrustedAccessRoleBinding, output *KubernetesClusterTrustedAccessRoleBindingResourceSchema) error {

	if err := r.mapTrustedAccessRoleBindingPropertiesToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input.Properties, output); err != nil {
		return fmt.Errorf("mapping SDK Field %q / Model %q to Schema: %+v", "TrustedAccessRoleBindingProperties", "Properties", err)
	}

	return nil
}
