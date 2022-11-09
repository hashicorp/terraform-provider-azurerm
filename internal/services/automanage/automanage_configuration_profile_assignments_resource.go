package automanage

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/automanage/2022-05-04/automanage"
)

func resourceAutomanageConfigurationProfileAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomanageConfigurationProfileAssignmentCreateUpdate,
		Read:   resourceAutomanageConfigurationProfileAssignmentRead,
		Update: resourceAutomanageConfigurationProfileAssignmentCreateUpdate,
		Delete: resourceAutomanageConfigurationProfileAssignmentDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutomanageConfigurationProfileAssignmentID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"vm_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"configuration_profile_id": {
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
func resourceAutomanageConfigurationProfileAssignmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Automanage.ConfigurationProfileAssignmentClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	vmName := d.Get("vm_name").(string)

	id := parse.NewAutomanageConfigurationProfileAssignmentID(subscriptionId, resourceGroup, vmName, name).ID()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, vmName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", name, resourceGroup, vmName, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automanage_configuration_profile_assignment", id)
		}
	}

	parameters := automanage.ConfigurationProfileAssignment{
		Properties: &automanage.ConfigurationProfileAssignmentProperties{
			ConfigurationProfile: utils.String(d.Get("configuration_profile_id").(string)),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, name, parameters, resourceGroup, vmName); err != nil {
		return fmt.Errorf("creating/updating Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", name, resourceGroup, vmName, err)
	}

	d.SetId(id)
	return resourceAutomanageConfigurationProfileAssignmentRead(d, meta)
}

func resourceAutomanageConfigurationProfileAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ConfigurationProfileAssignmentName, id.VirtualMachineName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", id.ConfigurationProfileAssignmentName, id.ResourceGroup, id.VirtualMachineName, err)
	}
	d.Set("name", id.ConfigurationProfileAssignmentName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("vm_name", id.VirtualMachineName)
	if props := resp.Properties; props != nil {
		d.Set("configuration_profile_id", props.ConfigurationProfile)
		d.Set("target_id", props.TargetID)
	}
	d.Set("managed_by", resp.ManagedBy)
	d.Set("type", resp.Type)
	return nil
}

func resourceAutomanageConfigurationProfileAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileAssignmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ConfigurationProfileAssignmentName, id.VirtualMachineName); err != nil {
		return fmt.Errorf("deleting Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", id.ConfigurationProfileAssignmentName, id.ResourceGroup, id.VirtualMachineName, err)
	}
	return nil
}
