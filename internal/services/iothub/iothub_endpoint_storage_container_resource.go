package iothub

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2021-03-31/devices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/parse"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotHubEndpointStorageContainer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubEndpointStorageContainerCreateUpdate,
		Read:   resourceIotHubEndpointStorageContainerRead,
		Update: resourceIotHubEndpointStorageContainerCreateUpdate,
		Delete: resourceIotHubEndpointStorageContainerDelete,

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

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IoTHubEndpointName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IoTHubName,
			},

			"container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.StorageContainerName,
			},

			"file_name_format": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  false,
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

			"connection_string": {
				Type:     pluginsdk.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
					accountKeyRegex := regexp.MustCompile("AccountKey=[^;]+")

					maskedNew := accountKeyRegex.ReplaceAllString(new, "AccountKey=****")
					return (new == d.Get(k).(string)) && (maskedNew == old)
				},
				Sensitive: true,
			},

			"encoding": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(devices.EncodingAvro),
					string(devices.EncodingAvroDeflate),
					string(devices.EncodingJSON),
				}, true),
			},
		},
	}
}

func resourceIotHubEndpointStorageContainerCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.ResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	id := parse.NewEndpointStorageContainerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	connectionStr := d.Get("connection_string").(string)
	containerName := d.Get("container_name").(string)
	fileNameFormat := d.Get("file_name_format").(string)
	batchFrequencyInSeconds := int32(d.Get("batch_frequency_in_seconds").(int))
	maxChunkSizeInBytes := int32(d.Get("max_chunk_size_in_bytes").(int))
	encoding := d.Get("encoding").(string)

	storageContainerEndpoint := devices.RoutingStorageContainerProperties{
		ConnectionString:        &connectionStr,
		Name:                    &id.EndpointName,
		SubscriptionID:          &subscriptionID,
		ResourceGroup:           &id.ResourceGroup,
		ContainerName:           &containerName,
		FileNameFormat:          &fileNameFormat,
		BatchFrequencyInSeconds: &batchFrequencyInSeconds,
		MaxChunkSizeInBytes:     &maxChunkSizeInBytes,
		Encoding:                devices.Encoding(encoding),
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
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
		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	d.Set("name", id.EndpointName)
	d.Set("iothub_name", id.IotHubName)
	d.Set("resource_group_name", id.ResourceGroup)

	if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		return nil
	}

	if endpoints := iothub.Properties.Routing.Endpoints.StorageContainers; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, id.EndpointName) {
					d.Set("connection_string", endpoint.ConnectionString)
					d.Set("container_name", endpoint.ContainerName)
					d.Set("file_name_format", endpoint.FileNameFormat)
					d.Set("batch_frequency_in_seconds", endpoint.BatchFrequencyInSeconds)
					d.Set("max_chunk_size_in_bytes", endpoint.MaxChunkSizeInBytes)
					d.Set("encoding", endpoint.Encoding)
				}
			}
		}
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
