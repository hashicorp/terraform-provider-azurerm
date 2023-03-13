package redis

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2022-06-01/patchschedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2022-06-01/redis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRedisPatchScheduleCache() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRedisCachePatchScheduleCreate,
		Read:   resourceRedisCachePatchScheduleRead,
		Update: resourceRedisCachePatchScheduleUpdate,
		Delete: resourceRedisCachePatchScheduleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := patchschedules.ParseRediID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"redis_cache_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: redis.ValidateRediID,
			},

			"patch_schedule": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day_of_week": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     validation.IsDayOfTheWeek(true),
						},

						"maintenance_window": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "PT5H",
							ValidateFunc: azValidate.ISO8601Duration,
						},

						"start_hour_utc": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 23),
						},
					},
				},
			},
		},
	}
}

func resourceRedisCachePatchScheduleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.PatchSchedules
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := patchschedules.ParseRediID(d.Get("redis_cache_id").(string))
	if err != nil {
		return err
	}
	existing, err := client.Get(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_redis_cache_patch_schedule", id.ID())
	}

	patchSchedule := expandRedisPatchSchedule(d)
	patchScheduleRedisId := patchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
	if _, err = client.CreateOrUpdate(ctx, patchScheduleRedisId, *patchSchedule); err != nil {
		return fmt.Errorf("setting Patch Schedule for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceRedisCachePatchScheduleRead(d, meta)
}

func resourceRedisCachePatchScheduleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.PatchSchedules
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := patchschedules.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	patchSchedule := expandRedisPatchSchedule(d)
	patchSchedulesRedisId := patchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
	_, err = client.CreateOrUpdate(ctx, patchSchedulesRedisId, *patchSchedule)
	if err != nil {
		return fmt.Errorf("setting Patch Schedule for %s: %+v", *id, err)
	}

	return resourceRedisCachePatchScheduleRead(d, meta)
}

func resourceRedisCachePatchScheduleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.PatchSchedules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := patchschedules.ParseRediID(d.Id())
	if err != nil {
		return err
	}
	d.Set("redis_cache_id", id.ID())

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	patchSchedulesRedisId := patchschedules.NewRediID(id.SubscriptionId, id.ResourceGroupName, id.RedisName)
	schedule, err := client.Get(ctx, patchSchedulesRedisId)
	if err != nil {
		return err
	}

	patchSchedules := flattenRedisPatchSchedules(*schedule.Model)
	if err = d.Set("patch_schedule", patchSchedules); err != nil {
		return fmt.Errorf("setting `patch_schedule`: %+v", err)
	}

	return nil
}

func resourceRedisCachePatchScheduleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Redis.PatchSchedules
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := patchschedules.ParseRediID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
