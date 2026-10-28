// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func TestCoalesceStorageSizeToGbs(t *testing.T) {
	tests := []struct {
		name string
		gbs  *int64
		tbs  *int64
		want int64
	}{
		{name: "neither unit"},
		{name: "sub TB storage", gbs: pointer.To(int64(40)), want: 40},
		{name: "TB only", tbs: pointer.To(int64(2)), want: 2048},
		{name: "consistent units", gbs: pointer.To(int64(2048)), tbs: pointer.To(int64(2)), want: 2048},
		{name: "GB takes precedence", gbs: pointer.To(int64(100)), tbs: pointer.To(int64(2)), want: 100},
		{name: "explicit zero GB", gbs: pointer.To(int64(0)), tbs: pointer.To(int64(2))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := coalesceStorageSizeToGbs(tt.gbs, tt.tbs); got != tt.want {
				t.Fatalf("got %d GB, want %d GB", got, tt.want)
			}
		})
	}
}

func TestAutonomousDatabaseStorageSchema(t *testing.T) {
	for _, mode := range []string{"false", "true"} {
		t.Run("sixPointOh="+mode, func(t *testing.T) {
			t.Setenv("ARM_SIXPOINTZERO_BETA", mode)
			wrapper := sdk.NewResourceWrapper(AutonomousDatabaseRegularResource{})
			if _, err := wrapper.Resource(); err != nil {
				t.Fatalf("building typed resource: %s", err)
			}
			args := (AutonomousDatabaseRegularResource{}).Arguments()
			gbs := args["data_storage_size_in_gbs"]
			tbs := args["data_storage_size_in_tbs"]
			if mode == "true" {
				if !gbs.Required || gbs.Optional || gbs.Computed || len(gbs.ExactlyOneOf) != 0 || tbs != nil {
					t.Fatal("6.0 must require GB and remove TB and cross-field constraints")
				}
				return
			}
			if tbs == nil || tbs.Deprecated == "" {
				t.Fatal("5.x must retain deprecated TB configuration")
			}
			for _, key := range []string{"data_storage_size_in_gbs", "data_storage_size_in_tbs"} {
				field := args[key]
				if field.Required || !field.Optional || !field.Computed || len(field.ExactlyOneOf) != 2 {
					t.Fatalf("5.x must allow exactly one configured storage unit: %s", key)
				}
			}
		})
	}
}
