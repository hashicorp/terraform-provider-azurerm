// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	sdkclient "github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestMysqlFlexibleServerRenameCLI(t *testing.T) {
	const resourceName = "azurerm_mysql_flexible_server.test"
	var mu sync.Mutex
	current := servers.NewFlexibleServerID(renameTestSubscription, "ExampleGroup", "source")
	destination := current
	exists := false
	var endpoint string
	var renameCount, deleteCount int
	client, url := mysqlRenameTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch {
		case req.Method == http.MethodGet && req.URL.Path == "/operations/rename":
			current = destination
			fmt.Fprint(w, `{"status":"Succeeded"}`)
		case req.Method == http.MethodGet:
			if !exists || req.URL.Path != current.ID() {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, `{"error":{"code":"ResourceNotFound","message":"not found"}}`)
				return
			}
			if err := json.NewEncoder(w).Encode(mysqlRenameTestServer(current)); err != nil {
				t.Error(err)
			}
		case req.Method == http.MethodPost:
			fmt.Fprint(w, `{"nameAvailable":true}`)
		case req.Method == http.MethodPut:
			renameCount++
			id, err := servers.ParseFlexibleServerID(req.URL.Path)
			if err != nil {
				t.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if req.URL.Query().Get("api-version") != "2025-12-30" {
				t.Errorf("unexpected rename API: %s", req.URL)
			}
			var payload servers.Server
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if payload.Properties == nil || pointer.From(payload.Properties.SourceServerResourceId) != current.ID() {
				t.Errorf("rename did not use current source ID")
			}
			destination = *id
			w.Header().Set("Azure-AsyncOperation", endpoint+"/operations/rename")
			w.Header().Set(sdkclient.SkipPollingDelayHeader, "true")
			w.WriteHeader(http.StatusAccepted)
		case req.Method == http.MethodDelete:
			if req.URL.Path != current.ID() {
				t.Errorf("cleanup used stale ID: %s", req.URL)
			}
			deleteCount++
			exists = false
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request: %s %s", req.Method, req.URL)
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	endpoint = url
	config := func(name string) string {
		return fmt.Sprintf(`
resource "azurerm_mysql_flexible_server" "test" {
  name                   = %q
  resource_group_name    = "ExampleGroup"
  location               = "westus2"
  administrator_login    = "exampleadmin"
  administrator_password = "ExamplePassword123!"
  sku_name               = "GP_Standard_D2ds_v4"
  rename_enabled         = true
}

output "server_id" {
  value = azurerm_mysql_flexible_server.test.id
}
`, name)
	}
	step := func(name string, action plancheck.ResourceActionType) resource.TestStep {
		id := servers.NewFlexibleServerID(renameTestSubscription, "ExampleGroup", name)
		planChecks := []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, action),
			plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("fqdn")),
			plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("id")),
		}
		return resource.TestStep{
			Config: config(name),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: planChecks,
			},
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(resourceName, "name", name),
				resource.TestCheckResourceAttr(resourceName, "fqdn", name+".mysql.database.azure.com"),
				resource.TestCheckResourceAttr(resourceName, "id", id.ID()),
				resource.TestCheckOutput("server_id", id.ID()),
			),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectIdentityValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(name)),
			},
		}
	}
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.12.0"))),
		},
		ProtoV5ProviderFactories: map[string]func() (tfprotov5.ProviderServer, error){
			"azurerm": func() (tfprotov5.ProviderServer, error) {
				r := resourceMysqlFlexibleServer()
				// Seed a local source; all later operations use the production callbacks.
				r.Create = func(d *pluginsdk.ResourceData, _ any) error {
					mu.Lock()
					defer mu.Unlock()
					exists = true
					d.SetId(current.ID())
					return resourceMysqlFlexibleServerFlatten(d, &current, pointer.To(mysqlRenameTestServer(current)))
				}
				provider := &schema.Provider{
					ResourcesMap: map[string]*schema.Resource{"azurerm_mysql_flexible_server": r},
					ConfigureFunc: func(_ *schema.ResourceData) (any, error) {
						return client, nil
					},
				}
				return NewRenameAwareProviderServer(provider.GRPCProvider()), nil
			},
		},
		Steps: []resource.TestStep{
			step("source", plancheck.ResourceActionCreate),
			step("destination", plancheck.ResourceActionUpdate),
			{Config: config("destination"), PlanOnly: true},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"administrator_password", "rename_enabled"},
			},
			step("destination-again", plancheck.ResourceActionUpdate),
			{Config: config("destination-again"), PlanOnly: true},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"administrator_password", "rename_enabled"},
			},
		},
	})
	mu.Lock()
	defer mu.Unlock()
	if renameCount != 2 || deleteCount != 1 {
		t.Fatalf("expected two renames and final-ID cleanup, got rename=%d delete=%d", renameCount, deleteCount)
	}
}
