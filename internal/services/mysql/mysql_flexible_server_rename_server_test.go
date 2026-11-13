// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type mysqlRenameTestProviderServer struct {
	tfprotov5.ProviderServer
	response    *tfprotov5.PlanResourceChangeResponse
	schema      *tfprotov5.Schema
	planError   error
	schemaError error
}

func (s *mysqlRenameTestProviderServer) PlanResourceChange(context.Context, *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	return s.response, s.planError
}

func (s *mysqlRenameTestProviderServer) GetProviderSchema(context.Context, *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	return &tfprotov5.GetProviderSchemaResponse{
		ResourceSchemas: map[string]*tfprotov5.Schema{mysqlFlexibleServerResourceName: s.schema},
	}, s.schemaError
}

func TestMysqlFlexibleServerRenameServer(t *testing.T) {
	resourceSchema := &tfprotov5.Schema{Block: &tfprotov5.SchemaBlock{
		Attributes: []*tfprotov5.SchemaAttribute{
			{Name: "id", Type: tftypes.String, Computed: true},
			{Name: "name", Type: tftypes.String, Required: true},
			{Name: "fqdn", Type: tftypes.String, Computed: true},
			{Name: "rename_enabled", Type: tftypes.Bool, Optional: true},
		},
	}}
	cases := []struct {
		name               string
		disabled           bool
		unknownOptIn       bool
		unknownName        bool
		destinationName    string
		knownID            bool
		knownFQDN          bool
		emptySourceID      bool
		create             bool
		destroy            bool
		otherResource      bool
		otherReplacement   bool
		noReplacement      bool
		errorDiagnostic    bool
		planError          bool
		schemaError        bool
		missingSchema      bool
		invalidPlannedData bool
		wantError          bool
		wantRename         bool
	}{
		{name: "opted-in rename", wantRename: true},
		{name: "disabled", disabled: true},
		{name: "unknown opt-in", unknownOptIn: true},
		{name: "unknown destination", unknownName: true},
		{name: "unchanged name", destinationName: "source"},
		{name: "case-only name", destinationName: "SOURCE"},
		{name: "known ID", knownID: true},
		{name: "known FQDN", knownFQDN: true},
		{name: "empty source ID", emptySourceID: true},
		{name: "create", create: true},
		{name: "destroy", destroy: true},
		{name: "other resource", otherResource: true},
		{name: "real replacement", otherReplacement: true},
		{name: "no replacement", noReplacement: true},
		{name: "error diagnostic", errorDiagnostic: true},
		{name: "plan error", planError: true, wantError: true},
		{name: "schema error", schemaError: true, wantError: true},
		{name: "missing schema", missingSchema: true, wantError: true},
		{name: "invalid planned data", invalidPlannedData: true, wantError: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value := func(attributes map[string]tftypes.Value) *tfprotov5.DynamicValue {
				t.Helper()
				var raw any = attributes
				if attributes == nil {
					raw = nil
				}
				result, err := tfprotov5.NewDynamicValue(resourceSchema.ValueType(), tftypes.NewValue(resourceSchema.ValueType(), raw))
				if err != nil {
					t.Fatal(err)
				}
				return &result
			}
			prior := map[string]tftypes.Value{
				"id":             tftypes.NewValue(tftypes.String, "source-id"),
				"name":           tftypes.NewValue(tftypes.String, "source"),
				"fqdn":           tftypes.NewValue(tftypes.String, "source.mysql.database.azure.com"),
				"rename_enabled": tftypes.NewValue(tftypes.Bool, true),
			}
			planned := map[string]tftypes.Value{
				"id":             tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"name":           tftypes.NewValue(tftypes.String, "destination"),
				"fqdn":           tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
				"rename_enabled": tftypes.NewValue(tftypes.Bool, !tc.disabled),
			}
			if tc.unknownOptIn {
				planned["rename_enabled"] = tftypes.NewValue(tftypes.Bool, tftypes.UnknownValue)
			}
			if tc.unknownName {
				planned["name"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
			}
			if tc.destinationName != "" {
				planned["name"] = tftypes.NewValue(tftypes.String, tc.destinationName)
			}
			if tc.knownID {
				planned["id"] = prior["id"]
			}
			if tc.knownFQDN {
				planned["fqdn"] = prior["fqdn"]
			}
			if tc.emptySourceID {
				prior["id"] = tftypes.NewValue(tftypes.String, "")
			}
			if tc.create {
				prior = nil
			}
			if tc.destroy {
				planned = nil
			}
			req := &tfprotov5.PlanResourceChangeRequest{
				TypeName:   mysqlFlexibleServerResourceName,
				PriorState: value(prior),
			}
			if tc.otherResource {
				req.TypeName = "azurerm_other_resource"
			}
			paths := []*tftypes.AttributePath{tftypes.NewAttributePath().WithAttributeName("id")}
			if tc.otherReplacement {
				paths = append(paths, tftypes.NewAttributePath().WithAttributeName("location"))
			}
			if tc.noReplacement {
				paths = nil
			}
			resp := &tfprotov5.PlanResourceChangeResponse{
				PlannedState:    value(planned),
				RequiresReplace: paths,
				PlannedPrivate:  []byte(`{"preserve":"private"}`),
			}
			if tc.errorDiagnostic {
				resp.Diagnostics = []*tfprotov5.Diagnostic{{Severity: tfprotov5.DiagnosticSeverityError, Summary: "invalid rename"}}
			}
			if tc.invalidPlannedData {
				resp.PlannedState = &tfprotov5.DynamicValue{JSON: []byte(`invalid`)}
			}
			server := &mysqlRenameTestProviderServer{response: resp, schema: resourceSchema}
			if tc.planError {
				server.planError = errors.New("plan failed")
			}
			if tc.schemaError {
				server.schemaError = errors.New("schema failed")
			}
			if tc.missingSchema {
				server.schema = nil
			}
			original := *resp
			got, err := NewRenameAwareProviderServer(server).PlanResourceChange(context.Background(), req)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v, want error = %t", err, tc.wantError)
			}
			if tc.wantRename {
				original.RequiresReplace = nil
			}
			if !reflect.DeepEqual(got, &original) {
				t.Fatalf("unexpected response changes: got %#v, want %#v", got, &original)
			}
		})
	}
}
