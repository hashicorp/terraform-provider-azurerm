package providerjson

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFromRaw(input *schema.Resource) (*ResourceJSON, error) {
	if input == nil {
		return nil, fmt.Errorf("resource not found")
	}

	result := &ResourceJSON{}
	translatedSchema := make(map[string]SchemaJSON)

	for k, s := range input.Schema {
		translatedSchema[k] = schemaFromRaw(s)
	}
	result.Schema = translatedSchema

	if input.Timeouts != nil {
		timeouts := &ResourceTimeoutJSON{
			Create: int(input.Timeouts.Create.Minutes()),
			Read:   int(input.Timeouts.Read.Minutes()),
			Delete: int(input.Timeouts.Delete.Minutes()),
			Update: int(input.Timeouts.Update.Minutes()),
		}
		result.Timeouts = timeouts
	}

	return result, nil
}

func schemaFromRaw(input *schema.Schema) SchemaJSON {

	return SchemaJSON{
		Type:        input.Type.String(),
		ConfigMode:  decodeConfigMode(input.ConfigMode),
		Optional:    input.Optional,
		Required:    input.Required,
		Default:     input.Default,
		Description: input.Description,
		Computed:    input.Computed,
		ForceNew:    input.ForceNew,
		Elem:        decodeElem(input.Elem),
		MaxItems:    input.MaxItems,
		MinItems:    input.MinItems,
	}
}

func decodeConfigMode(input schema.SchemaConfigMode) (out string) {
	switch input {
	case 1:
		out = "Auto"
	case 2:
		out = "Block"
	case 4:
		out = "Attribute"
	}
	return
}

func decodeElem(input interface{}) interface{} {
	switch t := input.(type) {
	case bool:
		return t
	case string:
		return t
	case int:
		return t
	case float32, float64:
		return t
	case *schema.Schema:
		return schemaFromRaw(t)
	case *schema.Resource:
		r, _ := resourceFromRaw(t)
		return r
	}
	return nil
}
