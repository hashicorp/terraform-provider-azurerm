// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubSharedAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubSharedAccessPolicyCreateUpdate,
		Read:   resourceIotHubSharedAccessPolicyRead,
		Update: resourceIotHubSharedAccessPolicyCreateUpdate,
		Delete: resourceIotHubSharedAccessPolicyDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubSharedAccessPolicyV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SharedAccessPolicyID(id)
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

			"resource_group_name": commonschema.ResourceGroupName(),

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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewSharedAccessPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	expandedAccessPolicy := devices.SharedAccessSignatureAuthorizationRule{
		KeyName: &id.IotHubKeyName,
		Rights:  devices.AccessRights(expandAccessRights(d)),
	}

	accessPolicies := make([]devices.SharedAccessSignatureAuthorizationRule, 0)

	alreadyExists := false
	for accessPolicyIterator, err := client.ListKeysComplete(ctx, id.ResourceGroup, id.IotHubName); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading %s: %+v", id, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if strings.EqualFold(*existingAccessPolicy.KeyName, id.IotHubKeyName) {
			if d.IsNewResource() {
				return tf.ImportAsExistsError("azurerm_iothub_shared_access_policy", id.ID())
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
		return fmt.Errorf("dnable to find %s", id)
	}

	iothub.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubSharedAccessPolicyRead(d, meta)
}

func resourceIotHubSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	accessPolicy, err := client.GetKeysForKeyName(ctx, id.ResourceGroup, id.IotHubName, id.IotHubKeyName)
	if err != nil {
		if utils.ResponseWasNotFound(accessPolicy.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("loading %s: %+v", id, err)
	}

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	d.Set("name", id.IotHubKeyName)
	d.Set("iothub_name", id.IotHubName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("primary_key", accessPolicy.PrimaryKey)
	if err := d.Set("primary_connection_string", getSharedAccessPolicyConnectionString(*iothub.Properties.HostName, id.IotHubKeyName, *accessPolicy.PrimaryKey)); err != nil {
		return fmt.Errorf("setting `primary_connection_string`: %v", err)
	}
	d.Set("secondary_key", accessPolicy.SecondaryKey)
	if err := d.Set("secondary_connection_string", getSharedAccessPolicyConnectionString(*iothub.Properties.HostName, id.IotHubKeyName, *accessPolicy.SecondaryKey)); err != nil {
		return fmt.Errorf("setting `secondary_connection_string`: %v", err)
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

	id, err := parse.SharedAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	accessPolicies := make([]devices.SharedAccessSignatureAuthorizationRule, 0)

	for accessPolicyIterator, err := client.ListKeysComplete(ctx, id.ResourceGroup, id.IotHubName); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
		if err != nil {
			return fmt.Errorf("loading %s: %+v", id, err)
		}
		existingAccessPolicy := accessPolicyIterator.Value()

		if !strings.EqualFold(*existingAccessPolicy.KeyName, id.IotHubKeyName) {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	iothub.Properties.AuthorizationPolicies = &accessPolicies

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
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
