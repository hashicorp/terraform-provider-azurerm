package eventgrid

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceEventGridSystemTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceEventGridSystemTopicRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringIsNotEmpty,
					validation.StringMatch(
						regexp.MustCompile("^[-a-zA-Z0-9]{3,128}$"),
						"EventGrid Topics name must be 3 - 128 characters long, contain only letters, numbers and hyphens.",
					),
				),
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"identity": commonschema.SystemOrUserAssignedIdentityComputed(),

			"location": commonschema.LocationComputed(),

			"source_arm_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"topic_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"metric_arm_resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceEventGridSystemTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.SystemTopicsClient
	subscriptionId := meta.(*clients.Client).EventGrid.DomainsClient.SubscriptionID
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSystemTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.SystemTopicProperties; props != nil {
		d.Set("metric_arm_resource_id", props.MetricResourceID)
		d.Set("source_arm_resource_id", props.Source)
		d.Set("topic_type", props.TopicType)
	}

	flattenedIdentity, err := flattenIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", flattenedIdentity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
