package iothub

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/iothub/mgmt/2020-03-01/devices"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotHubEnrichment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIotHubEnrichmentCreateUpdate,
		Read:   resourceArmIotHubEnrichmentRead,
		Update: resourceArmIotHubEnrichmentCreateUpdate,
		Delete: resourceArmIotHubEnrichmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-_.a-zA-Z0-9]{1,64}$"),
					"Enrichment Key name can only include alphanumeric characters, periods, underscores, hyphens, has a maximum length of 64 characters, and must be unique.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"iothub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IoTHubName,
			},

			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"endpoint_names": {
				Type:     schema.TypeList,
				MaxItems: 100,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Required: true,
			},
		},
	}
}

func resourceArmIotHubEnrichmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

	enrichmentKey := d.Get("key").(string)
	enrichmentValue := d.Get("value").(string)
	resourceId := fmt.Sprintf("%s/Enrichments/%s", *iothub.ID, enrichmentKey)
	endpointNamesRaw := d.Get("endpoint_names").([]interface{})

	enrichment := devices.EnrichmentProperties{
		Key:           &enrichmentKey,
		Value:         &enrichmentValue,
		EndpointNames: utils.ExpandStringSlice(endpointNamesRaw),
	}

	routing := iothub.Properties.Routing
	if routing == nil {
		routing = &devices.RoutingProperties{}
	}

	if routing.Enrichments == nil {
		enrichments := make([]devices.EnrichmentProperties, 0)
		routing.Enrichments = &enrichments
	}

	enrichments := make([]devices.EnrichmentProperties, 0)

	alreadyExists := false
	for _, existingEnrichment := range *routing.Enrichments {
		if existingEnrichment.Key != nil {
			if strings.EqualFold(*existingEnrichment.Key, enrichmentKey) {
				if d.IsNewResource() {
					return tf.ImportAsExistsError("azurerm_iothub_enrichment", resourceId)
				}
				enrichments = append(enrichments, enrichment)
				alreadyExists = true
			} else {
				enrichments = append(enrichments, existingEnrichment)
			}
		}
	}

	if d.IsNewResource() {
		enrichments = append(enrichments, enrichment)
	} else if !alreadyExists {
		return fmt.Errorf("Unable to find Enrichment %q defined for IotHub %q (Resource Group %q)", enrichmentKey, iothubName, resourceGroup)
	}
	routing.Enrichments = &enrichments

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the completion of the creating/updating of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.SetId(resourceId)

	return resourceArmIotHubEnrichmentRead(d, meta)
}

func resourceArmIotHubEnrichmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedEnrichmentId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedEnrichmentId.ResourceGroup
	iothubName := parsedEnrichmentId.Path["IotHubs"]
	enrichmentKey := parsedEnrichmentId.Path["Enrichments"]

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	d.Set("key", enrichmentKey)
	d.Set("iothub_name", iothubName)
	d.Set("resource_group_name", resourceGroup)

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		return nil
	}

	if enrichments := iothub.Properties.Routing.Enrichments; enrichments != nil {
		for _, enrichment := range *enrichments {
			if enrichment.Key != nil {
				if strings.EqualFold(*enrichment.Key, enrichmentKey) {
					d.Set("value", enrichment.Value)
					d.Set("endpoint_names", enrichment.EndpointNames)
				}
			}
		}
	}

	return nil
}

func resourceArmIotHubEnrichmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTHub.ResourceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parsedEnrichmentId, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := parsedEnrichmentId.ResourceGroup
	iothubName := parsedEnrichmentId.Path["IotHubs"]
	enrichmentKey := parsedEnrichmentId.Path["Enrichments"]

	locks.ByName(iothubName, IothubResourceName)
	defer locks.UnlockByName(iothubName, IothubResourceName)

	iothub, err := client.Get(ctx, resourceGroup, iothubName)
	if err != nil {
		if utils.ResponseWasNotFound(iothub.Response) {
			return fmt.Errorf("IotHub %q (Resource Group %q) was not found", iothubName, resourceGroup)
		}
		return fmt.Errorf("Error loading IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
	}

	if iothub.Properties == nil || iothub.Properties.Routing == nil {
		return nil
	}
	enrichments := iothub.Properties.Routing.Enrichments

	if enrichments == nil {
		return nil
	}

	updatedEnrichments := make([]devices.EnrichmentProperties, 0)
	for _, enrichment := range *enrichments {
		if enrichment.Key != nil {
			if !strings.EqualFold(*enrichment.Key, enrichmentKey) {
				updatedEnrichments = append(updatedEnrichments, enrichment)
			}
		}
	}
	iothub.Properties.Routing.Enrichments = &updatedEnrichments

	future, err := client.CreateOrUpdate(ctx, resourceGroup, iothubName, iothub, "")
	if err != nil {
		return fmt.Errorf("Error updating IotHub %q (Resource Group %q) with Enrichment %q: %+v", iothubName, resourceGroup, enrichmentKey, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for IotHub %q (Resource Group %q) to finish updating Enrichment %q: %+v", iothubName, resourceGroup, enrichmentKey, err)
	}

	return nil
}
