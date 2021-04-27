package iothub

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotHubEndpointServiceBusTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceIotHubEndpointServiceBusTopicCreateUpdate,
		Read:   resourceIotHubEndpointServiceBusTopicRead,
		Update: resourceIotHubEndpointServiceBusTopicCreateUpdate,
		Delete: resourceIotHubEndpointServiceBusTopicDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubEndpointName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"connection_string": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					sharedAccessKeyRegex := regexp.MustCompile("SharedAccessKey=[^;]+")
					sbProtocolRegex := regexp.MustCompile("sb://([^:]+)(:5671)?/;")

					maskedNew := sbProtocolRegex.ReplaceAllString(new, "sb://$1:5671/;")
					maskedNew = sharedAccessKeyRegex.ReplaceAllString(maskedNew, "SharedAccessKey=****")
					return (new == d.Get(k).(string)) && (maskedNew == old)
				},
				Sensitive: true,
			},
		},
	}
}

func resourceIotHubEndpointServiceBusTopicCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

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

	endpointName := d.Get("name").(string)
	resourceId := fmt.Sprintf("%s/Endpoints/%s", *iothub.ID, endpointName)

	topicEndpoint := devices.RoutingServiceBusTopicEndpointProperties{
		ConnectionString: utils.String(d.Get("connection_string").(string)),
		Name:             utils.String(endpointName),
		SubscriptionID:   utils.String(subscriptionID),
		ResourceGroup:    utils.String(resourceGroup),
	}

	routing := iothub.Properties.Routing
	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	if routing.Endpoints == nil {
		routing.Endpoints = &devices.RoutingEndpoints{}
	}

	if routing.Endpoints.EventHubs == nil {
		topics := make([]devices.RoutingServiceBusTopicEndpointProperties, 0)
		routing.Endpoints.ServiceBusTopics = &topics
	}
	endpoints := make([]devices.RoutingServiceBusTopicEndpointProperties, 0)

	alreadyExists := false
	for _, existingEndpoint := range *routing.Endpoints.ServiceBusTopics {
		if existingEndpointName := existingEndpoint.Name; existingEndpointName != nil {
			if strings.EqualFold(*existingEndpointName, endpointName) {
				if d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_iothub_endpoint_servicebus_topic", resourceId)
				}
				endpoints = append(endpoints, topicEndpoint)
				alreadyExists = true
			} else {
				endpoints = append(endpoints, existingEndpoint)
			}
		}
	}

	if d.IsNewResource() {
		endpoints = append(endpoints, topicEndpoint)
	} else if !alreadyExists {
		return fmt.Errorf("Unable to find ServiceBus Queue Endpoint %q defined for IotHub %q (Resource Group %q)", endpointName, iothubName, resourceGroup)
	}
	routing.Endpoints.ServiceBusTopics = &endpoints

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceIotHubEndpointServiceBusTopicRead(d, meta)
}

func resourceIotHubEndpointServiceBusTopicRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedIothubEndpointId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedIothubEndpointId.ResourceGroup
	iothubName := parsedIothubEndpointId.Path["IotHubs"]
	endpointName := parsedIothubEndpointId.Path["Endpoints"]

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.Set("name", endpointName)
	d.Set("iothub_name", iothubName)
	d.Set("resource_group_name", resourceGroup)

	if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		return nil
	}

	if endpoints := iothub.Properties.Routing.Endpoints.ServiceBusTopics; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, endpointName) {
					d.Set("connection_string", endpoint.ConnectionString)
				}
			}
		}
	}

	return nil
}

func resourceIotHubEndpointServiceBusTopicDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedIothubEndpointId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedIothubEndpointId.ResourceGroup
	iothubName := parsedIothubEndpointId.Path["IotHubs"]
	endpointName := parsedIothubEndpointId.Path["Endpoints"]

	locks.ByName(iothubName, IothubResourceName)
	defer locks.UnlockByName(iothubName, IothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if iothub.Properties == nil || iothub.Properties.Routing == nil || iothub.Properties.Routing.Endpoints == nil {
		return nil
	}
	endpoints := iothub.Properties.Routing.Endpoints.ServiceBusTopics

	if endpoints == nil {
		return nil
	}

	updatedEndpoints := make([]devices.RoutingServiceBusTopicEndpointProperties, 0)
	for _, endpoint := range *endpoints {
		if existingEndpointName := endpoint.Name; existingEndpointName != nil {
			if !strings.EqualFold(*existingEndpointName, endpointName) {
				updatedEndpoints = append(updatedEndpoints, endpoint)
			}
		}
	}
	iothub.Properties.Routing.Endpoints.ServiceBusTopics = &updatedEndpoints

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error updating IotHub %q (Resource Group %q) with ServiceBus Queue Endpoint %q: %+v", iothubName, resourceGroup, endpointName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub %q (Resource Group %q) to finish updating ServiceBus Queue Endpoint %q: %+v", iothubName, resourceGroup, endpointName, err)
	}

	return nil
}
