package iothub

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/hashicorp/go-multierror"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotHubSharedAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubSharedAccessPolicyCreateUpdate,
		Read:   resourceIotHubSharedAccessPolicyRead,
		Update: resourceIotHubSharedAccessPolicyCreateUpdate,
		Delete: resourceIotHubSharedAccessPolicyDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"registry_read": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"registry_write": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"service_connect": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"device_connect": {
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
		CustomizeDiff: pluginsdk.CustomizeDiffShim(iothubSharedAccessPolicyCustomizeDiff),
	}
}

func iothubSharedAccessPolicyCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ interface{}) (err error) {
	registryRead, hasRegistryRead := d.GetOk("registry_read")
	registryWrite, hasRegistryWrite := d.GetOk("registry_write")
	serviceConnect, hasServieConnect := d.GetOk("service_connect")
	deviceConnect, hasDeviceConnect := d.GetOk("device_connect")

	if !hasRegistryRead && !hasRegistryWrite && !hasServieConnect && !hasDeviceConnect {
		return fmt.Errorf("One of `registry_read`, `registry_write`, `service_connect` or `device_connect` properties must be set")
	}

	if !registryRead.(bool) && !registryWrite.(bool) && !serviceConnect.(bool) && !deviceConnect.(bool) {
		err = multierror.Append(err, fmt.Errorf("At least one of `registry_read`, `registry_write`, `service_connect` or `device_connect` properties must be set to true"))
	}

	if registryWrite.(bool) && !registryRead.(bool) {
		err = multierror.Append(err, fmt.Errorf("If `registry_write` is set to true, `registry_read` must also be set to true"))
	}

	return
}

func resourceIotHubSharedAccessPolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	locks.ByName(iothubName, IothubResourceName)
	defer locks.UnlockByName(iothubName, IothubResourceName)

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
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_iothub_shared_access_policy", resourceId)
			}

			if existingAccessPolicy.PrimaryKey != nil {
				expandedAccessPolicy.PrimaryKey = existingAccessPolicy.PrimaryKey
			}

			if existingAccessPolicy.SecondaryKey != nil {
				expandedAccessPolicy.SecondaryKey = existingAccessPolicy.SecondaryKey
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

	return resourceIotHubSharedAccessPolicyRead(d, meta)
}

func resourceIotHubSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedIothubSAPId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedIothubSAPId.ResourceGroup
	iothubName := parsedIothubSAPId.Path["IotHubs"]
	keyName := parsedIothubSAPId.Path["IotHubKeys"]

	accessPolicy, err := client.GetKeysForKeyName(ctx, resourceGroup, iothubName, keyName)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			log.Printf("[DEBUG] Shared Access Policy %q was not found on IotHub %q (Resource Group %q) - removing from state", keyName, iothubName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error loading IotHub Shared Access Policy %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.Set("name", keyName)
	d.Set("iothub_name", iothubName)
	d.Set("resource_group_name", resourceGroup)

	d.Set("primary_key", accessPolicy.PrimaryKey)
	if err := d.Set("primary_connection_string", getSharedAccessPolicyConnectionString(*iothub.Properties.HostName, keyName, *accessPolicy.PrimaryKey)); err != nil {
		return fmt.Errorf("error setting `primary_connection_string`: %v", err)
	}
	d.Set("secondary_key", accessPolicy.SecondaryKey)
	if err := d.Set("secondary_connection_string", getSharedAccessPolicyConnectionString(*iothub.Properties.HostName, keyName, *accessPolicy.SecondaryKey)); err != nil {
		return fmt.Errorf("error setting `secondary_connection_string`: %v", err)
	}

	rights := flattenAccessRights(accessPolicy.Rights)
	d.Set("registry_read", rights.registryRead)
	d.Set("registry_write", rights.registryWrite)
	d.Set("device_connect", rights.deviceConnect)
	d.Set("service_connect", rights.serviceConnect)

	return nil
}

func resourceIotHubSharedAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedIothubSAPId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedIothubSAPId.ResourceGroup
	iothubName := parsedIothubSAPId.Path["IotHubs"]
	keyName := parsedIothubSAPId.Path["IotHubKeys"]

	locks.ByName(iothubName, IothubResourceName)
	defer locks.UnlockByName(iothubName, IothubResourceName)

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

type accessRights struct {
	registryRead   bool
	registryWrite  bool
	serviceConnect bool
	deviceConnect  bool
}

func expandAccessRights(d *pluginsdk.ResourceData) string {
	possibleAccessRights := []struct {
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

func getSharedAccessPolicyConnectionString(iothubHostName string, keyName string, key string) string {
	return fmt.Sprintf("HostName=%s;SharedAccessKeyName=%s;SharedAccessKey=%s", iothubHostName, keyName, key)
}
