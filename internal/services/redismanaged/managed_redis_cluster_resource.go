// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redismanaged

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redismanaged/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedRedisClusterResource struct{}

var _ sdk.ResourceWithUpdate = ManagedRedisClusterResource{}

type ManagedRedisClusterResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	SkuName           string            `tfschema:"sku_name"`
	Zones             []string          `tfschema:"zones"`
	MinimumTlsVersion string            `tfschema:"minimum_tls_version"`
	Hostname          string            `tfschema:"hostname"`
	Tags              map[string]string `tfschema:"tags"`
}

func (r ManagedRedisClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedRedisName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ManagedRedisClusterSkuName,
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(redisenterprise.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(redisenterprise.TlsVersionOnePointTwo),
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r ManagedRedisClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ManagedRedisClusterResource) ModelObject() interface{} {
	return &ManagedRedisClusterResourceModel{}
}

func (r ManagedRedisClusterResource) ResourceType() string {
	return "azurerm_managed_redis_cluster"
}

func (r ManagedRedisClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return redisenterprise.ValidateRedisEnterpriseID
}

func (r ManagedRedisClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedisManaged.Client
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ManagedRedisClusterResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id := redisenterprise.NewRedisEnterpriseID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			sku := expandManagedRedisClusterSku(model.SkuName)

			// If the sku type is flash check to make sure that the sku is supported in that region
			if strings.Contains(string(sku.Name), "Flash") {
				if err := validate.ManagedRedisClusterLocationFlashSkuSupport(model.Location); err != nil {
					return fmt.Errorf("%s: %s", id, err)
				}
			}

			tlsVersion := redisenterprise.TlsVersion(model.MinimumTlsVersion)
			parameters := redisenterprise.Cluster{
				Location: model.Location,
				Sku:      sku,
				Properties: &redisenterprise.ClusterProperties{
					MinimumTlsVersion: &tlsVersion,
				},
				Tags: pointer.To(model.Tags),
			}

			if len(model.Zones) > 0 {
				// Zones are currently not supported in these regions
				if err := validate.ManagedRedisClusterLocationZoneSupport(model.Location); err != nil {
					return fmt.Errorf("%s: %s", id, err)
				}
				parameters.Zones = pointer.To(model.Zones)
			}

			if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			log.Printf("[DEBUG] Waiting for %s to become available..", id)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Creating", "Updating", "Enabling", "Deleting", "Disabling"},
				Target:     []string{"Running"},
				Refresh:    managedRedisClusterStateRefreshFunc(ctx, client, id),
				MinTimeout: 15 * time.Second,
				Timeout:    metadata.ResourceData.Timeout(pluginsdk.TimeoutCreate),
			}
			if _, err := stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagedRedisClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedisManaged.Client

			id, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
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

			state := ManagedRedisClusterResourceModel{
				Name:              id.RedisEnterpriseName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if skuName := flattenManagedRedisClusterSku(model.Sku); skuName != nil {
					state.SkuName = *skuName
				}

				state.Zones = pointer.From(model.Zones)

				if props := model.Properties; props != nil {
					state.Hostname = pointer.From(props.HostName)

					tlsVersion := ""
					if props.MinimumTlsVersion != nil {
						tlsVersion = string(*props.MinimumTlsVersion)
					}
					state.MinimumTlsVersion = tlsVersion
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagedRedisClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedisManaged.Client

			id, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ManagedRedisClusterResourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil || existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model or properties was nil", *id)
			}

			parameters := redisenterprise.ClusterUpdate{
				Properties: existing.Model.Properties,
			}

			if metadata.ResourceData.HasChange("minimum_tls_version") {
				parameters.Properties.MinimumTlsVersion = pointer.To(redisenterprise.TlsVersion(state.MinimumTlsVersion))
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(state.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			log.Printf("[DEBUG] Waiting for %s to become available", *id)
			stateConf := &pluginsdk.StateChangeConf{
				Pending:    []string{"Creating", "Updating", "Enabling", "Deleting", "Disabling"},
				Target:     []string{"Running"},
				Refresh:    managedRedisClusterStateRefreshFunc(ctx, client, *id),
				MinTimeout: 15 * time.Second,
				Timeout:    30 * time.Minute,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagedRedisClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RedisManaged.Client

			id, err := redisenterprise.ParseRedisEnterpriseID(metadata.ResourceData.Id())
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

func expandManagedRedisClusterSku(v string) redisenterprise.Sku {
	redisSku, _ := parse.RedisEnterpriseCacheSkuName(v)
	capacity, _ := strconv.ParseInt(redisSku.Capacity, 10, 32)

	return redisenterprise.Sku{
		Name:     redisenterprise.SkuName(redisSku.Name),
		Capacity: utils.Int64(capacity),
	}
}

func flattenManagedRedisClusterSku(input redisenterprise.Sku) *string {
	var name redisenterprise.SkuName
	var capacity int64

	if input.Name != "" {
		name = input.Name
	}

	if input.Capacity != nil {
		capacity = *input.Capacity
	}

	skuName := fmt.Sprintf("%s-%d", name, capacity)

	return &skuName
}

func managedRedisClusterStateRefreshFunc(ctx context.Context, client *redisenterprise.RedisEnterpriseClient, id redisenterprise.RedisEnterpriseId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.ResourceState == nil {
			return nil, "", fmt.Errorf("retrieving %s: model/resourceState was nil", id)
		}

		return res, string(*res.Model.Properties.ResourceState), nil
	}
}
