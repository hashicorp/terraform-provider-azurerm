package keyvault

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2021-10-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKeyVaultManagedHardwareSecurityModule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmKeyVaultManagedHardwareSecurityModuleCreate,
		Read:   resourceArmKeyVaultManagedHardwareSecurityModuleRead,
		Delete: resourceArmKeyVaultManagedHardwareSecurityModuleDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := managedhsms.ParseManagedHSMID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ManagedHardwareSecurityModuleName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(managedhsms.ManagedHsmSkuNameStandardBOne),
				}, false),
			},

			"admin_object_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"purge_protection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"soft_delete_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      90,
				ValidateFunc: validation.IntBetween(7, 90),
			},

			"hsm_uri": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				//Computed: true,
				Default:  true,
				ForceNew: true,
			},

			"network_acls": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"default_action": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(managedhsms.NetworkRuleActionAllow),
								string(managedhsms.NetworkRuleActionDeny),
							}, false),
						},
						"bypass": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(managedhsms.NetworkRuleBypassOptionsNone),
								string(managedhsms.NetworkRuleBypassOptionsAzureServices),
							}, false),
						},
					},
				},
			},

			// https://github.com/Azure/azure-rest-api-specs/issues/13365
			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceArmKeyVaultManagedHardwareSecurityModuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := managedhsms.NewManagedHSMID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_key_vault_managed_hardware_security_module", id.ID())
	}

	publicNetworkAccessEnabled := managedhsms.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccessEnabled = managedhsms.PublicNetworkAccessDisabled
	}
	hsm := managedhsms.ManagedHsm{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: &managedhsms.ManagedHsmProperties{
			InitialAdminObjectIds:     utils.ExpandStringSlice(d.Get("admin_object_ids").(*pluginsdk.Set).List()),
			CreateMode:                pointer.To(managedhsms.CreateModeDefault),
			EnableSoftDelete:          utils.Bool(true),
			SoftDeleteRetentionInDays: utils.Int64(int64(d.Get("soft_delete_retention_days").(int))),
			EnablePurgeProtection:     utils.Bool(d.Get("purge_protection_enabled").(bool)),
			PublicNetworkAccess:       pointer.To(publicNetworkAccessEnabled),
			NetworkAcls:               expandMHSMNetworkAcls(d.Get("network_acls").([]interface{})),
		},
		Sku: &managedhsms.ManagedHsmSku{
			Family: managedhsms.ManagedHsmSkuFamilyB,
			Name:   managedhsms.ManagedHsmSkuName(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if tenantId := d.Get("tenant_id").(string); tenantId != "" {
		hsm.Properties.TenantId = pointer.To(tenantId)
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, hsm); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmKeyVaultManagedHardwareSecurityModuleRead(d, meta)
}

func resourceArmKeyVaultManagedHardwareSecurityModuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[ERROR] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ManagedHSMName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			tenantId := ""
			if props.TenantId != nil {
				tenantId = *props.TenantId
			}
			d.Set("tenant_id", tenantId)
			d.Set("admin_object_ids", utils.FlattenStringSlice(props.InitialAdminObjectIds))
			d.Set("hsm_uri", props.HsmUri)
			d.Set("soft_delete_retention_days", props.SoftDeleteRetentionInDays)
			d.Set("purge_protection_enabled", props.EnablePurgeProtection)

			publicAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess != managedhsms.PublicNetworkAccessEnabled {
				publicAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicAccessEnabled)

			if err := d.Set("network_acls", flattenMHSMNetworkAcls(props.NetworkAcls)); err != nil {
				return fmt.Errorf("setting `network_acls`: %+v", err)
			}
		}

		skuName := ""
		if sku := model.Sku; sku != nil {
			skuName = string(sku.Name)
		}
		d.Set("sku_name", skuName)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceArmKeyVaultManagedHardwareSecurityModuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).KeyVault.ManagedHsmClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedhsms.ParseManagedHSMID(d.Id())
	if err != nil {
		return err
	}

	// We need to grab the keyvault hsm to see if purge protection is enabled prior to deletion
	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	loc := ""
	purgeProtectionEnabled := false
	if model := resp.Model; model != nil {
		loc = location.NormalizeNilable(model.Location)
		if props := model.Properties; props != nil && props.EnablePurgeProtection != nil {
			purgeProtectionEnabled = *props.EnablePurgeProtection
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if meta.(*clients.Client).Features.KeyVault.PurgeSoftDeletedHSMsOnDestroy {
		if purgeProtectionEnabled {
			log.Printf("[DEBUG] cannot purge %s because purge protection is enabled", id)
			return nil
		}
	}

	purgedId := managedhsms.NewDeletedManagedHSMID(id.SubscriptionId, loc, id.ManagedHSMName)
	if err := client.PurgeDeletedThenPoll(ctx, purgedId); err != nil {
		return fmt.Errorf("purging %s: %+v", id, err)
	}

	return nil
}

func expandMHSMNetworkAcls(input []interface{}) *managedhsms.MHSMNetworkRuleSet {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	res := &managedhsms.MHSMNetworkRuleSet{
		Bypass:        pointer.To(managedhsms.NetworkRuleBypassOptions(v["bypass"].(string))),
		DefaultAction: pointer.To(managedhsms.NetworkRuleAction(v["default_action"].(string))),
	}

	return res
}

func flattenMHSMNetworkAcls(acl *managedhsms.MHSMNetworkRuleSet) []interface{} {
	bypass := string(managedhsms.NetworkRuleBypassOptionsAzureServices)
	defaultAction := string(managedhsms.NetworkRuleActionAllow)

	if acl != nil {
		if acl.Bypass != nil {
			bypass = string(*acl.Bypass)
		}
		if acl.DefaultAction != nil {
			defaultAction = string(*acl.DefaultAction)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"bypass":         bypass,
			"default_action": defaultAction,
		},
	}
}
