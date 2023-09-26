// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2023-03-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	autoParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	autoVal "github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

func resourceArmStackHCICluster() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmStackHCIClusterCreate,
		Read:   resourceArmStackHCIClusterRead,
		Update: resourceArmStackHCIClusterUpdate,
		Delete: resourceArmStackHCIClusterDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := clusters.ParseClusterID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"client_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"tenant_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"automanage_configuration_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: autoVal.AutomanageConfigurationID,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceArmStackHCIClusterCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := clusters.NewClusterID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_stack_hci_cluster", id.ID())
	}

	cluster := clusters.Cluster{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &clusters.ClusterProperties{
			AadClientId: utils.String(d.Get("client_id").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		cluster.Properties.AadTenantId = utils.String(v.(string))
	} else {
		tenantId := meta.(*clients.Client).Account.TenantId
		cluster.Properties.AadTenantId = utils.String(tenantId)
	}

	if _, err := client.Create(ctx, id, cluster); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if v, ok := d.GetOk("automanage_configuration_id"); ok {
		hciAssignmentClient := meta.(*clients.Client).Automanage.HCIAssignmentClient
		autoConfigClient := meta.(*clients.Client).Automanage.ConfigurationClient

		automanageConfigId, err := autoParse.AutomanageConfigurationID(v.(string))
		if err != nil {
			return err
		}

		_, err = autoConfigClient.Get(ctx, automanageConfigId.ConfigurationProfileName, automanageConfigId.ResourceGroup)
		if err != nil {
			return fmt.Errorf("checking for existing %s: %+v", automanageConfigId, err)
		}

		hciAssignmentID := autoParse.NewAutomanageConfigurationHCIAssignmentID(subscriptionId, id.ResourceGroupName, id.ClusterName, "default")

		autoResp, err := hciAssignmentClient.Get(ctx, hciAssignmentID.ResourceGroup, hciAssignmentID.ClusterName, hciAssignmentID.ConfigurationProfileAssignmentName)
		if err != nil && !utils.ResponseWasNotFound(autoResp.Response) {
			return fmt.Errorf("checking for existing %s: %+v", hciAssignmentID, err)
		}

		if utils.ResponseWasNotFound(autoResp.Response) {
			properties := automanage.ConfigurationProfileAssignment{
				Properties: &automanage.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: utils.String(automanageConfigId.ID()),
				},
			}

			if _, err := hciAssignmentClient.CreateOrUpdate(ctx, properties, hciAssignmentID.ResourceGroup, hciAssignmentID.ClusterName, hciAssignmentID.ConfigurationProfileAssignmentName); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentID, err)
			}
		}
	}

	d.SetId(id.ID())

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	hciAssignmentClient := meta.(*clients.Client).Automanage.HCIAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("client_id", props.AadClientId)
			d.Set("tenant_id", props.AadTenantId)

			assignmentResp, err := hciAssignmentClient.Get(ctx, id.ResourceGroupName, id.ClusterName, "default")
			if err != nil && !utils.ResponseWasNotFound(assignmentResp.Response) {
				return err
			}
			configId := ""
			if !utils.ResponseWasNotFound(assignmentResp.Response) && assignmentResp.Properties != nil && assignmentResp.Properties.ConfigurationProfile != nil {
				automanageConfigId, err := autoParse.AutomanageConfigurationID(*assignmentResp.Properties.ConfigurationProfile)
				if err != nil {
					return err
				}
				configId = automanageConfigId.ID()
			}

			d.Set("automanage_configuration_id", configId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceArmStackHCIClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	cluster := clusters.ClusterPatch{}

	if d.HasChange("tags") {
		cluster.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.Update(ctx, *id, cluster); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("automanage_configuration_id") {
		if v, ok := d.GetOk("automanage_configuration_id"); ok {
			hciAssignmentClient := meta.(*clients.Client).Automanage.HCIAssignmentClient
			autoConfigClient := meta.(*clients.Client).Automanage.ConfigurationClient

			automanageConfigId, err := autoParse.AutomanageConfigurationID(v.(string))
			if err != nil {
				return err
			}

			_, err = autoConfigClient.Get(ctx, automanageConfigId.ConfigurationProfileName, automanageConfigId.ResourceGroup)
			if err != nil {
				return fmt.Errorf("checking for existing %s: %+v", automanageConfigId, err)
			}

			hciAssignmentID := autoParse.NewAutomanageConfigurationHCIAssignmentID(subscriptionId, id.ResourceGroupName, id.ClusterName, "default")

			properties := automanage.ConfigurationProfileAssignment{
				Properties: &automanage.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: utils.String(automanageConfigId.ID()),
				},
			}

			if _, err := hciAssignmentClient.CreateOrUpdate(ctx, properties, hciAssignmentID.ResourceGroup, hciAssignmentID.ClusterName, hciAssignmentID.ConfigurationProfileAssignmentName); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentID, err)
			}
		} else {
			hciAssignmentClient := meta.(*clients.Client).Automanage.HCIAssignmentClient
			assignmentResp, err := hciAssignmentClient.Get(ctx, id.ResourceGroupName, id.ClusterName, "default")
			if err != nil && !utils.ResponseWasNotFound(assignmentResp.Response) {
				return err
			}

			if !utils.ResponseWasNotFound(assignmentResp.Response) {
				if _, err := hciAssignmentClient.Delete(ctx, id.ResourceGroupName, id.ClusterName, "default"); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}
		}

	}

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	hciAssignmentClient := meta.(*clients.Client).Automanage.HCIAssignmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	assignmentResp, err := hciAssignmentClient.Get(ctx, id.ResourceGroupName, id.ClusterName, "default")
	if err != nil && !utils.ResponseWasNotFound(assignmentResp.Response) {
		return err
	}

	if !utils.ResponseWasNotFound(assignmentResp.Response) {
		if _, err := hciAssignmentClient.Delete(ctx, id.ResourceGroupName, id.ClusterName, "default"); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
