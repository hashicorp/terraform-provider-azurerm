// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package redis_linked_server

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-11-01/linkedserver"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisLinkedServerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	linkedRedisCacheId := d.Get("linked_redis_cache_id").(string)
	linkedRedisCacheLocation := d.Get("linked_redis_cache_location").(string)
	serverRole := linkedserver.ReplicationRole(d.Get("server_role").(string))

	// The name needs to match the linked_redis_cache_id
	cacheId, err := linkedserver.ParseRediID(linkedRedisCacheId)
	if err != nil {
		return err
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := linkedserver.NewLinkedServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("target_redis_cache_name").(string), cacheId.RedisName)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_redis_linked_server", id.ID())
		}
	}

	payload := linkedserver.RedisLinkedServerCreateParameters{
		Properties: linkedserver.RedisLinkedServerCreateProperties{
			LinkedRedisCacheId:       linkedRedisCacheId,
			LinkedRedisCacheLocation: location.Normalize(linkedRedisCacheLocation),
			ServerRole:               serverRole,
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	// TODO: is this still required now that this is using polling methods from `go-azure-sdk`?
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

	return resourceRedisLinkedServerRead(d, meta)
}
