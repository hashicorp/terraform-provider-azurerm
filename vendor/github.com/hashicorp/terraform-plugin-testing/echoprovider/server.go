// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package echoprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// NewProviderServer returns the "echo" provider, which is a protocol v6 Terraform provider meant only to be used for testing
// data which cannot be stored in Terraform artifacts (plan/state), such as an ephemeral resource. The "echo" provider can be included in
// an acceptance test with the `(resource.TestCase).ProtoV6ProviderFactories` field, for example:
//
//	resource.UnitTest(t, resource.TestCase{
//		// .. other TestCase fields
//		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
//			"echo": echoprovider.NewProviderServer(),
//		},
//
//		// .. TestSteps
//	})
//
// The "echo" provider configuration accepts in a dynamic "data" attribute, which will be stored in the "echo" managed resource "data" attribute, for example:
//
//	// Ephemeral resource that is under test
//	ephemeral "examplecloud_thing" "this" {
//		name = "thing-one"
//	}
//
//	provider "echo" {
//		data = ephemeral.examplecloud_thing.this
//	}
//
//	resource "echo" "test" {} // The `echo.test.data` attribute will contain the ephemeral data from `ephemeral.examplecloud_thing.this`
func NewProviderServer() func() (tfprotov6.ProviderServer, error) {
	return func() (tfprotov6.ProviderServer, error) {
		return &echoProviderServer{}, nil
	}
}

// echoProviderServer is a lightweight protocol version 6 provider server that saves data from the provider configuration (which is considered ephemeral)
// and then stores that data into state during ApplyResourceChange.
//
// As provider configuration is ephemeral, it's possible for the data to change between plan and apply. As a result of this, the echo provider
// will never propose new changes after it has been created, making it immutable (during plan, echo will always use prior state for it's plan,
// regardless of what the provider configuration is set to). This prevents the managed resource from continuously proposing new planned changes
// if the ephemeral data changes.
type echoProviderServer struct {
	// The value of the "data" attribute during provider configuration. Will be directly echoed to the echo.data attribute.
	providerConfigData tftypes.Value
}

const echoResourceType = "echo"

func (e *echoProviderServer) providerSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Description: "This provider is used to output the data attribute provided to the provider configuration into all resources instances of echo. " +
				"This is only useful for testing ephemeral resources where the data isn't stored to state.",
			DescriptionKind: tfprotov6.StringKindPlain,
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:            "data",
					Type:            tftypes.DynamicPseudoType,
					Description:     "Dynamic data to provide to the echo resource.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Optional:        true,
				},
			},
		},
	}
}

func (e *echoProviderServer) testResourceSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:            "data",
					Type:            tftypes.DynamicPseudoType,
					Description:     "Dynamic data that was provided to the provider configuration.",
					DescriptionKind: tfprotov6.StringKindPlain,
					Computed:        true,
				},
			},
		},
	}
}

func (e *echoProviderServer) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	resp := &tfprotov6.ApplyResourceChangeResponse{}

	if req.TypeName != echoResourceType {
		resp.Diagnostics = []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource",
				Detail:   fmt.Sprintf("ApplyResourceChange was called for a resource type that is not supported by this provider: %q", req.TypeName),
			},
		}

		return resp, nil
	}

	echoTestSchema := e.testResourceSchema()

	plannedState, diag := dynamicValueToValue(echoTestSchema, req.PlannedState)
	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	// Destroy Op, just return planned state, which is null
	if plannedState.IsNull() {
		resp.NewState = req.PlannedState
		return resp, nil
	}

	// Take the provider config "data" attribute verbatim and put back into state. It shares the same type (DynamicPseudoType)
	// as the echo "data" attribute.
	newVal := tftypes.NewValue(echoTestSchema.ValueType(), map[string]tftypes.Value{
		"data": e.providerConfigData,
	})

	newState, diag := valuetoDynamicValue(echoTestSchema, newVal)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.NewState = newState

	return resp, nil
}

