package network

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-05-01/network"
)

type commissionedAction string

const (
	commissionedActionProvision    commissionedAction = "Provision"
	commissionedActionCommission   commissionedAction = "Commission"
	commissionedActionDecommission commissionedAction = "Decommission"
	commissionedActionDeprovision  commissionedAction = "Deprovision"
)

type commissionedState string

const (
	commissionedStateDeprovisioned    commissionedState = "Deprovisioned"
	commissionedStateValidationFailed commissionedState = "ValidationFailed"
)

func resourceCustomIpPrefix() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCustomIpPrefixCreate,
		Read:   resourceCustomIpPrefixRead,
		Update: resourceCustomIpPrefixUpdate,
		Delete: resourceCustomIpPrefixDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CustomIpPrefixID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(1020 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(1020 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CustomIpPrefixName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"cidr": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"roa_expiration_date": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"wan_validation_signed_message": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"action": {
				Type:     pluginsdk.TypeString,
				Default:  string(commissionedActionProvision),
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					// When the operation is successful, the possible values of `CommissionedState` are :`Provisioned`,`Commissioned`,`Deprovisioned`.
					// The value of action is set as `Provision`, the commissionState returned `Provisioned` by API.
					// The value of action is set as `Decommission`, the commissionState returned `Provisioned` by API.
					// The value of action is set as `Commission`, the commissionState returned `Commissioned` by API.
					// The value of action is set as `Deprovision`, the commissionState returned `Deprovisioned` by API.
					// so we should suppress diff when the commissionState maps the action in the configuration.
					if old == string(network.CommissionedStateProvisioned) && (new == string(commissionedActionProvision) || new == string(commissionedActionDecommission)) {
						return true
					}
					if old == string(network.CommissionedStateCommissioned) && new == string(commissionedActionCommission) {
						return true
					}
					if old == string(commissionedStateDeprovisioned) && new == string(commissionedActionDeprovision) {
						return true
					}
					return false
				},
				ValidateFunc: validation.StringInSlice([]string{
					string(commissionedActionProvision),
					string(commissionedActionCommission),
					string(commissionedActionDecommission),
					string(commissionedActionDeprovision),
				}, false),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceCustomIpPrefixCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient
	subscriptionId := meta.(*clients.Client).Network.CustomIPPrefixesClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewCustomIpPrefixID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_custom_ip_prefix", id.ID())
	}

	// The `CommissionedState` can only be `Provisioning` when creating a resource.
	parameters := network.CustomIPPrefix{
		CustomIPPrefixPropertiesFormat: &network.CustomIPPrefixPropertiesFormat{
			Cidr:              utils.String(d.Get("cidr").(string)),
			CommissionedState: network.CommissionedStateProvisioning,
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = &zones
	}

	if v, ok := d.GetOk("roa_expiration_date"); ok {
		parameters.AuthorizationMessage = utils.String(subscriptionId + "|" + d.Get("cidr").(string) + "|" + v.(string))
	}

	if v, ok := d.GetOk("wan_validation_signed_message"); ok {
		parameters.SignedMessage = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{string(network.CommissionedStateProvisioning)},
		Target:                    []string{string(network.CommissionedStateProvisioned)},
		Refresh:                   customIpPrefixCreateRefreshFunc(ctx, client, id),
		PollInterval:              1 * time.Minute,
		ContinuousTargetOccurence: 3,
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCustomIpPrefixRead(d, meta)
}

func resourceCustomIpPrefixUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomIpPrefixID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	// Since the "Provisioning", "Commissioning", "Decommissioning" or `Deprovisioning` commissioned state of custom ip prefix can do nothing, it could not be updated.
	if existing.CommissionedState == network.CommissionedStateProvisioning || existing.CommissionedState == network.CommissionedStateCommissioning || existing.CommissionedState == network.CommissionedStateDecommissioning || existing.CommissionedState == network.CommissionedStateDeprovisioning {
		return fmt.Errorf("the commissioned state of %s is %s, it doesn't allow to be updated", id, string(existing.CommissionedState))
	}

	if existing.CustomIPPrefixPropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	props := existing.CustomIPPrefixPropertiesFormat

	if d.HasChange("action") {
		if err := prepareUpdate(ctx, d, meta, *id, props); err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		existing.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), existing)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	// The `commissionedstate` of custom ip prefix is in "Provisioning", "Commissioning", "Decommissioning" or `Deprovisioning` after waiting for completion.
	// We should poll the `commissionedstate` until the target state is complete. The detailed changes of the commissionedstate are as follows:
	// `Provisioning`      => `Provisioned`   - can do commission, or deprovision
	// `Commissioning`     => `Commissioned`  - can do decommission
	// `Decommissioning`   => `Provisioned`   - can do commission, or deprovision
	// `Deprovisioning`    => `Deprovisioned` - can do delete or provision
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{string(network.CommissionedStateProvisioning), string(network.CommissionedStateCommissioning), string(network.CommissionedStateDecommissioning), string(network.CommissionedStateDeprovisioning)},
		Target:                    []string{string(network.CommissionedStateCommissioned), string(network.CommissionedStateProvisioned), string(commissionedStateDeprovisioned)},
		Refresh:                   customIpPrefixCreateRefreshFunc(ctx, client, *id),
		PollInterval:              1 * time.Minute,
		ContinuousTargetOccurence: 3,
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceCustomIpPrefixRead(d, meta)
}

