package managedservices

import (
	"fmt"
	"time"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmRegistrationDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmRegistrationDefinitionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"registration_definition_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Computed: 	  true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew: 	  true,
				Computed: 	  true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed: 	  true,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed: 	  true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"managed_by_tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				Computed: 	  true,
				ValidateFunc: validation.IsUUID,
			},

			"managed_by_tenant_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed: 	  true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"authorizations": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"principal_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"role_definition_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmRegistrationDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.RegistrationDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureRegistrationDefinitionId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.scope, id.registrationDefinitionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Registration Definition '%s' was not found (Scope '%s')", id.registrationDefinitionId, id.scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Registration Definition %q (Scope %q): %+v", id.registrationDefinitionId, id.scope, err)
	}

	d.Set("registration_definition_id", resp.Name)
	d.Set("scope", id.scope)

	if props := resp.Properties; props != nil {
		if err := d.Set("authorization", flattenManagedServicesDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("name", props.RegistrationDefinitionName)
		d.Set("managed_by_tenant_id", props.ManagedByTenantID)
		d.Set("managed_by_tenant_name", props.ManagedByTenantName)
	}

	return nil
}