func (e *echoProviderServer) CallFunction(ctx context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	return &tfprotov6.CallFunctionResponse{}, nil
}

func (e *echoProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	resp := &tfprotov6.ConfigureProviderResponse{}

	configVal, diags := dynamicValueToValue(e.providerSchema(), req.Config)
	if diags != nil {
		resp.Diagnostics = append(resp.Diagnostics, diags)
		return resp, nil
	}

	objVal := map[string]tftypes.Value{}
	err := configVal.As(&objVal)
	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading Config",
			Detail:   err.Error(),
		}
		resp.Diagnostics = append(resp.Diagnostics, diag)
		return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
	}

	dynamicDataVal, ok := objVal["data"]
	if !ok {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  `Attribute "data" not found in config`,
		}
		resp.Diagnostics = append(resp.Diagnostics, diag)
		return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
	}

	e.providerConfigData = dynamicDataVal.Copy()

	return resp, nil
}

func (e *echoProviderServer) GetFunctions(ctx context.Context, req *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
	return &tfprotov6.GetFunctionsResponse{}, nil
}

func (e *echoProviderServer) GetMetadata(ctx context.Context, req *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
	return &tfprotov6.GetMetadataResponse{
		Resources: []tfprotov6.ResourceMetadata{
			{
				TypeName: echoResourceType,
			},
		},
	}, nil
}

func (e *echoProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	return &tfprotov6.GetProviderSchemaResponse{
		Provider: e.providerSchema(),
		// MAINTAINER NOTE: This provider is only really built to support a single special resource type ("echo"). In the future, if we want
		// to add more resource types to this provider, we'll likely need to refactor other RPCs in the provider server to handle that.
		ResourceSchemas: map[string]*tfprotov6.Schema{
			echoResourceType: e.testResourceSchema(),
		},
	}, nil
}

func (e *echoProviderServer) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	return &tfprotov6.ImportResourceStateResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource Operation",
				Detail:   "ImportResourceState is not supported by this provider.",
			},
		},
	}, nil
}

func (e *echoProviderServer) MoveResourceState(ctx context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
	return &tfprotov6.MoveResourceStateResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource Operation",
				Detail:   "MoveResourceState is not supported by this provider.",
			},
		},
	}, nil
}

func (e *echoProviderServer) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	resp := &tfprotov6.PlanResourceChangeResponse{}

	if req.TypeName != echoResourceType {
		resp.Diagnostics = []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource",
				Detail:   fmt.Sprintf("PlanResourceChange was called for a resource type that is not supported by this provider: %q", req.TypeName),
			},
		}

		return resp, nil
	}

	echoTestSchema := e.testResourceSchema()
	priorState, diag := dynamicValueToValue(echoTestSchema, req.PriorState)
	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	proposedNewState, diag := dynamicValueToValue(echoTestSchema, req.ProposedNewState)
	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	// Destroying the resource, just return proposed new state (which is null)
	if proposedNewState.IsNull() {
		return &tfprotov6.PlanResourceChangeResponse{
			PlannedState: req.ProposedNewState,
		}, nil
	}

	// If the echo resource has prior state, don't plan anything new as it's valid for the ephemeral data to change
	// between operations and we don't want to produce constant diffs. This resource is only for testing data, which a
	// single plan/apply should suffice.
	if !priorState.IsNull() {
		return &tfprotov6.PlanResourceChangeResponse{
			PlannedState: req.PriorState,
		}, nil
	}

	// If we are creating, mark data as unknown in the plan.
	//
	// We can't set the proposed new state to the provider config data because it could change between plan/apply (provider config is ephemeral).
	unknownVal := tftypes.NewValue(echoTestSchema.ValueType(), map[string]tftypes.Value{
		"data": tftypes.NewValue(tftypes.DynamicPseudoType, tftypes.UnknownValue),
	})

	plannedState, diag := valuetoDynamicValue(echoTestSchema, unknownVal)
	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.PlannedState = plannedState

	return resp, nil
}

