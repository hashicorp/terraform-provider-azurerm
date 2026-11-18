// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/webapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FunctionAppDapr struct {
	AppId                  string `tfschema:"app_id"`
	AppPort                int64  `tfschema:"app_port"`
	APILoggingEnabled      bool   `tfschema:"api_logging_enabled"`
	HTTPMaxRequestSizeInMB int64  `tfschema:"http_max_request_size_in_mb"`
	HTTPReadBufferSizeInKB int64  `tfschema:"http_read_buffer_size_in_kb"`
	LogLevel               string `tfschema:"log_level"`
}

func FunctionAppDaprSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeList,
		Optional:     true,
		MaxItems:     1,
		RequiredWith: []string{"container_app_environment_id"},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					Description:  "The Dapr Application Identifier.",
				},

				"app_port": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65535),
					Description:  "The port which the application is listening on.",
				},

				"api_logging_enabled": {
					Type:        pluginsdk.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Should API logging be enabled for the Dapr sidecar?",
				},

				"http_max_request_size_in_mb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntAtLeast(1),
					Description:  "The maximum size of the request body for the Dapr HTTP servers, in MB.",
				},

				"http_read_buffer_size_in_kb": {
					Type:         pluginsdk.TypeInt,
					Optional:     true,
					ValidateFunc: validation.IntBetween(1, 65),
					Description:  "The maximum size of the HTTP header read buffer for the Dapr sidecar, in KB.",
				},

				"log_level": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					Default:      string(webapps.DaprLogLevelInfo),
					ValidateFunc: validation.StringInSlice(webapps.PossibleValuesForDaprLogLevel(), false),
					Description:  "The log level for the Dapr sidecar.",
				},
			},
		},
	}
}

func FunctionAppDaprSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"app_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"app_port": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"api_logging_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"http_max_request_size_in_mb": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"http_read_buffer_size_in_kb": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},

				"log_level": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// ExpandFunctionAppDapr returns nil when no `dapr` block is configured. Azure disables Dapr when `daprConfig` is
// omitted from the request, so omitting it is how the block is removed.
func ExpandFunctionAppDapr(input []FunctionAppDapr) *webapps.DaprConfig {
	if len(input) == 0 {
		return nil
	}

	dapr := input[0]
	output := &webapps.DaprConfig{
		Enabled:          pointer.To(true),
		AppId:            pointer.To(dapr.AppId),
		EnableApiLogging: pointer.To(dapr.APILoggingEnabled),
		LogLevel:         pointer.ToEnum[webapps.DaprLogLevel](dapr.LogLevel),
	}

	if dapr.AppPort != 0 {
		output.AppPort = pointer.To(dapr.AppPort)
	}
	if dapr.HTTPMaxRequestSizeInMB != 0 {
		output.HTTPMaxRequestSize = pointer.To(dapr.HTTPMaxRequestSizeInMB)
	}
	if dapr.HTTPReadBufferSizeInKB != 0 {
		output.HTTPReadBufferSize = pointer.To(dapr.HTTPReadBufferSizeInKB)
	}

	return output
}

// FlattenFunctionAppDapr returns no block when Dapr is not enabled, since Azure keeps an `enabled: false` object
// for an app where Dapr was explicitly disabled.
func FlattenFunctionAppDapr(input *webapps.DaprConfig) []FunctionAppDapr {
	if input == nil || !pointer.From(input.Enabled) {
		return []FunctionAppDapr{}
	}

	logLevel := string(webapps.DaprLogLevelInfo)
	if input.LogLevel != nil && *input.LogLevel != "" {
		logLevel = string(*input.LogLevel)
	}

	return []FunctionAppDapr{{
		AppId:                  pointer.From(input.AppId),
		AppPort:                pointer.From(input.AppPort),
		APILoggingEnabled:      pointer.From(input.EnableApiLogging),
		HTTPMaxRequestSizeInMB: pointer.From(input.HTTPMaxRequestSize),
		HTTPReadBufferSizeInKB: pointer.From(input.HTTPReadBufferSize),
		LogLevel:               logLevel,
	}}
}
