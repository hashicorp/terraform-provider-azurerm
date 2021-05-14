package servicebus

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/servicebus/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ServiceBusNamespaceDisasterRecoveryConfigResource struct {
}

func resourceServiceBusNamespaceDisasterRecoveryConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceBusNamespaceDisasterRecoveryConfigCreate,
		Read:   resourceServiceBusNamespaceDisasterRecoveryConfigRead,
		Update: resourceServiceBusNamespaceDisasterRecoveryConfigUpdate,
		Delete: resourceServiceBusNamespaceDisasterRecoveryConfigDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NamespaceDisasterRecoveryConfigID(id)
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"primary_namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"partner_namespace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"alias_primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"alias_secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"default_secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceServiceBusNamespaceDisasterRecoveryConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for ServiceBus Namespace pairing create/update.")

	id, err := parse.NamespaceID(d.Get("primary_namespace_id").(string))
	if err != nil {
		return err
	}

	aliasName := d.Get("name").(string)
	partnerNamespaceId := d.Get("partner_namespace_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, aliasName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("error checking for presence of existing Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", aliasName, id.Name, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_servicebus_namespace_disaster_recovery_config", *existing.ID)
		}
	}

	parameters := servicebus.ArmDisasterRecovery{
		ArmDisasterRecoveryProperties: &servicebus.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(partnerNamespaceId),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, aliasName, parameters); err != nil {
		return fmt.Errorf("error creating/updating Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", aliasName, id.Name, id.ResourceGroup, err)
	}

	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id.ResourceGroup, id.Name, aliasName, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("error waiting for replication to complete for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", aliasName, id.Name, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name, aliasName)
	if err != nil {
		return fmt.Errorf("error reading Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %v", aliasName, id.Name, id.ResourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("got nil ID for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q)", aliasName, id.Name, id.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceServiceBusNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceServiceBusNamespaceDisasterRecoveryConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceDisasterRecoveryConfigID(d.State().ID)
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, serviceBusNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, serviceBusNamespaceResourceName)

	if d.HasChange("partner_namespace_id") {
		breakPair, err := client.BreakPairing(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName)
		if breakPair.StatusCode != http.StatusOK {
			return fmt.Errorf("error issuing break pairing request for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
		}

		if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("error waiting for break pairing request to complete for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
		}
	}

	parameters := servicebus.ArmDisasterRecovery{
		ArmDisasterRecoveryProperties: &servicebus.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(d.Get("partner_namespace_id").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName, parameters); err != nil {
		return fmt.Errorf("error creating/updating Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
	}

	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("error waiting for replication to complete for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
	}

	return resourceServiceBusNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceServiceBusNamespaceDisasterRecoveryConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	primaryId := parse.NewNamespaceID(id.SubscriptionId, id.ResourceGroup, id.NamespaceName)

	d.Set("name", id.DisasterRecoveryConfigName)
	d.Set("primary_namespace_id", primaryId.ID())
	d.Set("partner_namespace_id", resp.ArmDisasterRecoveryProperties.PartnerNamespace)

	keys, err := client.ListKeys(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName, serviceBusNamespaceDefaultAuthorizationRule)

	if err != nil {
		log.Printf("[WARN] listing default keys for %s: %+v", id, err)
	} else {
		d.Set("alias_primary_connection_string", keys.AliasPrimaryConnectionString)
		d.Set("alias_secondary_connection_string", keys.AliasSecondaryConnectionString)
		d.Set("default_primary_key", keys.PrimaryKey)
		d.Set("default_secondary_key", keys.SecondaryKey)
	}

	return nil
}

func resourceServiceBusNamespaceDisasterRecoveryConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ServiceBus.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceDisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	breakPair, err := client.BreakPairing(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName)
	if err != nil {
		return fmt.Errorf("breaking pairing %s: %+v", id, err)
	}

	if breakPair.StatusCode != http.StatusOK {
		return fmt.Errorf("error breaking pairing for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
	}

	if err := resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName, d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("error waiting for break pairing request to complete for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName); err != nil {
		return fmt.Errorf("error issuing delete request for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
	}

	// no future for deletion so wait for it to vanish
	deleteWait := &resource.StateChangeConf{
		Pending:    []string{"200"},
		Target:     []string{"404"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.DisasterRecoveryConfigName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return resp, strconv.Itoa(resp.StatusCode), nil
				}
				return nil, "nil", fmt.Errorf("error polling for the status of the Service Bus Namespace Disaster Recovery Configs %q deletion (Namespace %q / Resource Group %q): %v", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
			}

			return resp, strconv.Itoa(resp.StatusCode), nil
		},
	}

	if _, err := deleteWait.WaitForState(); err != nil {
		return fmt.Errorf("error waiting the deletion of Service Bus Namespace Disaster Recovery Configs %q deletion (Namespace %q / Resource Group %q): %v", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
	}

	// it can take some time for the name to become available again
	// this is mainly here 	to enable updating the resource in place
	nameFreeWait := &resource.StateChangeConf{
		Pending:    []string{"NameInUse"},
		Target:     []string{"None"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.CheckNameAvailabilityMethod(ctx, id.ResourceGroup, id.NamespaceName, servicebus.CheckNameAvailability{Name: utils.String(id.DisasterRecoveryConfigName)})
			if err != nil {
				return resp, "Error", fmt.Errorf("error checking if the Service Bus Namespace Disaster Recovery Configs %q name has been freed (Namespace %q / Resource Group %q): %v", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
			}

			return resp, string(resp.Reason), nil
		},
	}

	if _, err := nameFreeWait.WaitForState(); err != nil {
		return fmt.Errorf("error waiting the the Service Bus Namespace Disaster Recovery Configs %q name to be available (Namespace %q / Resource Group %q): %v", id.DisasterRecoveryConfigName, id.NamespaceName, id.ResourceGroup, err)
	}

	return nil
}

func resourceServiceBusNamespaceDisasterRecoveryConfigWaitForState(ctx context.Context, client *servicebus.DisasterRecoveryConfigsClient, resourceGroup, namespaceName, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(servicebus.Accepted)},
		Target:     []string{string(servicebus.Succeeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
		Refresh: func() (interface{}, string, error) {
			read, err := client.Get(ctx, resourceGroup, namespaceName, name)
			if err != nil {
				return nil, "error", fmt.Errorf("wait read Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %v", name, namespaceName, resourceGroup, err)
			}

			if props := read.ArmDisasterRecoveryProperties; props != nil {
				if props.ProvisioningState == servicebus.Failed {
					return read, "failed", fmt.Errorf("replication for Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q) failed", name, namespaceName, resourceGroup)
				}
				return read, string(props.ProvisioningState), nil
			}

			return read, "nil", fmt.Errorf("waiting for replication error Service Bus Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): provisioning state is nil", name, namespaceName, resourceGroup)
		},
	}

	_, err := stateConf.WaitForState()
	return err
}
