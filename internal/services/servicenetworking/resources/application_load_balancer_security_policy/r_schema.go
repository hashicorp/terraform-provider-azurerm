// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_load_balancer_security_policy

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2025-01-01/securitypoliciesinterface"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (f Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9_.-]{0,62}[a-zA-Z0-9])?$`),
				"`name` must begin with a letter or number, end with a letter and number, and must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"application_load_balancer_id": commonschema.ResourceIDReferenceRequiredForceNew(&securitypoliciesinterface.TrafficControllerId{}, "Application Load Balancer"),

		"location": commonschema.Location(),

		"web_application_firewall_policy_id": commonschema.ResourceIDReferenceRequiredForceNew(&webapplicationfirewallpolicies.ApplicationGatewayWebApplicationFirewallPolicyId{}, "Web Application Firewall Policy"),

		"tags": commonschema.Tags(),
	}
}

func (f Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
