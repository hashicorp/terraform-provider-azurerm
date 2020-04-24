package blueprints

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBlueprintPublishedVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBlueprintPublishedVersionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.UUID,
			},

			"management_group": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"blueprint_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			// Computed

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmBlueprintPublishedVersionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprints.PublishedBlueprintsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subscriptionRaw := d.Get("subscription_id")
	var subscription, managementGroup string
	if subscriptionRaw != nil {
		subscription = subscriptionRaw.(string)
	}
	managementGroupRaw := d.Get("managementGroup")
	if managementGroupRaw != nil {
		managementGroup = managementGroupRaw.(string)
	}

	var scope string

	if subscription == "" && managementGroup == "" {
		return fmt.Errorf("One of subscription or management group must be specified")
	}

	if subscription != "" {
		scope = fmt.Sprintf("subscriptions/%s", subscription)
	} else {
		scope = fmt.Sprintf("providers/Microsoft.Management/managementGroups/%s", managementGroup)
	}

	blueprintName := d.Get("blueprint_name").(string)
	versionID := d.Get("version").(string)

	resp, err := client.Get(ctx, scope, blueprintName, versionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Published Blueprint Version %q not found: %+v", versionID, err)
		}

		return fmt.Errorf("Read failed for Published Blueprint (%q) Version (%q): %+v", blueprintName, versionID, err)
	}

	d.SetId(*resp.ID)
	if resp.Type != nil {
		d.Set("type", resp.Type)
	}

	if resp.Status != nil {
		if resp.Status.TimeCreated != nil {
			d.Set("time_created", resp.Status.TimeCreated.String())
		}

		if resp.Status.LastModified != nil {
			d.Set("last_modified", resp.Status.LastModified.String())
		}
	}

	d.Set("target_scope", resp.TargetScope)

	if resp.DisplayName != nil {
		d.Set("display_name", resp.DisplayName)
	}

	if resp.Description != nil {
		d.Set("description", resp.Description)
	}

	return nil
}
