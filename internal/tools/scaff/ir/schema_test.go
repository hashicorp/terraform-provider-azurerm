// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ir

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"
)

// TestLive_SchemaRedHatOpenShift dumps and sanity-checks the resolved schema for
// the RedHatOpenShift cluster resource so the mechanical mapping can be compared
// against the hand-written gold resource.
func TestLive_SchemaRedHatOpenShift(t *testing.T) {
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

	fmt.Println("=== TOP LEVEL ===")
	for _, p := range res.TopLevel {
		fmt.Printf("  %-32s tf=%-10s go=%-16s req=%v opt=%v computed=%v block=%v enum=%v\n",
			p.TFName, p.TFType, p.GoType, p.Required, p.Optional, p.Computed, p.BlockName, p.IsEnum)
	}
	fmt.Println("=== BLOCKS ===")
	for _, b := range res.Blocks {
		fmt.Printf("  %s (from %s): %d fields\n", b.Name, b.SDKModel, len(b.Properties))
		for _, p := range b.Properties {
			fmt.Printf("      %-28s tf=%-10s go=%-14s block=%v enum=%v\n", p.TFName, p.TFType, p.GoType, p.BlockName, p.IsEnum)
		}
	}

	names := map[string]bool{}
	for _, p := range res.TopLevel {
		names[p.TFName] = true
	}
	for _, want := range []string{"name", "location", "resource_group_name", "tags", "cluster_profile", "network_profile", "master_profile", "worker_profiles"} {
		if !names[want] {
			t.Errorf("expected top-level property %q", want)
		}
	}
}
