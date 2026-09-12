// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// per https://learn.microsoft.com/en-us/rest/api/billing/invoice-sections/create-or-update
var invoiceSectionNamePattern = regexp.MustCompile(`^[a-zA-Z\d-_]{1,128}$`)

const invoiceSectionNameValidationMessage = "must be between 1 and 128 characters in length and may only contain letters, numbers, hyphens and underscores"

type BillingInvoiceSectionResource struct{}

var (
	_ sdk.ResourceWithUpdate   = BillingInvoiceSectionResource{}
	_ sdk.ResourceWithIdentity = BillingInvoiceSectionResource{}
)

type BillingInvoiceSectionResourceModel struct {
	Name             string            `tfschema:"name"`
	BillingProfileId string            `tfschema:"billing_profile_id"`
	DisplayName      string            `tfschema:"display_name"`
	Tags             map[string]string `tfschema:"tags"`

	SystemId    string `tfschema:"system_id"`
	State       string `tfschema:"state"`
	ReasonCode  string `tfschema:"reason_code"`
	TargetCloud string `tfschema:"target_cloud"`
}

func (BillingInvoiceSectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(invoiceSectionNamePattern, invoiceSectionNameValidationMessage),
		},

		"billing_profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: invoicesection.ValidateBillingProfileID,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (BillingInvoiceSectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"system_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reason_code": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"target_cloud": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (BillingInvoiceSectionResource) ModelObject() interface{} {
	return &BillingInvoiceSectionResourceModel{}
}

func (BillingInvoiceSectionResource) ResourceType() string {
	return "azurerm_billing_invoice_section"
}

func (BillingInvoiceSectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return invoicesection.ValidateInvoiceSectionID
}

func (BillingInvoiceSectionResource) Identity() resourceids.ResourceId {
	return &invoicesection.InvoiceSectionId{}
}

func (r BillingInvoiceSectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.InvoiceSectionClient

			var config BillingInvoiceSectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			billingProfileId, err := invoicesection.ParseBillingProfileID(config.BillingProfileId)
			if err != nil {
				return err
			}

			id := invoicesection.NewInvoiceSectionID(billingProfileId.BillingAccountName, billingProfileId.BillingProfileName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := expandBillingInvoiceSection(config)

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r BillingInvoiceSectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.InvoiceSectionClient

			id, err := invoicesection.ParseInvoiceSectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config BillingInvoiceSectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// the API has no `PATCH` operation, so the full payload is sent on each update
			payload := expandBillingInvoiceSection(config)

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r BillingInvoiceSectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.InvoiceSectionClient

			id, err := invoicesection.ParseInvoiceSectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (BillingInvoiceSectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.InvoiceSectionClient

			id, err := invoicesection.ParseInvoiceSectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// deletion is blocked whilst subscriptions or products are still associated with the
			// invoice section - checking upfront surfaces the reason rather than an opaque LRO
			// failure, but the check is advisory only and never blocks the deletion itself
			if eligibility, err := client.ValidateDeleteEligibility(ctx, *id); err != nil {
				log.Printf("[DEBUG] unable to validate delete eligibility for %s, continuing with the deletion: %+v", *id, err)
			} else if model := eligibility.Model; model != nil {
				if pointer.From(model.EligibilityStatus) == invoicesection.DeleteInvoiceSectionEligibilityStatusNotAllowed {
					return fmt.Errorf("%s cannot be deleted: %s", *id, flattenBillingInvoiceSectionEligibilityDetails(model.EligibilityDetails))
				}
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (BillingInvoiceSectionResource) flatten(metadata sdk.ResourceMetaData, id *invoicesection.InvoiceSectionId, model *invoicesection.InvoiceSection) error {
	state := BillingInvoiceSectionResourceModel{
		Name:             id.InvoiceSectionName,
		BillingProfileId: invoicesection.NewBillingProfileID(id.BillingAccountName, id.BillingProfileName).ID(),
	}

	if model != nil {
		state.Tags = flattenBillingInvoiceSectionTags(*model)

		if props := model.Properties; props != nil {
			state.DisplayName = pointer.From(props.DisplayName)
			state.SystemId = pointer.From(props.SystemId)
			state.State = string(pointer.From(props.State))
			state.ReasonCode = string(pointer.From(props.ReasonCode))
			state.TargetCloud = pointer.From(props.TargetCloud)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func expandBillingInvoiceSection(config BillingInvoiceSectionResourceModel) invoicesection.InvoiceSection {
	tags := pointer.To(config.Tags)
	if config.Tags == nil {
		tags = pointer.To(map[string]string{})
	}

	// the API documents `tags` both within `properties` and in the ARM envelope, so both are sent
	return invoicesection.InvoiceSection{
		Tags: tags,
		Properties: &invoicesection.InvoiceSectionProperties{
			DisplayName: pointer.To(config.DisplayName),
			Tags:        tags,
		},
	}
}

// flattenBillingInvoiceSectionTags prefers whichever of the two `tags` locations comes back
// populated, since expandBillingInvoiceSection writes to both
func flattenBillingInvoiceSectionTags(model invoicesection.InvoiceSection) map[string]string {
	if props := model.Properties; props != nil && len(pointer.From(props.Tags)) > 0 {
		return pointer.From(props.Tags)
	}

	return pointer.From(model.Tags)
}

func flattenBillingInvoiceSectionEligibilityDetails(input *[]invoicesection.DeleteInvoiceSectionEligibilityDetail) string {
	if input == nil || len(*input) == 0 {
		return "the Billing API didn't provide a reason"
	}

	reasons := make([]string, 0, len(*input))
	for _, detail := range *input {
		reasons = append(reasons, fmt.Sprintf("%s: %s", string(pointer.From(detail.Code)), pointer.From(detail.Message)))
	}

	return strings.Join(reasons, ", ")
}
