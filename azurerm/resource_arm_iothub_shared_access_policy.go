package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/hashicorp/terraform/helper/customdiff"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"regexp"
	"strings"
)

func resourceArmIotHubSharedAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubSharedAccessPolicyCreateUpdate,
		Read:   resourceArmIotHubSharedAccessPolicyRead,
		Update: resourceArmIotHubSharedAccessPolicyCreateUpdate,
		Delete: resourceArmIotHubSharedAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"iothub_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.IoTHubName,
			},
			"resource_group_name": resourceGroupNameSchema(),
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`\S{1,64}`), ""+
					"The shared access policy key name must not be empty, and must not exceed 64 characters in length.  The shared access policy key name must not contain whitespace characters."),
			},

			"registry_read": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"registry_write": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service_connect": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"device_connect": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"primary_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
		CustomizeDiff: customdiff.All(
			validateAtLeastOneIsTrue(
				"registry_read",
				"registry_write",
				"service_connect",
				"device_connect",
			),
			validateRegistryWriteImpliesRegistryRead,
		),
	}
}

func validateRegistryWriteImpliesRegistryRead(d *schema.ResourceDiff, meta interface{}) error {
	if d.Get("registry_write").(bool) && !d.Get("registry_read").(bool) {
		return fmt.Errorf("registry_read key must be set to true when registry_write is set to true")
	}
	return nil

}
func validateAtLeastOneIsTrue(parameterNames ...string) schema.CustomizeDiffFunc {
	return func(d *schema.ResourceDiff, meta interface{}) error {
		atLeastOneSet := false
		for _, param := range parameterNames {
			atLeastOneSet = atLeastOneSet || d.Get(param).(bool)
		}

		if !atLeastOneSet {
			return fmt.Errorf("at least one of the properties: %s must be set to true", strings.Join(parameterNames, ", "))
		}

		return nil
	}
}

func resourceArmIotHubSharedAccessPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	azureRMLockByName(iothubName, iothubResourceName)
	defer azureRMUnlockByName(iothubName, iothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	keyName := d.Get("name").(string)

	resourceId := fmt.Sprintf("%s/IotHubKeys/%s", *iothub.ID, keyName)

	expandedAccessPolicy := devices.SharedAccessSignatureAuthorizationRule{
		KeyName: &keyName,
		Rights:  devices.AccessRights(expandAccessRights(d)),
	}

	accessPolicies := make([]devices.SharedAccessSignatureAuthorizationRule, 0)

	alreadyExists := false
	for accessPolicyIterator, err := client.ListKeysComplete(ctx, resourceGroup, iothubName); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("Error loading Shared Access Profiles of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if strings.EqualFold(*existingAccessPolicy.KeyName, keyName) {
			if d.IsNewResource() && requireResourcesToBeImported {
				return tf.ImportAsExistsError("azurerm_iothub_shared_access_policy", resourceId)
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
		return fmt.Errorf("Unable to find Shared Access Policy %q defined for IotHub %q (Resource Group %q)", keyName, iothubName, resourceGroup)
	}

	iothub.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error updating IotHub %q (Resource Group %q) with Shared Access Profile %q: %+v", iothubName, resourceGroup, keyName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub %q (Resource Group %q) to finish updating Shared Access Profile %q: %+v", iothubName, resourceGroup, keyName, err)
	}

	d.SetId(resourceId)

	return resourceArmIotHubSharedAccessPolicyRead(d, meta)
}

type accessRights struct {
	registryRead   bool
	registryWrite  bool
	serviceConnect bool
	deviceConnect  bool
}

func expandAccessRights(d *schema.ResourceData) string {

	var possibleAccessRights = []struct {
		schema string
		right  string
	}{
		{"registry_read", "RegistryRead"},
		{"registry_write", "RegistryWrite"},
		{"service_connect", "ServiceConnect"},
		{"device_connect", "DeviceConnect"},
	}
	actualRights := make([]string, 0)
	// iteration order is important here, so we cannot use a map
	for _, possibleRight := range possibleAccessRights {
		if d.Get(possibleRight.schema).(bool) {
			actualRights = append(actualRights, possibleRight.right)
		}
	}
	strRights := strings.Join(actualRights, ", ")
	return strRights
}

func flattenAccessRights(r devices.AccessRights) accessRights {
	rights := accessRights{
		registryRead:   false,
		registryWrite:  false,
		deviceConnect:  false,
		serviceConnect: false,
	}

	actualAccessRights := strings.Split(string(r), ",")

	for _, right := range actualAccessRights {
		switch strings.ToLower(strings.Trim(right, " ")) {
		case "registrywrite":
			rights.registryWrite = true
			// RegistryWrite implies RegistryRead.
			// What's more, creating a Access Policy with both RegistryRead and RegistryWrite
			// only really sets RegistryWrite permission, which then also implies RedistryRead
			fallthrough
		case "registryread":
			rights.registryRead = true
		case "deviceconnect":
			rights.deviceConnect = true
		case "serviceconnect":
			rights.serviceConnect = true
		}
	}

	return rights
}

func resourceArmIotHubSharedAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	parsedIothubId, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	resourceGroup := parsedIothubId.ResourceGroup
	iothubName := parsedIothubId.Path["IotHubs"]
	keyName := parsedIothubId.Path["IotHubKeys"]

	accessPolicy, err := client.GetKeysForKeyName(ctx, resourceGroup, iothubName, keyName)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			log.Printf("[DEBUG] Shared Access Policy %q was not found on IotHub %q (Resource Group %q) - removing from state", keyName, iothubName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	_ = d.Set("name", keyName)
	_ = d.Set("iothub_name", iothubName)
	_ = d.Set("resource_group_name", resourceGroup)

	_ = d.Set("primary_key", accessPolicy.PrimaryKey)
	_ = d.Set("secondary_key", accessPolicy.SecondaryKey)

	rights := flattenAccessRights(accessPolicy.Rights)
	_ = d.Set("registry_read", rights.registryRead)
	_ = d.Set("registry_write", rights.registryWrite)
	_ = d.Set("device_connect", rights.deviceConnect)
	_ = d.Set("service_connect", rights.serviceConnect)

	return nil
}

func resourceArmIotHubSharedAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothubResourceClient
	ctx := meta.(*ArmClient).StopContext

	parsedIothubId, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	resourceGroup := parsedIothubId.ResourceGroup
	iothubName := parsedIothubId.Path["IotHubs"]
	keyName := parsedIothubId.Path["IotHubKeys"]

	azureRMLockByName(iothubName, iothubResourceName)
	defer azureRMUnlockByName(iothubName, iothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	accessPolicies := make([]devices.SharedAccessSignatureAuthorizationRule, 0)

	for accessPolicyIterator, err := client.ListKeysComplete(ctx, resourceGroup, iothubName); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("Error loading Shared Access Profiles of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if !strings.EqualFold(*existingAccessPolicy.KeyName, keyName) {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	iothub.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error updating IotHub %q (Resource Group %q) with Shared Access Profile %q: %+v", iothubName, resourceGroup, keyName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub %q (Resource Group %q) to finish updating Shared Access Profile %q: %+v", iothubName, resourceGroup, keyName, err)
	}

	return nil
}
