// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redismanaged

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redisenterprise/2025-04-01/redisenterprise"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redismanaged/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redismanaged/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redismanaged/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagedRedisClusterResource struct{}

var _ sdk.ResourceWithUpdate = ManagedRedisClusterResource{}

type ManagedRedisClusterResourceModel struct {
	Name                    string                                     `tfschema:"name"`
	ResourceGroupName       string                                     `tfschema:"resource_group_name"`
	Location                string                                     `tfschema:"location"`
	SkuName                 string                                     `tfschema:"sku_name"`
	CustomerManagedKey      []CustomerManagedKey                       `tfschema:"customer_managed_key"`
	HighAvailabilityEnabled bool                                       `tfschema:"high_availability_enabled"`
	Identity                []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	MinimumTlsVersion       string                                     `tfschema:"minimum_tls_version"`
	Tags                    map[string]string                          `tfschema:"tags"`
	Zones                   []string                                   `tfschema:"zones"`
	Hostname                string                                     `tfschema:"hostname"`
}

type CustomerManagedKey struct {
	EncryptionKeyUrl       string `tfschema:"encryption_key_url"`
	UserAssignedIdentityId string `tfschema:"user_assigned_identity_id"`
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

		"customer_managed_key": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"encryption_key_url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},

					"user_assigned_identity_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateUserAssignedIdentityID,
					},
				},
			},
		},

		"high_availability_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"minimum_tls_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(redisenterprise.TlsVersionOnePointTwo),
			ValidateFunc: validation.StringInSlice([]string{
				string(redisenterprise.TlsVersionOnePointTwo),
			}, false),
		},

		"tags": commonschema.Tags(),

		"zones": commonschema.ZonesMultipleOptionalForceNew(),
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
				return fmt.Errorf("decoding: %+v", err)
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

			sku, err := parse.ManagedRedisCacheSkuName(model.SkuName)
			if err != nil {
				return fmt.Errorf("parsing `sku_name`: %+v", err)
			}

			highAvailability := redisenterprise.HighAvailabilityEnabled
			if !model.HighAvailabilityEnabled {
				highAvailability = redisenterprise.HighAvailabilityDisabled
			}

			parameters := redisenterprise.Cluster{
				Location: location.Normalize(model.Location),
				Sku:      pointer.From(sku),
				Properties: &redisenterprise.ClusterProperties{
					MinimumTlsVersion: pointer.To(redisenterprise.TlsVersion(model.MinimumTlsVersion)),
					HighAvailability:  pointer.To(highAvailability),
				},
				Tags: pointer.To(model.Tags),
			}

			if len(model.CustomerManagedKey) > 0 {
				encryption := expandManagedRedisClusterCustomerManagedKey(model.CustomerManagedKey)
				parameters.Properties.Encryption = encryption
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}
			parameters.Identity = expandedIdentity

			if len(model.Zones) > 0 {
				parameters.Zones = pointer.To(model.Zones)
			}

			if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			pollerType := custompollers.NewClusterStatePoller(client, id)
			poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for cluster %s to become available: %+v", id, err)
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

				if skuName := flattenManagedRedisClusterSku(model.Sku); skuName != nil {
					state.SkuName = *skuName
				}

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = *flattenedIdentity

				state.Tags = pointer.From(model.Tags)
				state.Zones = pointer.From(model.Zones)

				if props := model.Properties; props != nil {
					state.CustomerManagedKey = flattenManagedRedisClusterCustomerManagedKey(props.Encryption)
					state.HighAvailabilityEnabled = strings.EqualFold(string(pointer.From(props.HighAvailability)), string(redisenterprise.HighAvailabilityEnabled))
					state.MinimumTlsVersion = string(pointer.From(props.MinimumTlsVersion))
					state.Hostname = pointer.From(props.HostName)
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
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", *id)
			}

			parameters := redisenterprise.ClusterUpdate{
				Properties: existing.Model.Properties,
			}

			if metadata.ResourceData.HasChange("customer_managed_key") {
				encryption := expandManagedRedisClusterCustomerManagedKey(state.CustomerManagedKey)
				parameters.Properties.Encryption = encryption
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(state.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				parameters.Identity = expandedIdentity
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

			pollerType := custompollers.NewClusterStatePoller(client, *id)
			poller := pollers.NewPoller(pollerType, 15*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for cluster %s to become available: %+v", id, err)
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

func flattenManagedRedisClusterSku(input redisenterprise.Sku) *string {
	var skuName string

	if input.Capacity != nil {
		capacity := *input.Capacity
		skuName = fmt.Sprintf("%s-%d", string(input.Name), capacity)
	} else {
		skuName = string(input.Name)
	}

	return pointer.To(skuName)
}

func expandManagedRedisClusterCustomerManagedKey(input []CustomerManagedKey) *redisenterprise.ClusterPropertiesEncryption {
	if len(input) == 0 {
		return nil
	}

	cmk := input[0]

	return &redisenterprise.ClusterPropertiesEncryption{
		CustomerManagedKeyEncryption: &redisenterprise.ClusterPropertiesEncryptionCustomerManagedKeyEncryption{
			KeyEncryptionKeyURL: pointer.To(cmk.EncryptionKeyUrl),
			KeyEncryptionKeyIdentity: &redisenterprise.ClusterPropertiesEncryptionCustomerManagedKeyEncryptionKeyEncryptionKeyIdentity{
				IdentityType:                   pointer.To(redisenterprise.CmkIdentityTypeUserAssignedIdentity),
				UserAssignedIdentityResourceId: pointer.To(cmk.UserAssignedIdentityId),
			},
		},
	}
}

func flattenManagedRedisClusterCustomerManagedKey(input *redisenterprise.ClusterPropertiesEncryption) []CustomerManagedKey {
	if input == nil || input.CustomerManagedKeyEncryption == nil {
		return make([]CustomerManagedKey, 0)
	}

	cmkEncryption := input.CustomerManagedKeyEncryption

	return []CustomerManagedKey{
		{
			EncryptionKeyUrl:       pointer.From(cmkEncryption.KeyEncryptionKeyURL),
			UserAssignedIdentityId: pointer.From(cmkEncryption.KeyEncryptionKeyIdentity.UserAssignedIdentityResourceId),
		},
	}
}
