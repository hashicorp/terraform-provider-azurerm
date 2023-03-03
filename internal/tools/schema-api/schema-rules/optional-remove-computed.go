package schema_rules

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

var _ BreakingChangeRule = optionalRemoveComputed{}

type optionalRemoveComputed struct {
}

func (optionalRemoveComputed) Check(base providerjson.SchemaJSON, current providerjson.SchemaJSON, propertyName string) *string {
	if (base.Optional && base.Computed) && (current.Optional && !current.Computed) {
		return pointer.To(fmt.Sprintf("cannot remove Computed from the Optional property %q as the user config may not supply the value, thus causing a diff", propertyName))
	}

	return nil
}
