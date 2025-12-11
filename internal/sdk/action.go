// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sdk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

const timeoutAttributeName = "timeout"

type ActionMetadata struct {
	Client *clients.Client

	SubscriptionId string

	Features features.UserFeatures
}

type BaseActionModel struct {
	Timeout types.String `tfsdk:"timeout"`
}

type ActionWrapper struct {
	ActionMetadata

	WrappedAction
}

var _ action.ActionWithConfigure = &ActionWrapper{}

func (a *ActionMetadata) Defaults(_ context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	if request.ProviderData == nil {
		// response.Diagnostics.AddError("Client Provider Data Error", "null provider data supplied")
		return
	}

	c, ok := request.ProviderData.(*clients.Client)
	if !ok {
		response.Diagnostics.AddError("Client Provider Data Error", "invalid provider data supplied")
		return
	}

	a.Client = c
	a.SubscriptionId = c.Account.SubscriptionId
	a.Features = c.Features
}

type WrappedAction interface {
	Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse)

	Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse)

	Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse, config any, metadata ActionMetadata)

	ModelObject() any

	// for consistent timeouts and ctx handling
	// working around the fact that action attributes have no defaults
	// or... request defaults support for action attributes?
	Timeout() time.Duration
}

type WrappedActionWithConfigure interface {
	WrappedAction

	Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse)
}

var _ action.Action = &ActionWrapper{}

func (a *ActionWrapper) Action() func() action.Action {
	return func() action.Action {
		return a
	}
}

func (a *ActionWrapper) Schema(ctx context.Context, request action.SchemaRequest, response *action.SchemaResponse) {
	a.WrappedAction.Schema(ctx, request, response)

	if _, ok := response.Schema.Attributes[timeoutAttributeName]; ok {
		// Ensures there are no conflicting `timeout` attributes
		panic("internal error - `timeout` is a reserved attribute and must not be included in the Action's schema")
	}

	response.Schema.Attributes[timeoutAttributeName] = schema.StringAttribute{
		Optional:            true,
		Description:         fmt.Sprintf("Timeout duration for the action to complete. Defaults to `%s`", a.WrappedAction.Timeout()),
		MarkdownDescription: fmt.Sprintf("Timeout duration for the action to complete. Defaults to `%s`", a.WrappedAction.Timeout()),
		Validators: []validator.String{
			Timeout(),
		},
	}
}

// TODO: move validation to go-azure-helpers

type timeoutValidator struct{}

var _ validator.String = timeoutValidator{}

func Timeout() timeoutValidator {
	return timeoutValidator{}
}

func (t timeoutValidator) Description(_ context.Context) string {
	return "Unable to Parse Timeout"
}

func (t timeoutValidator) MarkdownDescription(ctx context.Context) string {
	return t.Description(ctx)
}

func (t timeoutValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	if _, err := time.ParseDuration(request.ConfigValue.ValueString()); err != nil {
		response.Diagnostics.Append(
			diag.NewAttributeErrorDiagnostic(
				request.Path,
				t.Description(ctx),
				err.Error(),
			),
		)
	}
}

/////////////////////

func (a *ActionWrapper) Metadata(ctx context.Context, request action.MetadataRequest, response *action.MetadataResponse) {
	a.WrappedAction.Metadata(ctx, request, response)
}

func (a *ActionWrapper) Invoke(ctx context.Context, request action.InvokeRequest, response *action.InvokeResponse) {
	model := a.WrappedAction.ModelObject()
	response.Diagnostics.Append(request.Config.Get(ctx, model)...)
	if response.Diagnostics.HasError() {
		return
	}

	var timeoutAttr types.String
	response.Diagnostics.Append(request.Config.GetAttribute(ctx, path.Root(timeoutAttributeName), &timeoutAttr)...)
	if response.Diagnostics.HasError() {
		return
	}

	timeout := a.WrappedAction.Timeout()
	if !timeoutAttr.IsNull() {
		v, err := time.ParseDuration(timeoutAttr.ValueString())
		if err != nil {
			response.Diagnostics.AddError("parsing `timeout`", err.Error())
			return
		}
		timeout = v
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	a.WrappedAction.Invoke(ctx, request, response, model, a.ActionMetadata)
}

func (a *ActionWrapper) Configure(ctx context.Context, request action.ConfigureRequest, response *action.ConfigureResponse) {
	a.Defaults(ctx, request, response)

	if wa, ok := a.WrappedAction.(WrappedActionWithConfigure); ok {
		wa.Configure(ctx, request, response)
	}
}

///

func AssertActionModelType[T any](model any, response *action.InvokeResponse) *T {
	// TODO add to `AssertResourceModelType (and rename) or keep separate?
	result, ok := model.(*T)
	if !ok {
		response.Diagnostics.AddError("action had incorrect model type", fmt.Sprintf("got %T", model))

		return nil
	}

	return result
}
