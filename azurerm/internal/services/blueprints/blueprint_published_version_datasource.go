package blueprints

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	mgValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
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
			"blueprint_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"scope_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					azure.ValidateResourceID,
					mgValidate.ManagementGroupID,
				),
			},

			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
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

			"time_created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
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

	scope := d.Get("scope_id").(string)
	blueprintName := d.Get("blueprint_name").(string)
	versionID := d.Get("version").(string)

	resp, err := client.Get(ctx, scope, blueprintName, versionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Published Blueprint Version %q not found: %+v", versionID, err)
		}

		return fmt.Errorf("Read failed for Published Blueprint (%q) Version (%q): %+v", blueprintName, versionID, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Failed to retrieve ID for Blueprint %q Version %q", blueprintName, versionID)
	} else {
		d.SetId(*resp.ID)
	}

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
