package eventhub

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/checknameavailabilitydisasterrecoveryconfigs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/sdk/disasterrecoveryconfigs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceEventHubNamespaceDisasterRecoveryConfig() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceEventHubNamespaceDisasterRecoveryConfigCreate,
		Read:   resourceEventHubNamespaceDisasterRecoveryConfigRead,
		Update: resourceEventHubNamespaceDisasterRecoveryConfigUpdate,
		Delete: resourceEventHubNamespaceDisasterRecoveryConfigDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"partner_namespace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			// this property is broken and should not be reimplemented after 3.0 until this is addressed: https://github.com/Azure/azure-sdk-for-go/issues/5893
			"alternate_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.ValidateEventHubNamespaceName(),
				Deprecated:   "This property has been deprecated and will be removed in v3.0 of the provider as any DRC created with an alternate name cannot be deleted and the service is not going to change this. Please see: https://github.com/Azure/azure-sdk-for-go/issues/5893",
			},
		},
	}
}

func resourceEventHubNamespaceDisasterRecoveryConfigCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM EventHub Namespace Disaster Recovery Configs creation.")

	id := disasterrecoveryconfigs.NewDisasterRecoveryConfigID(subscriptionId, d.Get("resource_group_name").(string), d.Get("namespace_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_eventhub_namespace_disaster_recovery_config", id.ID())
		}
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	parameters := disasterrecoveryconfigs.ArmDisasterRecovery{
		Properties: &disasterrecoveryconfigs.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(d.Get("partner_namespace_id").(string)),
		},
	}

	if v, ok := d.GetOk("alternate_name"); ok {
		parameters.Properties.AlternateName = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, id); err != nil {
		return fmt.Errorf("waiting for replication of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceEventHubNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceEventHubNamespaceDisasterRecoveryConfigUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.DisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	if d.HasChange("partner_namespace_id") {
		// break pairing
		if _, err := client.BreakPairing(ctx, *id); err != nil {
			return fmt.Errorf("breaking the pairing for %s: %+v", *id, err)
		}

		if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
			return fmt.Errorf("waiting for the pairing to be broken for %s: %+v", *id, err)
		}
	}

	parameters := disasterrecoveryconfigs.ArmDisasterRecovery{
		Properties: &disasterrecoveryconfigs.ArmDisasterRecoveryProperties{
			PartnerNamespace: utils.String(d.Get("partner_namespace_id").(string)),
		},
	}

	if v, ok := d.GetOk("alternate_name"); ok {
		parameters.Properties.AlternateName = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for replication after update of %s: %+v", *id, err)
	}

	return resourceEventHubNamespaceDisasterRecoveryConfigRead(d, meta)
}

func resourceEventHubNamespaceDisasterRecoveryConfigRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.DisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil && model.Properties != nil {
		d.Set("partner_namespace_id", model.Properties.PartnerNamespace)
		d.Set("alternate_name", model.Properties.AlternateName)
	}

	return nil
}

func resourceEventHubNamespaceDisasterRecoveryConfigDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.DisasterRecoveryConfigsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := disasterrecoveryconfigs.DisasterRecoveryConfigID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NamespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, eventHubNamespaceResourceName)

	if _, err := client.BreakPairing(ctx, *id); err != nil {
		return fmt.Errorf("breaking pairing of %s: %+v", *id, err)
	}

	if err := resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx, client, *id); err != nil {
		return fmt.Errorf("waiting for pairing to break for %s: %+v", *id, err)
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	// no future for deletion so wait for it to vanish
	deleteWait := &pluginsdk.StateChangeConf{
		Pending:    []string{"200"},
		Target:     []string{"404"},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (interface{}, string, error) {
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
				}
				return nil, "nil", fmt.Errorf("polling to check the deletion state for %s: %+v", *id, err)
			}

			return resp, strconv.Itoa(resp.HttpResponse.StatusCode), nil
		},
	}

	if _, err := deleteWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting the deletion of %s: %+v", *id, err)
	}

	// it can take some time for the name to become available again
	// this is mainly here	to enable updating the resource in place
	parentNamespaceId := checknameavailabilitydisasterrecoveryconfigs.NewNamespaceID(id.SubscriptionId, id.ResourceGroup, id.NamespaceName)
	availabilityClient := meta.(*clients.Client).Eventhub.DisasterRecoveryNameAvailabilityClient
	nameFreeWait := &pluginsdk.StateChangeConf{
		Pending:    []string{"NameInUse"},
		Target:     []string{"None"},
		MinTimeout: 30 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
		Refresh: func() (interface{}, string, error) {
			input := checknameavailabilitydisasterrecoveryconfigs.CheckNameAvailabilityParameter{
				Name: id.Name,
			}
			resp, err := availabilityClient.DisasterRecoveryConfigsCheckNameAvailability(ctx, parentNamespaceId, input)
			if err != nil {
				return resp, "Error", fmt.Errorf("waiting for the name of %s to become free: %v", *id, err)
			}
			// TODO: new crash to handle here
			return resp, string(*resp.Model.Reason), nil
		},
	}

	if _, err := nameFreeWait.WaitForStateContext(ctx); err != nil {
		return err
	}

	return nil
}

func resourceEventHubNamespaceDisasterRecoveryConfigWaitForState(ctx context.Context, client *disasterrecoveryconfigs.DisasterRecoveryConfigsClient, id disasterrecoveryconfigs.DisasterRecoveryConfigId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(eventhub.ProvisioningStateDRAccepted)},
		Target:     []string{string(eventhub.ProvisioningStateDRSucceeded)},
		MinTimeout: 30 * time.Second,
		Timeout:    time.Until(deadline),
		Refresh: func() (interface{}, string, error) {
			read, err := client.Get(ctx, id)
			if err != nil {
				return nil, "error", fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := read.Model; model != nil {
				if props := model.Properties; props != nil {
					if props.ProvisioningState == nil {
						return read, "failed", fmt.Errorf("provisioningState was empty")
					}

					if *props.ProvisioningState == disasterrecoveryconfigs.ProvisioningStateDRFailed {
						return read, "failed", fmt.Errorf("replication failed for %s: %+v", id, err)
					}
					return read, string(*props.ProvisioningState), nil
				}
			}

			return read, "nil", fmt.Errorf("waiting for replication of %s: %+v", id, err)
		},
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
