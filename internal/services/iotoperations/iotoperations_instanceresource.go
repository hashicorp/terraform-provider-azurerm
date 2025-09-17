package iotoperations

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotoperations/armiotoperations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clients "github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func resourceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceInstanceCreate,
		Read:   resourceInstanceRead,
		Update: resourceInstanceUpdate,
		Delete: resourceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"extended_location_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extended_location_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		       "adr_namespace_ref": {
			       Type:     schema.TypeString,
			       Optional: true,
			       Description: "The Azure Device Registry Namespace used by Assets, Discovered Assets and devices",
		       },
		       "default_secret_provider_class_ref": {
			       Type:     schema.TypeString,
			       Optional: true,
			       Description: "The reference to the AIO Secret provider class.",
		       },
		       "description": {
			       Type:     schema.TypeString,
			       Optional: true,
			       Description: "Detailed description of the Instance.",
		       },
		       "features": {
			       Type:     schema.TypeMap,
			       Optional: true,
			       Elem:     &schema.Schema{Type: schema.TypeString},
			       Description: "The features of the AIO Instance.",
		       },
		       "provisioning_state": {
			       Type:     schema.TypeString,
			       Computed: true,
			       Description: "The status of the last operation.",
		       },
		       "schema_registry_ref": {
			       Type:     schema.TypeString,
			       Optional: true,
			       Description: "The reference to the Schema Registry for this AIO Instance.",
		       },
		       "version": {
			       Type:     schema.TypeString,
			       Optional: true,
			       Description: "The Azure IoT Operations version.",
		       },
		       "tags": {
			       Type:     schema.TypeMap,
			       Optional: true,
			       Elem:     &schema.Schema{Type: schema.TypeString},
		       },
		},
	}
}

func resourceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	top := meta.(*clients.Client)
	svc := top.IoTOperations
	ctx := context.Background()

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)

	// Build extended location if provided
	var extendedLocation *armiotoperations.ExtendedLocation
	if extName, ok := d.GetOk("extended_location_name"); ok {
		extType := d.Get("extended_location_type").(string)
		extendedLocation = &armiotoperations.ExtendedLocation{
			Name: toPtr(extName.(string)),
			Type: toPtr(armiotoperations.ExtendedLocationType(extType)),
		}
	}

	// Build properties
	props := &armiotoperations.InstanceProperties{
		Description: toPtr(d.Get("description").(string)),
		Version:     toPtr(d.Get("version").(string)),
	}

	// Build tags
	tags := make(map[string]*string)
	if v, ok := d.GetOk("tags"); ok {
		for k, v := range v.(map[string]interface{}) {
			tags[k] = toPtr(v.(string))
		}
	}

	instance := armiotoperations.InstanceResource{
		Location:         toPtr(location),
		ExtendedLocation: extendedLocation,
		Properties:       props,
		Tags:             tags,
	}

	poller, err := svc.InstanceClient.BeginCreateOrUpdate(ctx, rg, name, instance, nil)
	if err != nil {
		return fmt.Errorf("creating IoT Operations Instance: %+v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("waiting for IoT Operations Instance creation: %+v", err)
	}

	d.SetId(fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.IoTOperations/instances/%s", top.Account.SubscriptionId, rg, name))
	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	top := meta.(*clients.Client)
	svc := top.IoTOperations
	ctx := context.Background()

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := svc.InstanceClient.Get(ctx, rg, name, nil)
	if err != nil {
		// Handle not found, etc.
		d.SetId("")
		return nil
	}

	instance := resp.InstanceResource

	d.Set("name", name)
	d.Set("resource_group_name", rg)
	d.Set("location", instance.Location)
	if instance.ExtendedLocation != nil {
		d.Set("extended_location_name", instance.ExtendedLocation.Name)
		d.Set("extended_location_type", instance.ExtendedLocation.Type)
	}
	if instance.Properties != nil {
		d.Set("description", instance.Properties.Description)
		d.Set("version", instance.Properties.Version)
		d.Set("provisioning_state", instance.Properties.ProvisioningState)
	}
	if instance.Tags != nil {
		tags := make(map[string]string)
		for k, v := range instance.Tags {
			if v != nil {
				tags[k] = *v
			}
		}
		d.Set("tags", tags)
	}

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	top := meta.(*clients.Client)
	svc := top.IoTOperations
	ctx := context.Background()

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	// Build tags
	tags := make(map[string]*string)
	if v, ok := d.GetOk("tags"); ok {
		for k, v := range v.(map[string]interface{}) {
			tags[k] = toPtr(v.(string))
		}
	}

	patch := armiotoperations.InstancePatchModel{
		Tags: tags,
	}

	_, err := svc.InstanceClient.Update(ctx, rg, name, patch, nil)
	if err != nil {
		return fmt.Errorf("updating IoT Operations Instance: %+v", err)
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	top := meta.(*clients.Client)
	svc := top.IoTOperations
	ctx := context.Background()

	rg := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	poller, err := svc.InstanceClient.BeginDelete(ctx, rg, name, nil)
	if err != nil {
		return fmt.Errorf("deleting IoT Operations Instance: %+v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		return fmt.Errorf("waiting for IoT Operations Instance deletion: %+v", err)
	}

	d.SetId("")
	return nil
}

// Helper to get pointer to string
func toPtr[T any](v T) *T {
	return &v
}
