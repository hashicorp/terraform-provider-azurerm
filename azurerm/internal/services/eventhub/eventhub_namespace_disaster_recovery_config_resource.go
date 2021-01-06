package eventhub

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceEventHubNamespaceDisasterRecoveryConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventHubNamespaceDisasterRecoveryConfigCreate,
		Read:   resourceEventHubNamespaceDisasterRecoveryConfigRead,
		Update: resourceEventHubNamespaceDisasterRecoveryConfigUpdate,
		Delete: resourceEventHubNamespaceDisasterRecoveryConfigDelete,

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
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"partner_namespace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"alternate_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateEventHubNamespaceName(),
			},
		},
	}
}

func resourceEventHubNamespaceDisasterRecoveryConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace Disaster Recovery Configs creation.")

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, namespaceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_disaster_recovery_config", *existing.ID)
		}
	}

	parameters := eventhub.ArmDisasterRecovery{
		ArmDisasterRecoveryProperties: &eventhub.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(d.Get("partner_namespace_id").(string)),
		},
	}

	if v, ok := d.GetOk("alternate_name"); ok {
		parameters.ArmDisasterRecoveryProperties.AlternateName = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, resourceGroup, namespaceName, name, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("Error waiting for replication to complete for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error reading EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %v", name, namespaceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Got nil ID for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q)", name, namespaceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceEventHubNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceEventHubNamespaceDisasterRecoveryConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["disasterRecoveryConfigs"]
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]

	if d.HasChange("partner_namespace_id") {
		// break pairing
		breakPair, err := client.BreakPairing(ctx, resourceGroup, namespaceName, name)
		if breakPair.StatusCode != http.StatusOK {
			return fmt.Errorf("Error issuing break pairing request for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
		}

		if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, resourceGroup, namespaceName, name, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("Error waiting for break pairing request to complete for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
		}
	}

	parameters := eventhub.ArmDisasterRecovery{
		ArmDisasterRecoveryProperties: &eventhub.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(d.Get("partner_namespace_id").(string)),
		},
	}

	if v, ok := d.GetOk("alternate_name"); ok {
		parameters.ArmDisasterRecoveryProperties.AlternateName = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, namespaceName, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, resourceGroup, namespaceName, name, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return fmt.Errorf("Error waiting for replication to complete for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	return resourceEventHubNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceEventHubNamespaceDisasterRecoveryConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["disasterRecoveryConfigs"]
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]

	resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)

	if properties := resp.ArmDisasterRecoveryProperties; properties != nil {
		d.Set("partner_namespace_id", properties.PartnerNamespace)
		d.Set("alternate_name", properties.AlternateName)
	}

	return nil
}

func resourceEventHubNamespaceDisasterRecoveryConfigDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["disasterRecoveryConfigs"]
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]

	breakPair, err := client.BreakPairing(ctx, resourceGroup, namespaceName, name)
	if err != nil {
		return fmt.Errorf("Error issuing break pairing request for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}
	if breakPair.StatusCode != http.StatusOK {
		return fmt.Errorf("Error breaking pairing for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, resourceGroup, namespaceName, name, d.Timeout(schema.TimeoutDelete)); err != nil {
		return fmt.Errorf("Error waiting for break pairing request to complete for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	if _, err := client.Delete(ctx, resourceGroup, namespaceName, name); err != nil {
		return fmt.Errorf("Error issuing delete request for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %s", name, namespaceName, resourceGroup, err)
	}

	// no future for deletion so wait for it to vanish
	deleteWait := &resource.StateChangeConf{
		Pending:    []string{"200"},
		Target:     []string{"404"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, resourceGroup, namespaceName, name)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return resp, strconv.Itoa(resp.StatusCode), nil
				}
				return nil, "nil", fmt.Errorf("Error polling for the status of the EventHub Namespace Disaster Recovery Configs %q deletion (Namespace %q / Resource Group %q): %v", name, namespaceName, resourceGroup, err)
			}

			return resp, strconv.Itoa(resp.StatusCode), nil
		},
	}

	if _, err := deleteWait.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting the deletion of EventHub Namespace Disaster Recovery Configs %q deletion (Namespace %q / Resource Group %q): %v", name, namespaceName, resourceGroup, err)
	}

	// it can take some time for the name to become available again
	// this is mainly here 	to enable updating the resource in place
	nameFreeWait := &resource.StateChangeConf{
		Pending:    []string{"NameInUse"},
		Target:     []string{"None"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.CheckNameAvailability(ctx, resourceGroup, namespaceName, eventhub.CheckNameAvailabilityParameter{Name: utils.String(name)})
			if err != nil {
				return resp, "Error", fmt.Errorf("Error checking if the EventHub Namespace Disaster Recovery Configs %q name has been freed (Namespace %q / Resource Group %q): %v", name, namespaceName, resourceGroup, err)
			}

			return resp, string(resp.Reason), nil
		},
	}

	if _, err := nameFreeWait.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting the the EventHub Namespace Disaster Recovery Configs %q name to be available (Namespace %q / Resource Group %q): %v", name, namespaceName, resourceGroup, err)
	}

	return nil
}

func resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx context.Context, client *eventhub.DisasterRecoveryConfigsClient, resourceGroup, namespaceName, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(eventhub.ProvisioningStateDRAccepted)},
		Target:     []string{string(eventhub.ProvisioningStateDRSucceeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    timeout,
		Refresh: func() (interface{}, string, error) {
			read, err := client.Get(ctx, resourceGroup, namespaceName, name)
			if err != nil {
				return nil, "error", fmt.Errorf("Wait read EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): %v", name, namespaceName, resourceGroup, err)
			}

			if props := read.ArmDisasterRecoveryProperties; props != nil {
				if props.ProvisioningState == eventhub.ProvisioningStateDRFailed {
					return read, "failed", fmt.Errorf("Replication for EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q) failed!", name, namespaceName, resourceGroup)
				}
				return read, string(props.ProvisioningState), nil
			}

			return read, "nil", fmt.Errorf("Waiting for replication error EventHub Namespace Disaster Recovery Configs %q (Namespace %q / Resource Group %q): provisioning state is nil", name, namespaceName, resourceGroup)
		},
	}

	_, err := stateConf.WaitForState()
	return err
}
