package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceHybridConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceHybridConnectionCreateUpdate,
		Read:   resourceArmAppServiceHybridConnectionRead,
		Update: resourceArmAppServiceHybridConnectionCreateUpdate,
		Delete: resourceArmAppServiceHybridConnectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"app_service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(6, 50),
			},
			"relay_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"relay_arm_uri": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.PortNumberOrZero,
			},
			"service_bus_namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_bus_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"send_key_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"send_key_value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArmAppServiceHybridConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("app_service_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	relayName := d.Get("relay_name").(string)
	relayNameArmUri := d.Get("relay_arm_uri").(string)
	hostname := d.Get("hostname").(string)
	port := int32(d.Get("port").(int))
	serviceBusSuffix := d.Get("service_bus_suffix").(string)
	serviceBusNamespace := d.Get("service_bus_namespace").(string)
	sendKeyName := d.Get("send_key_name").(string)
	sendKeyValue := d.Get("send_key_value").(string)

	hybridConnectionProperties := web.HybridConnectionProperties{
		ServiceBusNamespace: &serviceBusNamespace,
		RelayName:           &relayName,
		RelayArmURI:         &relayNameArmUri,
		Hostname:            &hostname,
		Port:                &port,
		SendKeyName:         &sendKeyName,
		SendKeyValue:        &sendKeyValue,
		ServiceBusSuffix:    &serviceBusSuffix,
	}

	connectionEnvelope := web.HybridConnection{}
	connectionEnvelope.HybridConnectionProperties = &hybridConnectionProperties

	if _, err := client.CreateOrUpdateHybridConnection(ctx, resGroup, name, namespaceName, relayName, connectionEnvelope); err != nil {
		return fmt.Errorf("Error creating App Service Hybrid Connection %q (resource group %q): %s", name, resGroup, err)
	}
	return resourceArmAppServiceHybridConnectionRead(d, meta)
}

func resourceArmAppServiceHybridConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["sites"]
	namespaceName := id.Path["hybridConnectionNamespaces"]
	relayName := id.Path["relays"]

	resp, err := client.GetHybridConnection(ctx, resGroup, name, namespaceName, relayName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on App Service Hybrid Connection %q in Namespace %q, Resource Group %q: %s", name, namespaceName, resGroup, err)
	}
	d.Set("app_service_name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("namespace_name", namespaceName)
	d.Set("relay_name", relayName)

	if props := resp.HybridConnectionProperties; props != nil {
		d.Set("port", resp.Port)
		d.Set("service_bus_namespace", resp.ServiceBusNamespace)
		d.Set("send_key_name", resp.SendKeyName)
		d.Set("send_key_value", resp.SendKeyValue)
		d.Set("service_bus_suffix", resp.ServiceBusSuffix)
		d.Set("relay_arm_uri", resp.RelayArmURI)
		d.Set("hostname", resp.Hostname)
	}

	return nil
}

func resourceArmAppServiceHybridConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["sites"]
	namespaceName := id.Path["hybridConnectionNamespaces"]
	relayName := id.Path["relays"]

	resp, err := client.DeleteHybridConnection(ctx, resGroup, name, namespaceName, relayName)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return err
	}

	return nil
}
