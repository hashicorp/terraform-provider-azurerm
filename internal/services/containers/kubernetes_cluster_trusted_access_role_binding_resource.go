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
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-02-01/trustedaccess"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = KubernetesClusterTrustedAccessRoleBindingResource{}
	_ sdk.ResourceWithUpdate = KubernetesClusterTrustedAccessRoleBindingResource{}
)

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
			client := metadata.Client.Containers.TrustedAccessClient

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
			r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(config, &payload)

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
			client := metadata.Client.Containers.TrustedAccessClient
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
				r.mapTrustedAccessRoleBindingToKubernetesClusterTrustedAccessRoleBindingResourceSchema(*model, &schema)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.TrustedAccessClient

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
			client := metadata.Client.Containers.TrustedAccessClient

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

			r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(config, &payload)

			if _, err := client.RoleBindingsCreateOrUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBinding(input KubernetesClusterTrustedAccessRoleBindingResourceSchema, output *trustedaccess.TrustedAccessRoleBinding) {
	r.mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBindingProperties(input, &output.Properties)
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapTrustedAccessRoleBindingToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input trustedaccess.TrustedAccessRoleBinding, output *KubernetesClusterTrustedAccessRoleBindingResourceSchema) {
	r.mapTrustedAccessRoleBindingPropertiesToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input.Properties, output)
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapKubernetesClusterTrustedAccessRoleBindingResourceSchemaToTrustedAccessRoleBindingProperties(input KubernetesClusterTrustedAccessRoleBindingResourceSchema, output *trustedaccess.TrustedAccessRoleBindingProperties) {
	roles := append([]string{}, input.Roles...)
	output.Roles = roles

	output.SourceResourceId = input.SourceResourceId
}

func (r KubernetesClusterTrustedAccessRoleBindingResource) mapTrustedAccessRoleBindingPropertiesToKubernetesClusterTrustedAccessRoleBindingResourceSchema(input trustedaccess.TrustedAccessRoleBindingProperties, output *KubernetesClusterTrustedAccessRoleBindingResourceSchema) {
	roles := append([]string{}, input.Roles...)
	output.Roles = roles

	output.SourceResourceId = input.SourceResourceId
}
