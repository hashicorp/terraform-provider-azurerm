package web

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAppServiceEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAppServiceEnvironmentCreateOrUpdate,
		Read:   resourceArmAppServiceEnvironmentRead,
		Update: resourceArmAppServiceEnvironmentCreateOrUpdate,
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

			"internal_load_balancing_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(web.InternalLoadBalancingModeNone),
				ValidateFunc: validation.StringInSlice([]string{
					string(web.InternalLoadBalancingModeNone),
					string(web.InternalLoadBalancingModePublishing),
					string(web.InternalLoadBalancingModeWeb),
				}, false),
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: network.ValidateSubnetID,
			},

			// Note: This is currently 'multiSize' in the API for historic v1 reasons, may change in future?
			"pricing_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "I1",
				ValidateFunc: validate.AppServiceEnvironmentPricingTier,
			},

			"front_end_scale_factor": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      15,
				ValidateFunc: validation.IntBetween(5, 15),
			},

			"location": {
				Type:      schema.TypeString,
				Computed:  true,
				StateFunc: azure.NormalizeLocation,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceArmAppServiceEnvironmentCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	vnetClient := meta.(*clients.Client).Network.VnetClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	internalLoadBalancingMode := d.Get("internal_load_balancing_mode").(string)
	t := d.Get("tags").(map[string]interface{})

	subnetId := d.Get("subnet_id").(string)
	subnet, err := network.ParseSubnetID(subnetId)
	if err != nil {
		return fmt.Errorf("Error parsing subnet id %q: %+v", subnetId, err)
	}

	resourceGroup := subnet.ResourceGroup

	vnet, err := vnetClient.Get(ctx, resourceGroup, subnet.VirtualNetworkName, "")
	if err != nil {
		return fmt.Errorf("Error reading Virtual Network %q for App Service Environment %q: %+v", subnet.VirtualNetworkName, name, err)
	}

	var location string
	if vnetLoc := vnet.Location; vnetLoc != nil {
		location = azure.NormalizeLocation(*vnetLoc)
	} else {
		return fmt.Errorf("Error determining Location from Virtual Network %s", *vnet.Name)
	}

	frontEndScaleFactor := d.Get("front_end_scale_factor").(int)

	pricingTier := d.Get("pricing_tier").(string)

	// the SDK is coded primarily for v1, which needs a non-null entry for workerpool, so we construct an empty slice for it
	// TODO Submit change for SDK?
	wp := []web.WorkerPool{{}}

	envelope := web.AppServiceEnvironmentResource{
		Location: utils.String(location),
		Kind:     utils.String("ASEV2"),
		AppServiceEnvironment: &web.AppServiceEnvironment{
			FrontEndScaleFactor:       utils.Int32(int32(frontEndScaleFactor)),
			MultiSize:                 utils.String(convertFromIsolatedSKU(pricingTier)),
			Name:                      utils.String(name),
			Location:                  utils.String(location),
			InternalLoadBalancingMode: web.InternalLoadBalancingMode(internalLoadBalancingMode),
			VirtualNetwork: &web.VirtualNetworkProfile{
				ID:     utils.String(subnetId),
				Subnet: utils.String(subnet.Name),
			},
			WorkerPools: &wp,
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, envelope)
	if err != nil {
		return fmt.Errorf("Error creating App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the creation of App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving App Service Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*read.ID)

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

	resourceGroup := id.ResourceGroup
	name := id.Name

	appServiceEnvironment, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(appServiceEnvironment.Response) {
			log.Printf("[DEBUG] App Service Environmment %q (Resource Group %q) was not found!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving App Service Environmment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := appServiceEnvironment.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := tags.FlattenAndSet(d, appServiceEnvironment.Tags); err != nil {
		return fmt.Errorf("Error flattening and setting tags in App Service Environment %q (resource group %q): %+v", name, resourceGroup, err)
	}

	ase := appServiceEnvironment.AppServiceEnvironment
	if ase.InternalLoadBalancingMode != "" {
		d.Set("internal_load_balancing_mode", ase.InternalLoadBalancingMode)
	}
	if ase.VirtualNetwork.ID != nil {
		d.Set("subnet_id", ase.VirtualNetwork.ID)
	}
	if ase.FrontEndScaleFactor != nil {
		d.Set("front_end_scale_factor", int(*ase.FrontEndScaleFactor))
	}
	if ase.MultiSize != nil {
		d.Set("pricing_tier", convertToIsolatedSKU(*ase.MultiSize))
	}

	return nil
}

func resourceArmAppServiceEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServiceEnvironmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AppServiceEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Name

	log.Printf("[DEBUG] Deleting App Service Environment %q (Resource Group %q)", name, resGroup)

	// `true` below deletes any child resources (e.g. App Services / Plans / Certificates etc)
	// This potentially destroys resources outside of Terraform's state without the user knowing
	// It is set to true as this is consistent with other instances of this type of functionality in the provider.
	future, err := client.Delete(ctx, resGroup, name, utils.Bool(true))
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return err
	}

	return nil
}

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
