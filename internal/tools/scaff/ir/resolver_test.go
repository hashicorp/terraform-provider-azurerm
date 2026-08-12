// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ir

import (
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"
)

func serverAvailable(baseURL string) bool {
	c := &http.Client{Timeout: 2 * time.Second}
	resp, err := c.Get(baseURL + "/v1/resource-manager/services")
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// TestLive_ResolveRedHatOpenShift asserts the derived CRUD/ID metadata for the
// RedHatOpenShift cluster resource against a running Pandora Data API.
func TestLive_ResolveRedHatOpenShift(t *testing.T) {
	client := pandora.NewClient(pandora.DefaultBaseURL)
	if !serverAvailable(client.BaseURL) {
		t.Skipf("pandora data api not reachable at %s; skipping", client.BaseURL)
	}

	res, err := Resolve(client, Options{
		ARMType:        "Microsoft.RedHatOpenShift/openShiftClusters",
		APIVersion:     "2025-07-25",
		Name:           "redhat_openshift_cluster",
		GoName:         "RedHatOpenShiftCluster",
		ServicePackage: "redhatopenshift",
		ProviderName:   "azurerm",
	})
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	checks := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"ServiceName", res.ServiceName, "RedHatOpenShift"},
		{"ResourceProvider", res.ResourceProvider, "Microsoft.RedHatOpenShift"},
		{"PandoraResource", res.PandoraResource, "OpenShiftClusters"},
		{"SDKPackage", res.SDKPackage, "openshiftclusters"},
		{"ClientField", res.ClientField, "OpenShiftClustersClient"},
		{"SDKImportPath", res.SDKImportPath, "github.com/hashicorp/go-azure-sdk/resource-manager/redhatopenshift/2025-07-25/openshiftclusters"},
		{"Name", res.Name, "RedHatOpenShiftCluster"},
		{"TerraformType", res.TerraformType, "azurerm_redhat_openshift_cluster"},
		{"ModelStructName", res.ModelStructName, "RedHatOpenShiftClusterModel"},
		{"IDTypeName", res.IDTypeName, "OpenShiftClusterId"},
		{"IDParseFunc", res.IDParseFunc, "ParseOpenShiftClusterID"},
		{"IDNewFunc", res.IDNewFunc, "NewOpenShiftClusterID"},
		{"HasSubscription", res.HasSubscription, true},
		{"HasResourceGroup", res.HasResourceGroup, true},
		{"Updatable", res.Updatable, true},
		{"CreateOp", res.CreateOp, "CreateOrUpdate"},
		{"CreateLRO", res.CreateLRO, true},
		{"CreateModel", res.CreateModel, "OpenShiftCluster"},
		{"ReadOp", res.ReadOp, "Get"},
		{"ReadModel", res.ReadModel, "OpenShiftCluster"},
		{"UpdateOp", res.UpdateOp, "Update"},
		{"UpdateLRO", res.UpdateLRO, true},
		{"UpdateModel", res.UpdateModel, "OpenShiftClusterUpdate"},
		{"DeleteOp", res.DeleteOp, "Delete"},
		{"DeleteLRO", res.DeleteLRO, true},
	}
	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}

	wantArgs := []string{"subscriptionId", "config.ResourceGroup", "config.Name"}
	if len(res.IDNewArgs) != len(wantArgs) {
		t.Fatalf("IDNewArgs = %v, want %v", res.IDNewArgs, wantArgs)
	}
	for i := range wantArgs {
		if res.IDNewArgs[i] != wantArgs[i] {
			t.Errorf("IDNewArgs[%d] = %q, want %q", i, res.IDNewArgs[i], wantArgs[i])
		}
	}
}
