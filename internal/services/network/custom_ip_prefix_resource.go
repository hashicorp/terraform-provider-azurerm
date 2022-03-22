package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
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
		Create: resourceCustomIpPrefixCreateUpdate,
		Read:   resourceCustomIpPrefixRead,
		Update: resourceCustomIpPrefixCreateUpdate,
		Delete: resourceCustomIpPrefixDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CustomIpPrefixID(id)
			return err
		}),

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

			"zones": commonschema.ZonesMultipleRequiredForceNew(),

			"authorization_message": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"signed_message": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"action": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(commissionedActionProvision),
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

func resourceCustomIpPrefixCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient
	subscriptionId := meta.(*clients.Client).Network.CustomIPPrefixesClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewCustomIpPrefixID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	state := commissionedAction(d.Get("action").(string))

	existing, err := client.Get(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), "")
	if d.IsNewResource() {
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_custom_ip_prefix", id.ID())
		}
	} else {
		switch string(existing.CommissionedState) {
		case string(network.CommissionedStateProvisioned):
			if state != commissionedActionCommission && state != commissionedActionDeprovision {
				return fmt.Errorf("%s can do `Commission` or `Deprovision` when the commissioned state is `Provisioned`", id)
			}
		case string(network.CommissionedStateCommissioned):
			if state != commissionedActionDecommission {
				return fmt.Errorf("%s can do `Decommission` when the commissioned state is `Commissioned`", id)
			}
		case string(commissionedStateDeprovisioned):
			if state != commissionedActionProvision {
				return fmt.Errorf("%s can do delete or `Provision` when the commissioned state is `Deprovisioned`", id)
			}
		}
	}

	parameters := network.CustomIPPrefix{
		CustomIPPrefixPropertiesFormat: &network.CustomIPPrefixPropertiesFormat{
			Cidr: utils.String(d.Get("cidr").(string)),
		},
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("action"); ok {
		var state network.CommissionedState
		switch commissionedAction(v.(string)) {
		case commissionedActionProvision:
			state = network.CommissionedStateProvisioning
		case commissionedActionCommission:
			state = network.CommissionedStateCommissioning
		case commissionedActionDecommission:
			state = network.CommissionedStateDecommissioning
		case commissionedActionDeprovision:
			state = network.CommissionedStateDeprovisioning
		}
		parameters.CommissionedState = state
	}

	if v, ok := d.GetOk("authorization_message"); ok {
		parameters.AuthorizationMessage = utils.String(v.(string))
	}

	if v, ok := d.GetOk("signed_message"); ok {
		parameters.SignedMessage = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zones"); ok {
		zones := zones.Expand(v.(*schema.Set).List())
		if len(zones) > 0 {
			parameters.Zones = &zones
		}
	}

	future, err := client.CreateOrUpdate(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), parameters)
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
		Refresh:                   customIpPrefixCreateRefreshFunc(ctx, client, id),
		PollInterval:              1 * time.Minute,
		ContinuousTargetOccurence: 3,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
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
		d.Set("authorization_message", props.AuthorizationMessage)
		d.Set("signed_message", props.SignedMessage)
	}

	if v := resp.Zones; v != nil {
		d.Set("zones", zones.Flatten(v))
	}

	// The "action" will be changed by API automatically, hence setting it from config.
	d.Set("action", d.Get("action").(string))

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

	if err := prepareDelete(d, meta, *id); err != nil {
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

func prepareDelete(d *pluginsdk.ResourceData, meta interface{}, id parse.CustomIpPrefixId) error {
	client := meta.(*clients.Client).Network.CustomIPPrefixesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	existing, err := client.Get(ctx, id.ResourceGroup, id.CustomIpPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	state := existing.CommissionedState

	// Since the "Provisioning", "Commissioning", "Decommissioning" or `Deprovisioning` commissioned state of custom ip prefix can do nothing, should wait for the operation to complete.
	if state == network.CommissionedStateProvisioning || state == network.CommissionedStateCommissioning || state == network.CommissionedStateDecommissioning || state == network.CommissionedStateDeprovisioning {
		return fmt.Errorf("the commissioned state of %s in %s, it doesn't allow to be deleted", id, string(existing.CommissionedState))
	}

	// Since the custom ip prefix can be deleted when the commisionedstate is `Deprovisioned` or `ValidationFailed`, it is needed to update `commisionedstate` before deleting if the state does not allow.
	if string(state) != string(commissionedStateDeprovisioned) && string(state) != string(commissionedStateValidationFailed) {
		// For the `Commissioned` state should `Decommission` first, then update to `Deprovisioned` after the state changed to `Provided`.
		if state == network.CommissionedStateCommissioned {
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

	future, err := client.CreateOrUpdate(ctx, d.Get("resource_group_name").(string), d.Get("name").(string), parameters)
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

	// The "action" will be changed by API automatically, hence setting it from config.
	d.Set("action", d.Get("action").(string))

	return nil
}
