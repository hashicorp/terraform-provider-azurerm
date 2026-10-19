// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_shared_access_policy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIotHubSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SharedAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	accessPolicy, err := client.GetKeysForKeyName(ctx, id.ResourceGroup, id.IotHubName, id.IotHubKeyName)
	if err != nil {
		if response.WasNotFound(accessPolicy.Response.Response) {
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

func flattenAccessRights(r devices.AccessRights) accessRights {
	rights := accessRights{}

	actualAccessRights := strings.SplitSeq(string(r), ",")

	for right := range actualAccessRights {
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

func iothubSharedAccessPolicyCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) (err error) {
	registryRead, hasRegistryRead := d.GetOk("registry_read")
	registryWrite, hasRegistryWrite := d.GetOk("registry_write")
	serviceConnect, hasServiceConnect := d.GetOk("service_connect")
	deviceConnect, hasDeviceConnect := d.GetOk("device_connect")

	if !hasRegistryRead && !hasRegistryWrite && !hasServiceConnect && !hasDeviceConnect {
		return errors.New("one of `registry_read`, `registry_write`, `service_connect` or `device_connect` properties must be set")
	}

	if !registryRead.(bool) && !registryWrite.(bool) && !serviceConnect.(bool) && !deviceConnect.(bool) {
		err = multierror.Append(err, errors.New("at least one of `registry_read`, `registry_write`, `service_connect` or `device_connect` properties must be set to true"))
	}

	if registryWrite.(bool) && !registryRead.(bool) {
		err = multierror.Append(err, errors.New("if `registry_write` is set to true, `registry_read` must also be set to true"))
	}

	return
}

type accessRights struct {
	registryRead   bool
	registryWrite  bool
	serviceConnect bool
	deviceConnect  bool
}

func getSharedAccessPolicyConnectionString(iothubHostName string, keyName string, key string) string {
	return fmt.Sprintf("HostName=%s;SharedAccessKeyName=%s;SharedAccessKey=%s", iothubHostName, keyName, key)
}
