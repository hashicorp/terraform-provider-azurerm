package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	pricings_v2023_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/security/2023-01-01/pricings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSecurityCenterSubscriptionPricing() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterSubscriptionPricingUpdate,
		Read:   resourceSecurityCenterSubscriptionPricingRead,
		Update: resourceSecurityCenterSubscriptionPricingUpdate,
		Delete: resourceSecurityCenterSubscriptionPricingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PricingID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SubscriptionPricingV0ToV1{},
		}),

		Schema: map[string]*pluginsdk.Schema{
			"tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(security.PricingTierFree),
					string(security.PricingTierStandard),
				}, false),
			},
			"resource_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "VirtualMachines",
				ValidateFunc: validation.StringInSlice([]string{
					"AppServices",
					"ContainerRegistry",
					"KeyVaults",
					"KubernetesService",
					"SqlServers",
					"SqlServerVirtualMachines",
					"StorageAccounts",
					"VirtualMachines",
					"Arm",
					"Dns",
					"OpenSourceRelationalDatabases",
					"Containers",
					"CosmosDbs",
					"CloudPosture",
				}, false),
			},
			"subplan": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSecurityCenterSubscriptionPricingUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: add a requires import check ensuring this is != Free (meaning we should likely remove Free as a SKU option?)

	id := pricings_v2023_01_01.NewPricingID(subscriptionId, d.Get("resource_type").(string))
	pricing := pricings_v2023_01_01.Pricing{
		Properties: &pricings_v2023_01_01.PricingProperties{
			PricingTier: pricings_v2023_01_01.PricingTier(d.Get("tier").(string)),
		},
	}

	if v, ok := d.GetOk("subplan"); ok {
		pricing.Properties.SubPlan = utils.String(v.(string))
	}

	if _, err := client.Update(ctx, id, pricing); err != nil {
		return fmt.Errorf("setting %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSecurityCenterSubscriptionPricingRead(d, meta)
}

func resourceSecurityCenterSubscriptionPricingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.PricingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := pricings_v2023_01_01.ParsePricingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("resource_type", id.PricingName)
	if resp.Model != nil {
		if properties := resp.Model.Properties; properties != nil {
			d.Set("tier", properties.PricingTier)
			d.Set("subplan", properties.SubPlan)
		}
	}

	return nil
}

func resourceSecurityCenterSubscriptionPricingDelete(_ *pluginsdk.ResourceData, _ interface{}) error {
	// TODO: reset this back to Free

	log.Printf("[DEBUG] Security Center Subscription deletion invocation")
	return nil
}
