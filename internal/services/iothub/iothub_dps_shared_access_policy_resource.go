// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotHubDPSSharedAccessPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Read:   resourceIotHubDPSSharedAccessPolicyRead,
		Update: resourceIotHubDPSSharedAccessPolicyCreateUpdate,
		Delete: resourceIotHubDPSSharedAccessPolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := iotdpsresource.ParseKeyID(id)
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	iothubDpsId := commonids.NewProvisioningServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_dps_name").(string))

	locks.ByName(iothubDpsId.ProvisioningServiceName, IothubResourceName)
	defer locks.UnlockByName(iothubDpsId.ProvisioningServiceName, IothubResourceName)

	iothubDps, err := client.Get(ctx, iothubDpsId)
	if err != nil {
		if response.WasNotFound(iothubDps.HttpResponse) {
			return fmt.Errorf("IotHub DPS %s was not found", iothubDpsId.String())
		}

		return fmt.Errorf("retrieving IotHub DPS %s: %+v", iothubDpsId.String(), err)
	}

	if iothubDps.Model == nil {
		return fmt.Errorf("retrieving IotHub DPS %s: ID was nil", iothubDpsId.String())
	}

	id := iotdpsresource.NewKeyID(iothubDpsId.SubscriptionId, iothubDpsId.ResourceGroupName, iothubDpsId.ProvisioningServiceName, d.Get("name").(string))

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

	expandedAccessPolicy := iotdpsresource.SharedAccessSignatureAuthorizationRuleAccessRightsDescription{
		KeyName: id.KeyName,
		Rights:  iotdpsresource.AccessRightsDescription(expandDpsAccessRights(accessRights)),
	}

	accessPolicies := make([]iotdpsresource.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	alreadyExists := false
	existingAccessPolicies, err := client.ListKeysComplete(ctx, iothubDpsId)
	if err != nil {
		return fmt.Errorf("loading %s: %+v", id, err)
	}
	for _, existingAccessPolicy := range existingAccessPolicies.Items {
		if strings.EqualFold(existingAccessPolicy.KeyName, id.KeyName) {
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

	iothubDps.Model.Properties.AuthorizationPolicies = &accessPolicies

	if err := client.CreateOrUpdateThenPoll(ctx, iothubDpsId, *iothubDps.Model); err != nil {
		return fmt.Errorf("updating IotHub DPS %s with Shared Access Policy %s: %+v", iothubDpsId, id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubDPSSharedAccessPolicyRead(d, meta)
}

func resourceIotHubDPSSharedAccessPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceIotHubDPSSharedAccessPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.DPSResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotdpsresource.ParseKeyID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ProvisioningServiceName, IothubResourceName)
	defer locks.UnlockByName(id.ProvisioningServiceName, IothubResourceName)

	iothubDpsId := commonids.NewProvisioningServiceID(id.SubscriptionId, id.ResourceGroupName, id.ProvisioningServiceName)
	iothubDps, err := client.Get(ctx, iothubDpsId)
	if err != nil {
		if response.WasNotFound(iothubDps.HttpResponse) {
			return fmt.Errorf("IotHub DPS %q was not found", id)
		}

		return fmt.Errorf("loading IotHub DPS %q: %+v", id, err)
	}

	accessPolicies := make([]iotdpsresource.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, 0)

	existingAccessPolicies, err := client.ListKeysComplete(ctx, iothubDpsId)
	if err != nil {
		return fmt.Errorf("loading %s: %+v", id, err)
	}
	for _, existingAccessPolicy := range existingAccessPolicies.Items {
		if !strings.EqualFold(existingAccessPolicy.KeyName, id.KeyName) {
			accessPolicies = append(accessPolicies, existingAccessPolicy)
		}
	}

	iothubDps.Model.Properties.AuthorizationPolicies = &accessPolicies

	if err := client.CreateOrUpdateThenPoll(ctx, iothubDpsId, *iothubDps.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
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

func flattenDpsAccessRights(r iotdpsresource.AccessRightsDescription) dpsAccessRights {
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
