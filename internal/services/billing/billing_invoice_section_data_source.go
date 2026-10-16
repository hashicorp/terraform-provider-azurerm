// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package billing

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2024-04-01/invoicesection"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BillingInvoiceSectionDataSource struct{}

var _ sdk.DataSource = BillingInvoiceSectionDataSource{}

type BillingInvoiceSectionDataSourceModel struct {
	Name             string            `tfschema:"name"`
	BillingProfileId string            `tfschema:"billing_profile_id"`
	DisplayName      string            `tfschema:"display_name"`
	SystemId         string            `tfschema:"system_id"`
	State            string            `tfschema:"state"`
	ReasonCode       string            `tfschema:"reason_code"`
	TargetCloud      string            `tfschema:"target_cloud"`
	Tags             map[string]string `tfschema:"tags"`
}

func (BillingInvoiceSectionDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringMatch(invoiceSectionNamePattern, invoiceSectionNameValidationMessage),
		},

		"billing_profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: invoicesection.ValidateBillingProfileID,
		},
	}
}

func (BillingInvoiceSectionDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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

		"tags": commonschema.TagsDataSource(),
	}
}

func (BillingInvoiceSectionDataSource) ModelObject() interface{} {
	return &BillingInvoiceSectionDataSourceModel{}
}

func (BillingInvoiceSectionDataSource) ResourceType() string {
	return "azurerm_billing_invoice_section"
}

func (BillingInvoiceSectionDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Billing.InvoiceSectionClient

			var config BillingInvoiceSectionDataSourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			billingProfileId, err := invoicesection.ParseBillingProfileID(config.BillingProfileId)
			if err != nil {
				return err
			}

			id := invoicesection.NewInvoiceSectionID(billingProfileId.BillingAccountName, billingProfileId.BillingProfileName, config.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			state := BillingInvoiceSectionDataSourceModel{
				Name:             id.InvoiceSectionName,
				BillingProfileId: invoicesection.NewBillingProfileID(id.BillingAccountName, id.BillingProfileName).ID(),
				Tags:             flattenBillingInvoiceSectionTags(*resp.Model),
			}

			if props := resp.Model.Properties; props != nil {
				state.DisplayName = pointer.From(props.DisplayName)
				state.SystemId = pointer.From(props.SystemId)
				state.State = string(pointer.From(props.State))
				state.ReasonCode = string(pointer.From(props.ReasonCode))
				state.TargetCloud = pointer.From(props.TargetCloud)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
