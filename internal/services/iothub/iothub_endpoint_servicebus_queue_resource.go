// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iothub

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubEndpointServiceBusQueue() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubEndpointServiceBusQueueCreateUpdate,
		Read:   resourceIotHubEndpointServiceBusQueueRead,
		Update: resourceIotHubEndpointServiceBusQueueCreateUpdate,
		Delete: resourceIotHubEndpointServiceBusQueueDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubServiceBusQueueV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EndpointServiceBusQueueID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceIothubEndpointServicebusQueue(),
	}
}

func resourceIothubEndpointServicebusQueue() map[string]*pluginsdk.Schema {
	out := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IoTHubEndpointName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// lintignore: S013
		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IotHubID,
		},

		"authentication_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(devices.AuthenticationTypeKeyBased),
			ValidateFunc: validation.StringInSlice([]string{
				string(devices.AuthenticationTypeKeyBased),
				string(devices.AuthenticationTypeIdentityBased),
			}, false),
		},

		"identity_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			ConflictsWith: []string{"connection_string"},
		},

		"endpoint_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"entity_path"},
			ExactlyOneOf: []string{"endpoint_uri", "connection_string"},
		},

		"entity_path": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.QueueName(),
			RequiredWith: []string{"endpoint_uri"},
		},

		"connection_string": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				sharedAccessKeyRegex := regexp.MustCompile("SharedAccessKey=[^;]+")
				sbProtocolRegex := regexp.MustCompile("sb://([^:]+)(:5671)?/;")

				maskedNew := sbProtocolRegex.ReplaceAllString(new, "sb://$1:5671/;")
				maskedNew = sharedAccessKeyRegex.ReplaceAllString(maskedNew, "SharedAccessKey=****")
				return (new == d.Get(k).(string)) && (maskedNew == old)
			},
			Sensitive:     true,
			ConflictsWith: []string{"identity_id"},
			ExactlyOneOf:  []string{"endpoint_uri", "connection_string"},
		},
	}

	return out
}

func resourceIotHubEndpointServiceBusQueueCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.ResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	endpointRG := d.Get("resource_group_name").(string)

	iotHubId, err := parse.IotHubID(d.Get("iothub_id").(string))
	if err != nil {
		return err
	}
	iotHubName := iotHubId.Name
	iotHubRG := iotHubId.ResourceGroup

	id := parse.NewEndpointServiceBusQueueID(subscriptionId, iotHubRG, iotHubName, d.Get("name").(string))

	locks.ByName(iotHubName, IothubResourceName)
	defer locks.UnlockByName(iotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, iotHubRG, iotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iotHubName, iotHubRG)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", iotHubName, iotHubRG, err)
	}

	authenticationType := devices.AuthenticationType(d.Get("authentication_type").(string))

	queueEndpoint := devices.RoutingServiceBusQueueEndpointProperties{
		AuthenticationType: authenticationType,
		Name:               utils.String(id.EndpointName),
		SubscriptionID:     utils.String(subscriptionID),
		ResourceGroup:      utils.String(endpointRG),
	}

	if authenticationType == devices.AuthenticationTypeKeyBased {
		if v, ok := d.GetOk("connection_string"); ok {
			queueEndpoint.ConnectionString = utils.String(v.(string))
		} else {
			return fmt.Errorf("`connection_string` must be specified when `authentication_type` is `keyBased`")
		}
	} else {
		if v, ok := d.GetOk("endpoint_uri"); ok {
			queueEndpoint.EndpointURI = utils.String(v.(string))
			queueEndpoint.EntityPath = utils.String(d.Get("entity_path").(string))
		} else {
			return fmt.Errorf("`endpoint_uri` and `entity_path` must be specified when `authentication_type` is `identityBased`")
		}

		if v, ok := d.GetOk("identity_id"); ok {
			queueEndpoint.Identity = &devices.ManagedIdentity{
				UserAssignedIdentity: utils.String(v.(string)),
			}
		}
	}

	routing := iothub.Properties.Routing
	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	if routing.Endpoints == nil {
		routing.Endpoints = &devices.RoutingEndpoints{}
	}

	if routing.Endpoints.EventHubs == nil {
		queues := make([]devices.RoutingServiceBusQueueEndpointProperties, 0)
		routing.Endpoints.ServiceBusQueues = &queues
	}
	endpoints := make([]devices.RoutingServiceBusQueueEndpointProperties, 0)

	alreadyExists := false
	for _, existingEndpoint := range *routing.Endpoints.ServiceBusQueues {
		if existingEndpointName := existingEndpoint.Name; existingEndpointName != nil {
			if strings.EqualFold(*existingEndpointName, id.EndpointName) {
				if d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_iothub_endpoint_servicebus_queue", id.ID())
				}
				endpoints = append(endpoints, queueEndpoint)
				alreadyExists = true
			} else {
				endpoints = append(endpoints, existingEndpoint)
			}
		}
	}

	if d.IsNewResource() {
		endpoints = append(endpoints, queueEndpoint)
	} else if !alreadyExists {
		return fmt.Errorf("unable to find %s", id)
	}
	routing.Endpoints.ServiceBusQueues = &endpoints

	future, err := client.CreateOrUpdate(ctx, iotHubRG, iotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubEndpointServiceBusQueueRead(d, meta)
}

func resourceIotHubEndpointServiceBusQueueRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointServiceBusQueueID(d.Id())
	if err != nil {
		return err
	}

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	d.Set("name", id.EndpointName)

	iotHubId := parse.NewIotHubID(id.SubscriptionId, id.ResourceGroup, id.IotHubName)
	d.Set("iothub_id", iotHubId.ID())

	if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		d.SetId("")
		return nil
	}

	exist := false

	if endpoints := iothub.Properties.Routing.Endpoints.ServiceBusQueues; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, id.EndpointName) {
					exist = true
					d.Set("resource_group_name", endpoint.ResourceGroup)

					authenticationType := string(devices.AuthenticationTypeKeyBased)
					if string(endpoint.AuthenticationType) != "" {
						authenticationType = string(endpoint.AuthenticationType)
					}
					d.Set("authentication_type", authenticationType)

					connectionStr := ""
					if endpoint.ConnectionString != nil {
						connectionStr = *endpoint.ConnectionString
					}
					d.Set("connection_string", connectionStr)

					endpointUri := ""
					if endpoint.EndpointURI != nil {
						endpointUri = *endpoint.EndpointURI
					}
					d.Set("endpoint_uri", endpointUri)

					entityPath := ""
					if endpoint.EntityPath != nil {
						entityPath = *endpoint.EntityPath
					}
					d.Set("entity_path", entityPath)

					identityId := ""
					if endpoint.Identity != nil && endpoint.Identity.UserAssignedIdentity != nil {
						identityId = *endpoint.Identity.UserAssignedIdentity
					}
					d.Set("identity_id", identityId)
				}
			}
		}
	}

	if !exist {
		d.SetId("")
	}

	return nil
}

func resourceIotHubEndpointServiceBusQueueDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointServiceBusQueueID(d.Id())
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

	if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		return nil
	}
	endpoints := iothub.Properties.Routing.Endpoints.ServiceBusQueues

	if endpoints == nil {
		return nil
	}

	updatedEndpoints := make([]devices.RoutingServiceBusQueueEndpointProperties, 0)
	for _, endpoint := range *endpoints {
		if existingEndpointName := endpoint.Name; existingEndpointName != nil {
			if !strings.EqualFold(*existingEndpointName, id.EndpointName) {
				updatedEndpoints = append(updatedEndpoints, endpoint)
			}
		}
	}

	iothub.Properties.Routing.Endpoints.ServiceBusQueues = &updatedEndpoints

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}
