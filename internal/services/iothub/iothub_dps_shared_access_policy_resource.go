package iothub

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2021-10-15/iothub"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotHubDPSSharedAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Read:   resourceIotHubDPSSharedAccessPolicyRead,
		Update: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Delete: resourceIotHubDPSSharedAccessPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DpsSharedAccessPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IotHubSharedAccessPolicyName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_dps_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"enrollment_read": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enrollment_write": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"registration_read": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"registration_write": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_config": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func resourceIotHubDPSSharedAccessPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.DPSResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubDpsId := parse.NewIotHubDpsID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_dps_name").(string))

	locks.ByName(iothubDpsId.ProvisioningServiceName, IothubResourceName)
	defer locks.UnlockByName(iothubDpsId.ProvisioningServiceName, IothubResourceName)

	iothubDps, err := client.Get(ctx, iothubDpsId.ProvisioningServiceName, iothubDpsId.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(iothubDps.Response) {
			return fmt.Errorf("IotHub DPS %s was not found", iothubDpsId.String())
		}

		return fmt.Errorf("retrieving IotHub DPS %s: %+v", iothubDpsId.String(), err)
	}

	if iothubDps.ID == nil || *iothubDps.ID == "" {
		return fmt.Errorf("retrieving IotHub DPS %s: ID was nil", iothubDpsId.String())
	}

	id := parse.NewDpsSharedAccessPolicyID(iothubDpsId.SubscriptionId, iothubDpsId.ResourceGroup, iothubDpsId.ProvisioningServiceName, d.Get("name").(string))

	accessRights := dpsAccessRights{
		enrollmentRead:    d.Get("enrollment_read").(bool),
		enrollmentWrite:   d.Get("enrollment_write").(bool),
		registrationRead:  d.Get("registration_read").(bool),
		registrationWrite: d.Get("registration_write").(bool),
		serviceConfig:     d.Get("service_config").(bool),
	}

	if err := accessRights.validate(); err != nil {
		return fmt.Errorf("building Access Rights: %s", err)
	}

	expandedAccessPolicy := iothub.SharedAccessSignatureAuthorizationRuleAccessRightsDescription{
		KeyName: &id.KeyName,
		Rights:  iothub.AccessRightsDescription(expandDpsAccessRights(accessRights)),
	}

	accessPolicies := make([]iothub.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	alreadyExists := false
	for accessPolicyIterator, err := client.ListKeysComplete(ctx, id.ProvisioningServiceName, id.ResourceGroup); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading %s: %+v", id, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if strings.EqualFold(*existingAccessPolicy.KeyName, id.KeyName) {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_iothub_dps_shared_access_policy", id.ID())
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
		return fmt.Errorf("unable to find %s", id)
	}

	iothubDps.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ProvisioningServiceName, iothubDps)
	if err != nil {
		return fmt.Errorf("updating IotHub DPS %s with Shared Access Policy %s: %+v", iothubDpsId, id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for IotHub DPS %s to finish updating Shared Access Policy %s: %+v", iothubDpsId, id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubDPSSharedAccessPolicyRead(d, meta)
}

func resourceIotHubDPSSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DpsSharedAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	iothubDps, err := client.Get(ctx, id.ProvisioningServiceName, id.ResourceGroup)
	if err != nil {
		return fmt.Errorf("retrieving IotHub DPS %q (Resource Group %q): %+v", id.ProvisioningServiceName, id.ResourceGroup, err)
	}

	accessPolicy, err := client.ListKeysForKeyName(ctx, id.ProvisioningServiceName, id.KeyName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			log.Printf("[DEBUG] %s - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("loading %s: %+v", id, err)
	}

	d.Set("name", id.KeyName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("primary_key", accessPolicy.PrimaryKey)
	d.Set("secondary_key", accessPolicy.SecondaryKey)

	primaryConnectionString := ""
	secondaryConnectionString := ""
	if iothubDps.Properties != nil && iothubDps.Properties.ServiceOperationsHostName != nil {
		hostname := iothubDps.Properties.ServiceOperationsHostName
		if primary := accessPolicy.PrimaryKey; primary != nil {
			primaryConnectionString = getSAPConnectionString(*hostname, id.KeyName, *primary)
		}
		if secondary := accessPolicy.SecondaryKey; secondary != nil {
			secondaryConnectionString = getSAPConnectionString(*hostname, id.KeyName, *secondary)
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

func resourceIotHubDPSSharedAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DpsSharedAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ProvisioningServiceName, IothubResourceName)
	defer locks.UnlockByName(id.ProvisioningServiceName, IothubResourceName)

	iothubDps, err := client.Get(ctx, id.ProvisioningServiceName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(iothubDps.Response) {
			return fmt.Errorf("IotHub DPS %q (Resource Group %q) was not found", id.ProvisioningServiceName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub DPS %q (Resource Group %q): %+v", id.ProvisioningServiceName, id.ResourceGroup, err)
	}

	accessPolicies := make([]iothub.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	for accessPolicyIterator, err := client.ListKeysComplete(ctx, id.ProvisioningServiceName, id.ResourceGroup); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading %s: %+v", id, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if existingAccessPolicy.KeyName == nil {
			continue
		}

		if !strings.EqualFold(*existingAccessPolicy.KeyName, id.KeyName) {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	iothubDps.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ProvisioningServiceName, iothubDps)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
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
