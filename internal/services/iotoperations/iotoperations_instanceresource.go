package iotoperations

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotoperations/2024-11-01/instance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of instance.",
				ValidateFunc: validation.All(
					validation.StringLenBetween(3, 63),
					validation.StringMatch(regexp.MustCompile("^[a-z0-9][a-z0-9-]*[a-z0-9]$"), "Must match ^[a-z0-9][a-z0-9-]*[a-z0-9]$"),
				),
			},
			"resource_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The name of the resource group. The name is case insensitive.",
				ValidateFunc: validation.StringLenBetween(1, 90),
			}
				Description:  "The API version to use for this operation.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Detailed description of the Instance.",
			},
			"provisioning_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the last operation.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Azure IoT Operations version.",
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
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Azure Device Registry Namespace used by Assets, Discovered Assets and devices",
			},
			"default_secret_provider_class_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The reference to the AIO Secret provider class.",
			},
			"features": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The features of the AIO Instance.",
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
	client := meta.(*clients.Client)
	svc := client.IoTOperations
	ctx := context.Background()

	rg := d.Get("resourceGroupName").(string)
	name := d.Get("instanceName").(string)
	location := d.Get("location").(string)

	id := instance.NewInstanceID(client.Account.SubscriptionId, rg, name)

	// Build extended location if provided
	var extendedLocation *instance.ExtendedLocation
	if extName, ok := d.GetOk("extended_location_name"); ok {
		extType := d.Get("extended_location_type").(string)
		extendedLocation = &instance.ExtendedLocation{
			Name: toPtr(extName.(string)),
			Type: toPtr(instance.ExtendedLocationType(extType)),
		}
	}

	// Build properties
	props := &instance.InstanceProperties{
		Description: toPtr(d.Get("description").(string)),
		Version:     toPtr(d.Get("version").(string)),
	}

	// Build tags
	tags := make(map[string]string)
	if v, ok := d.GetOk("tags"); ok {
		for k, v := range v.(map[string]interface{}) {
			tags[k] = v.(string)
		}
	}

	payload := instance.InstanceResource{
		Location:         toPtr(location),
		ExtendedLocation: extendedLocation,
		Properties:       props,
		Tags:             &tags,
	}

	err := svc.InstanceClient.CreateOrUpdateThenPoll(ctx, id, payload)
	if err != nil {
		return fmt.Errorf("creating IoT Operations Instance: %+v", err)
	}

	d.SetId(id.ID())
	return resourceInstanceRead(d, meta)
}

func resourceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	svc := client.IoTOperations
	ctx := context.Background()

	id, err := instance.ParseInstanceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := svc.InstanceClient.Get(ctx, *id)
	if err != nil {
		// Handle not found, etc.
		d.SetId("")
		return nil
	}

	if model := resp.Model; model != nil {
		d.Set("instanceName", id.InstanceName)
		d.Set("resourceGroupName", id.ResourceGroupName)
		d.Set("location", model.Location)

		if model.ExtendedLocation != nil {
			d.Set("extended_location_name", model.ExtendedLocation.Name)
			d.Set("extended_location_type", model.ExtendedLocation.Type)
		}

		if model.Properties != nil {
			d.Set("description", model.Properties.Description)
			d.Set("version", model.Properties.Version)
			d.Set("provisioning_state", model.Properties.ProvisioningState)
		}

		if model.Tags != nil {
			d.Set("tags", *model.Tags)
		}
	}

	return nil
}

func resourceInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	svc := client.IoTOperations
	ctx := context.Background()

	id, err := instance.ParseInstanceID(d.Id())
	if err != nil {
		return err
	}

	// Build tags
	tags := make(map[string]string)
	if v, ok := d.GetOk("tags"); ok {
		for k, v := range v.(map[string]interface{}) {
			tags[k] = v.(string)
		}
	}

	payload := instance.InstancePatchModel{
		Tags: &tags,
	}

	_, err = svc.InstanceClient.Update(ctx, *id, payload)
	if err != nil {
		return fmt.Errorf("updating IoT Operations Instance: %+v", err)
	}

	return resourceInstanceRead(d, meta)
}

func resourceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client)
	svc := client.IoTOperations
	ctx := context.Background()

	id, err := instance.ParseInstanceID(d.Id())
	if err != nil {
		return err
	}

	err = svc.InstanceClient.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting IoT Operations Instance: %+v", err)
	}

	d.SetId("")
	return nil
}

// Helper to get pointer to string
func toPtr[T any](v T) *T {
	return &v
}
