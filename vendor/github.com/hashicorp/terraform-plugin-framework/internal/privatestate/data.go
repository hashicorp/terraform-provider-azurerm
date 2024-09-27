// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatestate

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
)

// Data contains private state data for the framework and providers.
type Data struct {
	// Potential future usage:
	// Framework contains private state data for framework usage.
	Framework map[string][]byte

	// Provider contains private state data for provider usage.
	Provider *ProviderData
}

// Bytes returns a JSON encoded slice of bytes containing the merged
// framework and provider private state data.
func (d *Data) Bytes(ctx context.Context) ([]byte, diag.Diagnostics) {
	var diags diag.Diagnostics

	if d == nil {
		return nil, nil
	}

	if (d.Provider == nil || len(d.Provider.data) == 0) && len(d.Framework) == 0 {
		return nil, nil
	}

	var providerData map[string][]byte

	if d.Provider != nil {
		providerData = d.Provider.data
	}

	mergedMap := make(map[string][]byte, len(d.Framework)+len(providerData))

	for _, m := range []map[string][]byte{d.Framework, providerData} {
		for k, v := range m {
			if len(v) == 0 {
				continue
			}

			// Values in FrameworkData and ProviderData should never be invalid UTF-8, but let's make sure.
			if !utf8.Valid(v) {
				diags.AddError(
					"Error Encoding Private State",
					"An error was encountered when validating private state value."+
						fmt.Sprintf("The value associated with key %q is is not valid UTF-8.\n\n", k)+
						"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.",
				)

				tflog.Error(ctx, "error encoding private state: invalid UTF-8 value", map[string]interface{}{"key": k, "value": v})

				continue
			}

			// Values in FrameworkData and ProviderData should never be invalid JSON, but let's make sure.
			if !json.Valid(v) {
				diags.AddError(
					"Error Encoding Private State",
					fmt.Sprintf("An error was encountered when validating private state value."+
						fmt.Sprintf("The value associated with key %q is is not valid JSON.\n\n", k)+
						"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer."),
				)

				tflog.Error(ctx, "error encoding private state: invalid JSON value", map[string]interface{}{"key": k, "value": v})

				continue
			}

			mergedMap[k] = v
		}
	}

	if diags.HasError() {
		return nil, diags
	}

	bytes, err := json.Marshal(mergedMap)
	if err != nil {
		diags.AddError(
			"Error Encoding Private State",
			fmt.Sprintf("An error was encountered when encoding private state: %s.\n\n"+
				"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.", err),
		)

		return nil, diags
	}

	return bytes, diags
}

// NewData creates a new Data based on the given slice of bytes.
// It must be a JSON encoded slice of bytes, that is map[string][]byte.
func NewData(ctx context.Context, data []byte) (*Data, diag.Diagnostics) {
	var (
		dataMap map[string][]byte
		diags   diag.Diagnostics
	)

	if len(data) == 0 {
		return nil, nil
	}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		// terraform-plugin-sdk stored private state by marshalling its data
		// as map[string]any, which is slightly incompatible with trying to
		// unmarshal it as map[string][]byte. If unmarshalling with
		// map[string]any works, we can ignore it for now, as provider
		// developers did not have access to managing the private state data.
		//
		// TODO: We can extract the terraform-plugin-sdk resource timeouts key
		// here to extract its prior data, if necessary.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/400
		if anyErr := json.Unmarshal(data, new(map[string]any)); anyErr == nil {
			logging.FrameworkWarn(ctx, "Discarding incompatible resource private state data", map[string]any{logging.KeyError: err.Error()})
			return nil, nil
		}

		diags.AddError(
			"Error Decoding Private State",
			fmt.Sprintf("An error was encountered when decoding private state: %s.\n\n"+
				"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.", err),
		)

		return nil, diags
	}

	output := Data{
		Framework: make(map[string][]byte),
		Provider: &ProviderData{
			make(map[string][]byte),
		},
	}

	for k, v := range dataMap {
		if !utf8.Valid(v) {
			diags.AddError(
				"Error Decoding Private State",
				"An error was encountered when validating private state value.\n"+
					fmt.Sprintf("The value being supplied for key %q is is not valid UTF-8.\n\n", k)+
					"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.",
			)

			tflog.Error(ctx, "error decoding private state: invalid UTF-8 value", map[string]interface{}{"key": k, "value": v})

			continue
		}

		if !json.Valid(v) {
			diags.AddError(
				"Error Decoding Private State",
				"An error was encountered when validating private state value.\n"+
					fmt.Sprintf("The value being supplied for key %q is is not valid JSON.\n\n", k)+
					"This is always a problem with Terraform or terraform-plugin-framework. Please report this to the provider developer.",
			)

			tflog.Error(ctx, "error decoding private state: invalid JSON value", map[string]interface{}{"key": k, "value": v})

			continue
		}

		if isInvalidProviderDataKey(ctx, k) {
			output.Framework[k] = v
			continue
		}

		output.Provider.data[k] = v
	}

	if diags.HasError() {
		return nil, diags
	}

	return &output, diags
}

// EmptyData creates an initialised but empty Data.
func EmptyData(ctx context.Context) *Data {
	return &Data{
		Provider: EmptyProviderData(ctx),
	}
}

