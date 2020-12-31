package iothub

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotHubDPSSharedAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Read:   resourceIotHubDPSSharedAccessPolicyRead,
		Update: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Delete: resourceIotHubDPSSharedAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IotHubSharedAccessPolicyName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_dps_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"enrollment_read": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enrollment_write": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"registration_read": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"registration_write": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_config": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"primary_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceIotHubDPSSharedAccessPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubDpsName := d.Get("iothub_dps_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(iothubDpsName, IothubResourceName)
	defer locks.UnlockByName(iothubDpsName, IothubResourceName)

	iothubDps, err := client.Get(ctx, iothubDpsName, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(iothubDps.Response) {
			return fmt.Errorf("IotHub DPS %q (Resource Group %q) was not found", iothubDpsName, resourceGroup)
		}

		return fmt.Errorf("Error retrieving IotHub DPS %q (Resource Group %q): %+v", iothubDpsName, resourceGroup, err)
	}

	if iothubDps.ID == nil || *iothubDps.ID == "" {
		return fmt.Errorf("Error retrieving IotHub DPS %q (Resource Group %q): ID was nil", iothubDpsName, resourceGroup)
	}

	keyName := d.Get("name").(string)
	resourceID := fmt.Sprintf("%s/keys/%s", *iothubDps.ID, keyName)

	accessRights := dpsAccessRights{
		enrollmentRead:    d.Get("enrollment_read").(bool),
		enrollmentWrite:   d.Get("enrollment_write").(bool),
		registrationRead:  d.Get("registration_read").(bool),
		registrationWrite: d.Get("registration_write").(bool),
		serviceConfig:     d.Get("service_config").(bool),
	}

	if err := accessRights.validate(); err != nil {
		return fmt.Errorf("Error building Access Rights: %s", err)
	}

	expandedAccessPolicy := iothub.SharedAccessSignatureAuthorizationRuleAccessRightsDescription{
		KeyName: &keyName,
		Rights:  iothub.AccessRightsDescription(expandDpsAccessRights(accessRights)),
	}

	accessPolicies := make([]iothub.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	alreadyExists := false
	for accessPolicyIterator, err := client.ListKeysComplete(ctx, iothubDpsName, resourceGroup); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("Error loading Shared Access Policies of IotHub DPS %q (Resource Group %q): %+v", iothubDpsName, resourceGroup, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if strings.EqualFold(*existingAccessPolicy.KeyName, keyName) {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_iothub_dps_shared_access_policy", resourceID)
			}
			accessPolicies = append(accessPolicies, expandedAccessPolicy)
			alreadyExists = true
		} else {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	if d.IsNewResource() {
		accessPolicies = append(accessPolicies, expandedAccessPolicy)
	} else if !alreadyExists {
		return fmt.Errorf("Unable to find Shared Access Policy %q defined for IotHub DPS %q (Resource Group %q)", keyName, iothubDpsName, resourceGroup)
	}

	iothubDps.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubDpsName, iothubDps)
	if err != nil {
		return fmt.Errorf("Error updating IotHub DPS %q (Resource Group %q) with Shared Access Policy %q: %+v", iothubDpsName, resourceGroup, keyName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub DPS %q (Resource Group %q) to finish updating Shared Access Policy %q: %+v", iothubDpsName, resourceGroup, keyName, err)
	}

	d.SetId(resourceID)

	return resourceIotHubDPSSharedAccessPolicyRead(d, meta)
}

func resourceIotHubDPSSharedAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	iothubDpsName := id.Path["provisioningServices"]
	keyName := id.Path["keys"]

	iothubDps, err := client.Get(ctx, iothubDpsName, resourceGroup)
	if err != nil {
		return fmt.Errorf("Error retrieving IotHub DPS %q (Resource Group %q): %+v", iothubDpsName, resourceGroup, err)
	}

	accessPolicy, err := client.ListKeysForKeyName(ctx, iothubDpsName, keyName, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			log.Printf("[DEBUG] Shared Access Policy %q was not found on IotHub DPS %q (Resource Group %q) - removing from state", keyName, iothubDpsName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading Shared Access Policy %q (IotHub DPS %q / Resource Group %q): %+v", keyName, iothubDpsName, resourceGroup, err)
	}

	d.Set("name", keyName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("primary_key", accessPolicy.PrimaryKey)
	d.Set("secondary_key", accessPolicy.SecondaryKey)

	primaryConnectionString := ""
	secondaryConnectionString := ""
	if iothubDps.Properties != nil && iothubDps.Properties.ServiceOperationsHostName != nil {
		hostname := iothubDps.Properties.ServiceOperationsHostName
		if primary := accessPolicy.PrimaryKey; primary != nil {
			primaryConnectionString = getSAPConnectionString(*hostname, keyName, *primary)
		}
		if secondary := accessPolicy.SecondaryKey; secondary != nil {
			secondaryConnectionString = getSAPConnectionString(*hostname, keyName, *secondary)
		}
	}
	d.Set("primary_connection_string", primaryConnectionString)
	d.Set("secondary_connection_string", secondaryConnectionString)

	rights := flattenDpsAccessRights(accessPolicy.Rights)
	d.Set("enrollment_read", rights.enrollmentRead)
	d.Set("enrollment_write", rights.enrollmentWrite)
	d.Set("registration_read", rights.registrationRead)
	d.Set("registration_write", rights.registrationWrite)
	d.Set("service_config", rights.serviceConfig)

	return nil
}

func resourceIotHubDPSSharedAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	iothubDpsName := id.Path["provisioningServices"]
	keyName := id.Path["keys"]

	locks.ByName(iothubDpsName, IothubResourceName)
	defer locks.UnlockByName(iothubDpsName, IothubResourceName)

	iothubDps, err := client.Get(ctx, iothubDpsName, resourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(iothubDps.Response) {
			return fmt.Errorf("IotHub DPS %q (Resource Group %q) was not found", iothubDpsName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub DPS %q (Resource Group %q): %+v", iothubDpsName, resourceGroup, err)
	}

	accessPolicies := make([]iothub.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	for accessPolicyIterator, err := client.ListKeysComplete(ctx, iothubDpsName, resourceGroup); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("Error loading Shared Access Policies of IotHub DPS %q (Resource Group %q): %+v", iothubDpsName, resourceGroup, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if existingAccessPolicy.KeyName == nil {
			continue
		}

		if !strings.EqualFold(*existingAccessPolicy.KeyName, keyName) {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	iothubDps.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubDpsName, iothubDps)
	if err != nil {
		return fmt.Errorf("Error updating IotHub DPS %q (Resource Group %q) with Shared Access Policy %q: %+v", iothubDpsName, resourceGroup, keyName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub DPS %q (Resource Group %q) to finish updating Shared Access Policy %q: %+v", iothubDpsName, resourceGroup, keyName, err)
	}

	return nil
}

type dpsAccessRights struct {
	enrollmentRead    bool
	enrollmentWrite   bool
	registrationRead  bool
	registrationWrite bool
	serviceConfig     bool
}

func (r dpsAccessRights) validate() error {
	var err error

	if !r.enrollmentRead && !r.enrollmentWrite && !r.registrationRead && !r.registrationWrite && !r.serviceConfig {
		err = multierror.Append(err, fmt.Errorf("At least one of `enrollment_read`, `enrollment_write`, `registration_read`, `registration_write` , or `service_config` properties must be set to true"))
	}

	if r.enrollmentRead && !r.registrationRead {
		err = multierror.Append(err, fmt.Errorf("If `enrollment_read` is set to true, `registration_read` must also be set to true"))
	}

	if r.registrationWrite && !r.registrationRead {
		err = multierror.Append(err, fmt.Errorf("If `registration_write` is set to true, `registration_read` must also be set to true"))
	}

	if r.enrollmentWrite && !r.enrollmentRead && !r.registrationRead && !r.registrationWrite {
		err = multierror.Append(err, fmt.Errorf("If `enrollment_write` is set to true, `enrollment_read`, `registration_read`, and `registration_write` must also be set to true"))
	}

	return err
}

func expandDpsAccessRights(input dpsAccessRights) string {
	actualRights := make([]string, 0)

	// NOTE: the iteration order is important here
	if input.enrollmentRead {
		actualRights = append(actualRights, "EnrollmentRead")
	}

	if input.enrollmentWrite {
		actualRights = append(actualRights, "EnrollmentWrite")
	}

	if input.registrationRead {
		actualRights = append(actualRights, "RegistrationStatusRead")
	}

	if input.registrationWrite {
		actualRights = append(actualRights, "RegistrationStatusWrite")
	}

	if input.serviceConfig {
		actualRights = append(actualRights, "ServiceConfig")
	}

	return strings.Join(actualRights, ", ")
}

func flattenDpsAccessRights(r iothub.AccessRightsDescription) dpsAccessRights {
	rights := dpsAccessRights{
		enrollmentRead:    false,
		enrollmentWrite:   false,
		registrationRead:  false,
		registrationWrite: false,
		serviceConfig:     false,
	}

	actualAccessRights := strings.Split(string(r), ",")

	for _, right := range actualAccessRights {
		switch strings.ToLower(strings.Trim(right, " ")) {
		case "enrollmentread":
			rights.enrollmentRead = true
		case "enrollmentwrite":
			rights.enrollmentWrite = true
		case "registrationstatusread":
			rights.registrationRead = true
		case "registrationstatuswrite":
			rights.registrationWrite = true
		case "serviceconfig":
			rights.serviceConfig = true
		}
	}

	return rights
}

func getSAPConnectionString(iothubDpsHostName string, keyName string, key string) string {
	return fmt.Sprintf("HostName=%s;SharedAccessKeyName=%s;SharedAccessKey=%s", iothubDpsHostName, keyName, key)
}