func (e *echoProviderServer) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	return &tfprotov6.ReadDataSourceResponse{}, nil
}

func (e *echoProviderServer) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	// Just return current state, since the data doesn't need to be refreshed.
	return &tfprotov6.ReadResourceResponse{
		NewState: req.CurrentState,
	}, nil
}

func (e *echoProviderServer) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	return &tfprotov6.StopProviderResponse{}, nil
}

func (e *echoProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	resp := &tfprotov6.UpgradeResourceStateResponse{}

	if req.TypeName != echoResourceType {
		resp.Diagnostics = []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource",
				Detail:   fmt.Sprintf("UpgradeResourceState was called for a resource type that is not supported by this provider: %q", req.TypeName),
			},
		}

		return resp, nil
	}

	// Define options to be used when unmarshalling raw state.
	// IgnoreUndefinedAttributes will silently skip over fields in the JSON
	// that do not have a matching entry in the schema.
	unmarshalOpts := tfprotov6.UnmarshalOpts{
		ValueFromJSONOpts: tftypes.ValueFromJSONOpts{
			IgnoreUndefinedAttributes: true,
		},
	}

	providerSchema := e.providerSchema()

	if req.Version != providerSchema.Version {
		resp.Diagnostics = []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported Resource",
				Detail:   "UpgradeResourceState was called for echo, which does not support multiple schema versions",
			},
		}

		return resp, nil
	}

	// Terraform CLI can call UpgradeResourceState even if the stored state
	// version matches the current schema. Presumably this is to account for
	// the previous terraform-plugin-sdk implementation, which handled some
	// state fixups on behalf of Terraform CLI. This will attempt to roundtrip
	// the prior RawState to a state matching the current schema.
	rawStateValue, err := req.RawState.UnmarshalWithOpts(providerSchema.ValueType(), unmarshalOpts)

	if err != nil {
		diag := &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to Read Previously Saved State for UpgradeResourceState",
			Detail:   "There was an error reading the saved resource state using the current resource schema: " + err.Error(),
		}

		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil //nolint:nilerr // error via diagnostic, not gRPC
	}

	upgradedState, diag := valuetoDynamicValue(providerSchema, rawStateValue)

	if diag != nil {
		resp.Diagnostics = append(resp.Diagnostics, diag)

		return resp, nil
	}

	resp.UpgradedState = upgradedState

	return resp, nil
}

func (e *echoProviderServer) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	return &tfprotov6.ValidateDataResourceConfigResponse{}, nil
}

func (e *echoProviderServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	return &tfprotov6.ValidateProviderConfigResponse{}, nil
}

func (e *echoProviderServer) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	return &tfprotov6.ValidateResourceConfigResponse{}, nil
}

func (e *echoProviderServer) OpenEphemeralResource(ctx context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	return &tfprotov6.OpenEphemeralResourceResponse{}, nil
}

func (e *echoProviderServer) RenewEphemeralResource(ctx context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	return &tfprotov6.RenewEphemeralResourceResponse{}, nil
}

func (e *echoProviderServer) CloseEphemeralResource(ctx context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	return &tfprotov6.CloseEphemeralResourceResponse{}, nil
}

func (e *echoProviderServer) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	return &tfprotov6.ValidateEphemeralResourceConfigResponse{}, nil
}

func (e *echoProviderServer) GetResourceIdentitySchemas(context.Context, *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	return &tfprotov6.GetResourceIdentitySchemasResponse{}, nil
}

func (e *echoProviderServer) UpgradeResourceIdentity(context.Context, *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
	return &tfprotov6.UpgradeResourceIdentityResponse{
		Diagnostics: []*tfprotov6.Diagnostic{
			{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Unsupported UpgradeResourceIdentity Operation",
				Detail:   "Resource Identity is not supported by this provider.",
			},
		},
	}, nil
}
