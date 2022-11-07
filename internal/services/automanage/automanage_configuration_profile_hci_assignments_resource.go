package automanage

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

func resourceAutomanageConfigurationProfileHCIAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomanageConfigurationProfileHCIAssignmentCreateUpdate,
		Read:   resourceAutomanageConfigurationProfileHCIAssignmentRead,
		Update: resourceAutomanageConfigurationProfileHCIAssignmentCreateUpdate,
		Delete: resourceAutomanageConfigurationProfileHCIAssignmentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomanageConfigurationProfileHCIAssignmentID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"cluster_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"configuration_profile": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"managed_by": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"target_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceAutomanageConfigurationProfileHCIAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	clusterName := d.Get("cluster_name").(string)

	id := parse.NewAutomanageConfigurationProfileHCIAssignmentID(subscriptionId, resourceGroup, clusterName, name).ID()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, clusterName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q): %+v", name, resourceGroup, clusterName, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automanage_configuration_profile_hci_assignment", id)
		}
	}

	parameters := automanage.ConfigurationProfileAssignment{
		Properties: &automanage.ConfigurationProfileAssignmentProperties{
			ConfigurationProfile: utils.String(d.Get("configuration_profile").(string)),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, parameters, resourceGroup, clusterName, name); err != nil {
		return fmt.Errorf("creating/updating Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q): %+v", name, resourceGroup, clusterName, err)
	}

	d.SetId(id)
	return resourceAutomanageConfigurationProfileHCIAssignmentRead(d, meta)
}

func resourceAutomanageConfigurationProfileHCIAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileHCIAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q): %+v", id.ConfigurationProfileAssignmentName, id.ResourceGroup, id.ClusterName, err)
	}
	d.Set("name", id.ConfigurationProfileAssignmentName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("cluster_name", id.ClusterName)
	if props := resp.Properties; props != nil {
		d.Set("configuration_profile", props.ConfigurationProfile)
		d.Set("target_id", props.TargetID)
	}
	d.Set("managed_by", resp.ManagedBy)
	d.Set("type", resp.Type)
	return nil
}

func resourceAutomanageConfigurationProfileHCIAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCIAssignmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileHCIAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.ConfigurationProfileAssignmentName); err != nil {
		return fmt.Errorf("deleting Automanage ConfigurationProfileHCIAssignment %q (Resource Group %q / clusterName %q): %+v", id.ConfigurationProfileAssignmentName, id.ResourceGroup, id.ClusterName, err)
	}
	return nil
}
