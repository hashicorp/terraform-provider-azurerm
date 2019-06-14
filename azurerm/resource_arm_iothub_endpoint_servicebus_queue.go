package azurerm

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/iothub/mgmt/2018-12-01-preview/devices"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIotHubEndpointServiceBusQueue() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubEndpointServiceBusQueueCreateUpdate,
		Read:   resourceArmIotHubEndpointServiceBusQueueRead,
		Update: resourceArmIotHubEndpointServiceBusQueueCreateUpdate,
		Delete: resourceArmIotHubEndpointServiceBusQueueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceArmIotHubEndpointServiceBusQueueCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext
	subscriptionID := meta.(*ArmClient).subscriptionId

	iothubName := d.Get("iothub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	azureRMLockByName(iothubName, iothubResourceName)
	defer azureRMUnlockByName(iothubName, iothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}

		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	endpointName := d.Get("name").(string)

	resourceId := fmt.Sprintf("%s/Endpoints/%s", *iothub.ID, endpointName)

	connectionStr := d.Get("connection_string").(string)

	queueEndpoint := devices.RoutingServiceBusQueueEndpointProperties{
		ConnectionString: &connectionStr,
		Name:             &endpointName,
		SubscriptionID:   &subscriptionID,
		ResourceGroup:    &resourceGroup,
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
		if strings.EqualFold(*existingEndpoint.Name, endpointName) {
			if d.IsNewResource() && requireResourcesToBeImported {
				return tf.ImportAsExistsError("azurerm_iothub_endpoint_servicebus_queue", resourceId)
			}
			endpoints = append(endpoints, queueEndpoint)
			alreadyExists = true

		} else {
			endpoints = append(endpoints, existingEndpoint)
		}
	}

	if d.IsNewResource() {
		endpoints = append(endpoints, queueEndpoint)
	} else if !alreadyExists {
		return fmt.Errorf("Unable to find ServiceBus Queue Endpoint %q defined for IotHub %q (Resource Group %q)", endpointName, iothubName, resourceGroup)
	}

	routing.Endpoints.ServiceBusQueues = &endpoints

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceArmIotHubEndpointServiceBusQueueRead(d, meta)
}

func resourceArmIotHubEndpointServiceBusQueueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext

	parsedIothubEndpointId, err := parseAzureResourceID(d.Id())

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

	if endpoints := iothub.Properties.Routing.Endpoints.ServiceBusQueues; endpoints != nil {
		for _, endpoint := range *endpoints {
			if strings.EqualFold(*endpoint.Name, endpointName) {
				d.Set("connection_string", endpoint.ConnectionString)
			}
		}
	}

	return nil
}

func resourceArmIotHubEndpointServiceBusQueueDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).iothub.ResourceClient
	ctx := meta.(*ArmClient).StopContext

	parsedIothubEndpointId, err := parseAzureResourceID(d.Id())

	if err != nil {
		return err
	}

	resourceGroup := parsedIothubEndpointId.ResourceGroup
	iothubName := parsedIothubEndpointId.Path["IotHubs"]
	endpointName := parsedIothubEndpointId.Path["Endpoints"]

	azureRMLockByName(iothubName, iothubResourceName)
	defer azureRMUnlockByName(iothubName, iothubResourceName)

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
	endpoints := iothub.Properties.Routing.Endpoints.ServiceBusQueues

	if endpoints == nil {
		return nil
	}

	updatedEndpoints := make([]devices.RoutingServiceBusQueueEndpointProperties, 0)
	for _, endpoint := range *endpoints {
		if !strings.EqualFold(*endpoint.Name, endpointName) {
			updatedEndpoints = append(updatedEndpoints, endpoint)
		}
	}

	iothub.Properties.Routing.Endpoints.ServiceBusQueues = &updatedEndpoints

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error updating IotHub %q (Resource Group %q) with ServiceBus Queue Endpoint %q: %+v", iothubName, resourceGroup, endpointName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub %q (Resource Group %q) to finish updating ServiceBus Queue Endpoint %q: %+v", iothubName, resourceGroup, endpointName, err)
	}

	return nil
}
