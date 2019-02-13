package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventgrid/mgmt/2018-09-15-preview/eventgrid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmEventGridDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmEventGridDomainCreateUpdate,
		Read:   resourceArmEventGridDomainRead,
		Update: resourceArmEventGridDomainCreateUpdate,
		Delete: resourceArmEventGridDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),

			"input_schema": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(eventgrid.InputSchemaEventGridSchema),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(eventgrid.InputSchemaCloudEventV01Schema),
					string(eventgrid.InputSchemaCustomEventSchema),
					string(eventgrid.InputSchemaEventGridSchema),
				}, true),
			},

			"input_mapping_fields": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"topic": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"event_time": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"event_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"subject": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"data_version": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},

			"input_mapping_default_values": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"subject": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"data_version": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmEventGridDomainCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGridDomainsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventGrid Domain %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventgrid_domain", *existing.ID)
		}
	}

	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	domainProperties := &eventgrid.DomainProperties{
		InputSchemaMapping: expandAzureRmEventgridDomainInputMapping(d),
		InputSchema:        eventgrid.InputSchema(d.Get("input_schema").(string)),
	}

	domain := eventgrid.Domain{
		Location:         &location,
		DomainProperties: domainProperties,
		Tags:             expandTags(tags),
	}

	log.Printf("[INFO] preparing arguments for AzureRM EventGrid Domain creation with Properties: %+v", domain)

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, domain)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read EventGrid Domain %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmEventGridDomainRead(d, meta)
}

func resourceArmEventGridDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGridDomainsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["domains"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] EventGrid Domain '%s' was not found (resource group '%s')", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on EventGrid Domain '%s': %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.DomainProperties; props != nil {
		d.Set("endpoint", props.Endpoint)

		d.Set("input_schema", string(props.InputSchema))

		if err := d.Set("input_mapping_fields", flattenAzureRmEventgridDomainInputMapping(props.InputSchemaMapping)); err != nil {
			return fmt.Errorf("Error setting `input_schema_mapping_fields`: %+v", err)
		}

		if err := d.Set("input_mapping_default_values", flattenAzureRmEventgridDomainInputMappingDefaultValues(props.InputSchemaMapping)); err != nil {
			return fmt.Errorf("Error setting `input_schema_mapping_fields`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmEventGridDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).eventGridDomainsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["domains"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Domain %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Event Grid Domain %q: %+v", name, err)
	}

	return nil
}

func expandAzureRmEventgridDomainInputMapping(d *schema.ResourceData) *eventgrid.JSONInputSchemaMapping {
	imf, imfok := d.GetOk("input_mapping_fields")

	imdv, imdvok := d.GetOk("input_mapping_default_values")

	if !imfok && !imdvok {
		return nil
	}

	jismp := eventgrid.JSONInputSchemaMappingProperties{}

	if imfok {
		mappings := imf.([]interface{})
		mapping := mappings[0].(map[string]interface{})

		if id := mapping["id"].(string); id != "" {
			jismp.ID = &eventgrid.JSONField{SourceField: &id}
		}

		if eventTime := mapping["event_time"].(string); eventTime != "" {
			jismp.EventTime = &eventgrid.JSONField{SourceField: &eventTime}
		}

		if topic := mapping["topic"].(string); topic != "" {
			jismp.Topic = &eventgrid.JSONField{SourceField: &topic}
		}

		if dataVersion := mapping["data_version"].(string); dataVersion != "" {
			jismp.DataVersion = &eventgrid.JSONFieldWithDefault{SourceField: &dataVersion}
		}

		if subject := mapping["subject"].(string); subject != "" {
			jismp.Subject = &eventgrid.JSONFieldWithDefault{SourceField: &subject}
		}

		if eventType := mapping["event_type"].(string); eventType != "" {
			jismp.EventType = &eventgrid.JSONFieldWithDefault{SourceField: &eventType}
		}
	}

	if imdvok {
		mappings := imdv.([]interface{})
		mapping := mappings[0].(map[string]interface{})

		if dataVersion := mapping["data_version"].(string); dataVersion != "" {
			jismp.DataVersion = &eventgrid.JSONFieldWithDefault{DefaultValue: &dataVersion}
		}

		if subject := mapping["subject"].(string); subject != "" {
			jismp.Subject = &eventgrid.JSONFieldWithDefault{DefaultValue: &subject}
		}

		if eventType := mapping["event_type"].(string); eventType != "" {
			jismp.EventType = &eventgrid.JSONFieldWithDefault{DefaultValue: &eventType}
		}
	}

	jsonMapping := eventgrid.JSONInputSchemaMapping{
		JSONInputSchemaMappingProperties: &jismp,
		InputSchemaMappingType:           eventgrid.InputSchemaMappingTypeJSON,
	}

	return &jsonMapping
}

func flattenAzureRmEventgridDomainInputMapping(input eventgrid.BasicInputSchemaMapping) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	jsonValues := input.(eventgrid.JSONInputSchemaMapping).JSONInputSchemaMappingProperties

	if jsonValues.EventTime != nil && jsonValues.EventTime.SourceField != nil {
		result["event_time"] = *jsonValues.EventTime.SourceField
	}

	if jsonValues.ID != nil && jsonValues.ID.SourceField != nil {
		result["id"] = *jsonValues.ID.SourceField
	}

	if jsonValues.Topic != nil && jsonValues.Topic.SourceField != nil {
		result["topic"] = *jsonValues.Topic.SourceField
	}

	if jsonValues.DataVersion != nil && jsonValues.DataVersion.SourceField != nil {
		result["data_version"] = *jsonValues.DataVersion.SourceField
	}

	if jsonValues.EventType != nil && jsonValues.EventType.SourceField != nil {
		result["event_type"] = *jsonValues.EventType.SourceField
	}

	if jsonValues.Subject != nil && jsonValues.Subject.SourceField != nil {
		result["subject"] = *jsonValues.Subject.SourceField
	}

	return []interface{}{result}
}

func flattenAzureRmEventgridDomainInputMappingDefaultValues(input eventgrid.BasicInputSchemaMapping) []interface{} {
	if input == nil {
		return nil
	}
	result := make(map[string]interface{})

	jsonValues := input.(eventgrid.JSONInputSchemaMapping).JSONInputSchemaMappingProperties

	if jsonValues.DataVersion != nil && jsonValues.DataVersion.DefaultValue != nil {
		result["data_version"] = *jsonValues.DataVersion.DefaultValue
	}

	if jsonValues.EventType != nil && jsonValues.EventType.DefaultValue != nil {
		result["event_type"] = *jsonValues.EventType.DefaultValue
	}

	if jsonValues.Subject != nil && jsonValues.Subject.DefaultValue != nil {
		result["subject"] = *jsonValues.Subject.DefaultValue
	}

	return []interface{}{result}
}
