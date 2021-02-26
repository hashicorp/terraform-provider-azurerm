package policy

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/guestconfiguration/mgmt/2020-06-25/guestconfiguration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	computeParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	computeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceGuestConfigurationAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceGuestConfigurationAssignmentCreateUpdate,
		Read:   resourceGuestConfigurationAssignmentRead,
		Update: resourceGuestConfigurationAssignmentCreateUpdate,
		Delete: resourceGuestConfigurationAssignmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.GuestConfigurationAssignmentID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"virtual_machine_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: computeValidate.VirtualMachineID,
			},

			"guest_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"parameter": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"content_hash": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"content_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"assignment_hash": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"compliance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGuestConfigurationAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmId, err := computeParse.VirtualMachineID(d.Get("virtual_machine_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewGuestConfigurationAssignmentID(subscriptionId, vmId.ResourceGroup, vmId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, id.VMName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing GuestConfiguration GuestConfigurationAssignment %q (Virtual Machine ID %q): %+v", id.Name, vmId.ID(), err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_guest_configuration_assignment", id.ID())
		}
	}

	parameter := guestconfiguration.Assignment{
		Name:     utils.String(d.Get("name").(string)),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Properties: &guestconfiguration.AssignmentProperties{
			GuestConfiguration: expandGuestConfigurationAssignmentNavigation(d.Get("guest_configuration").([]interface{})),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.Name, parameter, id.ResourceGroup, id.VMName)
	if err != nil {
		return fmt.Errorf("creating/updating GuestConfigurationAssignment %q (Virtual Machine ID %q): %+v", id.Name, vmId.ID(), err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for GuestConfigurationAssignment %q (Virtual Machine ID %q): %+v", id.Name, vmId.ID(), err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, id.VMName)
	if err != nil {
		return fmt.Errorf("retrieving GuestConfigurationAssignment %q (Virtual Machine ID %q): %+v", id.Name, vmId.ID(), err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for GuestConfigurationAssignment %q (Virtual Machine ID %q) ID", id.Name, vmId.ID())
	}

	d.SetId(id.ID())

	return resourceGuestConfigurationAssignmentRead(d, meta)
}

func resourceGuestConfigurationAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GuestConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, id.VMName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] guestConfiguration %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving GuestConfigurationAssignment %q (Resource Group %q / vmName %q): %+v", id.Name, id.ResourceGroup, id.VMName, err)
	}
	vmId := computeParse.NewVirtualMachineID(subscriptionId, id.ResourceGroup, id.VMName)
	d.Set("name", id.Name)
	d.Set("virtual_machine_id", vmId.ID())
	d.Set("location", location.NormalizeNilable(resp.Location))
	if props := resp.Properties; props != nil {
		if err := d.Set("guest_configuration", flattenGuestConfigurationAssignmentNavigation(props.GuestConfiguration)); err != nil {
			return fmt.Errorf("setting `guest_configuration`: %+v", err)
		}
		d.Set("assignment_hash", props.AssignmentHash)
		d.Set("compliance_status", props.ComplianceStatus)
	}
	return nil
}

func resourceGuestConfigurationAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.GuestConfigurationAssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GuestConfigurationAssignmentID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name, id.VMName)
	if err != nil {
		return fmt.Errorf("deleting GuestConfiguration GuestConfigurationAssignment %q (Resource Group %q / vmName %q): %+v", id.Name, id.ResourceGroup, id.VMName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for GuestConfiguration GuestConfigurationAssignment %q (Resource Group %q / vmName %q): %+v", id.Name, id.ResourceGroup, id.VMName, err)
	}
	return nil
}

func expandGuestConfigurationAssignmentNavigation(input []interface{}) *guestconfiguration.Navigation {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &guestconfiguration.Navigation{
		Name:                   utils.String(v["name"].(string)),
		Version:                utils.String(v["version"].(string)),
		ConfigurationParameter: expandGuestConfigurationAssignmentConfigurationParameterArray(v["parameter"].(*schema.Set).List()),
	}
}

func expandGuestConfigurationAssignmentConfigurationParameterArray(input []interface{}) *[]guestconfiguration.ConfigurationParameter {
	results := make([]guestconfiguration.ConfigurationParameter, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		results = append(results, guestconfiguration.ConfigurationParameter{
			Name:  utils.String(v["name"].(string)),
			Value: utils.String(v["value"].(string)),
		})
	}
	return &results
}

func expandGuestConfigurationAssignmentConfigurationSetting(input []interface{}) *guestconfiguration.ConfigurationSetting {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	rebootIfNeeded := guestconfiguration.RebootIfNeededFalse
	if v["reboot_if_needed"].(bool) {
		rebootIfNeeded = guestconfiguration.RebootIfNeededTrue
	}
	allowModuleOverwrite := guestconfiguration.False
	if v["allow_module_overwrite"].(bool) {
		allowModuleOverwrite = guestconfiguration.True
	}
	return &guestconfiguration.ConfigurationSetting{
		ConfigurationMode:              guestconfiguration.ConfigurationMode(v["configuration_mode"].(string)),
		AllowModuleOverwrite:           allowModuleOverwrite,
		ActionAfterReboot:              guestconfiguration.ActionAfterReboot(v["action_after_reboot"].(string)),
		RefreshFrequencyMins:           utils.Float(v["refresh_frequency_in_minute"].(float64)),
		RebootIfNeeded:                 rebootIfNeeded,
		ConfigurationModeFrequencyMins: utils.Float(v["configuration_mode_frequency_mins"].(float64)),
	}
}

func flattenGuestConfigurationAssignmentNavigation(input *guestconfiguration.Navigation) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var name string
	if input.Name != nil {
		name = *input.Name
	}
	var version string
	if input.Version != nil {
		version = *input.Version
	}
	var contentHash string
	if input.ContentHash != nil {
		contentHash = *input.ContentHash
	}
	var contentUri string
	if input.ContentURI != nil {
		contentUri = *input.ContentURI
	}
	return []interface{}{
		map[string]interface{}{
			"name":         name,
			"parameter":    flattenGuestConfigurationAssignmentConfigurationParameterArray(input.ConfigurationParameter),
			"version":      version,
			"content_hash": contentHash,
			"content_uri":  contentUri,
		},
	}
}

func flattenGuestConfigurationAssignmentConfigurationParameterArray(input *[]guestconfiguration.ConfigurationParameter) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}
		var value string
		if item.Value != nil {
			value = *item.Value
		}
		results = append(results, map[string]interface{}{
			"name":  name,
			"value": value,
		})
	}
	return results
}

func flattenGuestConfigurationAssignmentConfigurationSetting(input *guestconfiguration.ConfigurationSetting) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	configurationMode := input.ConfigurationMode

	var configurationModeFrequencyMins float64
	if input.ConfigurationModeFrequencyMins != nil {
		configurationModeFrequencyMins = *input.ConfigurationModeFrequencyMins
	}

	var refreshFrequencyMinute float64
	if input.RefreshFrequencyMins != nil {
		refreshFrequencyMinute = *input.RefreshFrequencyMins
	}
	return []interface{}{
		map[string]interface{}{
			"action_after_reboot":               input.ActionAfterReboot,
			"allow_module_overwrite":            input.AllowModuleOverwrite == guestconfiguration.True,
			"configuration_mode":                configurationMode,
			"configuration_mode_frequency_mins": configurationModeFrequencyMins,
			"reboot_if_needed":                  input.RebootIfNeeded == guestconfiguration.RebootIfNeededTrue,
			"refresh_frequency_in_minute":       refreshFrequencyMinute,
		},
	}
}
