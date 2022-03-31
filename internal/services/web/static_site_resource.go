package web

// This file is a workaround for issue https://github.com/Azure/azure-rest-api-specs/issues/18501. This can be removed once it is resolved.

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceStaticSite() *pluginsdk.Resource {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.StaticSiteName,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"default_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"api_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.Schema(),
	}

	if features.FourPointOhBeta() {
		schema["sku_name"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(web.SkuNameFree),
			ValidateFunc: validation.StringInSlice([]string{
				string(web.SkuNameStandard),
				string(web.SkuNameFree),
			}, false),
		}
	} else {
		schema["sku_size"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(web.SkuNameFree),
			ValidateFunc: validation.StringInSlice([]string{
				string(web.SkuNameStandard),
				string(web.SkuNameFree),
			}, false),
		}
		schema["sku_tier"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(web.SkuNameFree),
			ValidateFunc: validation.StringInSlice([]string{
				string(web.SkuNameStandard),
				string(web.SkuNameFree),
			}, false),
		}
	}

	return &pluginsdk.Resource{
		Create: resourceStaticSiteCreate,
		Read:   resourceStaticSiteRead,
		Update: resourceStaticSiteUpdate,
		Delete: resourceStaticSiteDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.StaticSiteID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: schema,
	}
}

func resourceStaticSiteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Static Site creation.")

	id := parse.NewStaticSiteID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.GetStaticSite(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("failed checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_static_site", id.ID())
	}

	loc := location.Normalize(d.Get("location").(string))

	var skuName string
	if features.FourPointOhBeta() {
		skuName = d.Get("sku_name").(string)
	} else {
		skuName = d.Get("sku_size").(string)
	}

	identity, err := expandStaticSiteIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	// See: https://github.com/Azure/azure-rest-api-specs/issues/17525
	if skuName == string(web.SkuNameFree) && identity != nil {
		return fmt.Errorf("a Managed Identity cannot be used when tier is set to `Free`")
	}

	siteEnvelope := web.StaticSiteARMResource{
		Sku: &web.SkuDescription{
			Name: &skuName,
		},
		StaticSite: &web.StaticSite{},
		Location:   &loc,
		Identity:   identity,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if features.FourPointOhBeta() {
		siteEnvelope.Sku.Tier = utils.String(skuName)
	} else {
		siteEnvelope.Sku.Tier = utils.String(d.Get("sku_tier").(string))
	}

	future, err := client.CreateOrUpdateStaticSite(ctx, id.ResourceGroup, id.Name, siteEnvelope)
	if err != nil {
		return fmt.Errorf("failed creating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed waiting for creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceStaticSiteRead(d, meta)
}

func resourceStaticSiteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StaticSiteID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.GetStaticSite(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving existing %q: %+v", id, err)
	}

	var skuName string
	if features.FourPointOhBeta() {
		skuName = d.Get("sku_name").(string)
	} else {
		skuName = d.Get("sku_size").(string)
	}

	if d.HasChange("identity") {
		identity, err := expandStaticSiteIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		// In case the expanded identity is nil, it means the users have removed the identity block from the config. We need to explicitly set the identity type to "None" here
		// (as ignoring the identity field in the API request will ends up not changing it at all).
		// Untill issue https://github.com/Azure/azure-rest-api-specs/issues/17525 is addressed, we can put this logic back to the expand function instead.
		if identity == nil && skuName != string(web.SkuNameFree) {
			identity = &web.ManagedServiceIdentity{Type: web.ManagedServiceIdentityTypeNone}
		}
		existing.Identity = identity
	}

	if existing.Sku != nil {
		if features.FourPointOhBeta() {
			if d.HasChange("sku_name") {
				existing.Sku.Name = utils.String(skuName)
				existing.Sku.Tier = utils.String(skuName)
			}
		} else {
			if d.HasChange("sku_tier") {
				existing.Sku.Tier = utils.String(d.Get("sku_tier").(string))
			}
			if d.HasChange("sku_size") {
				existing.Sku.Name = &skuName
			}
		}
	}

	if d.HasChange("tags") {
		existing.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := azuresdkhacks.CreateOrUpdateStaticSite(ctx, client, id.ResourceGroup, id.Name, existing)
	if err != nil {
		return fmt.Errorf("failed updating %s: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed waiting for updating %s: %+v", id, err)
	}

	return resourceStaticSiteRead(d, meta)
}

func resourceStaticSiteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StaticSiteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetStaticSite(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Static Site %q (resource group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed making Read request on %s: %+v", id, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	identity, err := flattenStaticSiteIdentity(resp.Identity)
	if err != nil {
		return err
	}
	d.Set("identity", identity)

	d.Set("location", location.NormalizeNilable(resp.Location))

	if prop := resp.StaticSite; prop != nil {
		defaultHostname := ""
		if prop.DefaultHostname != nil {
			defaultHostname = *prop.DefaultHostname
		}
		d.Set("default_host_name", defaultHostname)
	}

	if features.FourPointOhBeta() {
		skuName := ""
		skuTier := ""
		if sku := resp.Sku; sku != nil {
			if v := sku.Name; v != nil {
				skuName = *v
			}

			if v := sku.Tier; v != nil {
				skuTier = *v
			}
		}
		d.Set("sku_size", skuName)
		d.Set("sku_tier", skuTier)
	} else {
		skuName := ""
		if sku := resp.Sku; sku != nil {
			if v := sku.Name; v != nil {
				skuName = *v
			}
		}
		d.Set("sku_name", skuName)
	}

	secretResp, err := client.ListStaticSiteSecrets(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("listing secretes for %s: %v", id, err)
	}

	apiKey := ""
	if pkey := secretResp.Properties["apiKey"]; pkey != nil {
		apiKey = *pkey
	}
	d.Set("api_key", apiKey)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceStaticSiteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.StaticSitesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.StaticSiteID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting Static Site %q (resource group %q)", id.Name, id.ResourceGroup)

	future, err := client.DeleteStaticSite(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("failed deleting %q: %v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("failed waiting for deleting %q: %v", id, err)
	}

	return nil
}

func expandStaticSiteIdentity(input []interface{}) (*web.ManagedServiceIdentity, error) {
	config, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	if config.Type == identity.TypeNone {
		return nil, nil
	}

	out := web.ManagedServiceIdentity{
		Type: web.ManagedServiceIdentityType(config.Type),
	}

	if len(config.IdentityIds) != 0 {
		out.UserAssignedIdentities = make(map[string]*web.UserAssignedIdentity)
		for id := range config.IdentityIds {
			out.UserAssignedIdentities[id] = &web.UserAssignedIdentity{}
		}
	}

	return &out, nil
}

func flattenStaticSiteIdentity(input *web.ManagedServiceIdentity) (*[]interface{}, error) {
	var transform *identity.SystemAndUserAssignedMap

	if input != nil {
		transform = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}

		if input.PrincipalID != nil {
			transform.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transform.TenantId = *input.TenantID
		}
		for k, v := range input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transform)
}
