// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_endpoint

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Manages a Private Endpoint.",
		Create:      resourcePrivateEndpointCreate,
		Read:        resourcePrivateEndpointRead,
		Update:      resourcePrivateEndpointUpdate,
		Delete:      resourcePrivateEndpointDelete,
		Importer:    pluginsdk.ImporterValidatingIdentity(&privateendpoints.PrivateEndpointId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&privateendpoints.PrivateEndpointId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: privateEndpointSchema(),
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta any) error {
			rawConfig := d.GetRawConfig()
			if d.Id() != "" && d.HasChange("ip_version_type") && rawConfig.IsKnown() && !rawConfig.IsNull() && rawConfig.GetAttr("ip_version_type").IsKnown() {
				old, new := d.GetChange("ip_version_type")
				current := old.(string)
				if current == "" {
					current = string(privateendpoints.PrivateEndpointIPVersionTypeIPvFour)
				}
				next := new.(string)
				// Azure only supports updating the IP version type in place from IPv4 to DualStack.
				if current != next && (current != string(privateendpoints.PrivateEndpointIPVersionTypeIPvFour) || next != string(privateendpoints.PrivateEndpointIPVersionTypeDualStack)) {
					return d.ForceNew("ip_version_type")
				}
			}

			privateServiceConnections := d.Get("private_service_connection").([]any)
			for i, psc := range privateServiceConnections {
				privateServiceConnection := psc.(map[string]any)
				name := privateServiceConnection["name"].(string)

				isManual := privateServiceConnection["manual_connection"].(bool)
				if !features.SixPointOh() && !pluginsdk.IsExplicitlyNullInConfigAt(d, fmt.Sprintf("private_service_connection.%d.is_manual_connection", i)) {
					isManual = privateServiceConnection["is_manual_connection"].(bool)
				}

				// If this is not a manual connection and the message is set return an error since this does not make sense.
				if !isManual && privateServiceConnection["request_message"].(string) != "" {
					return fmt.Errorf(`"private_service_connection":%q is invalid, the "request_message" attribute cannot be set if the "manual_connection" attribute is "false"`, name)
				}

				// If this is a manual connection and the message isn't set return an error.
				if isManual && strings.TrimSpace(privateServiceConnection["request_message"].(string)) == "" {
					return fmt.Errorf(`"private_service_connection":%q is invalid, the "request_message" attribute must not be empty`, name)
				}
			}
			return nil
		},
	}
}
