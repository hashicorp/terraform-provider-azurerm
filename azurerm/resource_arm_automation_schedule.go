package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutomationSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutomationScheduleCreateUpdate,
		Read:   resourceArmAutomationScheduleRead,
		Update: resourceArmAutomationScheduleCreateUpdate,
		Delete: resourceArmAutomationScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"start_time": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareDataAsUTCSuppressFunc,
				ValidateFunc:     validateStartTime,
			},

			"expiry_time": {
				Type:             schema.TypeString,
				DiffSuppressFunc: compareDataAsUTCSuppressFunc,
				Optional:         true,
				Computed:         true,
			},

			"frequency": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.Day),
					string(automation.Hour),
					string(automation.Month),
					string(automation.OneTime),
					string(automation.Week),
				}, true),
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
			},
		},
	}
}

func compareDataAsUTCSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	ot, oerr := time.Parse(time.RFC3339, old)
	nt, nerr := time.Parse(time.RFC3339, new)
	if oerr != nil || nerr != nil {
		return false
	}

	return nt.Equal(ot)
}

func validateStartTime(v interface{}, k string) (ws []string, errors []error) {
	starttime, tperr := time.Parse(time.RFC3339, v.(string))
	if tperr != nil {
		errors = append(errors, fmt.Errorf("Cannot parse %q", k))
	}

	u := time.Until(starttime)
	if u < 5*time.Minute {
		errors = append(errors, fmt.Errorf("%q should be at least 5 minutes in the future", k))
	}

	return
}

func resourceArmAutomationScheduleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationScheduleClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Automation Schedule creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	client.ResourceGroupName = resGroup

	accName := d.Get("account_name").(string)
	freqstr := d.Get("frequency").(string)
	freq := automation.ScheduleFrequency(freqstr)

	cst := d.Get("start_time").(string)
	starttime, tperr := time.Parse(time.RFC3339, cst)
	if tperr != nil {
		return fmt.Errorf("Cannot parse start_time: %q", cst)
	}

	expirytime, teperr := time.Parse(time.RFC3339, cst)
	if teperr != nil {
		return fmt.Errorf("Cannot parse expiry_time: %q", cst)
	}

	stdt := date.Time{Time: starttime}
	etdt := date.Time{Time: expirytime}

	description := d.Get("description").(string)
	timezone := d.Get("timezone").(string)

	parameters := automation.ScheduleCreateOrUpdateParameters{
		Name: &name,
		ScheduleCreateOrUpdateProperties: &automation.ScheduleCreateOrUpdateProperties{
			Description: &description,
			Frequency:   freq,
			StartTime:   &stdt,
			ExpiryTime:  &etdt,
			TimeZone:    &timezone,
		},
	}

	_, err := client.CreateOrUpdate(ctx, accName, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, accName, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Automation Schedule '%s' (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmAutomationScheduleRead(d, meta)
}

func resourceArmAutomationScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationScheduleClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	client.ResourceGroupName = resGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["schedules"]

	resp, err := client.Get(ctx, accName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on AzureRM Automation Schedule '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("account_name", accName)
	d.Set("frequency", resp.Frequency)
	d.Set("description", resp.Description)
	d.Set("start_time", string(resp.StartTime.Format(time.RFC3339)))
	d.Set("expiry_time", string(resp.ExpiryTime.Format(time.RFC3339)))
	d.Set("timezone", resp.TimeZone)
	return nil
}

func resourceArmAutomationScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).automationScheduleClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	client.ResourceGroupName = resGroup
	accName := id.Path["automationAccounts"]
	name := id.Path["schedules"]

	resp, err := client.Delete(ctx, accName, name)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Automation Schedule '%s': %+v", name, err)
	}

	return nil
}
