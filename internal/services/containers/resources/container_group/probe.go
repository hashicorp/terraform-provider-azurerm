// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_group

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerinstance/2025-09-01/containerinstance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaContainerGroupProbe() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Description: "The definition of a readiness probe for this container.",
		Type:        pluginsdk.TypeList,
		Optional:    true,
		ForceNew:    true,
		MaxItems:    1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"exec": {
					Description: "Commands to be run to validate container readiness.",
					Type:        pluginsdk.TypeList,
					Optional:    true,
					ForceNew:    true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.NoZeroValues,
					},
				},

				// lintignore:XS003
				"http_get": {
					Description: "The definition of the http_get for this container.",
					Type:        pluginsdk.TypeList,
					Optional:    true,
					ForceNew:    true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"path": {
								Description:  "Path to access on the HTTP server.",
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"port": {
								Description:  "Number of the port to access on the container.",
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.IsPortNumber,
							},
							"scheme": {
								Description:  "Scheme to use for connecting to the host.",
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice(containerinstance.PossibleValuesForScheme(), false),
							},
							"http_headers": {
								Description: "A map of HTTP headers used to access on the container.",
								Type:        pluginsdk.TypeMap,
								Optional:    true,
								ForceNew:    true,
								Elem: &pluginsdk.Schema{
									Type: pluginsdk.TypeString,
								},
							},
						},
					},
				},

				"initial_delay_seconds": {
					Description: "Number of seconds after the container has started before liveness or readiness probes are initiated.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					ForceNew:    true,
				},

				"period_seconds": {
					Description: "How often (in seconds) to perform the probe.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					ForceNew:    true,
				},

				"failure_threshold": {
					Description: "How many times to try the probe before restarting the container (liveness probe) or marking the container as unhealthy (readiness probe).",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					ForceNew:    true,
				},

				"success_threshold": {
					Description: "Minimum consecutive successes for the probe to be considered successful after having failed.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					ForceNew:    true,
				},

				"timeout_seconds": {
					Description: "Number of seconds after which the probe times out.",
					Type:        pluginsdk.TypeInt,
					Optional:    true,
					ForceNew:    true,
				},
			},
		},
	}
}