func customIpPrefixCreateRefreshFunc(ctx context.Context, client *network.CustomIPPrefixesClient, id parse.CustomIpPrefixId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.CustomIpPrefixeName, "")
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return nil, "NotFound", nil
			}
			return nil, "", fmt.Errorf("polling for the commissionedstate of the Custom Ip Prefix %q (Resource Group: %q): %+v", id.CustomIpPrefixeName, id.ResourceGroup, err)
		}

		if res.CustomIPPrefixPropertiesFormat == nil {
			return nil, "", fmt.Errorf("unexpected nil properties format of Custom Ip Prefix %q (Resource Group %q)", id.CustomIpPrefixeName, id.ResourceGroup)
		}
		return res, string(res.CustomIPPrefixPropertiesFormat.CommissionedState), nil
	}
}

func resourceCustomIpPrefixRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomIpPrefixID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.CustomIpPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", resp.Name)
	d.Set("location", azure.NormalizeLocation(*resp.Location))
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.CustomIPPrefixPropertiesFormat; props != nil {
		d.Set("cidr", props.Cidr)
		if v := props.AuthorizationMessage; v != nil {
			authMessage := strings.Split(*v, "|")
			if len(authMessage) == 3 {
				d.Set("roa_expiration_date", authMessage[2])
			}
		}
		d.Set("wan_validation_signed_message", props.SignedMessage)
		d.Set("action", props.CommissionedState)
	}

	if v := resp.Zones; v != nil {
		d.Set("zones", zones.Flatten(v))
	}

	if err := d.Set("tags", tags.Flatten(resp.Tags)); err != nil {
		return fmt.Errorf("setting `tags`: %+v", err)
	}

	return nil
}

func resourceCustomIpPrefixDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomIpPrefixID(d.Id())
	if err != nil {
		return err
	}

	if err := prepareDelete(ctx, d, meta, *id); err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.CustomIpPrefixeName)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}

	return nil
}

