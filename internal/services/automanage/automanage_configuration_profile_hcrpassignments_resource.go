package automanage

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/mgmt/2022-05-04/automanage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"log"
	"time"
)

func resourceAutomanageConfigurationProfileHCRPAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomanageConfigurationProfileHCRPAssignmentCreateUpdate,
		Read:   resourceAutomanageConfigurationProfileHCRPAssignmentRead,
		Update: resourceAutomanageConfigurationProfileHCRPAssignmentCreateUpdate,
		Delete: resourceAutomanageConfigurationProfileHCRPAssignmentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomanageConfigurationProfileHCRPAssignmentID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"machine_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"configuration_profile": {
				Type:     pluginsdk.TypeString,
				Optional: true,
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
func resourceAutomanageConfigurationProfileHCRPAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCRPAssignmentClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	machineName := d.Get("machine_name").(string)

	id := parse.NewAutomanageConfigurationProfileHCRPAssignmentID(subscriptionId, resourceGroup, machineName, name).ID()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, machineName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Automanage ConfigurationProfileHCRPAssignment %q (Resource Group %q / machineName %q): %+v", name, resourceGroup, machineName, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automanage_configuration_profile_hcrpassignment", id)
		}
	}

	parameters := automanage.ConfigurationProfileAssignment{
		Properties: &automanage.ConfigurationProfileAssignmentProperties{
			ConfigurationProfile: utils.String(d.Get("configuration_profile").(string)),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, parameters, resourceGroup, machineName, name); err != nil {
		return fmt.Errorf("creating/updating Automanage ConfigurationProfileHCRPAssignment %q (Resource Group %q / machineName %q): %+v", name, resourceGroup, machineName, err)
	}

	d.SetId(id)
	return resourceAutomanageConfigurationProfileHCRPAssignmentRead(d, meta)
}

func resourceAutomanageConfigurationProfileHCRPAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCRPAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileHCRPAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.MachineName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfileHCRPAssignment %q (Resource Group %q / machineName %q): %+v", id.Name, id.ResourceGroup, id.MachineName, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("machine_name", id.MachineName)
	if props := resp.Properties; props != nil {
		d.Set("configuration_profile", props.ConfigurationProfile)
		d.Set("target_id", props.TargetID)
	}
	d.Set("managed_by", resp.ManagedBy)
	d.Set("type", resp.Type)
	return nil
}

func resourceAutomanageConfigurationProfileHCRPAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileHCRPAssignmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileHCRPAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.MachineName, id.Name); err != nil {
		return fmt.Errorf("deleting Automanage ConfigurationProfileHCRPAssignment %q (Resource Group %q / machineName %q): %+v", id.Name, id.ResourceGroup, id.MachineName, err)
	}
	return nil
}
