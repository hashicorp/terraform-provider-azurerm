// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestExpandFunctionAppDapr(t *testing.T) {
	testCases := []struct {
		name     string
		input    []FunctionAppDapr
		expected *webapps.DaprConfig
	}{
		{
			// Azure disables Dapr when `daprConfig` is omitted, so this is how the block is removed
			name:     "no block",
			input:    nil,
			expected: nil,
		},
		{
			name:  "minimal",
			input: []FunctionAppDapr{{AppId: "app", LogLevel: "info"}},
			expected: &webapps.DaprConfig{
				Enabled:          pointer.To(true),
				AppId:            pointer.To("app"),
				EnableApiLogging: pointer.To(false),
				LogLevel:         pointer.To(webapps.DaprLogLevelInfo),
			},
		},
		{
			name: "all fields",
			input: []FunctionAppDapr{{
				AppId:                  "app",
				AppPort:                3001,
				APILoggingEnabled:      true,
				HTTPMaxRequestSizeInMB: 8,
				HTTPReadBufferSizeInKB: 32,
				LogLevel:               "debug",
			}},
			expected: &webapps.DaprConfig{
				Enabled:            pointer.To(true),
				AppId:              pointer.To("app"),
				AppPort:            pointer.To(int64(3001)),
				EnableApiLogging:   pointer.To(true),
				HTTPMaxRequestSize: pointer.To(int64(8)),
				HTTPReadBufferSize: pointer.To(int64(32)),
				LogLevel:           pointer.To(webapps.DaprLogLevelDebug),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := ExpandFunctionAppDapr(testCase.input); !reflect.DeepEqual(actual, testCase.expected) {
				t.Fatalf("expected %+v but got %+v", testCase.expected, actual)
			}
		})
	}
}

func TestFlattenFunctionAppDapr(t *testing.T) {
	testCases := []struct {
		name     string
		input    *webapps.DaprConfig
		expected []FunctionAppDapr
	}{
		{name: "nil", input: nil, expected: []FunctionAppDapr{}},
		{
			// Azure keeps an object with `enabled: false` once Dapr has been disabled
			name:     "disabled",
			input:    &webapps.DaprConfig{Enabled: pointer.To(false), EnableApiLogging: pointer.To(false), LogLevel: pointer.To(webapps.DaprLogLevelInfo)},
			expected: []FunctionAppDapr{},
		},
		{
			name:     "enabled without enabled being returned",
			input:    &webapps.DaprConfig{AppId: pointer.To("app")},
			expected: []FunctionAppDapr{},
		},
		{
			name:     "minimal, log level defaulted by Azure",
			input:    &webapps.DaprConfig{Enabled: pointer.To(true), AppId: pointer.To("app"), LogLevel: pointer.To(webapps.DaprLogLevelInfo)},
			expected: []FunctionAppDapr{{AppId: "app", LogLevel: "info"}},
		},
		{
			name:     "missing log level uses the default",
			input:    &webapps.DaprConfig{Enabled: pointer.To(true), AppId: pointer.To("app")},
			expected: []FunctionAppDapr{{AppId: "app", LogLevel: "info"}},
		},
		{
			name: "all fields",
			input: &webapps.DaprConfig{
				Enabled:            pointer.To(true),
				AppId:              pointer.To("app"),
				AppPort:            pointer.To(int64(3001)),
				EnableApiLogging:   pointer.To(true),
				HTTPMaxRequestSize: pointer.To(int64(8)),
				HTTPReadBufferSize: pointer.To(int64(32)),
				LogLevel:           pointer.To(webapps.DaprLogLevelWarn),
			},
			expected: []FunctionAppDapr{{
				AppId:                  "app",
				AppPort:                3001,
				APILoggingEnabled:      true,
				HTTPMaxRequestSizeInMB: 8,
				HTTPReadBufferSizeInKB: 32,
				LogLevel:               "warn",
			}},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if actual := FlattenFunctionAppDapr(testCase.input); !reflect.DeepEqual(actual, testCase.expected) {
				t.Fatalf("expected %+v but got %+v", testCase.expected, actual)
			}
		})
	}
}

// A configuration round-tripped through Azure must produce no diff, so what is expanded must flatten back to it.
func TestFunctionAppDaprRoundTrip(t *testing.T) {
	input := []FunctionAppDapr{{
		AppId:                  "app",
		AppPort:                3001,
		APILoggingEnabled:      true,
		HTTPMaxRequestSizeInMB: 8,
		HTTPReadBufferSizeInKB: 32,
		LogLevel:               "error",
	}}

	if actual := FlattenFunctionAppDapr(ExpandFunctionAppDapr(input)); !reflect.DeepEqual(actual, input) {
		t.Fatalf("expected %+v but got %+v", input, actual)
	}
}

func TestFunctionAppDaprSchema(t *testing.T) {
	block := FunctionAppDaprSchema()
	if len(block.RequiredWith) != 1 || block.RequiredWith[0] != "container_app_environment_id" {
		t.Fatalf("expected `dapr` to require `container_app_environment_id` but got %+v", block.RequiredWith)
	}
}

func TestFunctionAppDaprSchemaHTTPReadBufferSize(t *testing.T) {
	elem := FunctionAppDaprSchema().Elem.(*pluginsdk.Resource)
	validate := elem.Schema["http_read_buffer_size_in_kb"].ValidateFunc

	// Azure rejects values above 65 with ContainerAppDaprHttpReadBufferSizeOutOfRange
	for value, valid := range map[int]bool{0: false, 1: true, 65: true, 66: false, 128: false} {
		_, errs := validate(value, "http_read_buffer_size_in_kb")
		if (len(errs) == 0) != valid {
			t.Fatalf("expected %d to be valid: %t but got errors %+v", value, valid, errs)
		}
	}
}
