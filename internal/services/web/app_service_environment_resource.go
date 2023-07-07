// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	helpersValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	LoadBalancingModeWebPublishing web.LoadBalancingMode = "Web, Publishing"
)

func resourceAppServiceEnvironment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceEnvironmentCreate,
		Read:   resourceAppServiceEnvironmentRead,
		Update: resourceAppServiceEnvironmentUpdate,
		Delete: resourceAppServiceEnvironmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AppServiceEnvironmentID(id)
			return err
		}),

		// Need to find sane values for below, some operations on this resource can take an exceptionally long time
		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: resourceAppServiceEnvironmentSchema(),
	}
}

func resourceAppServiceEnvironmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	networksClient := meta.(*clients.Client).Network.VnetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	internalLoadBalancingMode := d.Get("internal_load_balancing_mode").(string)
	internalLoadBalancingMode = strings.ReplaceAll(internalLoadBalancingMode, " ", "")
	t := d.Get("tags").(map[string]interface{})

	var userWhitelistedIPRangesRaw []interface{}
	if v, ok := d.GetOk("allowed_user_ip_cidrs"); ok {
		userWhitelistedIPRangesRaw = v.(*pluginsdk.Set).List()
	}

	subnetId := d.Get("subnet_id").(string)
	subnet, err := commonids.ParseSubnetID(subnetId)
	if err != nil {
		return err
	}

	// TODO: Remove the implicit behaviour in new major version.
	// Discrepancy of resource group between ASE and Subnet is allowed. While for the sake of
	// compatibility, we still allow user to use the resource group of Subnet to be the one for
	// ASE implicitly. While allow user to explicitly specify the resource group, which takes higher
	// precedence.
	resourceGroup := subnet.ResourceGroupName
	if v, ok := d.GetOk("resource_group_name"); ok {
		resourceGroup = v.(string)
	}
	id := parse.NewAppServiceEnvironmentID(subscriptionId, resourceGroup, d.Get("name").(string))

	vnet, err := networksClient.Get(ctx, subnet.ResourceGroupName, subnet.VirtualNetworkName, "")
	if err != nil {
		return fmt.Errorf("retrieving Virtual Network %q (Resource Group %q): %+v", subnet.VirtualNetworkName, subnet.ResourceGroupName, err)
	}

	// the App Service Environment has to be in the same location as the Virtual Network
	var location string
	if loc := vnet.Location; loc != nil {
		location = azure.NormalizeLocation(*loc)
	} else {
		return fmt.Errorf("determining Location from Virtual Network %q (Resource Group %q): `location` was nil", subnet.VirtualNetworkName, subnet.ResourceGroupName)
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_app_service_environment", id.ID())
	}

	frontEndScaleFactor := d.Get("front_end_scale_factor").(int)
	pricingTier := d.Get("pricing_tier").(string)

	envelope := web.AppServiceEnvironmentResource{
		Location: utils.String(location),
		Kind:     utils.String("ASEV2"),
		AppServiceEnvironment: &web.AppServiceEnvironment{
			InternalLoadBalancingMode: web.LoadBalancingMode(internalLoadBalancingMode),
			FrontEndScaleFactor:       utils.Int32(int32(frontEndScaleFactor)),
			MultiSize:                 utils.String(convertFromIsolatedSKU(pricingTier)),
			VirtualNetwork: &web.VirtualNetworkProfile{
				ID:     utils.String(subnetId),
				Subnet: utils.String(subnet.SubnetName),
			},
			UserWhitelistedIPRanges: utils.ExpandStringSlice(userWhitelistedIPRangesRaw),
		},
		Tags: tags.Expand(t),
	}

	if clusterSettingsRaw, ok := d.GetOk("cluster_setting"); ok {
		envelope.AppServiceEnvironment.ClusterSettings = expandAppServiceEnvironmentClusterSettings(clusterSettingsRaw)
	}

	// whilst this returns a future go-autorest has a max number of retries
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.HostingEnvironmentName, envelope)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	createWait := pluginsdk.StateChangeConf{
		Pending: []string{
			string(web.ProvisioningStateInProgress),
		},
		Target: []string{
			string(web.ProvisioningStateSucceeded),
		},
		MinTimeout: 1 * time.Minute,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
		Refresh:    appServiceEnvironmentRefresh(ctx, client, id.ResourceGroup, id.HostingEnvironmentName),
	}

	// as such we'll ignore it and use a custom poller instead
	if _, err := createWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceEnvironmentRead(d, meta)
}

func resourceAppServiceEnvironmentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	e := web.AppServiceEnvironmentPatchResource{
		AppServiceEnvironment: &web.AppServiceEnvironment{},
	}

	if d.HasChange("internal_load_balancing_mode") {
		v := d.Get("internal_load_balancing_mode").(string)
		v = strings.ReplaceAll(v, " ", "")
		e.AppServiceEnvironment.InternalLoadBalancingMode = web.LoadBalancingMode(v)
	}

	if d.HasChange("front_end_scale_factor") {
		v := d.Get("front_end_scale_factor").(int)
		e.AppServiceEnvironment.FrontEndScaleFactor = utils.Int32(int32(v))
	}

	if d.HasChange("pricing_tier") {
		v := d.Get("pricing_tier").(string)
		v = convertFromIsolatedSKU(v)
		e.AppServiceEnvironment.MultiSize = utils.String(v)
	}

	if d.HasChanges("allowed_user_ip_cidrs") {
		if v, ok := d.GetOk("allowed_user_ip_cidrs"); ok {
			e.UserWhitelistedIPRanges = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
		}
	}

	if d.HasChange("cluster_setting") {
		e.ClusterSettings = expandAppServiceEnvironmentClusterSettings(d.Get("cluster_setting"))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.HostingEnvironmentName, e); err != nil {
		return fmt.Errorf("updating App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	updateWait := pluginsdk.StateChangeConf{
		Pending: []string{
			string(web.ProvisioningStateInProgress),
		},
		Target: []string{
			string(web.ProvisioningStateSucceeded),
		},
		MinTimeout: 1 * time.Minute,
		Timeout:    d.Timeout(pluginsdk.TimeoutUpdate),
		Refresh:    appServiceEnvironmentRefresh(ctx, client, id.ResourceGroup, id.HostingEnvironmentName),
	}

	if _, err := updateWait.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Update of App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	return resourceAppServiceEnvironmentRead(d, meta)
}

func resourceAppServiceEnvironmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.HostingEnvironmentName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			log.Printf("[DEBUG] App Service Environmment %q (Resource Group %q) was not found - removing from state!", id.HostingEnvironmentName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving App Service Environmment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	d.Set("name", id.HostingEnvironmentName)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := existing.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := existing.AppServiceEnvironment; props != nil {
		d.Set("internal_load_balancing_mode", string(props.InternalLoadBalancingMode))

		subnetId := ""
		if props.VirtualNetwork != nil && props.VirtualNetwork.ID != nil {
			subnetId = *props.VirtualNetwork.ID
		}
		d.Set("subnet_id", subnetId)

		frontendScaleFactor := 0
		if props.FrontEndScaleFactor != nil {
			frontendScaleFactor = int(*props.FrontEndScaleFactor)
		}
		d.Set("front_end_scale_factor", frontendScaleFactor)

		pricingTier := ""
		if props.MultiSize != nil {
			pricingTier = convertToIsolatedSKU(*props.MultiSize)
		}
		d.Set("pricing_tier", pricingTier)
		d.Set("allowed_user_ip_cidrs", props.UserWhitelistedIPRanges)
		d.Set("cluster_setting", flattenClusterSettings(props.ClusterSettings))
	}

	// Get IP attributes for ASE.
	vipInfo, err := client.GetVipInfo(ctx, id.ResourceGroup, id.HostingEnvironmentName)
	if err != nil {
		if utils.ResponseWasNotFound(vipInfo.Response) {
			return fmt.Errorf("retrieving VIP info: App Service Environment %q (Resource Group %q) was not found", id.HostingEnvironmentName, id.ResourceGroup)
		}
		return fmt.Errorf("retrieving VIP info App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	d.Set("internal_ip_address", vipInfo.InternalIPAddress)
	d.Set("service_ip_address", vipInfo.ServiceIPAddress)
	d.Set("outbound_ip_addresses", vipInfo.OutboundIPAddresses)

	return tags.FlattenAndSet(d, existing.Tags)
}

func resourceAppServiceEnvironmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Environment %q (Resource Group %q)", id.HostingEnvironmentName, id.ResourceGroup)

	forceDeleteAllChildren := utils.Bool(false)
	future, err := client.Delete(ctx, id.ResourceGroup, id.HostingEnvironmentName, forceDeleteAllChildren)
	if err != nil {
		return fmt.Errorf("deleting App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
		}
	}

	return nil
}

func appServiceEnvironmentRefresh(ctx context.Context, client *web.AppServiceEnvironmentsClient, resourceGroup string, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		read, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return "", "", err
		}

		if read.AppServiceEnvironment == nil {
			return "", "", fmt.Errorf("`properties` was nil")
		}

		state := read.AppServiceEnvironment.ProvisioningState
		return state, string(state), nil
	}
}

