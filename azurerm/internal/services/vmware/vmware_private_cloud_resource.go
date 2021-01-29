package vmware

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceVmwarePrivateCloud() *schema.Resource {
	return &schema.Resource{
		Create: resourceVmwarePrivateCloudCreate,
		Read:   resourceVmwarePrivateCloudRead,
		Update: resourceVmwarePrivateCloudUpdate,
		Delete: resourceVmwarePrivateCloudDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Hour),
			Delete: schema.DefaultTimeout(10 * time.Hour),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PrivateCloudID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"av20",
					"av36",
					"av36t",
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"management_cluster": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3, 16),
						},

						"hosts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"network_subnet_cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},

			"internet_connection_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"nsxt_password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vcenter_password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"circuit": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"express_route_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"express_route_private_peering_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"primary_subnet_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"secondary_subnet_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"hcx_cloud_manager_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"management_subnet_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nsxt_certificate_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"nsxt_manager_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provisioning_subnet_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vcenter_certificate_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vcsa_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vmotion_subnet_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}
func resourceVmwarePrivateCloudCreate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewPrivateCloudID(subscriptionId, resourceGroup, name).ID()

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Vmware Private Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_vmware_private_cloud", id)
	}

	internet := avs.Disabled
	if d.Get("internet_connection_enabled").(bool) {
		internet = avs.Enabled
	}

	privateCloud := avs.PrivateCloud{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Sku: &avs.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		PrivateCloudProperties: &avs.PrivateCloudProperties{
			ManagementCluster: &avs.ManagementCluster{
				ClusterSize: utils.Int32(int32(d.Get("management_cluster.0.size").(int))),
			},
			NetworkBlock:    utils.String(d.Get("network_subnet_cidr").(string)),
			Internet:        internet,
			NsxtPassword:    utils.String(d.Get("nsxt_password").(string)),
			VcenterPassword: utils.String(d.Get("vcenter_password").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, privateCloud)
	if err != nil {
		return fmt.Errorf("creating Vmware Private Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of the Vmware Private Cloud %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)

	return resourceVmwarePrivateCloudRead(d, meta)
}

func resourceVmwarePrivateCloudRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Vmware Private Cloud %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Vmware Private Cloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.PrivateCloudProperties; props != nil {
		if err := d.Set("management_cluster", flattenArmPrivateCloudManagementCluster(props.ManagementCluster)); err != nil {
			return fmt.Errorf("setting `management_cluster`: %+v", err)
		}
		d.Set("network_subnet_cidr", props.NetworkBlock)
		if err := d.Set("circuit", flattenArmPrivateCloudCircuit(props.Circuit)); err != nil {
			return fmt.Errorf("setting `circuit`: %+v", err)
		}

		d.Set("internet_connection_enabled", props.Internet == avs.Enabled)
		d.Set("hcx_cloud_manager_endpoint", props.Endpoints.HcxCloudManager)
		d.Set("nsxt_manager_endpoint", props.Endpoints.NsxtManager)
		d.Set("vcsa_endpoint", props.Endpoints.Vcsa)
		d.Set("management_subnet_cidr", props.ManagementNetwork)
		d.Set("nsxt_certificate_thumbprint", props.NsxtCertificateThumbprint)
		d.Set("provisioning_subnet_cidr", props.ProvisioningNetwork)
		d.Set("vcenter_certificate_thumbprint", props.VcenterCertificateThumbprint)
		d.Set("vmotion_subnet_cidr", props.VmotionNetwork)
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVmwarePrivateCloudUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	privateCloudUpdate := avs.PrivateCloudUpdate{
		PrivateCloudUpdateProperties: &avs.PrivateCloudUpdateProperties{},
	}

	clusterSizeChange := false
	if d.HasChange("management_cluster") {
		privateCloudUpdate.PrivateCloudUpdateProperties.ManagementCluster = &avs.ManagementCluster{
			ClusterSize: utils.Int32(int32(d.Get("management_cluster.0.size").(int))),
		}
		clusterSizeChange = true
	}
	if d.HasChange("internet_connection_enabled") {
		if clusterSizeChange {
			return fmt.Errorf("`management_cluster.0.size` and `internet_connection_enabled` could not be changed together")
		}

		internet := avs.Disabled
		if d.Get("internet_connection_enabled").(bool) {
			internet = avs.Enabled
		}
		privateCloudUpdate.PrivateCloudUpdateProperties.Internet = internet
	}

	if d.HasChange("tags") {
		privateCloudUpdate.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, privateCloudUpdate)
	if err != nil {
		return fmt.Errorf("updating Vmware Private Cloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Vmware Private Cloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return resourceVmwarePrivateCloudRead(d, meta)
}

func resourceVmwarePrivateCloudDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Vmware Private Cloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of the Vmware Private Cloud %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func flattenArmPrivateCloudManagementCluster(input *avs.ManagementCluster) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var clusterSize int32
	if input.ClusterSize != nil {
		clusterSize = *input.ClusterSize
	}
	var clusterId int32
	if input.ClusterID != nil {
		clusterId = *input.ClusterID
	}
	return []interface{}{
		map[string]interface{}{
			"size":  clusterSize,
			"id":    clusterId,
			"hosts": utils.FlattenStringSlice(input.Hosts),
		},
	}
}

func flattenArmPrivateCloudCircuit(input *avs.Circuit) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var expressRouteId string
	if input.ExpressRouteID != nil {
		expressRouteId = *input.ExpressRouteID
	}
	var expressRoutePrivatePeeringId string
	if input.ExpressRoutePrivatePeeringID != nil {
		expressRoutePrivatePeeringId = *input.ExpressRoutePrivatePeeringID
	}
	var primarySubnet string
	if input.PrimarySubnet != nil {
		primarySubnet = *input.PrimarySubnet
	}
	var secondarySubnet string
	if input.SecondarySubnet != nil {
		secondarySubnet = *input.SecondarySubnet
	}
	return []interface{}{
		map[string]interface{}{
			"express_route_id":                 expressRouteId,
			"express_route_private_peering_id": expressRoutePrivatePeeringId,
			"primary_subnet_cidr":              primarySubnet,
			"secondary_subnet_cidr":            secondarySubnet,
		},
	}
}
