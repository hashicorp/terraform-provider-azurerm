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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	devices "github.com/tombuildsstuff/kermit/sdk/iothub/2022-04-30-preview/iothub"
)

func resourceIotHubEndpointStorageContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubEndpointStorageContainerCreateUpdate,
		Read:   resourceIotHubEndpointStorageContainerRead,
		Update: resourceIotHubEndpointStorageContainerCreateUpdate,
		Delete: resourceIotHubEndpointStorageContainerDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.IoTHubEndPointStorageContainerV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EndpointStorageContainerID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceIothubEndpointStorageContainerSchema(),
	}
}

func resourceIothubEndpointStorageContainerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IoTHubEndpointName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		//lintignore: S013
		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IotHubID,
		},

		"container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageContainerName,
		},

		"file_name_format": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}",
			ValidateFunc: iothubValidate.FileNameFormat,
		},

		"batch_frequency_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      300,
			ValidateFunc: validation.IntBetween(60, 720),
		},

		"max_chunk_size_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      314572800,
			ValidateFunc: validation.IntBetween(10485760, 524288000),
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
			ExactlyOneOf: []string{"endpoint_uri", "connection_string"},
		},

		"connection_string": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				accountKeyRegex := regexp.MustCompile("AccountKey=[^;]+")

				maskedNew := accountKeyRegex.ReplaceAllString(new, "AccountKey=****")
				return (new == d.Get(k).(string)) && (maskedNew == old)
			},
			Sensitive:     true,
			ConflictsWith: []string{"identity_id"},
			ExactlyOneOf:  []string{"endpoint_uri", "connection_string"},
		},

		"encoding": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			Default:          string(devices.EncodingAvro),
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc: validation.StringInSlice([]string{
				string(devices.EncodingAvro),
				string(devices.EncodingAvroDeflate),
				string(devices.EncodingJSON),
			}, true),
		},
	}
}

func resourceIotHubEndpointStorageContainerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	id := parse.NewEndpointStorageContainerID(subscriptionId, iotHubRG, iotHubName, d.Get("name").(string))

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
	containerName := d.Get("container_name").(string)
	fileNameFormat := d.Get("file_name_format").(string)
	batchFrequencyInSeconds := int32(d.Get("batch_frequency_in_seconds").(int))
	maxChunkSizeInBytes := int32(d.Get("max_chunk_size_in_bytes").(int))
	encoding := d.Get("encoding").(string)

	storageContainerEndpoint := devices.RoutingStorageContainerProperties{
		AuthenticationType:      authenticationType,
		Name:                    &id.EndpointName,
		SubscriptionID:          &subscriptionID,
		ResourceGroup:           &endpointRG,
		ContainerName:           &containerName,
		FileNameFormat:          &fileNameFormat,
		BatchFrequencyInSeconds: &batchFrequencyInSeconds,
		MaxChunkSizeInBytes:     &maxChunkSizeInBytes,
		Encoding:                devices.Encoding(encoding),
	}

	if authenticationType == devices.AuthenticationTypeKeyBased {
		if v, ok := d.GetOk("connection_string"); ok {
			storageContainerEndpoint.ConnectionString = utils.String(v.(string))
		} else {
			return fmt.Errorf("`connection_string` must be specified when `authentication_type` is `keyBased`")
		}
	} else {
		if v, ok := d.GetOk("endpoint_uri"); ok {
			storageContainerEndpoint.EndpointURI = utils.String(v.(string))
		} else {
			return fmt.Errorf("`endpoint_uri` must be specified when `authentication_type` is `identityBased`")
		}

		if v, ok := d.GetOk("identity_id"); ok {
			storageContainerEndpoint.Identity = &devices.ManagedIdentity{
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

	if routing.Endpoints.StorageContainers == nil {
		storageContainers := make([]devices.RoutingStorageContainerProperties, 0)
		routing.Endpoints.StorageContainers = &storageContainers
	}

	endpoints := make([]devices.RoutingStorageContainerProperties, 0)

	alreadyExists := false
	for _, existingEndpoint := range *routing.Endpoints.StorageContainers {
		if existingEndpointName := existingEndpoint.Name; existingEndpointName != nil {
			if strings.EqualFold(*existingEndpointName, id.EndpointName) {
				if d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_iothub_endpoint_storage_container", id.ID())
				}
				endpoints = append(endpoints, storageContainerEndpoint)
				alreadyExists = true
			} else {
				endpoints = append(endpoints, existingEndpoint)
			}
		}
	}

	if d.IsNewResource() {
		endpoints = append(endpoints, storageContainerEndpoint)
	} else if !alreadyExists {
		return fmt.Errorf("unable to find %s", id)
	}
	routing.Endpoints.StorageContainers = &endpoints

	future, err := client.CreateOrUpdate(ctx, iotHubRG, iotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubEndpointStorageContainerRead(d, meta)
}

func resourceIotHubEndpointStorageContainerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointStorageContainerID(d.Id())
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

	if endpoints := iothub.Properties.Routing.Endpoints.StorageContainers; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, id.EndpointName) {
					exist = true
					d.Set("container_name", endpoint.ContainerName)
					d.Set("file_name_format", endpoint.FileNameFormat)
					d.Set("batch_frequency_in_seconds", endpoint.BatchFrequencyInSeconds)
					d.Set("max_chunk_size_in_bytes", endpoint.MaxChunkSizeInBytes)
					d.Set("encoding", endpoint.Encoding)
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

func resourceIotHubEndpointStorageContainerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointStorageContainerID(d.Id())
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
	endpoints := iothub.Properties.Routing.Endpoints.StorageContainers

	if endpoints == nil {
		return nil
	}

	updatedEndpoints := make([]devices.RoutingStorageContainerProperties, 0)
	for _, endpoint := range *endpoints {
		if existingEndpointName := endpoint.Name; existingEndpointName != nil {
			if !strings.EqualFold(*existingEndpointName, id.EndpointName) {
				updatedEndpoints = append(updatedEndpoints, endpoint)
			}
		}
	}
	iothub.Properties.Routing.Endpoints.StorageContainers = &updatedEndpoints

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}
