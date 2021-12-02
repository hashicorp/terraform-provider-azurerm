package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var appServiceResourceName = "azurerm_app_service"

func resourceAppServiceIpRestriction() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceIpRestrictionCreate,
		Read:   resourceAppServiceIpRestrictionRead,
		Update: resourceAppServiceIpRestrictionUpdate,
		Delete: resourceAppServiceIpRestrictionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"app_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"ip_restriction": {
				Type:       pluginsdk.TypeList,
				Required:   true,
				MinItems:   1,
				MaxItems:   1,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Elem: &pluginsdk.Resource{
					Schema: schemaAppServiceIpRestrictionElement(),
				},
			},
		},
	}
}

func resourceAppServiceIpRestrictionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appServiceId := d.Get("app_service_id").(string)
	id, err := parse.AppServiceID(appServiceId)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("error checking for presence of existing App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
		}
	}

	if resp.SiteConfig == nil || resp.SiteConfig.IPSecurityRestrictions == nil {
		return fmt.Errorf("failed reading IP Restrictions for %q (resource group %q)", id.SiteName, id.ResourceGroup)
	}

	// Locking to prevent parallel changes causing issues
	locks.ByName(id.SiteName, appServiceResourceName)
	defer locks.UnlockByName(id.SiteName, appServiceResourceName)

	ipRestrictionArr, err := expandAppServiceIpRestriction(d.Get("ip_restriction"))
	if err != nil {
		return err
	}

	name := ipRestrictionArr[0].Name
	if name == nil || *name == "" {
		return fmt.Errorf("no name specified for IP restriction for App Service %q (Resource Group %q)", id.SiteName, id.ResourceGroup)
	}

	// This is because azure doesn't have an 'id' for single app service ip restriction
	// In order to compensate for this and allow importing of this resource we are artificially
	// creating an identity for an app service ip restriction object
	// /subscriptions/<guid>/resourceGroups/<rg-name>/providers/Microsoft.Web/sites/<site-name>/ipRestriction/<restriction-name>
	resourceId := fmt.Sprintf("%s/ipRestriction/%s", *resp.ID, *name)
	_, ipRestriction := FindIPRestriction(resp.SiteConfig.IPSecurityRestrictions, *name)
	if ipRestriction != nil {
		return tf.ImportAsExistsError("azurerm_app_service_ip_restriction", resourceId)
	}

	restrictions := append(*resp.SiteConfig.IPSecurityRestrictions, ipRestrictionArr...)
	resp.SiteConfig.IPSecurityRestrictions = &restrictions

	siteConfigResource := web.SiteConfigResource{
		SiteConfig: resp.SiteConfig,
	}

	if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.SiteName, siteConfigResource); err != nil {
		return fmt.Errorf("updating Configuration for App Service %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	d.SetId(resourceId)
	return resourceAppServiceIpRestrictionRead(d, meta)
}

func resourceAppServiceIpRestrictionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	siteName := id.Path["sites"]
	restrictionName := id.Path["ipRestriction"]

	resp, err := client.Get(ctx, id.ResourceGroup, siteName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("error checking for presence of existing App Service %q (Resource Group %q): %s", siteName, id.ResourceGroup, err)
		}
	}

	if resp.SiteConfig == nil || resp.SiteConfig.IPSecurityRestrictions == nil {
		return fmt.Errorf("failed reading IP Restrictions for %q (resource group %q)", siteName, id.ResourceGroup)
	}

	_, restriction := FindIPRestriction(resp.SiteConfig.IPSecurityRestrictions, restrictionName)

	if restriction == nil {
		log.Printf("[INFO] IP Restriction %q was not found in App Service %q (Resource Group %q) - removing from state", restrictionName, siteName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	restrictionArr := []web.IPSecurityRestriction{*restriction}
	appServiceIpRestriction := flattenAppServiceIpRestriction(&restrictionArr)
	if len(appServiceIpRestriction) != 1 {
		return fmt.Errorf("failed to flatten IP Restriction %q for App Service %q (resource group %q)", restrictionName, siteName, id.ResourceGroup)
	}
	d.Set("ip_restriction", appServiceIpRestriction)

	return nil
}

func resourceAppServiceIpRestrictionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appServiceId := d.Get("app_service_id").(string)
	id, err := parse.AppServiceID(appServiceId)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("error checking for presence of existing App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
		}
	}

	if resp.SiteConfig == nil || resp.SiteConfig.IPSecurityRestrictions == nil {
		return fmt.Errorf("failed reading IP Restrictions for %q (resource group %q)", id.SiteName, id.ResourceGroup)
	}

	// Locking to prevent parallel changes causing issues
	locks.ByName(id.SiteName, appServiceResourceName)
	defer locks.UnlockByName(id.SiteName, appServiceResourceName)

	if d.HasChange("ip_restriction") {

		ipRestrictionArr, err := expandAppServiceIpRestriction(d.Get("ip_restriction"))
		if err != nil {
			return err
		}

		name := ipRestrictionArr[0].Name
		if name == nil || *name == "" {
			return fmt.Errorf("no name specified for IP restriction for App Service %q (Resource Group %q)", id.SiteName, id.ResourceGroup)
		}

		idx, _ := FindIPRestriction(resp.SiteConfig.IPSecurityRestrictions, *name)

		if idx < 0 {
			d.SetId("")
			return fmt.Errorf(" IP Restriction %q was not found in App Service %q (Resource Group %q) - removing from state", *name, id.SiteName, id.ResourceGroup)
		}

		(*resp.SiteConfig.IPSecurityRestrictions)[idx] = ipRestrictionArr[0]

		siteConfigResource := web.SiteConfigResource{
			SiteConfig: resp.SiteConfig,
		}

		if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.SiteName, siteConfigResource); err != nil {
			return fmt.Errorf("updating Configuration for App Service %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
		}
	}

	return nil
}

func resourceAppServiceIpRestrictionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	appServiceId := d.Get("app_service_id").(string)
	id, err := parse.AppServiceID(appServiceId)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SiteName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("error checking for presence of existing App Service %q (Resource Group %q): %s", id.SiteName, id.ResourceGroup, err)
		}
	}

	if resp.SiteConfig == nil || resp.SiteConfig.IPSecurityRestrictions == nil {
		return fmt.Errorf("failed reading IP Restrictions for %q (resource group %q)", id.SiteName, id.ResourceGroup)
	}

	// Locking to prevent parallel changes causing issues
	locks.ByName(id.SiteName, appServiceResourceName)
	defer locks.UnlockByName(id.SiteName, appServiceResourceName)

	ipRestrictionArr, err := expandAppServiceIpRestriction(d.Get("ip_restriction"))
	if err != nil {
		return err
	}

	name := ipRestrictionArr[0].Name
	if name == nil || *name == "" {
		return fmt.Errorf("no name specified for IP restriction for App Service %q (Resource Group %q)", id.SiteName, id.ResourceGroup)
	}

	restrictions, itemToRemove := removeIPRestriction(resp.SiteConfig.IPSecurityRestrictions, *name)

	if itemToRemove == nil {
		log.Printf("[INFO] IP Restriction %q was not found in App Service %q (Resource Group %q) - removing from state", *name, id.SiteName, id.ResourceGroup)
		d.SetId("")
		return nil
	}

	resp.SiteConfig.IPSecurityRestrictions = restrictions
	siteConfigResource := web.SiteConfigResource{
		SiteConfig: resp.SiteConfig,
	}

	if _, err := client.CreateOrUpdateConfiguration(ctx, id.ResourceGroup, id.SiteName, siteConfigResource); err != nil {
		return fmt.Errorf("updating Configuration for App Service %q (Resource Group %q): %+v", id.SiteName, id.ResourceGroup, err)
	}

	d.SetId("")
	return nil
}

func removeIPRestriction(restrictions *[]web.IPSecurityRestriction, name string) (*[]web.IPSecurityRestriction, *web.IPSecurityRestriction) {
	if restrictions == nil || len(*restrictions) == 0 {
		return nil, nil
	}
	for i, item := range *restrictions {
		if item.Name != nil && strings.EqualFold(*item.Name, name) {
			arr := append((*restrictions)[:i], (*restrictions)[i+1:]...)
			return &arr, &item
		}
	}
	return restrictions, nil
}

func FindIPRestriction(restrictions *[]web.IPSecurityRestriction, name string) (int, *web.IPSecurityRestriction) {
	if restrictions == nil || len(*restrictions) == 0 {
		return -1, nil
	}
	for idx, item := range *restrictions {
		if item.Name != nil && strings.EqualFold(*item.Name, name) {
			return idx, &item
		}
	}
	return -1, nil
}
