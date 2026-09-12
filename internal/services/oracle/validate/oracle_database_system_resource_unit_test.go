package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle"
	"testing"
)

func TestDatabaseSystemResourcePluggableDatabaseNameDiffSuppress(t *testing.T) {
	schema := oracle.DatabaseSystemResource{}.Arguments()["pluggable_database_name"]
	if schema == nil {
		t.Fatalf("expected pluggable_database_name schema to be defined")
	}

	diff := schema.DiffSuppressFunc
	if diff == nil {
		t.Fatalf("expected DiffSuppressFunc to be configured")
	}

	if !diff("pluggable_database_name", "", "example", nil) {
		t.Fatalf("expected diff suppress when old value is empty and new value is set")
	}

	if diff("pluggable_database_name", "old", "new", nil) {
		t.Fatalf("expected diff not to be suppressed when both values are set")
	}
}
