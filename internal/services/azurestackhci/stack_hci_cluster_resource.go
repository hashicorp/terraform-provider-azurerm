// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				Optional:     true,
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
				// TODO: this field should be removed in 4.0 - there's an "association" API specifically for this purpose
				// so we should be outputting this as an association resource.
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: configurationprofiles.ValidateConfigurationProfileID,
			},

			"cloud_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"resource_provider_object_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityOptional(),

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

	if v, ok := d.GetOk("identity"); ok {
		cluster.Identity = expandSystemAssigned(v.([]interface{}))
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
		configurationProfilesClient := meta.(*clients.Client).Automanage.ConfigurationProfilesClient
		hciAssignmentsClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient

		configurationProfileId, err := configurationprofiles.ParseConfigurationProfileID(v.(string))
		if err != nil {
			return err
		}

		if _, err = configurationProfilesClient.Get(ctx, *configurationProfileId); err != nil {
			return fmt.Errorf("checking for existing %s: %+v", configurationProfileId, err)
		}

		hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(subscriptionId, id.ResourceGroupName, id.ClusterName, "default")
		assignmentsResp, err := hciAssignmentsClient.Get(ctx, hciAssignmentId)
		if err != nil && !response.WasNotFound(assignmentsResp.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", hciAssignmentId, err)
		}

		if response.WasNotFound(assignmentsResp.HttpResponse) {
			properties := configurationprofilehciassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehciassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: utils.String(configurationProfileId.ID()),
				},
			}

			if _, err := hciAssignmentsClient.CreateOrUpdate(ctx, hciAssignmentId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentId, err)
			}
		}
	}

	d.SetId(id.ID())

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	hciAssignmentsClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
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
		d.Set("identity", flattenSystemAssigned(model.Identity))

		if props := model.Properties; props != nil {
			d.Set("client_id", props.AadClientId)
			d.Set("tenant_id", props.AadTenantId)
			d.Set("cloud_id", props.CloudId)
			d.Set("service_endpoint", props.ServiceEndpoint)
			d.Set("resource_provider_object_id", props.ResourceProviderObjectId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	hclAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	assignmentResp, err := hciAssignmentsClient.Get(ctx, hclAssignmentId)
	if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
		return err
	}
	configId := ""
	if model := assignmentResp.Model; model != nil && model.Properties != nil && model.Properties.ConfigurationProfile != nil {
		parsed, err := configurationprofiles.ParseConfigurationProfileIDInsensitively(*model.Properties.ConfigurationProfile)
		if err != nil {
			return err
		}
		configId = parsed.ID()
	}
	d.Set("automanage_configuration_id", configId)

	return nil
}

func resourceArmStackHCIClusterUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
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

	if d.HasChange("identity") {
		cluster.Identity = expandSystemAssigned(d.Get("identity").([]interface{}))
	}

	if _, err := client.Update(ctx, *id, cluster); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("automanage_configuration_id") {
		hciAssignmentClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
		configurationProfilesClient := meta.(*clients.Client).Automanage.ConfigurationProfilesClient
		hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")

		if v, ok := d.GetOk("automanage_configuration_id"); ok {
			configurationProfileId, err := configurationprofiles.ParseConfigurationProfileID(v.(string))
			if err != nil {
				return err
			}

			if _, err = configurationProfilesClient.Get(ctx, *configurationProfileId); err != nil {
				return fmt.Errorf("checking for existing %s: %+v", configurationProfileId, err)
			}

			properties := configurationprofilehciassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehciassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: utils.String(configurationProfileId.ID()),
				},
			}

			if _, err := hciAssignmentClient.CreateOrUpdate(ctx, hciAssignmentId, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", hciAssignmentId, err)
			}
		} else {
			assignmentResp, err := hciAssignmentClient.Get(ctx, hciAssignmentId)
			if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
				return err
			}

			if !response.WasNotFound(assignmentResp.HttpResponse) {
				if _, err := hciAssignmentClient.Delete(ctx, hciAssignmentId); err != nil {
					return fmt.Errorf("deleting %s: %+v", id, err)
				}
			}
		}
	}

	return resourceArmStackHCIClusterRead(d, meta)
}

func resourceArmStackHCIClusterDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AzureStackHCI.Clusters
	hciAssignmentClient := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := clusters.ParseClusterID(d.Id())
	if err != nil {
		return err
	}

	hciAssignmentId := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, "default")
	assignmentResp, err := hciAssignmentClient.Get(ctx, hciAssignmentId)
	if err != nil && !response.WasNotFound(assignmentResp.HttpResponse) {
		return err
	}

	if !response.WasNotFound(assignmentResp.HttpResponse) {
		if _, err := hciAssignmentClient.Delete(ctx, hciAssignmentId); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

// API does not accept userAssignedIdentity as in swagger https://github.com/Azure/azure-rest-api-specs/issues/28260
func expandSystemAssigned(input []interface{}) *identity.SystemAndUserAssignedMap {
	if len(input) == 0 || input[0] == nil {
		return &identity.SystemAndUserAssignedMap{
			Type: identity.TypeNone,
		}
	}

	return &identity.SystemAndUserAssignedMap{
		Type: identity.TypeSystemAssigned,
	}
}

func flattenSystemAssigned(input *identity.SystemAndUserAssignedMap) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	if input.Type == identity.TypeNone {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"type":         input.Type,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}
}
