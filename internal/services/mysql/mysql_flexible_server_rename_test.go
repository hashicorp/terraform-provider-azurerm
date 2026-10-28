// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	mysqlclient "github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const renameTestSubscription = "00000000-0000-0000-0000-000000000000"

func mysqlRenameTestConfig(name string, enabled bool) map[string]any {
	return map[string]any{
		"name": name, "rename_enabled": enabled,
		"resource_group_name": "ExampleGroup", "location": "westus2",
		"administrator_login": "exampleadmin", "administrator_password": "ExamplePassword123!",
		"sku_name": "GP_Standard_D2ds_v4",
	}
}

func mysqlRenameTestState(t *testing.T, name string, enabled bool) *terraform.InstanceState {
	t.Helper()
	r := resourceMysqlFlexibleServer()
	d := schema.TestResourceDataRaw(t, r.Schema, mysqlRenameTestConfig(name, enabled))
	id := servers.NewFlexibleServerID(renameTestSubscription, "ExampleGroup", name)
	d.SetId(id.ID())
	if err := d.Set("fqdn", name+".mysql.database.azure.com"); err != nil {
		t.Fatal(err)
	}
	d = r.Data(d.State())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		t.Fatal(err)
	}
	return d.State()
}

func TestMysqlFlexibleServerRenameDiff(t *testing.T) {
	if err := resourceMysqlFlexibleServer().InternalValidate(nil, true); err != nil {
		t.Fatalf("invalid server schema: %v", err)
	}
	cases := []struct {
		name           string
		enabled        bool
		oldEnabled     bool
		newName        string
		other          map[string]any
		replace        bool
		wantError      bool
		unknownEnabled bool
		unknownName    bool
	}{
		{name: "default replacement", newName: "destination", replace: true},
		{name: "explicit rename", enabled: true, oldEnabled: true, newName: "destination"},
		{name: "enable and rename", enabled: true, newName: "destination"},
		{name: "disable and replace", oldEnabled: true, newName: "destination", replace: true},
		{name: "enable only", enabled: true, newName: "source"},
		{name: "same name", enabled: true, oldEnabled: true, newName: "source"},
		{name: "case only", enabled: true, newName: "SOURCE", wantError: true},
		{name: "rename with tags", enabled: true, newName: "destination", other: map[string]any{"tags": map[string]any{"env": "test"}}, wantError: true},
		{name: "rename with location", enabled: true, newName: "destination", other: map[string]any{"location": "eastus"}, wantError: true},
		{name: "rename with resource group", enabled: true, newName: "destination", other: map[string]any{"resource_group_name": "another"}, wantError: true},
		{name: "unknown opt-in", newName: "destination", unknownEnabled: true, wantError: true},
		{name: "unknown rename name", enabled: true, unknownName: true, wantError: true},
		{name: "unknown replacement name", unknownName: true, replace: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceMysqlFlexibleServer()
			if !r.ResourceBehavior.MutableIdentity {
				t.Fatal("rename must allow mutable resource identity")
			}
			state := mysqlRenameTestState(t, "source", tc.oldEnabled)
			config := mysqlRenameTestConfig(tc.newName, tc.enabled)
			for k, v := range tc.other {
				config[k] = v
			}
			rawConfig := terraform.NewResourceConfigRaw(config)
			if tc.unknownEnabled || tc.unknownName {
				values := make(map[string]cty.Value)
				for name, valueType := range r.CoreConfigSchema().ImpliedType().AttributeTypes() {
					values[name] = cty.NullVal(valueType)
				}
				for name, value := range config {
					if str, ok := value.(string); ok {
						values[name] = cty.StringVal(str)
					}
				}
				values["rename_enabled"] = cty.BoolVal(tc.enabled)
				if tc.unknownEnabled {
					values["rename_enabled"] = cty.UnknownVal(cty.Bool)
				}
				if tc.unknownName {
					values["name"] = cty.UnknownVal(cty.String)
				}
				rawConfig = terraform.NewResourceConfigShimmed(cty.ObjectVal(values), r.CoreConfigSchema())
			}
			diff, err := r.Diff(context.Background(), state, rawConfig, nil)
			if tc.wantError {
				if err == nil {
					t.Fatal("expected a planning error")
				}
				if (tc.unknownEnabled || tc.unknownName) && !strings.Contains(err.Error(), "must be known") {
					t.Fatalf("expected unknown input diagnostic, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if diff == nil {
				if tc.newName != "source" || tc.oldEnabled != tc.enabled {
					t.Fatal("expected a diff")
				}
				return
			}
			if diff.RequiresNew() != tc.replace {
				t.Fatalf("RequiresNew = %t, want %t", diff.RequiresNew(), tc.replace)
			}
			if tc.enabled && tc.newName != "source" {
				for _, key := range []string{"id", "fqdn"} {
					if attr := diff.Attributes[key]; attr == nil || !attr.NewComputed {
						t.Errorf("%s must become unknown in the rename plan: %#v", key, attr)
					}
				}
			}
		})
	}
}

func mysqlRenameTestClient(t *testing.T, handler http.Handler) (*clients.Client, string) {
	t.Helper()
	t.Setenv("GO_AZURE_SDK_SKIP_POLLING_DELAY", "true")
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	api := environments.ResourceManagerAPI(server.URL)
	mysql, err := mysqlclient.NewClient(&common.ClientOptions{
		Environment:                 environments.Environment{ResourceManager: api},
		Authorizers:                 &common.Authorizers{},
		DisableCorrelationRequestID: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, client := range []*resourcemanager.Client{
		mysql.FlexibleServers.Servers.Client,
		mysql.FlexibleServers.CheckNameAvailability.Client,
		mysql.FlexibleServerRename.Client,
	} {
		client.AuthorizeRequest = nil
	}
	return &clients.Client{
		StopContext: context.Background(),
		MySQL:       mysql,
	}, server.URL
}

func mysqlRenameTestServer(id servers.FlexibleServerId) servers.Server {
	return servers.Server{
		Id: pointer.To(id.ID()), Location: "westus2",
		Properties: &servers.ServerProperties{
			State:                    pointer.To(servers.ServerStateReady),
			FullyQualifiedDomainName: pointer.To(id.FlexibleServerName + ".mysql.database.azure.com"),
			AdministratorLogin:       pointer.To("exampleadmin"),
			Backup: &servers.Backup{
				BackupRetentionDays: pointer.To(int64(7)),
				GeoRedundantBackup:  pointer.To(servers.EnableStatusEnumDisabled),
			},
		},
		Sku: &servers.MySQLServerSku{Name: "Standard_D2ds_v4", Tier: servers.ServerSkuTierGeneralPurpose},
	}
}

func TestMysqlFlexibleServerRenameLifecycle(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusAccepted} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			source := servers.NewFlexibleServerID(renameTestSubscription, "ExampleGroup", "source")
			current := source
			var destination servers.FlexibleServerId
			var url string
			var putCount, pollCount, patchCount, deleteCount int
			var availabilityCount int
			deleted := false
			client, endpoint := mysqlRenameTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case req.Method == http.MethodGet && req.URL.Path == "/operations/rename":
					pollCount++
					current = destination
					fmt.Fprint(w, `{"status":"Succeeded"}`)
				case req.Method == http.MethodGet:
					if req.URL.Query().Get("api-version") != "2023-12-30" {
						t.Errorf("ordinary GET changed API version: %s", req.URL)
					}
					if req.URL.Path != current.ID() || deleted {
						w.WriteHeader(http.StatusNotFound)
						fmt.Fprint(w, `{"error":{"code":"ResourceNotFound","message":"not found"}}`)
						return
					}
					model := mysqlRenameTestServer(current)
					if err := json.NewEncoder(w).Encode(model); err != nil {
						t.Error(err)
					}
				case req.Method == http.MethodPost:
					availabilityCount++
					if status == http.StatusCreated && strings.Contains(req.URL.Path, "/locations/") {
						w.WriteHeader(http.StatusForbidden)
						fmt.Fprint(w, `{"error":{"code":"AuthorizationFailed","message":"use the subscription endpoint"}}`)
						return
					}
					if status != http.StatusCreated && !strings.Contains(req.URL.Path, "/locations/westus2/checkNameAvailability") {
						t.Errorf("availability must use source location: %s", req.URL)
					}
					var payload map[string]any
					if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
						t.Error(err)
					}
					if payload["type"] != "Microsoft.DBforMySQL/flexibleServers" {
						t.Errorf("incorrect availability resource type: %#v", payload)
					}
					fmt.Fprint(w, `{"nameAvailable":true}`)
				case req.Method == http.MethodPut:
					putCount++
					if req.URL.Query().Get("api-version") != "2025-12-30" {
						t.Errorf("incorrect rename API version: %s", req.URL)
					}
					parsed, err := servers.ParseFlexibleServerID(req.URL.Path)
					if err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					destination = *parsed
					var payload map[string]any
					if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
						t.Error(err)
					}
					expected := map[string]any{
						"location":   "westus2",
						"properties": map[string]any{"createMode": "Rename", "sourceServerResourceId": current.ID()},
					}
					if !reflect.DeepEqual(payload, expected) {
						t.Errorf("rename must contain only the minimal payload: got %#v, want %#v", payload, expected)
					}
					w.Header().Set("Azure-AsyncOperation", url+"/operations/rename")
					w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
					w.WriteHeader(status)
				case req.Method == http.MethodPatch:
					patchCount++
					if req.URL.Path != current.ID() {
						t.Errorf("update used stale ID: %s", req.URL.Path)
					}
					fmt.Fprint(w, `{"properties":{"provisioningState":"Succeeded"}}`)
				case req.Method == http.MethodDelete:
					deleteCount++
					if req.URL.Path != current.ID() {
						t.Errorf("delete used stale ID: %s", req.URL.Path)
					}
					deleted = true
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("unexpected request: %s %s", req.Method, req.URL)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			url = endpoint
			r := resourceMysqlFlexibleServer()
			state := mysqlRenameTestState(t, source.FlexibleServerName, true)
			for _, name := range []string{"destination", "destination-again"} {
				config := terraform.NewResourceConfigRaw(mysqlRenameTestConfig(name, true))
				diff, err := r.Diff(context.Background(), state, config, nil)
				if err != nil {
					t.Fatal(err)
				}
				next, diagnostics := r.Apply(context.Background(), state, diff, client)
				if diagnostics.HasError() {
					t.Fatalf("rename failed: %+v", diagnostics)
				}
				if next.ID != current.ID() || next.Attributes["name"] != name || next.Attributes["fqdn"] != name+".mysql.database.azure.com" {
					t.Fatalf("incorrect renamed state: %#v", next)
				}
				if next.Identity["name"] != name {
					t.Fatalf("identity not updated: %#v", next.Identity)
				}
				diff, err = r.Diff(context.Background(), next, config, nil)
				if err != nil {
					t.Fatal(err)
				}
				if diff != nil && (len(diff.Attributes) != 0 || diff.Destroy || !reflect.DeepEqual(diff.Identity, next.Identity)) {
					t.Fatalf("rename did not converge: attributes=%#v identity=%#v", diff.Attributes, diff.Identity)
				}
				state = next
			}
			if putCount != 2 || pollCount != 2 || patchCount != 0 || deleteCount != 0 {
				t.Fatalf("unexpected requests: put=%d poll=%d patch=%d delete=%d", putCount, pollCount, patchCount, deleteCount)
			}
			expectedAvailability := 2
			if status == http.StatusCreated {
				expectedAvailability = 4
			}
			if availabilityCount != expectedAvailability {
				t.Fatalf("availability requests = %d, want %d", availabilityCount, expectedAvailability)
			}
			config := mysqlRenameTestConfig(current.FlexibleServerName, true)
			config["tags"] = map[string]any{"env": "test"}
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
			if err != nil {
				t.Fatal(err)
			}
			state, diagnostics := r.Apply(context.Background(), state, diff, client)
			if diagnostics.HasError() {
				t.Fatalf("post-rename update failed: %+v", diagnostics)
			}
			d := r.Data(state)
			if err := resourceMysqlFlexibleServerDelete(d, client); err != nil {
				t.Fatal(err)
			}
			if patchCount != 1 || deleteCount != 1 {
				t.Fatalf("expected update and delete against destination: patch=%d delete=%d", patchCount, deleteCount)
			}
		})
	}
}

