package vmware

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/avs/mgmt/2020-03-20/avs"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/vmware/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceVmwarePrivateCloud() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVmwarePrivateCloudRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"management_cluster": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"size": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"hosts": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"id": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"network_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"internet_connection_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"circuit": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"express_route_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"express_route_private_peering_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"primary_subnet_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"secondary_subnet_cidr": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"hcx_cloud_manager_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"management_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"nsxt_certificate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"nsxt_manager_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vcenter_certificate_thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vcsa_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"vmotion_subnet_cidr": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceVmwarePrivateCloudRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewPrivateCloudID(subscriptionId, resourceGroup, name).ID()

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("VMware PrivateCloud %q does not exist", name)
		}
		return fmt.Errorf("retrieving VMware PrivateCloud %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(id)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
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
