package managedservices

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/managedservices/mgmt/2019-06-01/managedservices"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLighthouseDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLighthouseDefinitionCreateUpdate,
		Read:   resourceArmLighthouseDefinitionRead,
		Update: resourceArmLighthouseDefinitionCreateUpdate,
		Delete: resourceArmLighthouseDefinitionDelete,
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
			"registration_definition_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"registration_definition_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"managed_by_tenant_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
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
		},
	}
}

func resourceArmLighthouseDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.LighthouseDefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	lighthouseDefinitionID := d.Get("registration_definition_id").(string)
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

	scope := buildScopeForLighthouseDefinition(subscriptionID)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, scope, lighthouseDefinitionID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Lighthouse Definition %q (Scope %q): %+v", lighthouseDefinitionID, scope, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_lighthouse_definition", *existing.ID)
		}
	}

	parameters := managedservices.RegistrationDefinition{
		Properties: &managedservices.RegistrationDefinitionProperties{
			Description:                utils.String(d.Get("description").(string)),
			Authorizations:             expandManagedServicesDefinitionAuthorization(d.Get("authorization").(*schema.Set).List()),
			RegistrationDefinitionName: utils.String(d.Get("registration_definition_name").(string)),
			ManagedByTenantID:          utils.String(d.Get("managed_by_tenant_id").(string)),
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

	return resourceArmLighthouseDefinitionRead(d, meta)
}

func resourceArmLighthouseDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.LighthouseDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureLighthouseDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.scope, id.lighthouseDefinitionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Lighthouse Definition '%s' was not found (Scope '%s')", id.lighthouseDefinitionID, id.scope)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Lighthouse Definition %q (Scope %q): %+v", id.lighthouseDefinitionID, id.scope, err)
	}

	d.Set("registration_definition_id", resp.Name)
	d.Set("scope", id.scope)

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

func resourceArmLighthouseDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ManagedServices.LighthouseDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parseAzureLighthouseDefinitionID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.lighthouseDefinitionID, id.scope)
	if err != nil {
		return fmt.Errorf("Error deleting Lighthouse Definition %q at Scope %q: %+v", id.lighthouseDefinitionID, id.scope, err)
	}

	return nil
}

type lighthouseDefinitionID struct {
	scope                  string
	lighthouseDefinitionID string
}

func parseAzureLighthouseDefinitionID(id string) (*lighthouseDefinitionID, error) {
	segments := strings.Split(id, "/providers/Microsoft.ManagedServices/registrationDefinitions/")

	if len(segments) != 2 {
		return nil, fmt.Errorf("Expected ID to be in the format `{scope}/providers/Microsoft.ManagedServices/registrationDefinitions/{name} - got %d segments", len(segments))
	}

	azurelighthouseDefinitionID := lighthouseDefinitionID{
		scope:                  segments[0],
		lighthouseDefinitionID: segments[1],
	}

	return &azurelighthouseDefinitionID, nil
}

func buildScopeForLighthouseDefinition(subscriptionID string) string {
	return fmt.Sprintf("/subscriptions/%s", subscriptionID)
}

func flattenManagedServicesDefinitionAuthorization(input *[]managedservices.Authorization) []interface{} {
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

func expandManagedServicesDefinitionAuthorization(input []interface{}) *[]managedservices.Authorization {
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
