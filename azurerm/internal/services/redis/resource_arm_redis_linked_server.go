package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2018-03-01/redis"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRedisLinkedServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRedisLinkedServerCreate,
		Read:   resourceArmRedisLinkedServerRead,
		Delete: resourceArmRedisLinkedServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"target_redis_cache_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"linked_redis_cache_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"linked_redis_cache_location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(redis.ReplicationRolePrimary),
					string(redis.ReplicationRoleSecondary),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmRedisLinkedServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Redis Linked Server creation.")

	redisCacheName := d.Get("target_redis_cache_name").(string)
	linkedRedisCacheId := d.Get("linked_redis_cache_id").(string)
	linkedRedisCacheLocation := d.Get("linked_redis_cache_location").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverRole := redis.ReplicationRole(d.Get("server_role").(string))

	// The name needs to match the linked_redis_cache_id
	id, err := azure.ParseAzureResourceID(linkedRedisCacheId)
	if err != nil {
		return err
	}
	name := id.Path["Redis"]

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, redisCacheName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Redis Linked Server %q (cache %q / resource group %q) ID", name, redisCacheName, resourceGroup)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_redis_linked_server", *existing.ID)
		}
	}

	parameters := redis.LinkedServerCreateParameters{
		LinkedServerCreateProperties: &redis.LinkedServerCreateProperties{
			LinkedRedisCacheID:       utils.String(linkedRedisCacheId),
			LinkedRedisCacheLocation: utils.String(linkedRedisCacheLocation),
			ServerRole:               serverRole,
		},
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		future, err := client.Create(ctx, resourceGroup, redisCacheName, name, parameters)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error issuing for the create of Redis Linked Server %s (resource group %s): %v", name, resourceGroup, err))
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error waiting for the create of Redis Linked Server %s (resource group %s): %v", name, resourceGroup, err))
		}

		read, err := client.Get(ctx, resourceGroup, redisCacheName, name)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Expected instance to be created but was in non existent state, retrying"))
		}
		if read.ID == nil {
			return resource.NonRetryableError(fmt.Errorf("Cannot read Redis Linked Server %q (cache %q / resource group %q) ID", name, redisCacheName, resourceGroup))
		}

		log.Printf("[DEBUG] Waiting for Redis Linked Server (%s) to become available", d.Get("name"))
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"Linking", "Updating", "Creating", "Syncing"},
			Target:     []string{"Succeeded"},
			Refresh:    redisLinkedServerStateRefreshFunc(ctx, client, resourceGroup, redisCacheName, name),
			MinTimeout: 15 * time.Second,
			Timeout:    d.Timeout(schema.TimeoutCreate),
		}

		if _, err = stateConf.WaitForState(); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error waiting for Redis Linked Server (%s) to become available: %s", d.Get("name"), err))
		}

		d.SetId(*read.ID)

		return resource.NonRetryableError(resourceArmRedisLinkedServerRead(d, meta))
	})
}

func resourceArmRedisLinkedServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	redisCacheName := id.Path["Redis"]
	name := id.Path["linkedServers"]

	resp, err := client.Get(ctx, resourceGroup, redisCacheName, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Redis Linked Server %q was not found in Cache %q / Resource Group %q - removing from state", name, redisCacheName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure Redis Linked Server %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("target_redis_cache_name", redisCacheName)
	d.Set("resource_group_name", resourceGroup)
	if props := resp.LinkedServerProperties; props != nil {
		d.Set("linked_redis_cache_id", props.LinkedRedisCacheID)
		d.Set("linked_redis_cache_location", props.LinkedRedisCacheLocation)
		d.Set("server_role", string(props.ServerRole))
	}

	return nil
}

func resourceArmRedisLinkedServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	redisCacheName := id.Path["Redis"]
	name := id.Path["linkedServers"]

	resp, err := client.Delete(ctx, resourceGroup, redisCacheName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing AzureRM delete request of Redis Linked Server %q (cache %q / resource group %q): %+v", name, redisCacheName, resourceGroup, err)
		}
	}

	// No LinkedServerDeleteFuture
	// https://github.com/Azure/azure-sdk-for-go/issues/12159
	log.Printf("[DEBUG] Waiting for Redis Linked Server %q (cache %q / Resource Group %q) to be eventually deleted", name, redisCacheName, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   redisLinkedServerDeleteStateRefreshFunc(ctx, client, resourceGroup, redisCacheName, name),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("failed to wait for Redis Linked Server %q (cache %q / resource group %q) to be deleted: %+v", name, redisCacheName, resourceGroup, err)
	}

	return nil
}

func redisLinkedServerStateRefreshFunc(ctx context.Context, client *redis.LinkedServerClient, resourceGroupName string, redisCacheName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, redisCacheName, name)
		if err != nil {
			return nil, "", fmt.Errorf("Error issuing read request in redisStateRefreshFunc to Azure ARM for Redis Linked Server Instance '%s' (RG: '%s'): %s", name, resourceGroupName, err)
		}

		return res, *res.LinkedServerProperties.ProvisioningState, nil
	}
}

func redisLinkedServerDeleteStateRefreshFunc(ctx context.Context, client *redis.LinkedServerClient, resourceGroupName string, redisCacheName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, redisCacheName, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("failed to poll to check if the Linked Server has been deleted: %+v", err)
		}

		return res, "Exists", nil
	}
}
