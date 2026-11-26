// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_cache

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redispatchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/redisresources"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisCacheCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Redis.RedisResourcesClient
	patchClient := meta.(*clients.Client).Redis.PatchSchedulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := redisresources.NewRediID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.RedisGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_redis_cache", id.ID())
		}
	}

	patchSchedule := expandRedisPatchSchedule(d)
	redisConfiguration, err := expandRedisConfiguration(d)
	if err != nil {
		return fmt.Errorf("parsing Redis Configuration: %+v", err)
	}

	publicNetworkAccess := redisresources.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = redisresources.PublicNetworkAccessDisabled
	}

	redisIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf(`expanding "identity": %v`, err)
	}

	enableNonSslPort := d.Get("non_ssl_port_enabled")

	parameters := redisresources.RedisCreateParameters{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: redisresources.RedisCreateProperties{
			DisableAccessKeyAuthentication: pointer.To(!d.Get("access_keys_authentication_enabled").(bool)),
			EnableNonSslPort:               pointer.To(enableNonSslPort.(bool)),
			Sku: redisresources.Sku{
				Capacity: int64(d.Get("capacity").(int)),
				Family:   redisresources.SkuFamily(d.Get("family").(string)),
				Name:     redisresources.SkuName(d.Get("sku_name").(string)),
			},
			MinimumTlsVersion:   pointer.ToEnum[redisresources.TlsVersion](d.Get("minimum_tls_version").(string)),
			RedisConfiguration:  redisConfiguration,
			PublicNetworkAccess: pointer.To(publicNetworkAccess),
		},
		Identity: redisIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("shard_count"); ok {
		parameters.Properties.ShardCount = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("replicas_per_master"); ok {
		parameters.Properties.ReplicasPerMaster = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("replicas_per_primary"); ok {
		parameters.Properties.ReplicasPerPrimary = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("redis_version"); ok {
		parameters.Properties.RedisVersion = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("tenant_settings"); ok {
		parameters.Properties.TenantSettings = expandTenantSettings(v.(map[string]any))
	}

	if v, ok := d.GetOk("private_static_ip_address"); ok {
		parameters.Properties.StaticIP = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		parsed, parseErr := commonids.ParseSubnetID(v.(string))
		if parseErr != nil {
			return err
		}

		locks.ByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)
		defer locks.UnlockByName(parsed.VirtualNetworkName, network.VirtualNetworkResourceName)

		locks.ByName(parsed.SubnetName, network.SubnetResourceName)
		defer locks.UnlockByName(parsed.SubnetName, network.SubnetResourceName)

		parameters.Properties.SubnetId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := zones.ExpandUntyped(v.(*schema.Set).List())
		if len(zones) > 0 {
			parameters.Zones = &zones
		}
	}

	if err := client.RedisCreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// TODO: is this still required now that this is using polling methods from `go-azure-sdk`?
	log.Printf("[DEBUG] Waiting for %s to become available", id)
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Scaling", "Updating", "Creating", "ConfiguringAAD"},
		Target:     []string{"Succeeded"},
		Refresh:    redisStateRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
		Timeout:    time.Until(deadline),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	if patchSchedule != nil {
		patchScheduleRedisId := redispatchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
		if _, err = patchClient.PatchSchedulesCreateOrUpdate(ctx, patchScheduleRedisId, *patchSchedule); err != nil {
			return fmt.Errorf("setting Patch Schedule for %s: %+v", id, err)
		}
	}

	return resourceRedisCacheRead(d, meta)
}
