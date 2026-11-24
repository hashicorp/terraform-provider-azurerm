// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/subscription/2021-10-01/subscriptions"
	billingValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func subscriptionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subscription_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Description:  "The Display Name for the Subscription.",
			ValidateFunc: validate.SubscriptionName,
		},

		"alias": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// Note: O+C - This value is supplied by the provider if omitted so must remain `Computed`
			Computed:     true,
			ForceNew:     true,
			Description:  "The Alias Name of the subscription. If omitted a new UUID will be generated for this property.",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"billing_scope_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ExactlyOneOf: []string{
				"subscription_id",
				"billing_scope_id",
			},
			ValidateFunc: validation.Any(
				billingValidate.MicrosoftCustomerAccountBillingScopeID,
				billingValidate.EnrollmentBillingScopeID,
				billingValidate.MicrosoftPartnerAccountBillingScopeID,
			),
		},

		"workload": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "The workload type for the Subscription.",
			// Other RP's have updated Constants with contextual prefixes so these are likely to change
			ValidateFunc: validation.StringInSlice(subscriptions.PossibleValuesForWorkload(), false),
			// Workload is not exposed in any way, so must be ignored if the resource is imported.
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				return new == ""
			},
		},

		"subscription_id": {
			Type:        pluginsdk.TypeString,
			Description: "The GUID of the Subscription.",
			ForceNew:    true,
			Optional:    true,
			// Note: O+C because Azure returns a computed value when a new subscription is created
			Computed: true,
			ExactlyOneOf: []string{
				"subscription_id",
				"billing_scope_id",
			},
			ValidateFunc: validation.IsUUID,
		},

		"tenant_id": {
			Type:        pluginsdk.TypeString,
			Description: "The Tenant ID to which the subscription belongs.",
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