// NewProviderData creates a new ProviderData based on the given slice of bytes.
// It must be a JSON encoded slice of bytes, that is map[string][]byte.
func NewProviderData(ctx context.Context, data []byte) (*ProviderData, diag.Diagnostics) {
	providerData := EmptyProviderData(ctx)

	if len(data) == 0 {
		return providerData, nil
	}

	var (
		dataMap map[string][]byte
		diags   diag.Diagnostics
	)

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		diags.AddError(
			"Error Decoding Provider Data",
			fmt.Sprintf("An error was encountered when decoding provider data: %s.\n\n"+
				"Please check that the data you are supplying is a byte representation of valid JSON.", err),
		)

		return nil, diags
	}

	for k, v := range dataMap {
		diags.Append(providerData.SetKey(ctx, k, v)...)
	}

	if diags.HasError() {
		return nil, diags
	}

	return providerData, diags
}

// EmptyProviderData creates a ProviderData containing initialised but empty data.
func EmptyProviderData(ctx context.Context) *ProviderData {
	return &ProviderData{
		data: make(map[string][]byte),
	}
}

// ProviderData contains private state data for provider usage.
type ProviderData struct {
	data map[string][]byte
}

// Equal returns true if the given ProviderData is exactly equivalent. The
// internal data is compared byte-for-byte, not accounting for semantic
// equivalency such as JSON whitespace or property reordering.
func (d *ProviderData) Equal(o *ProviderData) bool {
	if d == nil && o == nil {
		return true
	}

	if d == nil || o == nil {
		return false
	}

	if !reflect.DeepEqual(d.data, o.data) {
		return false
	}

	return true
}

// GetKey returns the private state data associated with the given key.
//
// If the key is reserved for framework usage, an error diagnostic
// is returned. If the key is valid, but private state data is not found,
// nil is returned.
//
// The naming of keys only matters in context of a single resource,
// however care should be taken that any historical keys are not reused
// without accounting for older resource instances that may still have
// older data at the key.
func (d *ProviderData) GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics) {
	if d == nil || d.data == nil {
		return nil, nil
	}

	diags := ValidateProviderDataKey(ctx, key)

	if diags.HasError() {
		return nil, diags
	}

	value, ok := d.data[key]
	if !ok {
		return nil, nil
	}

	return value, nil
}

// SetKey sets the private state data at the given key.
//
// If the key is reserved for framework usage, an error diagnostic
// is returned. The data must be valid JSON and UTF-8 safe or an error
// diagnostic is returned.
//
// The naming of keys only matters in context of a single resource,
// however care should be taken that any historical keys are not reused
// without accounting for older resource instances that may still have
// older data at the key.
func (d *ProviderData) SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics {
	var diags diag.Diagnostics

	if d == nil {
		tflog.Error(ctx, "error calling SetKey on uninitialized ProviderData")

		diags.AddError("Uninitialized ProviderData",
			"ProviderData must be initialized before it is used.\n\n"+
				"Call privatestate.NewProviderData to obtain an initialized instance of ProviderData.",
		)

		return diags
	}

	if d.data == nil {
		d.data = make(map[string][]byte)
	}

	diags.Append(ValidateProviderDataKey(ctx, key)...)

	if diags.HasError() {
		return diags
	}

	// Support removing keys by setting them to nil or zero-length value.
	if len(value) == 0 {
		delete(d.data, key)

		return diags
	}

	if !utf8.Valid(value) {
		tflog.Error(ctx, "invalid UTF-8 value", map[string]interface{}{"key": key, "value": value})

		diags.AddError("UTF-8 Invalid",
			"Values stored in private state must be valid UTF-8.\n\n"+
				fmt.Sprintf("The value being supplied for key %q is invalid. Please verify that the value is valid UTF-8.", key),
		)

		return diags
	}

	if !json.Valid(value) {
		tflog.Error(ctx, "invalid JSON value", map[string]interface{}{"key": key, "value": value})

		diags.AddError("JSON Invalid",
			"Values stored in private state must be valid JSON.\n\n"+
				fmt.Sprintf("The value being supplied for key %q is invalid. Please verify that the value is valid JSON.", key),
		)

		return diags
	}

	d.data[key] = value

	return nil
}

// ValidateProviderDataKey determines whether the key supplied is allowed on the basis of any
// restrictions that are in place, such as key prefixes that are reserved for use with
// framework private state data.
func ValidateProviderDataKey(ctx context.Context, key string) diag.Diagnostics {
	if isInvalidProviderDataKey(ctx, key) {
		return diag.Diagnostics{
			diag.NewErrorDiagnostic(
				"Restricted Resource Private State Namespace",
				"Using a period ('.') as a prefix for a key used in private state is not allowed.\n\n"+
					fmt.Sprintf("The key %q is invalid. Please check the key you are supplying does not use a a period ('.') as a prefix.", key),
			),
		}
	}

	return nil
}

// isInvalidProviderDataKey determines whether the supplied key has a prefix that is reserved for
// keys in Data.Framework
func isInvalidProviderDataKey(_ context.Context, key string) bool {
	return strings.HasPrefix(key, ".")
}

// MustMarshalToJson is for use in tests and panics if input cannot be marshalled to JSON.
func MustMarshalToJson(input map[string][]byte) []byte {
	output, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	return output
}

// MustProviderData is for use in tests and panics if the underlying call to NewProviderData
// returns diag.Diagnostics that contains any errors.
func MustProviderData(ctx context.Context, data []byte) *ProviderData {
	providerData, diags := NewProviderData(ctx, data)

	if diags.HasError() {
		var diagMsgs []string

		for _, v := range diags {
			diagMsgs = append(diagMsgs, fmt.Sprintf("%s: %s", v.Summary(), v.Detail()))
		}

		panic(fmt.Sprintf("error creating new provider data: %s", strings.Join(diagMsgs, ", ")))
	}

	return providerData
}