func prepareUpdate(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}, id parse.CustomIpPrefixId, props *network.CustomIPPrefixPropertiesFormat) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient

	existing, err := client.Get(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	currentState := existing.CommissionedState

	// Since the "Provisioning", "Commissioning", "Decommissioning" or `Deprovisioning` commissioned state of custom ip prefix can do nothing, it could not be updated.
	if currentState == network.CommissionedStateProvisioning || currentState == network.CommissionedStateCommissioning || currentState == network.CommissionedStateDecommissioning || currentState == network.CommissionedStateDeprovisioning {
		return fmt.Errorf("the commissioned state of %s is %s, it doesn't allow to be updated", id, string(existing.CommissionedState))
	}

	action := commissionedAction(d.Get("action").(string))

	if (currentState == network.CommissionedStateProvisioned && action == commissionedActionProvision) ||
		(currentState == network.CommissionedStateCommissioned && action == commissionedActionCommission) ||
		(currentState == network.CommissionedStateProvisioned && action == commissionedActionDecommission) ||
		(string(currentState) == string(commissionedStateDeprovisioned) && action == commissionedActionDeprovision) {
		return fmt.Errorf("the current commissioned state of %s is %s, no need to update", id, string(existing.CommissionedState))
	}

	switch action {
	case commissionedActionProvision:
		if string(currentState) == string(network.CommissionedStateCommissioned) {
			// When the CommissionedState of existed resource is `Commissioned`, the CommissionedState can do `Decommission`.
			// After `Decommission` is done, the CommissionedState changes back to `Provisioned`.
			// So, when current CommissionedState is Commissioned, Decommission is executed instead of Provision.
			props.CommissionedState = network.CommissionedStateDecommissioning
		} else {
			props.CommissionedState = network.CommissionedStateProvisioning
		}
	case commissionedActionCommission:
		if string(currentState) == string(commissionedStateDeprovisioned) {
			// When the CommissionedState of existed resource is `Deprovisioned`, the CommissionedState can do `Provision`.
			// Since `Provisioned` of CommissionedState can do `Commission`, when current CommissionedState is Deprovisioned, do Provision first, then do Commission.
			if updateCommissionedState(d, meta, id, existing, network.CommissionedStateProvisioning) != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		}
		props.CommissionedState = network.CommissionedStateCommissioning

	case commissionedActionDecommission:
		if string(currentState) == string(commissionedStateDeprovisioned) {
			// 	Since the commissionState changes back to `Provisioned` once the `Decommission` is done, `Provision` is executed instead of `Decommission` when current state is `Deprovisioned`.
			props.CommissionedState = network.CommissionedStateProvisioning
		} else {
			props.CommissionedState = network.CommissionedStateDecommissioning
		}
	case commissionedActionDeprovision:
		if string(currentState) == string(network.CommissionedStateCommissioned) {
			// 	Since the commissionState changes back to `Provisioned` once the `Decommission` is done and `Provisioned` of CommissionedState can do `Deprovision`, `Decommission` is executed first and then do `Deprovision`.
			if updateCommissionedState(d, meta, id, existing, network.CommissionedStateDecommissioning) != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		}
		props.CommissionedState = network.CommissionedStateDeprovisioning
	}
	return nil
}

func prepareDelete(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}, id parse.CustomIpPrefixId) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient

	existing, err := client.Get(ctx, id.ResourceGroup, id.CustomIpPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	currentState := existing.CommissionedState

	// Since the "Provisioning", "Commissioning", "Decommissioning" or `Deprovisioning` commissioned state of custom ip prefix can do nothing, should wait for the operation to complete.
	if currentState == network.CommissionedStateProvisioning || currentState == network.CommissionedStateCommissioning || currentState == network.CommissionedStateDecommissioning || currentState == network.CommissionedStateDeprovisioning {
		return fmt.Errorf("the commissioned state of %s in %s, it doesn't allow to be deleted", id, string(existing.CommissionedState))
	}

	// Since the custom ip prefix can be deleted when the commisionedstate is `Deprovisioned` or `ValidationFailed`, it is needed to update `commisionedstate` before deleting if the state does not allow.
	if string(currentState) != string(commissionedStateDeprovisioned) && string(currentState) != string(commissionedStateValidationFailed) {
		// For the `Commissioned` state should `Decommission` first, then update to `Deprovisioned` after the state changed to `Provided`.
		if currentState == network.CommissionedStateCommissioned {
			if updateCommissionedState(d, meta, id, existing, network.CommissionedStateDecommissioning) != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
		}

		if updateCommissionedState(d, meta, id, existing, network.CommissionedStateDeprovisioning) != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return nil
}

func updateCommissionedState(d *pluginsdk.ResourceData, meta interface{}, id parse.CustomIpPrefixId, existing network.CustomIPPrefix, state network.CommissionedState) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parameters := network.CustomIPPrefix{
		CustomIPPrefixPropertiesFormat: &network.CustomIPPrefixPropertiesFormat{
			Cidr:              existing.Cidr,
			CommissionedState: state,
		},
		Location: existing.Location,
	}

	if v := existing.Tags; v != nil {
		parameters.Tags = v
	}

	if v := existing.Zones; v != nil {
		parameters.Zones = v
	}

	if v := existing.AuthorizationMessage; v != nil {
		parameters.AuthorizationMessage = v
	}

	if v := existing.SignedMessage; v != nil {
		parameters.SignedMessage = v
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.CustomIpPrefixeName, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the updating completion of %s: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{string(state)},
		Target:                    []string{string(network.CommissionedStateProvisioned), string(commissionedStateDeprovisioned)},
		Refresh:                   customIpPrefixCreateRefreshFunc(ctx, client, id),
		PollInterval:              1 * time.Minute,
		ContinuousTargetOccurence: 3,
		Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to become ready for updating: %+v", id, err)
	}

	return nil
}
