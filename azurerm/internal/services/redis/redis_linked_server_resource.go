package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/redis/mgmt/2020-06-01/redis"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/redis/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceRedisLinkedServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisLinkedServerCreate,
		Read:   resourceRedisLinkedServerRead,
		Delete: resourceRedisLinkedServerDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.LinkedServerID(id)
			return err
		}),

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
				ValidateFunc: validate.CacheID,
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
					// TODO: make this case-sensitive in 3.0
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

func resourceRedisLinkedServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for Redis Linked Server creation.")

	linkedRedisCacheId := d.Get("linked_redis_cache_id").(string)
	linkedRedisCacheLocation := d.Get("linked_redis_cache_location").(string)
	serverRole := redis.ReplicationRole(d.Get("server_role").(string))

	// The name needs to match the linked_redis_cache_id
	cacheId, err := parse.CacheID(linkedRedisCacheId)
	if err != nil {
		return err
	}

	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resourceId := parse.NewLinkedServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("target_redis_cache_name").(string), cacheId.RediName)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.RediName, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Linked Server %q (Redis Cache %q / Resource Group %q): %+v", resourceId.Name, resourceId.RediName, resourceId.ResourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_redis_linked_server", resourceId.ID())
		}
	}

	parameters := redis.LinkedServerCreateParameters{
		LinkedServerCreateProperties: &redis.LinkedServerCreateProperties{
			LinkedRedisCacheID:       utils.String(linkedRedisCacheId),
			LinkedRedisCacheLocation: utils.String(linkedRedisCacheLocation),
			ServerRole:               serverRole,
		},
	}

	future, err := client.Create(ctx, resourceId.ResourceGroup, resourceId.RediName, resourceId.Name, parameters)
	if err != nil {
		return fmt.Errorf("waiting for creation of Linked Server %q (Redis Cache %q / Resource Group %q): %+v", resourceId.Name, resourceId.RediName, resourceId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of Linked Server %q (Redis Cache %q / Resource Group %q): %+v", resourceId.Name, resourceId.RediName, resourceId.ResourceGroup, err)
	}

	log.Printf("[DEBUG] Waiting for Linked Server %q (Redis Cache %q / Resource Group %q) to become available", resourceId.Name, resourceId.RediName, resourceId.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Linking", "Updating", "Creating", "Syncing"},
		Target:     []string{"Succeeded"},
		Refresh:    redisLinkedServerStateRefreshFunc(ctx, client, resourceId),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutCreate),
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Linked Server %q (Redis Cache %q / Resource Group %q) to become available: %+v", resourceId.Name, resourceId.RediName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceRedisLinkedServerRead(d, meta)
}

func resourceRedisLinkedServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.RediName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Linked Server %q (Redis Cache %q / Resource Group %q) was not found - removing from state!", id.Name, id.RediName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Linked Server %q (Redis Cache %q / Resource Group %q): %+v", id.Name, id.RediName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("target_redis_cache_name", id.RediName)
	d.Set("resource_group_name", id.ResourceGroup)
	if props := resp.LinkedServerProperties; props != nil {
		linkedRedisCacheId := ""
		if props.LinkedRedisCacheID != nil {
			cacheId, err := parse.CacheID(*props.LinkedRedisCacheID)
			if err != nil {
				return err
			}

			linkedRedisCacheId = cacheId.ID()
		}
		d.Set("linked_redis_cache_id", linkedRedisCacheId)

		d.Set("linked_redis_cache_location", props.LinkedRedisCacheLocation)
		d.Set("server_role", string(props.ServerRole))
	}

	return nil
}

func resourceRedisLinkedServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.LinkedServerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LinkedServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.RediName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Linked Server %q (Redis Cache %q / Resource Group %q): %+v", id.Name, id.RediName, id.ResourceGroup, err)
		}
	}

	// No LinkedServerDeleteFuture
	// https://github.com/Azure/azure-sdk-for-go/issues/12159
	log.Printf("[DEBUG] Waiting for Linked Server %q (Redis Cache %q / Resource Group %q) to be eventually deleted", id.Name, id.RediName, id.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   redisLinkedServerDeleteStateRefreshFunc(ctx, client, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
		Timeout:                   d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Linked Server %q (Redis Cache %q / Resource Group %q) to be deleted: %+v", id.Name, id.RediName, id.ResourceGroup, err)
	}

	return nil
}

func redisLinkedServerStateRefreshFunc(ctx context.Context, client *redis.LinkedServerClient, id parse.LinkedServerId) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.RediName, id.Name)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving status of Linked Server %q (Redis Cache %q / Resource Group %q): %+v", id.Name, id.RediName, id.ResourceGroup, err)
		}

		return res, *res.LinkedServerProperties.ProvisioningState, nil
	}
}

func redisLinkedServerDeleteStateRefreshFunc(ctx context.Context, client *redis.LinkedServerClient, id parse.LinkedServerId) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.RediName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("retrieving status of Linked Server %q (Redis Cache %q / Resource Group %q): %+v", id.Name, id.RediName, id.ResourceGroup, err)
		}

		return res, "Exists", nil
	}
}