func TestMysqlFlexibleServerRenameFailuresRetainSource(t *testing.T) {
	cases := []struct {
		name               string
		sourceStatus       int
		sourceState        servers.ServerState
		existing           bool
		unavailable        bool
		putStatus          int
		operationState     string
		missingHeader      bool
		badDestination     bool
		cancelPolling      bool
		pollStatus         int
		destinationStatus  int
		availabilityStatus int
		badSourceID        bool
		wantPut            int
	}{
		{name: "source missing", sourceStatus: http.StatusNotFound},
		{name: "source forbidden", sourceStatus: http.StatusForbidden},
		{name: "source not ready", sourceState: servers.ServerStateUpdating},
		{name: "unexpected source ID", badSourceID: true},
		{name: "destination exists", existing: true},
		{name: "destination forbidden", destinationStatus: http.StatusForbidden},
		{name: "name unavailable", unavailable: true},
		{name: "availability forbidden", availabilityStatus: http.StatusForbidden},
		{name: "feature disabled", putStatus: http.StatusBadRequest, wantPut: 1},
		{name: "unauthorized", putStatus: http.StatusForbidden, wantPut: 1},
		{name: "conflict", putStatus: http.StatusConflict, wantPut: 1},
		{name: "failed operation", operationState: "Failed", wantPut: 1},
		{name: "canceled operation", operationState: "Canceled", wantPut: 1},
		{name: "polling forbidden", pollStatus: http.StatusForbidden, wantPut: 1},
		{name: "interrupted after acceptance", cancelPolling: true, wantPut: 1},
		{name: "missing operation URL", missingHeader: true, wantPut: 1},
		{name: "unexpected destination", operationState: "Succeeded", badDestination: true, wantPut: 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := servers.NewFlexibleServerID(renameTestSubscription, "ExampleGroup", "source")
			destination := servers.NewFlexibleServerID(renameTestSubscription, "ExampleGroup", "destination")
			var url string
			var putCount int
			completed := false
			var cancel context.CancelFunc
			client, endpoint := mysqlRenameTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				switch {
				case req.Method == http.MethodGet && req.URL.Path == "/operations/rename":
					if tc.cancelPolling {
						cancel()
						return
					}
					if tc.pollStatus != 0 {
						w.WriteHeader(tc.pollStatus)
						fmt.Fprint(w, `{"error":{"code":"PollError","message":"cannot read operation"}}`)
						return
					}
					completed = tc.operationState == "Succeeded"
					fmt.Fprintf(w, `{"status":%q,"error":{"code":"RenameFailed","message":"service error"}}`, tc.operationState)
				case req.Method == http.MethodGet && req.URL.Path == source.ID():
					if tc.sourceStatus != 0 {
						w.WriteHeader(tc.sourceStatus)
						fmt.Fprint(w, `{"error":{"code":"SourceError","message":"source error"}}`)
						return
					}
					model := mysqlRenameTestServer(source)
					if tc.badSourceID {
						model.Id = pointer.To(destination.ID())
					}
					if tc.sourceState != "" {
						model.Properties.State = &tc.sourceState
					}
					json.NewEncoder(w).Encode(model)
				case req.Method == http.MethodGet && req.URL.Path == destination.ID():
					if tc.destinationStatus != 0 {
						w.WriteHeader(tc.destinationStatus)
						fmt.Fprint(w, `{"error":{"code":"DestinationError","message":"cannot read destination"}}`)
						return
					}
					if !tc.existing && !completed {
						w.WriteHeader(http.StatusNotFound)
						fmt.Fprint(w, `{"error":{"code":"ResourceNotFound","message":"not found"}}`)
						return
					}
					model := mysqlRenameTestServer(destination)
					if tc.badDestination {
						model.Id = pointer.To(source.ID())
					}
					json.NewEncoder(w).Encode(model)
				case req.Method == http.MethodPost:
					if tc.availabilityStatus != 0 {
						w.WriteHeader(tc.availabilityStatus)
						fmt.Fprint(w, `{"error":{"code":"Forbidden","message":"cannot check name"}}`)
						return
					}
					fmt.Fprintf(w, `{"nameAvailable":%t,"message":"already taken"}`, !tc.unavailable)
				case req.Method == http.MethodPut:
					putCount++
					if tc.putStatus != 0 {
						w.WriteHeader(tc.putStatus)
						fmt.Fprint(w, `{"error":{"code":"RenameRejected","message":"service rejected rename"}}`)
						return
					}
					if !tc.missingHeader {
						w.Header().Set("Azure-AsyncOperation", url+"/operations/rename")
					}
					w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
					w.WriteHeader(http.StatusAccepted)
				default:
					t.Errorf("unexpected mutation or request: %s %s", req.Method, req.URL)
					w.WriteHeader(http.StatusBadRequest)
				}
			}))
			url = endpoint
			client.StopContext, cancel = context.WithCancel(context.Background())
			defer cancel()
			r := resourceMysqlFlexibleServer()
			state := mysqlRenameTestState(t, "source", true)
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(mysqlRenameTestConfig("destination", true)), nil)
			if err != nil {
				t.Fatal(err)
			}
			next, diagnostics := r.Apply(context.Background(), state, diff, client)
			if !diagnostics.HasError() {
				t.Fatal("expected rename failure")
			}
			if next.ID != source.ID() || next.Attributes["name"] != "source" || next.Identity["name"] != "source" ||
				next.Attributes["fqdn"] != "source.mysql.database.azure.com" {
				t.Fatalf("failed rename changed source identity: %#v", next)
			}
			if putCount != tc.wantPut {
				t.Fatalf("PUT count = %d, want %d", putCount, tc.wantPut)
			}
		})
	}
}

