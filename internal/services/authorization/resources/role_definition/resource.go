// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package role_definition

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_role_definition"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input any, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := parse.RoleDefinitionId(v); err != nil {
			errors = append(errors, err)
		}

		return
	}
}
