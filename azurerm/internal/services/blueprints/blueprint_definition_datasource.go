package blueprints

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/blueprints/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmBlueprintDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmBlueprintDefinitionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.BlueprintName,
			},

			"scope_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"subscriptions",
					"managementGroup",
				}, false),
			},

			"scope_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceArmBlueprintDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprints.BlueprintsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	var scope string
	scopeName := d.Get("scope_name").(string)
	scopeType := d.Get("scope_type").(string)
	switch scopeType {
	case "subscriptions":
		scope = fmt.Sprintf("subscriptions/%s", scopeName)
	case "managementGroup":
		scope = fmt.Sprintf("providers/Microsoft.Management/managementGroups/%s", scopeName)
	}

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Blueprint Definition %q not found in Scope (%q): %+v", name, scope, err)
		}

		return fmt.Errorf("Read failed for Blueprint Definition (%q) in Sccope (%q): %+v", name, scope, err)
	}

	if resp.ID != nil {
		d.SetId(*resp.ID)
	} else {
		return fmt.Errorf("Failed to retrieve ID for Blueprint %q", name)
	}

	if resp.Description != nil {
		d.Set("description", resp.Description)
	}

	if resp.DisplayName != nil {
		d.Set("display_name", resp.DisplayName)
	}

	d.Set("last_modified", resp.Status.LastModified.String())

	d.Set("time_created", resp.Status.TimeCreated.String())

	d.Set("target_scope", resp.TargetScope)

	if resp.Versions != nil {
		d.Set("versions", resp.Versions)
	}

	return nil
}
