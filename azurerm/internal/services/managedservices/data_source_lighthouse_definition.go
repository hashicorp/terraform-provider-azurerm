package managedservices

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmLighthouseDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmLighthouseDefinitionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"registration_definition_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"registration_definition_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"managed_by_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"authorization": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_definition_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmLighthouseDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.LighthouseDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	registrationDefinitionID := d.Get("registration_definition_id").(string)
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	if subscriptionID == "" {
		return fmt.Errorf("Error reading Subscription for Registration Definition %q", registrationDefinitionID)
	}

	scope := buildScopeForLighthouseDefinition(subscriptionID)

	resp, err := client.Get(ctx, scope, registrationDefinitionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Registration Definition '%s' was not found (Scope '%s')", registrationDefinitionID, scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Registration Definition %q (Scope %q): %+v", registrationDefinitionID, scope, err)
	}

	d.SetId(*resp.ID)
	d.Set("scope", scope)

	if props := resp.Properties; props != nil {
		if err := d.Set("authorization", flattenManagedServicesDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("registration_definition_name", props.RegistrationDefinitionName)
		d.Set("managed_by_tenant_id", props.ManagedByTenantID)
	}

	return nil
}
