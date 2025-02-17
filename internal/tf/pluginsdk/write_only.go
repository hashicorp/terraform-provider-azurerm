package pluginsdk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-cty/cty"
)

func GetWriteOnly(d *ResourceData, name string, attributeType cty.Type) (*cty.Value, error) {
	value, diags := d.GetRawConfigAt(cty.GetAttrPath(name))
	if diags.HasError() {
		return nil, fmt.Errorf("blah")
	}

	if !value.Type().Equals(attributeType) {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: value is not of type %v", name, attributeType)
	}
	return pointer.To(value), nil
}
