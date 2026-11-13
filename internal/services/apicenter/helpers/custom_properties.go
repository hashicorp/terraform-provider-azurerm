// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CustomPropertiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:             pluginsdk.TypeString,
		Description:      "A JSON object containing custom properties.",
		Optional:         true,
		Default:          "{}",
		ValidateFunc:     validation.StringIsJSON,
		DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
	}
}

func ExpandCustomProperties(input string) (*any, error) {
	properties, err := pluginsdk.ExpandJsonFromString(input)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling `custom_properties`: %+v", err)
	}
	if properties == nil {
		return nil, fmt.Errorf("`custom_properties` must be a JSON object")
	}

	return pointer.To[any](properties), nil
}

func FlattenCustomProperties(input *any) (string, error) {
	if input == nil || *input == nil {
		return "{}", nil
	}

	properties, err := json.Marshal(*input)
	if err != nil {
		return "", fmt.Errorf("marshalling `custom_properties`: %+v", err)
	}

	return string(properties), nil
}
