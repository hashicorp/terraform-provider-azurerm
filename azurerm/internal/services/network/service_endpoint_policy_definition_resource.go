package network

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	mgValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceEndpointPolicyDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceEndpointPolicyDefinitionCreateUpdate,
		Read:   resourceArmServiceEndpointPolicyDefinitionRead,
		Update: resourceArmServiceEndpointPolicyDefinitionCreateUpdate,
		Delete: resourceArmServiceEndpointPolicyDefinitionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServiceEndpointPolicyDefinitionID(id)
			return err
		}),

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
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"policy_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServiceEndpointPolicyID,
			},

			"service_endpoint_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"service_resources": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.Any(
						azure.ValidateResourceID,
						mgValidate.ManagementGroupID,
					),
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 140),
			},
		},
	}
}

func resourceArmServiceEndpointPolicyDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicyDefinitionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	policyId, err := parse.ServiceEndpointPolicyID(d.Get("policy_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		resp, err := client.Get(ctx, policyId.ResourceGroup, policyId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Service Endpoint Policy Definition %q (Policy %q / Resource Group %q): %+v", name, policyId.Name, policyId.ResourceGroup, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_service_endpoint_policy_definition", *resp.ID)
		}
	}

	param := network.ServiceEndpointPolicyDefinition{
		ServiceEndpointPolicyDefinitionPropertiesFormat: &network.ServiceEndpointPolicyDefinitionPropertiesFormat{
			Description:      utils.String(d.Get("description").(string)),
			Service:          utils.String(d.Get("service_endpoint_name").(string)),
			ServiceResources: utils.ExpandStringSlice(d.Get("service_resources").(*schema.Set).List()),
		},
		Name: &name,
	}

	future, err := client.CreateOrUpdate(ctx, policyId.ResourceGroup, policyId.Name, name, param)
	if err != nil {
		return fmt.Errorf("creating Service Endpoint Policy Definition %q (Policy %q / Resource Group %q): %+v", name, policyId.Name, policyId.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Service Endpoint Policy Definition %q (Policy %q / Resource Group %q): %+v", name, policyId.Name, policyId.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, policyId.ResourceGroup, policyId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Service Endpoint Policy Definition %q (Policy %q / Resource Group %q): %+v", name, policyId.Name, policyId.ResourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Service Endpoint Policy Definition %q (Policy %q / Resource Group %q) ID", name, policyId.Name, policyId.ResourceGroup)
	}

	id, err := parse.ServiceEndpointPolicyDefinitionID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(subscriptionId))

	return resourceArmServiceEndpointPolicyDefinitionRead(d, meta)
}

func resourceArmServiceEndpointPolicyDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.ServiceEndpointPolicyDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceEndpointPolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Policy, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Service Endpoint Policy Definition %q was not found in Policy %q / Resource Group %q - removing from state!", id.Name, id.Policy, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Service Endpoint Policy Definition %q (Policy %q / Resource Group %q): %+v", id.Name, id.Policy, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)

	policyId := parse.NewServiceEndpointPolicyID(id.ResourceGroup, id.Policy)
	d.Set("policy_id", policyId.ID(subscriptionId))

	if prop := resp.ServiceEndpointPolicyDefinitionPropertiesFormat; prop != nil {
		d.Set("description", prop.Description)
		d.Set("service_endpoint_name", prop.Service)
		if err := d.Set("service_resources", utils.FlattenStringSlice(prop.ServiceResources)); err != nil {
			return fmt.Errorf("setting `service_resources`")
		}
	}

	return nil
}

func resourceArmServiceEndpointPolicyDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicyDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceEndpointPolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Policy, id.Name); err != nil {
		return fmt.Errorf("deleting Service Endpoint Policy Definition %q (Policy %q / Resource Group %q): %+v", id.Name, id.Policy, id.ResourceGroup, err)
	}

	return nil
}
