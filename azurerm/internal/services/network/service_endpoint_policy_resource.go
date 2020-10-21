package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmServiceEndpointPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmServiceEndpointPolicyCreateUpdate,
		Read:   resourceArmServiceEndpointPolicyRead,
		Update: resourceArmServiceEndpointPolicyCreateUpdate,
		Delete: resourceArmServiceEndpointPolicyDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ServiceEndpointPolicyID(id)
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
				ValidateFunc: validate.ServiceEndpointPolicyName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": location.Schema(),

			"tags": tags.Schema(),
		},
	}
}

func resourceArmServiceEndpointPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Service Endpoint Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_service_endpoint_policy", *resp.ID)
		}
	}

	param := network.ServiceEndpointPolicy{
		Location: &location,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, param)
	if err != nil {
		return fmt.Errorf("creating Service Endpoint Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Service Endpoint Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Service Endpoint Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Service Endpoint Policy %q (Resource Group %q) ID", name, resourceGroup)
	}

	id, err := parse.ServiceEndpointPolicyID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(subscriptionId))

	return resourceArmServiceEndpointPolicyRead(d, meta)
}

func resourceArmServiceEndpointPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceEndpointPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Service Endpoint Policy %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Service Endpoint Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmServiceEndpointPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ServiceEndpointPolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting Service Endpoint Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}
