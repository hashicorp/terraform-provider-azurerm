// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"encoding/json"
	"fmt"
	"strconv"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-testing/internal/addrs"
	"github.com/hashicorp/terraform-plugin-testing/internal/tfdiags"
)

type shimmedState struct {
	state *terraform.State
}

func shimStateFromJson(jsonState *tfjson.State) (*terraform.State, error) {
	state := terraform.NewState() //nolint:staticcheck // legacy usage
	state.TFVersion = jsonState.TerraformVersion

	if jsonState.Values == nil {
		// the state is empty
		return state, nil
	}

	for key, output := range jsonState.Values.Outputs {
		os, err := shimOutputState(output)
		if err != nil {
			return nil, err
		}
		state.RootModule().Outputs[key] = os
	}

	ss := &shimmedState{state}
	err := ss.shimStateModule(jsonState.Values.RootModule)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func shimOutputState(so *tfjson.StateOutput) (*terraform.OutputState, error) {
	os := &terraform.OutputState{
		Sensitive: so.Sensitive,
	}

	switch v := so.Value.(type) {
	case string:
		os.Type = "string"
		os.Value = v
		return os, nil
	case []interface{}:
		os.Type = "list"
		if len(v) == 0 {
			os.Value = v
			return os, nil
		}

		switch firstElem := v[0].(type) {
		case string:
			elements := make([]interface{}, len(v))
			for i, el := range v {
				strElement, ok := el.(string)
				// If the type of the element doesn't match the first elem, it's a tuple, return the original value
				if !ok {
					os.Value = v
					return os, nil
				}
				elements[i] = strElement
			}
			os.Value = elements
		case bool:
			elements := make([]interface{}, len(v))
			for i, el := range v {
				boolElement, ok := el.(bool)
				// If the type of the element doesn't match the first elem, it's a tuple, return the original value
				if !ok {
					os.Value = v
					return os, nil
				}

				elements[i] = boolElement
			}
			os.Value = elements
		// unmarshalled number from JSON will always be json.Number
		case json.Number:
			elements := make([]interface{}, len(v))
			for i, el := range v {
				numberElement, ok := el.(json.Number)
				// If the type of the element doesn't match the first elem, it's a tuple, return the original value
				if !ok {
					os.Value = v
					return os, nil
				}

				elements[i] = numberElement
			}
			os.Value = elements
		case []interface{}:
			os.Value = v
		case map[string]interface{}:
			os.Value = v
		default:
			return nil, fmt.Errorf("unexpected output list element type: %T", firstElem)
		}
		return os, nil
	case map[string]interface{}:
		os.Type = "map"
		os.Value = v
		return os, nil
	case bool:
		os.Type = "string"
		os.Value = strconv.FormatBool(v)
		return os, nil
	// unmarshalled number from JSON will always be json.Number
	case json.Number:
		os.Type = "string"
		os.Value = v.String()
		return os, nil
	}

	return nil, fmt.Errorf("unexpected output type: %T", so.Value)
}

func (ss *shimmedState) shimStateModule(sm *tfjson.StateModule) error {
	var path addrs.ModuleInstance

	if sm.Address == "" {
		path = addrs.RootModuleInstance
	} else {
		var diags tfdiags.Diagnostics
		path, diags = addrs.ParseModuleInstanceStr(sm.Address)
		if diags.HasErrors() {
			return diags.Err()
		}
	}

	mod := ss.state.AddModule(path) //nolint:staticcheck // legacy usage
	for _, res := range sm.Resources {
		resourceState, err := shimResourceState(res)
		if err != nil {
			return err
		}

		key, err := shimResourceStateKey(res)
		if err != nil {
			return err
		}

		mod.Resources[key] = resourceState
	}

	if len(sm.ChildModules) > 0 {
		return fmt.Errorf("Modules are not supported. Found %d modules.",
			len(sm.ChildModules))
	}
	return nil
}

func shimResourceStateKey(res *tfjson.StateResource) (string, error) {
	if res.Index == nil {
		return res.Address, nil
	}

	var mode terraform.ResourceMode
	switch res.Mode {
	case tfjson.DataResourceMode:
		mode = terraform.DataResourceMode
	case tfjson.ManagedResourceMode:
		mode = terraform.ManagedResourceMode
	default:
		return "", fmt.Errorf("unexpected resource mode for %q", res.Address)
	}

	var index int
	switch idx := res.Index.(type) {
	case json.Number:
		i, err := idx.Int64()
		if err != nil {
			return "", fmt.Errorf("unexpected index value (%q) for %q, ",
				idx, res.Address)
		}
		index = int(i)
	default:
		return "", fmt.Errorf("unexpected index type (%T) for %q, "+
			"for_each is not supported", res.Index, res.Address)
	}

	rsk := &terraform.ResourceStateKey{
		Mode:  mode,
		Type:  res.Type,
		Name:  res.Name,
		Index: index,
	}

	return rsk.String(), nil
}

func shimResourceState(res *tfjson.StateResource) (*terraform.ResourceState, error) {
	sf := &shimmedFlatmap{}
	err := sf.FromMap(res.AttributeValues)
	if err != nil {
		return nil, err
	}
	attributes := sf.Flatmap()

	// The instance state identifier was a Terraform versions 0.11 and earlier
	// concept which helped core and the then SDK determine if the resource
	// should be removed and as an identifier value in the human readable
	// output. This concept unfortunately carried over to the testing logic when
	// the testing logic was mostly changed to use the public, machine-readable
	// JSON interface with Terraform, rather than reusing prior internal logic
	// from Terraform. Using the "id" attribute value for this identifier was
	// the default implementation and therefore those older versions of
	// Terraform required the attribute. This is no longer necessary after
	// Terraform versions 0.12 and later.
	//
	// If the "id" attribute is not found, set the instance state identifier to
	// a synthetic value that can hopefully lead someone encountering the value
	// to these comments. The prior logic used to raise an error if the
	// attribute was not present, but this value should now only be present in
	// legacy logic of this Go module, such as unintentionally exported logic in
	// the terraform package, and not encountered during normal testing usage.
	//
	// Reference: https://github.com/hashicorp/terraform-plugin-testing/issues/84
	instanceStateID, ok := attributes["id"]

	if !ok {
		instanceStateID = "id-attribute-not-set"
	}

	return &terraform.ResourceState{
		Provider: res.ProviderName,
		Type:     res.Type,
		Primary: &terraform.InstanceState{
			ID:         instanceStateID,
			Attributes: attributes,
			Meta: map[string]interface{}{
				"schema_version": int(res.SchemaVersion),
			},
			Tainted: res.Tainted,
		},
		Dependencies: res.DependsOn,
	}, nil
}

type shimmedFlatmap struct {
	m map[string]string
}

func (sf *shimmedFlatmap) FromMap(attributes map[string]interface{}) error {
	if sf.m == nil {
		sf.m = make(map[string]string, len(attributes))
	}

	return sf.AddMap("", attributes)
}

func (sf *shimmedFlatmap) AddMap(prefix string, m map[string]interface{}) error {
	for key, value := range m {
		k := key
		if prefix != "" {
			k = fmt.Sprintf("%s.%s", prefix, key)
		}

		err := sf.AddEntry(k, value)
		if err != nil {
			return fmt.Errorf("unable to add map key %q entry: %w", k, err)
		}
	}

	mapLength := "%"
	if prefix != "" {
		mapLength = fmt.Sprintf("%s.%s", prefix, "%")
	}

	if err := sf.AddEntry(mapLength, strconv.Itoa(len(m))); err != nil {
		return fmt.Errorf("unable to add map length %q entry: %w", mapLength, err)
	}

	return nil
}

func (sf *shimmedFlatmap) AddSlice(name string, elements []interface{}) error {
	for i, elem := range elements {
		key := fmt.Sprintf("%s.%d", name, i)
		err := sf.AddEntry(key, elem)
		if err != nil {
			return fmt.Errorf("unable to add slice key %q entry: %w", key, err)
		}
	}

	sliceLength := fmt.Sprintf("%s.#", name)
	if err := sf.AddEntry(sliceLength, strconv.Itoa(len(elements))); err != nil {
		return fmt.Errorf("unable to add slice length %q entry: %w", sliceLength, err)
	}

	return nil
}

func (sf *shimmedFlatmap) AddEntry(key string, value interface{}) error {
	switch el := value.(type) {
	case nil:
		// omit the entry
		return nil
	case bool:
		sf.m[key] = strconv.FormatBool(el)
	case json.Number:
		sf.m[key] = el.String()
	case string:
		sf.m[key] = el
	case map[string]interface{}:
		err := sf.AddMap(key, el)
		if err != nil {
			return err
		}
	case []interface{}:
		err := sf.AddSlice(key, el)
		if err != nil {
			return err
		}
	default:
		// This should never happen unless terraform-json
		// changes how attributes (types) are represented.
		//
		// We handle all types which the JSON unmarshaler
		// can possibly produce
		// https://golang.org/pkg/encoding/json/#Unmarshal

		return fmt.Errorf("%q: unexpected type (%T)", key, el)
	}
	return nil
}

func (sf *shimmedFlatmap) Flatmap() map[string]string {
	return sf.m
}
