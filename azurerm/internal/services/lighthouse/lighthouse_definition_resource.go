package lighthouse

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/validate"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/lighthouse/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLighthouseDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceLighthouseDefinitionCreateUpdate,
		Read:   resourceLighthouseDefinitionRead,
		Update: resourceLighthouseDefinitionCreateUpdate,
		Delete: resourceLighthouseDefinitionDelete,
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
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"managing_tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"scope": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SubscriptionID,
			},

			"authorization": {
				Type:     schema.TypeSet,
				Required: true,
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"lighthouse_definition_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceLighthouseDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseDefinitionID := d.Get("lighthouse_definition_id").(string)
	if lighthouseDefinitionID == "" {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return fmt.Errorf("Error generating UUID for Lighthouse Definition: %+v", err)
		}

		lighthouseDefinitionID = uuid
	}

	subscriptionID := meta.(*clients.Client).Account.SubscriptionId
	if subscriptionID == "" {
		return fmt.Errorf("Error reading Subscription for Lighthouse Definition %q", lighthouseDefinitionID)
	}

	scope := d.Get("scope").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, lighthouseDefinitionID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Lighthouse Definition %q (Scope %q): %+v", lighthouseDefinitionID, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_lighthouse_definition", *existing.ID)
		}
	}

	parameters := managedservices.RegistrationDefinition{
		Properties: &managedservices.RegistrationDefinitionProperties{
			Description:                utils.String(d.Get("description").(string)),
			Authorizations:             expandLighthouseDefinitionAuthorization(d.Get("authorization").(*schema.Set).List()),
			RegistrationDefinitionName: utils.String(d.Get("name").(string)),
			ManagedByTenantID:          utils.String(d.Get("managing_tenant_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, lighthouseDefinitionID, scope, parameters); err != nil {
		return fmt.Errorf("Error Creating/Updating Lighthouse Definition %q (Scope %q): %+v", lighthouseDefinitionID, scope, err)
	}

	read, err := client.Get(ctx, scope, lighthouseDefinitionID)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Lighthouse Definition %q ID (scope %q) ID", lighthouseDefinitionID, scope)
	}

	d.SetId(*read.ID)

	return resourceLighthouseDefinitionRead(d, meta)
}

func resourceLighthouseDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LighthouseDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Scope, id.LighthouseDefinitionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Lighthouse Definition %q was not found (Scope %q)", id.LighthouseDefinitionID, id.Scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Lighthouse Definition %q (Scope %q): %+v", id.LighthouseDefinitionID, id.Scope, err)
	}

	d.Set("lighthouse_definition_id", resp.Name)
	d.Set("scope", id.Scope)

	if props := resp.Properties; props != nil {
		if err := d.Set("authorization", flattenLighthouseDefinitionAuthorization(props.Authorizations)); err != nil {
			return fmt.Errorf("setting `authorization`: %+v", err)
		}
		d.Set("description", props.Description)
		d.Set("name", props.RegistrationDefinitionName)
		d.Set("managing_tenant_id", props.ManagedByTenantID)
	}

	return nil
}

func resourceLighthouseDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Lighthouse.DefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.LighthouseDefinitionID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.LighthouseDefinitionID, id.Scope)
	if err != nil {
		return fmt.Errorf("Error deleting Lighthouse Definition %q at Scope %q: %+v", id.LighthouseDefinitionID, id.Scope, err)
	}

	return nil
}

func flattenLighthouseDefinitionAuthorization(input *[]managedservices.Authorization) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		principalID := ""
		if item.PrincipalID != nil {
			principalID = *item.PrincipalID
		}

		roleDefinitionID := ""
		if item.RoleDefinitionID != nil {
			roleDefinitionID = *item.RoleDefinitionID
		}

		results = append(results, map[string]interface{}{
			"role_definition_id": roleDefinitionID,
			"principal_id":       principalID,
		})
	}

	return results
}

func expandLighthouseDefinitionAuthorization(input []interface{}) *[]managedservices.Authorization {
	results := make([]managedservices.Authorization, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		result := managedservices.Authorization{
			RoleDefinitionID: utils.String(v["role_definition_id"].(string)),
			PrincipalID:      utils.String(v["principal_id"].(string)),
		}
		results = append(results, result)
	}
	return &results
}
