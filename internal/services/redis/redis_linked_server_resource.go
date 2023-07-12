// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2023-04-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisLinkedServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRedisLinkedServerCreate,
		Read:   resourceRedisLinkedServerRead,
		Delete: resourceRedisLinkedServerDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := redis.ParseLinkedServerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.LinkedServerV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"target_redis_cache_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"linked_redis_cache_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: redis.ValidateRediID,
			},

			"linked_redis_cache_location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"server_role": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.ReplicationRolePrimary),
					string(redis.ReplicationRoleSecondary),
				}, false),
			},

			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRedisLinkedServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	linkedRedisCacheId := d.Get("linked_redis_cache_id").(string)
	linkedRedisCacheLocation := d.Get("linked_redis_cache_location").(string)
	serverRole := redis.ReplicationRole(d.Get("server_role").(string))

	// The name needs to match the linked_redis_cache_id
	cacheId, err := redis.ParseRediID(linkedRedisCacheId)
	if err != nil {
		return err
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := redis.NewLinkedServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("target_redis_cache_name").(string), cacheId.RedisName)
	if d.IsNewResource() {
		existing, err := client.LinkedServerGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_redis_linked_server", id.ID())
		}
	}

	payload := redis.RedisLinkedServerCreateParameters{
		Properties: redis.RedisLinkedServerCreateProperties{
			LinkedRedisCacheId:       linkedRedisCacheId,
			LinkedRedisCacheLocation: location.Normalize(linkedRedisCacheLocation),
			ServerRole:               serverRole,
		},
	}

	if err := client.LinkedServerCreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	log.Printf("[DEBUG] Waiting for %s to become available", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Linking", "Updating", "Creating", "Syncing"},
		Target:     []string{"Succeeded"},
		Refresh:    redisLinkedServerStateRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
		Timeout:    time.Until(deadline),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceRedisLinkedServerRead(d, meta)
}

func resourceRedisLinkedServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redis.ParseLinkedServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.LinkedServerGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.LinkedServerName)
	d.Set("target_redis_cache_name", id.RedisName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			cacheId, err := redis.ParseRediIDInsensitively(props.LinkedRedisCacheId)
			if err != nil {
				return fmt.Errorf("parsing `linkedRedisCacheId` %q: %+v", props.LinkedRedisCacheId, err)
			}
			d.Set("linked_redis_cache_id", cacheId.ID())

			d.Set("linked_redis_cache_location", location.Normalize(props.LinkedRedisCacheLocation))
			d.Set("server_role", string(props.ServerRole))
		}
	}

	return nil
}

func resourceRedisLinkedServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.Redis
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := redis.ParseLinkedServerID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.LinkedServerDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// No LinkedServerDeleteFuture
	// https://github.com/Azure/azure-sdk-for-go/issues/12159
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	log.Printf("[DEBUG] Waiting for %s to be eventually deleted", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   redisLinkedServerDeleteStateRefreshFunc(ctx, client, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}

func redisLinkedServerStateRefreshFunc(ctx context.Context, client *redis.RedisClient, id redis.LinkedServerId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.LinkedServerGet(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving status of %s: %+v", id, err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil && props.ProvisioningState != nil {
				return resp, *props.ProvisioningState, nil
			}
		}

		return nil, "", fmt.Errorf("retrieving %s: `model` was nil", id)
	}
}

func redisLinkedServerDeleteStateRefreshFunc(ctx context.Context, client *redis.RedisClient, id redis.LinkedServerId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.LinkedServerGet(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("retrieving status of %s: %+v", id, err)
		}

		return res, "Exists", nil
	}
}
