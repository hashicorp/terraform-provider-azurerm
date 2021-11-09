package web

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	relayParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/sdk/2017-04-01/namespaces"
	relayValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceHybridConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceHybridConnectionCreateUpdate,
		Read:   resourceAppServiceHybridConnectionRead,
		Update: resourceAppServiceHybridConnectionCreateUpdate,
		Delete: resourceAppServiceHybridConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.HybridConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"app_service_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"relay_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: relayValidate.HybridConnectionID,
			},

			"hostname": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"port": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: azValidate.PortNumberOrZero,
			},

			"send_key_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      "RootManageSharedAccessKey",
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"namespace_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"relay_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_bus_namespace": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_bus_suffix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"send_key_value": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceAppServiceHybridConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("app_service_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	relayArmURI := d.Get("relay_id").(string)
	relayId, err := relayParse.HybridConnectionID(relayArmURI)
	if err != nil {
		return fmt.Errorf("parsing relay ID %q: %s", relayArmURI, err)
	}
	namespaceName := relayId.NamespaceName
	relayName := relayId.Name

	if d.IsNewResource() {
		existing, err := client.GetHybridConnection(ctx, resourceGroup, name, namespaceName, relayName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing App Service Hybrid Connection %q (Resource Group %q, Namespace %q, Relay Name %q): %s", name, resourceGroup, namespaceName, relayName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_app_service_hybrid_connection", *existing.ID)
		}
	}

	port := int32(d.Get("port").(int))

	connectionEnvelope := web.HybridConnection{
		HybridConnectionProperties: &web.HybridConnectionProperties{
			RelayArmURI:  &relayArmURI,
			Hostname:     utils.String(d.Get("hostname").(string)),
			Port:         &port,
			SendKeyName:  utils.String(d.Get("send_key_name").(string)),
			SendKeyValue: utils.String(""), // The service creates this no matter what is sent, but the API requires the field to be set
		},
	}

	hybridConnection, err := client.CreateOrUpdateHybridConnection(ctx, resourceGroup, name, namespaceName, relayName, connectionEnvelope)
	if err != nil {
		return fmt.Errorf("failed creating App Service Hybrid Connection %q (resource group %q): %s", name, resourceGroup, err)
	}

	if hybridConnection.ID == nil && *hybridConnection.ID == "" {
		return fmt.Errorf("failed to read ID for Hybrid Connection %q", name)
	}
	d.SetId(*hybridConnection.ID)

	return resourceAppServiceHybridConnectionRead(d, meta)
}

func resourceAppServiceHybridConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Hybrid Connection %q for App Service %q (Resource Group %q) was not found - removing from state", id.HybridConnectionNamespaceName, id.SiteName, id.RelayName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving App Service Hybrid Connection %q in Namespace %q, Resource Group %q: %s", id.SiteName, id.HybridConnectionNamespaceName, id.ResourceGroup, err)
	}
	d.Set("app_service_name", id.SiteName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("namespace_name", id.HybridConnectionNamespaceName)
	d.Set("relay_name", id.RelayName)

	if props := resp.HybridConnectionProperties; props != nil {
		d.Set("port", resp.Port)
		d.Set("service_bus_namespace", resp.ServiceBusNamespace)
		d.Set("send_key_name", resp.SendKeyName)
		d.Set("service_bus_suffix", resp.ServiceBusSuffix)
		d.Set("relay_id", resp.RelayArmURI)
		d.Set("hostname", resp.Hostname)
	}

	// key values are not returned in the response, so we get the primary key from the relay namespace ListKeys func
	if resp.ServiceBusNamespace != nil && resp.SendKeyName != nil {
		relayNamespacesClient := meta.(*clients.Client).Relay.NamespacesClient
		relayNamespaceRG, err := findRelayNamespace(ctx, relayNamespacesClient, id.SubscriptionId, *resp.ServiceBusNamespace)
		if err != nil {
			return err
		}
		authRuleId := namespaces.NewAuthorizationRuleID(id.SubscriptionId, *relayNamespaceRG, *resp.ServiceBusNamespace, *resp.SendKeyName)
		accessKeys, err := relayNamespacesClient.ListKeys(ctx, authRuleId)
		if err != nil {
			return fmt.Errorf("unable to List Access Keys for Namespace %q (Resource Group %q): %+v", *resp.ServiceBusNamespace, id.ResourceGroup, err)
		}

		if model := accessKeys.Model; model != nil {
			d.Set("send_key_value", model.PrimaryKey)
		}
	}

	return nil
}

func resourceAppServiceHybridConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteHybridConnection(ctx, id.ResourceGroup, id.SiteName, id.HybridConnectionNamespaceName, id.RelayName)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting App Service Hybrid Connection %q (Resource Group %q, Relay %q): %+v", id.SiteName, id.ResourceGroup, id.RelayName, err)
		}
	}

	return nil
}

func findRelayNamespace(ctx context.Context, client *namespaces.NamespacesClient, subscriptionId, name string) (*string, error) {
	subId := namespaces.NewSubscriptionID(subscriptionId)
	relayNSIterator, err := client.ListComplete(ctx, subId)
	if err != nil {
		return nil, fmt.Errorf("listing Relay Namespaces: %+v", err)
	}

	var found *namespaces.RelayNamespace
	for _, item := range relayNSIterator.Items {
		if item.Name != nil && *item.Name == name {
			found = &item
			break
		}
	}

	if found == nil || found.Id == nil {
		return nil, fmt.Errorf("could not find Relay Namespace with name: %q", name)
	}

	id, err := namespaces.ParseNamespaceID(*found.Id)
	if err != nil {
		return nil, fmt.Errorf("relay Namespace id not valid: %+v", err)
	}
	return &id.ResourceGroup, nil
}
