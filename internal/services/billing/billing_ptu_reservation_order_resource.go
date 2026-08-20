// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	billingClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/billing/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BillingPtuReservationOrderResource struct{}

var _ sdk.Resource = BillingPtuReservationOrderResource{}

type BillingPtuReservationOrderResourceModel struct {
	Name             string `tfschema:"name"`
	Location         string `tfschema:"location"`
	Capacity         int64  `tfschema:"capacity"`
	BillingScopeId   string `tfschema:"billing_scope_id"`
	SkuName          string `tfschema:"sku_name"`
	Term             string `tfschema:"term"`
	BillingPlan      string `tfschema:"billing_plan"`
	AppliedScopeType string `tfschema:"applied_scope_type"`
	Renew            bool   `tfschema:"renew"`
	OrderId          string `tfschema:"order_id"`
}

func (r BillingPtuReservationOrderResource) ResourceType() string {
	return "azurerm_billing_ptu_reservation_order"
}

func (r BillingPtuReservationOrderResource) ModelObject() interface{} {
	return &BillingPtuReservationOrderResourceModel{}
}

func (r BillingPtuReservationOrderResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.ValidatePtuReservationOrderID
}

func (r BillingPtuReservationOrderResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Display name for the PTU reservation order.",
		},

		"location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Azure region where the PTU capacity is reserved (e.g. 'eastus').",
		},

		"capacity": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Number of PTUs (Provisioned Throughput Units) to reserve.",
		},

		"billing_scope_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Billing scope for the reservation (e.g. '/subscriptions/{subscriptionId}').",
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "DataZoneProvisionedManaged",
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "SKU name for the PTU reservation (e.g. 'DataZoneProvisionedManaged', 'openai_provisioned_managed_datazone').",
		},

		"term": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "P1Y",
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "Reservation term in ISO 8601 duration format (e.g. 'P1M', 'P1Y', 'P3Y').",
		},

		"billing_plan": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "Upfront",
			ValidateFunc: validation.StringInSlice([]string{
				"Upfront",
				"Monthly",
			}, false),
			Description: "Billing plan for the reservation.",
		},

		"applied_scope_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "Shared",
			ValidateFunc: validation.StringInSlice([]string{
				"Shared",
				"Single",
				"ManagementGroup",
			}, false),
			Description: "Scope type to which the reservation benefit is applied.",
		},

		"renew": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     true,
			Description: "Whether the reservation auto-renews at the end of the term. Changing this forces a new resource to be created.",
		},
	}
}

func (r BillingPtuReservationOrderResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"order_id": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "UUID of the reservation order in Azure.",
		},
	}
}

func (r BillingPtuReservationOrderResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing

			var model BillingPtuReservationOrderResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Build the payload without a reservation order ID first — CalculatePrice
			// will open a purchase session and return the order ID we must use.
			payload := billingClient.ReservationOrderPurchaseRequest{
				Sku:      &billingClient.ReservationOrderSku{Name: model.SkuName},
				Location: model.Location,
				Properties: &billingClient.ReservationOrderPurchaseProperties{
					DisplayName:          model.Name,
					ReservedResourceType: "OpenAIPTU",
					ReservedResourceProperties: &billingClient.ReservedResourceProperties{
						InstanceFlexibility: "On",
					},
					Term:             model.Term,
					BillingPlan:      model.BillingPlan,
					BillingScopeId:   model.BillingScopeId,
					AppliedScopeType: model.AppliedScopeType,
					Quantity:         model.Capacity,
					Renew:            model.Renew,
				},
			}

			// Step 1: CalculatePrice opens the purchase session and returns the
			// reservation order ID. This ID must be used in the subsequent PUT and
			// both calls must happen within the session window (a few minutes).
			orderID, err := client.CalculatePrice(ctx, payload)
			if err != nil {
				return fmt.Errorf("calculating price for PTU reservation %q: %+v", model.Name, err)
			}

			id := parse.NewPtuReservationOrderId(orderID)

			// Step 2: Purchase — use the session-approved order ID immediately.
			payload.Properties.ReservationOrderId = orderID
			if _, err = client.CreateOrUpdate(ctx, id.ID(), payload); err != nil {
				return fmt.Errorf("purchasing %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BillingPtuReservationOrderResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing

			id, err := parse.PtuReservationOrderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ID())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing == nil {
				return metadata.MarkAsGone(id)
			}

			var state BillingPtuReservationOrderResourceModel
			if err = metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state.OrderId = id.OrderId
			if existing.Properties != nil {
				state.Name = existing.Properties.DisplayName
				state.Capacity = existing.Properties.OriginalQuantity
				state.Renew = existing.Properties.Renew
				// Azure omits certain fields (e.g. appliedScopeType, term, billingPlan,
				// billingScopeId) when they are at their default / "Shared" values. Only
				// overwrite state when the API actually returns a non-empty value so we
				// don't create a phantom diff against the config.
				if existing.Properties.Term != "" {
					state.Term = existing.Properties.Term
				}
				if existing.Properties.BillingPlan != "" {
					state.BillingPlan = existing.Properties.BillingPlan
				}
				if existing.Properties.AppliedScopeType != "" {
					state.AppliedScopeType = existing.Properties.AppliedScopeType
				}
				if existing.Properties.BillingScopeId != "" {
					state.BillingScopeId = existing.Properties.BillingScopeId
				}
			}
			if existing.Sku != nil {
				state.SkuName = existing.Sku.Name
			}
			if existing.Location != "" {
				state.Location = existing.Location
			}

			return metadata.Encode(&state)
		},
	}
}

func (r BillingPtuReservationOrderResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (r BillingPtuReservationOrderResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing

			id, err := parse.PtuReservationOrderID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// Azure only allows reservation returns within 30 days of purchase.
			// If the return window has passed the API returns an error. We treat
			// this as a warning rather than a hard failure so the resource can be
			// removed from state when it can no longer be returned — operators
			// should rely on `lifecycle { prevent_destroy = true }` to guard
			// against accidental deletions during the active reservation term.
			if returnErr := client.Return(ctx, id.ID()); returnErr != nil {
				if IsReturnNotAllowed(returnErr) {
					log.Printf("[WARN] azurerm_billing_ptu_reservation_order: cannot return %s (return window may have passed or reservation type does not support returns): %+v", id, returnErr)
					log.Printf("[WARN] The reservation remains active in Azure until it expires. Only the Terraform state entry is being removed.")
					return nil
				}
				return fmt.Errorf("returning %s: %+v", id, returnErr)
			}

			return nil
		},
	}
}

// IsReturnNotAllowed returns true for error messages that indicate a reservation
// cannot be returned (expired window, unsupported type, etc.).
func IsReturnNotAllowed(err error) bool {
	msg := strings.ToLower(err.Error())
	for _, kw := range []string{
		"returnfailed",
		"returnnotallowed",
		"cannot be returned",
		"refund period",
		"is not eligible",
	} {
		if strings.Contains(msg, kw) {
			return true
		}
	}
	return false
}