// Note: These are abstractions and possibly subject to change if Azure changes the underlying SKU for Isolated instances.
func convertFromIsolatedSKU(isolated string) (vmSKU string) {
	switch isolated {
	case "I1":
		vmSKU = "Standard_D1_V2"
	case "I2":
		vmSKU = "Standard_D2_V2"
	case "I3":
		vmSKU = "Standard_D3_V2"
	}
	return vmSKU
}

func convertToIsolatedSKU(vmSKU string) (isolated string) {
	switch vmSKU {
	case "Standard_D1_V2":
		isolated = "I1"
	case "Standard_D2_V2":
		isolated = "I2"
	case "Standard_D3_V2":
		isolated = "I3"
	}
	return isolated
}

func loadBalancingModeDiffSuppress(k, old, new string, d *pluginsdk.ResourceData) bool {
	return strings.ReplaceAll(old, " ", "") == strings.ReplaceAll(new, " ", "")
}

func expandAppServiceEnvironmentClusterSettings(input interface{}) *[]web.NameValuePair {
	var clusterSettings []web.NameValuePair
	if input == nil {
		return &clusterSettings
	}

	clusterSettingsRaw := input.([]interface{})
	for _, v := range clusterSettingsRaw {
		setting := v.(map[string]interface{})
		clusterSettings = append(clusterSettings, web.NameValuePair{
			Name:  utils.String(setting["name"].(string)),
			Value: utils.String(setting["value"].(string)),
		})
	}
	return &clusterSettings
}

func flattenClusterSettings(input *[]web.NameValuePair) interface{} {
	if input == nil || len(*input) == 0 {
		return []map[string]interface{}{}
	}

	settings := make([]map[string]interface{}, 0)
	for _, v := range *input {
		name := ""
		if v.Name == nil {
			continue
		} else {
			name = *v.Name
		}

		value := ""
		if v.Value != nil {
			value = *v.Value
		}

		settings = append(settings, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}
	return settings
}

func resourceAppServiceEnvironmentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AppServiceEnvironmentName,
		},

		"subnet_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSubnetID,
		},

		"cluster_setting": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"internal_load_balancing_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(web.LoadBalancingModeNone),
			ValidateFunc: validation.StringInSlice([]string{
				string(web.LoadBalancingModeNone),
				string(web.LoadBalancingModePublishing),
				string(web.LoadBalancingModeWeb),
				string(web.LoadBalancingModeWebPublishing),
				// (@jackofallops) breaking change in SDK - Enum for internal_load_balancing_mode changed from Web, Publishing to Web,Publishing
				string(LoadBalancingModeWebPublishing),
			}, false),
			DiffSuppressFunc: loadBalancingModeDiffSuppress,
		},

		"front_end_scale_factor": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      15,
			ValidateFunc: validation.IntBetween(5, 15),
		},

		"pricing_tier": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "I1",
			ValidateFunc: validation.StringInSlice([]string{
				"I1",
				"I2",
				"I3",
			}, false),
		},

		"allowed_user_ip_cidrs": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: helpersValidate.CIDR,
			},
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"tags": tags.ForceNewSchema(),

		// Computed

		// VipInfo
		"internal_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"outbound_ip_addresses": {
			Type: pluginsdk.TypeList,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Computed: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
