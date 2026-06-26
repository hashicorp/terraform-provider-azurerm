package statecheck

import (
	"context"
	"fmt"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

var _ statecheck.StateCheck = expectAllIdentityFieldsAreChecked{}

type expectAllIdentityFieldsAreChecked struct {
	resourceAddress string
	checkedFields   map[string]struct{}
}

func (e expectAllIdentityFieldsAreChecked) CheckState(_ context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	var resource *tfjson.StateResource

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")

		return
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")

		return
	}

	if req.State.Values.RootModule == nil {
		resp.Error = fmt.Errorf("state does not contain a root module")

		return
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}

	if resource.IdentitySchemaVersion == nil || len(resource.IdentityValues) == 0 {
		resp.Error = fmt.Errorf("%s - Identity not found in state. Either the resource does not support identity or the Terraform version running the test does not support identity. (must be v1.12+)", e.resourceAddress)

		return
	}

	uncheckedFields := make([]string, 0)
	for k := range resource.IdentityValues {
		if _, ok := e.checkedFields[k]; !ok {
			uncheckedFields = append(uncheckedFields, k)
		}
	}

	if len(uncheckedFields) > 0 {
		resp.Error = fmt.Errorf("%s - Identity Schema contained fields that were not explicitly checked (%s)", e.resourceAddress, strings.Join(uncheckedFields, ", "))
	}
}

// ExpectAllIdentityFieldsAreChecked is a check for the generated resource identity tests and schema.
// it ensures that each schema field is explicitly checked by passing the schema fields using the `-properties`, `known-values`, or `compare-values` flags.
func ExpectAllIdentityFieldsAreChecked(resourceAddress string, fields map[string]struct{}) statecheck.StateCheck {
	return expectAllIdentityFieldsAreChecked{
		resourceAddress: resourceAddress,
		checkedFields:   fields,
	}
}
