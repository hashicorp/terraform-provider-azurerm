// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package authorization

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = ServicePrincipalAttributeSetResource{}
	_ sdk.ResourceWithUpdate = ServicePrincipalAttributeSetResource{}
)

type ServicePrincipalAttributeSetResource struct{}

type ServicePrincipalAttributeSetModel struct {
	ServicePrincipalObjectId string                 `tfschema:"service_principal_object_id"`
	AttributeSetName         string                 `tfschema:"attribute_set_name"`
	Attributes               map[string]interface{} `tfschema:"attributes"`
}

func (r ServicePrincipalAttributeSetResource) ModelObject() interface{} {
	return &ServicePrincipalAttributeSetModel{}
}

func (r ServicePrincipalAttributeSetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidateServicePrincipalAttributeSetID
}

func (r ServicePrincipalAttributeSetResource) ResourceType() string {
	return "azurerm_service_principal_attribute_set"
}

func (r ServicePrincipalAttributeSetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"service_principal_object_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
		"attribute_set_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"attributes": {
			Type:     pluginsdk.TypeMap,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r ServicePrincipalAttributeSetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ServicePrincipalAttributeSetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ServicePrincipalCustomSecurityAttributesClient

			var config ServicePrincipalAttributeSetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := parse.NewServicePrincipalAttributeSetId(config.ServicePrincipalObjectId, config.AttributeSetName)
			existing, err := client.Get(ctx, config.ServicePrincipalObjectId)
			if err != nil {
				return err
			}
			if existing == nil {
				return fmt.Errorf("retrieving service principal %q: not found", config.ServicePrincipalObjectId)
			}

			if _, exists := existing.CustomSecurityAttributes[config.AttributeSetName]; exists {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			customSecurityAttributes, err := upsertManagedAttributeSet(existing.CustomSecurityAttributes, config.AttributeSetName, config.Attributes)
			if err != nil {
				return err
			}

			if err := client.Update(ctx, config.ServicePrincipalObjectId, customSecurityAttributes); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ServicePrincipalAttributeSetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ServicePrincipalCustomSecurityAttributesClient

			id, err := parse.ServicePrincipalAttributeSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ServicePrincipalObjectId)
			if err != nil {
				return err
			}
			if existing == nil {
				return metadata.MarkAsGone(*id)
			}

			attributes, ok := existing.CustomSecurityAttributes[id.AttributeSetName]
			if !ok {
				return metadata.MarkAsGone(*id)
			}

			flattenedAttributes, err := flattenManagedAttributeSet(attributes)
			if err != nil {
				return err
			}

			state := ServicePrincipalAttributeSetModel{
				ServicePrincipalObjectId: id.ServicePrincipalObjectId,
				AttributeSetName:         id.AttributeSetName,
				Attributes:               flattenedAttributes,
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ServicePrincipalAttributeSetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ServicePrincipalCustomSecurityAttributesClient

			var config ServicePrincipalAttributeSetModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, config.ServicePrincipalObjectId)
			if err != nil {
				return err
			}
			if existing == nil {
				return fmt.Errorf("retrieving service principal %q: not found", config.ServicePrincipalObjectId)
			}

			customSecurityAttributes, err := upsertManagedAttributeSet(existing.CustomSecurityAttributes, config.AttributeSetName, config.Attributes)
			if err != nil {
				return err
			}

			return client.Update(ctx, config.ServicePrincipalObjectId, customSecurityAttributes)
		},
	}
}

func (r ServicePrincipalAttributeSetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Authorization.ServicePrincipalCustomSecurityAttributesClient

			id, err := parse.ServicePrincipalAttributeSetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ServicePrincipalObjectId)
			if err != nil {
				return err
			}
			if existing == nil {
				return nil
			}
			if _, exists := existing.CustomSecurityAttributes[id.AttributeSetName]; !exists {
				return nil
			}

			customSecurityAttributes := cloneCustomSecurityAttributes(existing.CustomSecurityAttributes)
			delete(customSecurityAttributes, id.AttributeSetName)

			return client.Update(ctx, id.ServicePrincipalObjectId, customSecurityAttributes)
		},
	}
}

func upsertManagedAttributeSet(existing map[string]map[string]interface{}, attributeSetName string, attributes map[string]interface{}) (map[string]map[string]interface{}, error) {
	customSecurityAttributes := cloneCustomSecurityAttributes(existing)

	updatedManagedSet := make(map[string]interface{}, len(attributes))
	for key, value := range attributes {
		stringValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected attribute %q in set %q to be a string but got %T", key, attributeSetName, value)
		}
		updatedManagedSet[key] = stringValue
	}

	customSecurityAttributes[attributeSetName] = updatedManagedSet
	return customSecurityAttributes, nil
}

func flattenManagedAttributeSet(attributes map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{}, len(attributes))
	for key, value := range attributes {
		if key == "@odata.type" {
			continue
		}

		stringValue, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected attribute %q to be a string but got %T", key, value)
		}
		out[key] = stringValue
	}

	return out, nil
}

func cloneCustomSecurityAttributes(input map[string]map[string]interface{}) map[string]map[string]interface{} {
	if input == nil {
		return map[string]map[string]interface{}{}
	}

	out := make(map[string]map[string]interface{}, len(input))
	for setName, values := range input {
		clonedValues := make(map[string]interface{}, len(values))
		for key, value := range values {
			clonedValues[key] = value
		}
		out[setName] = clonedValues
	}

	return out
}
