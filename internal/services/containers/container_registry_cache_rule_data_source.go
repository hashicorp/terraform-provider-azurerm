// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/registries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-07-01/cacherules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerRegistryCacheRuleDataSource struct{}

func (ContainerRegistryCacheRuleDataSource) ResourceType() string {
	return "azurerm_container_registry_cache_rule"
}

func (ContainerRegistryCacheRuleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.CacheRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var state ContainerRegistryCacheRuleDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			registryId, err := registries.ParseRegistryID(metadata.ResourceData.Get("container_registry_id").(string))
			if err != nil {
				return err
			}

			id := cacherules.NewCacheRuleID(subscriptionId, registryId.ResourceGroupName, registryId.RegistryName, metadata.ResourceData.Get("name").(string))

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)
			state.Name = id.CacheRuleName
			state.ContainerRegistryId = registryId.ID()

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.SourceRepo = pointer.From(props.SourceRepository)
					state.TargetRepo = pointer.From(props.TargetRepository)
					state.CredentialSetId = pointer.From(props.CredentialSetResourceId)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

type ContainerRegistryCacheRuleDataSourceModel struct {
	Name                string `tfschema:"name"`
	ContainerRegistryId string `tfschema:"container_registry_id"`
	CredentialSetId     string `tfschema:"credential_set_id"`
	SourceRepo          string `tfschema:"source_repo"`
	TargetRepo          string `tfschema:"target_repo"`
}

func (ContainerRegistryCacheRuleDataSource) ModelObject() interface{} {
	return &ContainerRegistryCacheRuleDataSourceModel{}
}

func (ContainerRegistryCacheRuleDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ContainerRegistryCacheRuleName,
		},
		"container_registry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: registries.ValidateRegistryID,
		},
	}
}

func (ContainerRegistryCacheRuleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"credential_set_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"source_repo": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"target_repo": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