func TestMysqlFlexibleServerRenameReadMissing(t *testing.T) {
	client, _ := mysqlRenameTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			t.Errorf("unexpected mutation: %s", req.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"error":{"code":"ResourceNotFound","message":"not found"}}`)
	}))
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			r := resourceMysqlFlexibleServer()
			d := r.Data(mysqlRenameTestState(t, "source", enabled))
			err := resourceMysqlFlexibleServerRead(d, client)
			if enabled {
				if err == nil || d.Id() == "" {
					t.Fatal("a missing rename-enabled server must retain its identity and report an error")
				}
			} else if err != nil || d.Id() != "" {
				t.Fatalf("ordinary missing-resource behavior changed: id=%s err=%v", d.Id(), err)
			}
		})
	}
}

func TestMysqlFlexibleServerRenameToggleDoesNotMutate(t *testing.T) {
	id := servers.NewFlexibleServerID(renameTestSubscription, "ExampleGroup", "source")
	client, _ := mysqlRenameTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.Path != id.ID() {
			t.Errorf("toggling rename must only read the server: %s %s", req.Method, req.URL)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(mysqlRenameTestServer(id)); err != nil {
			t.Error(err)
		}
	}))
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			r := resourceMysqlFlexibleServer()
			state := mysqlRenameTestState(t, "source", !enabled)
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(mysqlRenameTestConfig("source", enabled)), nil)
			if err != nil {
				t.Fatal(err)
			}
			next, diagnostics := r.Apply(context.Background(), state, diff, client)
			if diagnostics.HasError() || next.ID != state.ID {
				t.Fatalf("toggle failed: %+v, state: %#v", diagnostics, next)
			}
		})
	}
}

func TestMysqlFlexibleServerRenameCanceledContext(t *testing.T) {
	client, _ := mysqlRenameTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Errorf("canceled operation made a request: %s %s", req.Method, req.URL)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	cancel()
	r := resourceMysqlFlexibleServer()
	d := r.Data(mysqlRenameTestState(t, "source", true))
	if err := d.Set("name", "destination"); err != nil {
		t.Fatal(err)
	}
	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		t.Fatal(err)
	}
	if err := resourceMysqlFlexibleServerRename(ctx, d, client, *id); err == nil || d.Id() != id.ID() {
		t.Fatalf("canceled rename must fail without changing ID: %v", err)
	}
}
