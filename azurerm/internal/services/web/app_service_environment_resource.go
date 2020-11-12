package web

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/go-azure-helpers/response"
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
	InternalLoadBalancingModeWebPublishing web.InternalLoadBalancingMode = "Web, Publishing"
)

func resourceArmAppServiceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceEnvironmentCreate,
		Read:   resourceArmAppServiceEnvironmentRead,
		Update: resourceArmAppServiceEnvironmentUpdate,
		Delete: resourceArmAppServiceEnvironmentDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AppServiceEnvironmentID(id)
			return err
		}),

		// Need to find sane values for below, some operations on this resource can take an exceptionally long time
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Hour),
			Delete: schema.DefaultTimeout(4 * time.Hour),
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

			"internal_load_balancing_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(web.InternalLoadBalancingModeNone),
				ValidateFunc: validation.StringInSlice([]string{
					string(web.InternalLoadBalancingModeNone),
					string(web.InternalLoadBalancingModePublishing),
					string(web.InternalLoadBalancingModeWeb),
					string(InternalLoadBalancingModeWebPublishing),
				}, false),
			},

			"front_end_scale_factor": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntBetween(5, 15),
			},

			"pricing_tier": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "I1",
				ValidateFunc: validation.StringInSlice([]string{
					"I1",
					"I2",
					"I3",
				}, false),
			},

			"allowed_user_ip_cidrs": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true, // remove in 3.0
				ConflictsWith: []string{"user_whitelisted_ip_ranges"},
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: helpersValidate.CIDR,
				},
			},

			"user_whitelisted_ip_ranges": {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true, // remove in 3.0
				ConflictsWith: []string{"allowed_user_ip_cidrs"},
				Deprecated:    "this property has been renamed to `allowed_user_ip_cidrs` better reflect the expected ip range format",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: helpersValidate.CIDR,
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

func resourceArmAppServiceEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	networksClient := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	internalLoadBalancingMode := d.Get("internal_load_balancing_mode").(string)
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
		return fmt.Errorf("Error retrieving Virtual Network %q (Resource Group %q): %+v", subnet.VirtualNetworkName, subnet.ResourceGroup, err)
	}

	// the App Service Environment has to be in the same location as the Virtual Network
	var location string
	if loc := vnet.Location; loc != nil {
		location = azure.NormalizeLocation(*loc)
	} else {
		return fmt.Errorf("Error determining Location from Virtual Network %q (Resource Group %q): `location` was nil", subnet.VirtualNetworkName, subnet.ResourceGroup)
	}

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("Error checking for presence of existing App Service Environment %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_app_service_environment", *existing.ID)
	}

	frontEndScaleFactor := d.Get("front_end_scale_factor").(int)
	pricingTier := d.Get("pricing_tier").(string)

	envelope := web.AppServiceEnvironmentResource{
		Location: utils.String(location),
		Kind:     utils.String("ASEV2"),
		AppServiceEnvironment: &web.AppServiceEnvironment{
			Name:                      utils.String(name),
			Location:                  utils.String(location),
			InternalLoadBalancingMode: web.InternalLoadBalancingMode(internalLoadBalancingMode),
			FrontEndScaleFactor:       utils.Int32(int32(frontEndScaleFactor)),
			MultiSize:                 utils.String(convertFromIsolatedSKU(pricingTier)),
			VirtualNetwork: &web.VirtualNetworkProfile{
				ID:     utils.String(subnetId),
				Subnet: utils.String(subnet.Name),
			},
			UserWhitelistedIPRanges: utils.ExpandStringSlice(userWhitelistedIPRangesRaw),

			// the SDK is coded primarily for v1, which needs a non-null entry for workerpool, so we construct an empty slice for it
			// TODO: remove this hack once https://github.com/Azure/azure-rest-api-specs/pull/8433 has been merged
			WorkerPools: &[]web.WorkerPool{{}},
		},
		Tags: tags.Expand(t),
	}

	// whilst this returns a future go-autorest has a max number of retries
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, envelope); err != nil {
		return fmt.Errorf("Error creating App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// as such we'll ignore it and use a custom poller instead
	if err := waitForAppServiceEnvironmentToStabilize(ctx, client, resourceGroup, name); err != nil {
		return fmt.Errorf("Error waiting for the creation of App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmAppServiceEnvironmentRead(d, meta)
}

func resourceArmAppServiceEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
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
		e.AppServiceEnvironment.InternalLoadBalancingMode = web.InternalLoadBalancingMode(v)
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

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, e); err != nil {
		return fmt.Errorf("Error updating App Service Environment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := waitForAppServiceEnvironmentToStabilize(ctx, client, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("Error waiting for Update of App Service Environment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return resourceArmAppServiceEnvironmentRead(d, meta)
}

func resourceArmAppServiceEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			log.Printf("[DEBUG] App Service Environmment %q (Resource Group %q) was not found - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving App Service Environmment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
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
		d.Set("user_whitelisted_ip_ranges", props.UserWhitelistedIPRanges)
		d.Set("allowed_user_ip_cidrs", props.UserWhitelistedIPRanges)
	}

	return tags.FlattenAndSet(d, existing.Tags)
}

func resourceArmAppServiceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Environment %q (Resource Group %q)", id.Name, id.ResourceGroup)

	// TODO: should this behaviour be added to the `features` block?
	forceDeleteAllChildren := utils.Bool(false)
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, forceDeleteAllChildren)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting App Service Environment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of App Service Environment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func waitForAppServiceEnvironmentToStabilize(ctx context.Context, client *web.AppServiceEnvironmentsClient, resourceGroup string, name string) error {
	for {
		time.Sleep(1 * time.Minute)

		read, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return err
		}

		if read.AppServiceEnvironment == nil {
			return fmt.Errorf("`properties` was nil")
		}

		state := read.AppServiceEnvironment.ProvisioningState
		if state == web.ProvisioningStateSucceeded {
			return nil
		}

		if state == web.ProvisioningStateInProgress {
			continue
		}

		return fmt.Errorf("Unexpected ProvisioningState: %q", state)
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
