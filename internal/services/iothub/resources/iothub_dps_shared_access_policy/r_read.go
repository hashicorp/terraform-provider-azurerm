// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_shared_access_policy

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotdpsresource.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	iothubDpsId := commonids.NewProvisioningServiceID(id.SubscriptionId, id.ResourceGroupName, id.ProvisioningServiceName)
	iothubDps, err := client.Get(ctx, iothubDpsId)
	if err != nil {
		return fmt.Errorf("retrieving IotHub DPS %q: %+v", id, err)
	}

	accessPolicy, err := client.ListKeysForKeyName(ctx, *id)
	if err != nil {
		if response.WasNotFound(accessPolicy.HttpResponse) {
			log.Printf("[DEBUG] %s - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("loading %s: %+v", id, err)
	}

	d.Set("name", id.KeyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := accessPolicy.Model; model != nil {
		d.Set("primary_key", model.PrimaryKey)
		d.Set("secondary_key", model.SecondaryKey)

		if iothubDpsModel := iothubDps.Model; iothubDpsModel != nil {
			primaryConnectionString := ""
			secondaryConnectionString := ""
			properties := iothubDpsModel.Properties
			if properties.ServiceOperationsHostName != nil {
				hostname := properties.ServiceOperationsHostName
				if primary := model.PrimaryKey; primary != nil {
					primaryConnectionString = getSAPConnectionString(*hostname, id.KeyName, *primary)
				}
				if secondary := model.SecondaryKey; secondary != nil {
					secondaryConnectionString = getSAPConnectionString(*hostname, id.KeyName, *secondary)
				}
			}
			d.Set("primary_connection_string", primaryConnectionString)
			d.Set("secondary_connection_string", secondaryConnectionString)
		}

		rights := flattenDpsAccessRights(model.Rights)
		d.Set("enrollment_read", rights.enrollmentRead)
		d.Set("enrollment_write", rights.enrollmentWrite)
		d.Set("registration_read", rights.registrationRead)
		d.Set("registration_write", rights.registrationWrite)
		d.Set("service_config", rights.serviceConfig)
	}

	return nil
}

func flattenDpsAccessRights(r iotdpsresource.AccessRightsDescription) dpsAccessRights {
	rights := dpsAccessRights{}

	actualAccessRights := strings.SplitSeq(string(r), ",")

	for right := range actualAccessRights {
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
		err = multierror.Append(err, errors.New("at least one of `enrollment_read`, `enrollment_write`, `registration_read`, `registration_write` , or `service_config` properties must be set to true"))
	}

	if r.enrollmentRead && !r.registrationRead {
		err = multierror.Append(err, errors.New("if `enrollment_read` is set to true, `registration_read` must also be set to true"))
	}

	if r.registrationWrite && !r.registrationRead {
		err = multierror.Append(err, errors.New("if `registration_write` is set to true, `registration_read` must also be set to true"))
	}

	if r.enrollmentWrite && !r.enrollmentRead && !r.registrationRead && !r.registrationWrite {
		err = multierror.Append(err, errors.New("if `enrollment_write` is set to true, `enrollment_read`, `registration_read`, and `registration_write` must also be set to true"))
	}

	return err
}

func getSAPConnectionString(iothubDpsHostName string, keyName string, key string) string {
	return fmt.Sprintf("HostName=%s;SharedAccessKeyName=%s;SharedAccessKey=%s", iothubDpsHostName, keyName, key)
}
