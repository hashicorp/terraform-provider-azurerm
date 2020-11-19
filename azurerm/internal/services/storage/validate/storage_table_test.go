package validate

import (
	"strings"
	"testing"
)

func TestValidateTableName(t *testing.T) {
	validNames := []string{
		"mytable01",
		"mytable",
		"myTable",
		"MYTABLE",
		"tbl",
		strings.Repeat("w", 63),
	}
	for _, v := range validNames {
		_, errors := TableName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Storage Table Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"table",
		"-invalidname1",
		"invalid_name",
		"invalid!",
		"ww",
		strings.Repeat("w", 64),
	}
	for _, v := range invalidNames {
		_, errors := TableName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Storage Table Name", v)
		}
	}
}
