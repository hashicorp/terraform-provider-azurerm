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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotHubEndpointServiceBusTopic() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotHubEndpointServiceBusTopicCreateUpdate,
		Read:   resourceIotHubEndpointServiceBusTopicRead,
		Update: resourceIotHubEndpointServiceBusTopicCreateUpdate,
		Delete: resourceIotHubEndpointServiceBusTopicDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.EndpointServiceBusTopicID(id)
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
				ValidateFunc: validate.IoTHubEndpointName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"connection_string": {
				Type:     pluginsdk.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
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

func resourceIotHubEndpointServiceBusTopicCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	subscriptionId := meta.(*clients.Client).IoTHub.ResourceClient.SubscriptionID
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	subscriptionID := meta.(*clients.Client).Account.SubscriptionId

	id := parse.NewEndpointServiceBusTopicID(subscriptionId, d.Get("resource_group_name").(string), d.Get("iothub_name").(string), d.Get("name").(string))

	locks.ByName(id.IotHubName, IothubResourceName)
	defer locks.UnlockByName(id.IotHubName, IothubResourceName)

	iothub, err := client.Get(ctx, id.ResourceGroup, id.IotHubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", id.IotHubName, id.ResourceGroup)
		}

		return fmt.Errorf("loading IotHub %q (Resource Group %q): %+v", id.IotHubName, id.ResourceGroup, err)
	}

	topicEndpoint := devices.RoutingServiceBusTopicEndpointProperties{
		ConnectionString: utils.String(d.Get("connection_string").(string)),
		Name:             utils.String(id.EndpointName),
		SubscriptionID:   utils.String(subscriptionID),
		ResourceGroup:    utils.String(id.ResourceGroup),
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
			if strings.EqualFold(*existingEndpointName, id.EndpointName) {
				if d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_iothub_endpoint_servicebus_topic", id.ID())
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
		return fmt.Errorf("unable to find %s", id)
	}
	routing.Endpoints.ServiceBusTopics = &endpoints

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the completion of the creating/updating of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIotHubEndpointServiceBusTopicRead(d, meta)
}

func resourceIotHubEndpointServiceBusTopicRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointServiceBusTopicID(d.Id())
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

	if endpoints := iothub.Properties.Routing.Endpoints.ServiceBusTopics; endpoints != nil {
		for _, endpoint := range *endpoints {
			if existingEndpointName := endpoint.Name; existingEndpointName != nil {
				if strings.EqualFold(*existingEndpointName, id.EndpointName) {
					d.Set("connection_string", endpoint.ConnectionString)
				}
			}
		}
	}

	return nil
}

func resourceIotHubEndpointServiceBusTopicDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointServiceBusTopicID(d.Id())
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
	endpoints := iothub.Properties.Routing.Endpoints.ServiceBusTopics

	if endpoints == nil {
		return nil
	}

	updatedEndpoints := make([]devices.RoutingServiceBusTopicEndpointProperties, 0)
	for _, endpoint := range *endpoints {
		if existingEndpointName := endpoint.Name; existingEndpointName != nil {
			if !strings.EqualFold(*existingEndpointName, id.EndpointName) {
				updatedEndpoints = append(updatedEndpoints, endpoint)
			}
		}
	}
	iothub.Properties.Routing.Endpoints.ServiceBusTopics = &updatedEndpoints

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IotHubName, iothub, "")
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to finish updating: %+v", id, err)
	}

	return nil
}
