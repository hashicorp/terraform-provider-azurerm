package web

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2020-06-01/web"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	helpersValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	networkParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const (
	LoadBalancingModeWebPublishing web.LoadBalancingMode = "Web, Publishing"
)

func resourceAppServiceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppServiceEnvironmentCreate,
		Read:   resourceAppServiceEnvironmentRead,
		Update: resourceAppServiceEnvironmentUpdate,
		Delete: resourceAppServiceEnvironmentDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppServiceEnvironmentID(id)
			return err
		}),

		// Need to find sane values for below, some operations on this resource can take an exceptionally long time
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(6 * time.Hour),
			Delete: schema.DefaultTimeout(6 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AppServiceEnvironmentName,
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
			},

			"cluster_setting": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"internal_load_balancing_mode": {
				Type:     schema.TypeString,
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
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntBetween(5, 15),
			},

			// TODO - Not allowed in V3, but a value is returned
			"pricing_tier": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"I1",
					"I2",
					"I3",
				}, false),
				ConflictsWith: []string{
					"version",
				},
			},

			// TODO - Not allowed in V3
			"allowed_user_ip_cidrs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // remove in 3.0
				ConflictsWith: []string{
					"user_whitelisted_ip_ranges",
					"version",
				},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: helpersValidate.CIDR,
				},
			},

			// TODO - Not allowed in V3
			"user_whitelisted_ip_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true, // remove in 3.0
				ConflictsWith: []string{
					"allowed_user_ip_cidrs",
					"version",
				},
				Deprecated: "this property has been renamed to `allowed_user_ip_cidrs` better reflect the expected ip range format",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: helpersValidate.CIDR,
				},
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ASEV2",
				ValidateFunc: validation.StringInSlice([]string{
					"ASEV2",
					"ASEV3",
				}, false),
				ConflictsWith: []string{
					"pricing_tier",
					"allowed_user_ip_cidrs",
					"user_whitelisted_ip_ranges",
				},
			},

			// TODO in 3.0 Make it "Required"
			"resource_group_name": azure.SchemaResourceGroupNameOptionalComputed(),

			"tags": tags.ForceNewSchema(),

			// Computed
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppServiceEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	networksClient := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	internalLoadBalancingMode := d.Get("internal_load_balancing_mode").(string)
	internalLoadBalancingMode = strings.ReplaceAll(internalLoadBalancingMode, " ", "")
	t := d.Get("tags").(map[string]interface{})
	userWhitelistedIPRangesRaw := d.Get("user_whitelisted_ip_ranges").(*schema.Set).List()
	if v, ok := d.GetOk("allowed_user_ip_cidrs"); ok {
		userWhitelistedIPRangesRaw = v.(*schema.Set).List()
	}

	subnetId := d.Get("subnet_id").(string)
	subnet, err := networkParse.SubnetID(subnetId)
	if err != nil {
		return err
	}

	// TODO: Remove the implicit behaviour in new major version.
	// Discrepancy of resource group between ASE and Subnet is allowed. While for the sake of
	// compatibility, we still allow user to use the resource group of Subnet to be the one for
	// ASE implicitly. While allow user to explicitly specify the resource group, which takes higher
	// precedence.
	resourceGroup := subnet.ResourceGroup
	if v, ok := d.GetOk("resource_group_name"); ok {
		resourceGroup = v.(string)
	}

	vnet, err := networksClient.Get(ctx, subnet.ResourceGroup, subnet.VirtualNetworkName, "")
	if err != nil {
		return fmt.Errorf("retrieving Virtual Network %q (Resource Group %q): %+v", subnet.VirtualNetworkName, subnet.ResourceGroup, err)
	}

	// the App Service Environment has to be in the same location as the Virtual Network
	var location string
	if loc := vnet.Location; loc != nil {
		location = azure.NormalizeLocation(*loc)
	} else {
		return fmt.Errorf("determining Location from Virtual Network %q (Resource Group %q): `location` was nil", subnet.VirtualNetworkName, subnet.ResourceGroup)
	}

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing App Service Environment %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_app_service_environment", *existing.ID)
	}

	frontEndScaleFactor := d.Get("front_end_scale_factor").(int)
	pricingTier := d.Get("pricing_tier").(string)
	kind := d.Get("version").(string)
	if kind == "ASEV2" && pricingTier == "" {
		pricingTier = "I1"
	}

	envelope := web.AppServiceEnvironmentResource{
		Location: utils.String(location),
		Kind:     utils.String(kind),
		AppServiceEnvironment: &web.AppServiceEnvironment{
			Name:                      utils.String(name),
			Location:                  utils.String(location),
			InternalLoadBalancingMode: web.LoadBalancingMode(internalLoadBalancingMode),
			FrontEndScaleFactor:       utils.Int32(int32(frontEndScaleFactor)),
			VirtualNetwork: &web.VirtualNetworkProfile{
				ID:     utils.String(subnetId),
				Subnet: utils.String(subnet.Name),
			},
			// the SDK is coded primarily for v1, which needs a non-null entry for workerpool, so we construct an empty slice for it
			// TODO: remove this hack once https://github.com/Azure/azure-rest-api-specs/pull/8433 has been merged
			WorkerPools: &[]web.WorkerPool{{}},
		},
		Tags: tags.Expand(t),
	}

	if clusterSettingsRaw, ok := d.GetOk("cluster_setting"); ok {
		envelope.AppServiceEnvironment.ClusterSettings = expandAppServiceEnvironmentClusterSettings(clusterSettingsRaw)
	}

	if kind == "ASEV2" {
		envelope.AppServiceEnvironment.MultiSize = utils.String(convertFromIsolatedSKU(pricingTier))
		envelope.AppServiceEnvironment.UserWhitelistedIPRanges = utils.ExpandStringSlice(userWhitelistedIPRangesRaw)
	}

	// whilst this returns a future go-autorest has a max number of retries
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, envelope); err != nil {
		return fmt.Errorf("creating App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	createWait := resource.StateChangeConf{
		Pending: []string{
			string(web.ProvisioningStateInProgress),
		},
		Target: []string{
			string(web.ProvisioningStateSucceeded),
		},
		MinTimeout: 1 * time.Minute,
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Refresh:    appServiceEnvironmentRefresh(ctx, client, resourceGroup, name),
	}

	// as such we'll ignore it and use a custom poller instead
	if _, err := createWait.WaitForState(); err != nil {
		return fmt.Errorf("waiting for the creation of App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceAppServiceEnvironmentRead(d, meta)
}

func resourceAppServiceEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if d.HasChanges("user_whitelisted_ip_ranges", "allowed_user_ip_cidrs") {
		e.UserWhitelistedIPRanges = utils.ExpandStringSlice(d.Get("user_whitelisted_ip_ranges").(*schema.Set).List())
		if v, ok := d.GetOk("user_whitelisted_ip_ranges"); ok {
			e.UserWhitelistedIPRanges = utils.ExpandStringSlice(v.(*schema.Set).List())
		}
	}

	if d.HasChange("cluster_setting") {
		e.ClusterSettings = expandAppServiceEnvironmentClusterSettings(d.Get("cluster_setting"))
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.HostingEnvironmentName, e); err != nil {
		return fmt.Errorf("updating App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	updateWait := resource.StateChangeConf{
		Pending: []string{
			string(web.ProvisioningStateInProgress),
		},
		Target: []string{
			string(web.ProvisioningStateSucceeded),
		},
		MinTimeout: 1 * time.Minute,
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Refresh:    appServiceEnvironmentRefresh(ctx, client, id.ResourceGroup, id.HostingEnvironmentName),
	}

	if _, err := updateWait.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Update of App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	return resourceAppServiceEnvironmentRead(d, meta)
}

func resourceAppServiceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
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

	kind := ""
	if existing.Kind != nil {
		kind = *existing.Kind
	}
	d.Set("version", kind)

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
		d.Set("user_whitelisted_ip_ranges", props.UserWhitelistedIPRanges)
		d.Set("allowed_user_ip_cidrs", props.UserWhitelistedIPRanges)
		d.Set("cluster_setting", flattenClusterSettings(props.ClusterSettings))
	}

	return tags.FlattenAndSet(d, existing.Tags)
}

func resourceAppServiceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Environment %q (Resource Group %q)", id.HostingEnvironmentName, id.ResourceGroup)

	// TODO: should this behaviour be added to the `features` block?
	forceDeleteAllChildren := utils.Bool(false)
	future, err := client.Delete(ctx, id.ResourceGroup, id.HostingEnvironmentName, forceDeleteAllChildren)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("deleting App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("waiting for deletion of App Service Environment %q (Resource Group %q): %+v", id.HostingEnvironmentName, id.ResourceGroup, err)
	}

	return nil
}

func appServiceEnvironmentRefresh(ctx context.Context, client *web.AppServiceEnvironmentsClient, resourceGroup string, name string) resource.StateRefreshFunc {
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
	default:
		vmSKU = isolated
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
	default:
		isolated = vmSKU
	}
	return isolated
}

func loadBalancingModeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
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
	if len(*input) == 0 {
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
