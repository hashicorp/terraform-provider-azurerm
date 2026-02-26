package sdk

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// TODO: make generic for all JSON types or is this sufficient for our use-case?
type PrivateStateValue struct {
	Value string `json:"value"`
}

func NewPrivateStateValue(value string) PrivateStateValue {
	return PrivateStateValue{Value: value}
}

func (v PrivateStateValue) Equals(data []byte, diags *diag.Diagnostics) bool {
	if data == nil {
		return false
	}

	value := NewPrivateStateValue("")
	if err := json.Unmarshal(data, &value); err != nil {
		// This shouldn't happen given the provider is in control of setting the private state data
		diags.Append(diag.NewErrorDiagnostic("Internal-error: parsing private state value", err.Error()))
		return false
	}

	return v.Value == value.Value
}

func (v PrivateStateValue) Bytes() []byte {
	return []byte(v.String())
}

func (v PrivateStateValue) String() string {
	return fmt.Sprintf(`{"value": "%s"}`, v.Value)
}
