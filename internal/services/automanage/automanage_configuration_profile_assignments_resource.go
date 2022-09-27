package automanage

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/mgmt/2022-05-04/automanage"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
)

func resourceAutomanageConfigurationProfileAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAutomanageConfigurationProfileAssignmentCreateUpdate,
		Read:   resourceAutomanageConfigurationProfileAssignmentRead,
		Update: resourceAutomanageConfigurationProfileAssignmentCreateUpdate,
		Delete: resourceAutomanageConfigurationProfileAssignmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.AutomanageConfigurationProfileAssignmentID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"vm_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"configuration_profile": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"managed_by": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"target_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceAutomanageConfigurationProfileAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
			ConfigurationProfile: utils.String(d.Get("configuration_profile").(string)),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, name, parameters, resourceGroup, vmName); err != nil {
		return fmt.Errorf("creating/updating Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", name, resourceGroup, vmName, err)
	}

	d.SetId(id)
	return resourceAutomanageConfigurationProfileAssignmentRead(d, meta)
}

func resourceAutomanageConfigurationProfileAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileAssignmentClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, id.VMName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] automanage %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", id.Name, id.ResourceGroup, id.VMName, err)
	}
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("vm_name", id.VMName)
	if props := resp.Properties; props != nil {
		d.Set("configuration_profile", props.ConfigurationProfile)
		d.Set("target_id", props.TargetID)
	}
	d.Set("managed_by", resp.ManagedBy)
	d.Set("type", resp.Type)
	return nil
}

func resourceAutomanageConfigurationProfileAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automanage.ConfigurationProfileAssignmentClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomanageConfigurationProfileAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name, id.VMName); err != nil {
		return fmt.Errorf("deleting Automanage ConfigurationProfileAssignment %q (Resource Group %q / vmName %q): %+v", id.Name, id.ResourceGroup, id.VMName, err)
	}
	return nil
}
