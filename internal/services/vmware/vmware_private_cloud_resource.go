package vmware

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/sdk/2020-03-20/privateclouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVmwarePrivateCloud() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVmwarePrivateCloudCreate,
		Read:   resourceVmwarePrivateCloudRead,
		Update: resourceVmwarePrivateCloudUpdate,
		Delete: resourceVmwarePrivateCloudDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(10 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(10 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(10 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privateclouds.ParsePrivateCloudID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"av20",
					"av36",
					"av36t",
				}, false),
			},

			"management_cluster": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"size": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(3, 16),
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
			},

			"internet_connection_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"nsxt_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vcenter_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceVmwarePrivateCloudCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privateclouds.NewPrivateCloudID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_vmware_private_cloud", id.ID())
	}

	internet := privateclouds.InternetEnumDisabled
	if d.Get("internet_connection_enabled").(bool) {
		internet = privateclouds.InternetEnumEnabled
	}

	privateCloud := privateclouds.PrivateCloud{
		Location: location.Normalize(d.Get("location").(string)),
		Sku: privateclouds.Sku{
			Name: d.Get("sku_name").(string),
		},
		Properties: privateclouds.PrivateCloudProperties{
			ManagementCluster: privateclouds.ManagementCluster{
				ClusterSize: int64(d.Get("management_cluster.0.size").(int)),
			},
			NetworkBlock:    d.Get("network_subnet_cidr").(string),
			Internet:        &internet,
			NsxtPassword:    utils.String(d.Get("nsxt_password").(string)),
			VcenterPassword: utils.String(d.Get("vcenter_password").(string)),
		},
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, privateCloud); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVmwarePrivateCloudRead(d, meta)
}

func resourceVmwarePrivateCloudRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateclouds.ParsePrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		props := model.Properties

		if err := d.Set("management_cluster", flattenPrivateCloudManagementCluster(props.ManagementCluster)); err != nil {
			return fmt.Errorf("setting `management_cluster`: %+v", err)
		}
		d.Set("network_subnet_cidr", props.NetworkBlock)
		if err := d.Set("circuit", flattenPrivateCloudCircuit(props.Circuit)); err != nil {
			return fmt.Errorf("setting `circuit`: %+v", err)
		}

		internetConnectionEnabled := false
		if props.Internet != nil {
			internetConnectionEnabled = *props.Internet == privateclouds.InternetEnumEnabled
		}
		d.Set("internet_connection_enabled", internetConnectionEnabled)

		d.Set("hcx_cloud_manager_endpoint", props.Endpoints.HcxCloudManager)
		d.Set("nsxt_manager_endpoint", props.Endpoints.NsxtManager)
		d.Set("vcsa_endpoint", props.Endpoints.Vcsa)
		d.Set("management_subnet_cidr", props.ManagementNetwork)
		d.Set("nsxt_certificate_thumbprint", props.NsxtCertificateThumbprint)
		d.Set("provisioning_subnet_cidr", props.ProvisioningNetwork)
		d.Set("vcenter_certificate_thumbprint", props.VcenterCertificateThumbprint)
		d.Set("vmotion_subnet_cidr", props.VmotionNetwork)

		d.Set("sku_name", model.Sku.Name)

		if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceVmwarePrivateCloudUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateclouds.ParsePrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	privateCloudUpdate := privateclouds.PrivateCloudUpdate{
		Properties: &privateclouds.PrivateCloudUpdateProperties{},
	}

	if d.HasChange("management_cluster") && d.HasChange("internet_connection_enabled") {
		return fmt.Errorf("`management_cluster.0.size` and `internet_connection_enabled` could not be changed together")
	}

	if d.HasChange("management_cluster") {
		privateCloudUpdate.Properties.ManagementCluster = &privateclouds.ManagementCluster{
			ClusterSize: int64(d.Get("management_cluster.0.size").(int)),
		}
	}

	if d.HasChange("internet_connection_enabled") {
		internet := privateclouds.InternetEnumDisabled
		if d.Get("internet_connection_enabled").(bool) {
			internet = privateclouds.InternetEnumEnabled
		}
		privateCloudUpdate.Properties.Internet = &internet
	}

	if d.HasChange("tags") {
		privateCloudUpdate.Tags = expandTags(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, privateCloudUpdate); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVmwarePrivateCloudRead(d, meta)
}

func resourceVmwarePrivateCloudDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Vmware.PrivateCloudClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateclouds.ParsePrivateCloudID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
