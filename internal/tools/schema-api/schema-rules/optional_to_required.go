package schema_rules

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

type optionalToRequired struct {
}

var _ BreakingChangeRule = optionalToRequired{}

func (o optionalToRequired) Check(base providerjson.SchemaJSON, current providerjson.SchemaJSON, propertyName string) *string {
	if base.Optional && current.Required {
		return pointer.To(fmt.Sprintf("Cannot change property %q from Optional to Required", propertyName))
	}

	return nil
}
