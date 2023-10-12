// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"testing"
)

func TestPostgresqlDatabaseCollation_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			name:  "Empty",
			input: "",
			valid: false,
		},
		{
			name:  "Invalid Characters",
			input: "en_US%",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := PostgresqlDatabaseCollation(tt.input, "collation"); err == nil {
				t.Errorf("expected an error for value %q but didn't get one", tt.input)
			}
		})
	}
}

func TestPostgresqlDatabaseCollation_Language(t *testing.T) {
	for languageCode := range languageCodes {
		t.Run(fmt.Sprintf("Language Code %q", languageCode), func(t *testing.T) {
			fmtTemplates := []string{
				"%s",
				"%s.utf8",
				"%s.UTF8",
			}
			for _, template := range fmtTemplates {
				value := fmt.Sprintf(template, languageCode)
				if _, err := PostgresqlDatabaseCollation(value, "collation"); err != nil {
					t.Errorf("Expected no error for %q but got %+v", value, err)
				}
			}
		})
	}
}

func TestPostgresqlDatabaseCollation_LanguagePlusNumbers(t *testing.T) {
	validValues := map[string]struct{}{
		"ar_001":  {},
		"en_001":  {},
		"en_029":  {},
		"en_150":  {},
		"eo_001":  {},
		"fr_029":  {},
		"ia_001":  {},
		"la_001":  {},
		"pap_029": {},
	}
	for value := range validValues {
		t.Run(fmt.Sprintf("Value %q", value), func(t *testing.T) {
			fmtTemplates := []string{
				"%s",
				"%s.utf8",
				"%s.UTF8",
			}
			for _, template := range fmtTemplates {
				val := fmt.Sprintf(template, value)
				if _, err := PostgresqlDatabaseCollation(val, "collation"); err != nil {
					t.Errorf("Expected no error for %q but got %+v", val, err)
				}
			}
		})
	}
}

func TestPostgresqlDatabaseCollation_Locales(t *testing.T) {
	for languageCode := range languageCodes {
		t.Run(fmt.Sprintf("Language Code %q", languageCode), func(t *testing.T) {
			fmtTemplates := []string{
				"%s",
				"%s-en.utf8",
				"%s-EN.utf8",
				"%s-de.UTF8",
				"%s-DE.UTF8",
			}
			for _, template := range fmtTemplates {
				value := fmt.Sprintf(template, languageCode)
				if _, err := PostgresqlDatabaseCollation(value, "collation"); err != nil {
					t.Errorf("Expected no error for %q but got %+v", value, err)
				}
			}
		})
	}
}

func TestPostgresqlDatabaseCollation_LocalesThirdComponent(t *testing.T) {
	validValues := map[string]struct{}{
		"bm_Latn_ML":     {},
		"ca_ES_valencia": {},
		"chr_Cher_US":    {},
		"ff_Latn_SN":     {},
		"ha_Latn_GH":     {},
		"ha_Latn_NE":     {},
		"ha_Latn_NG":     {},
		"iu_Cans_CA":     {},
		"iu_Latn_CA":     {},
		"ks_Deva_IN":     {},
		"jv_Java_ID":     {},
	}
	for value := range validValues {
		t.Run(fmt.Sprintf("Value %q", value), func(t *testing.T) {
			fmtTemplates := []string{
				"%s",
				"%s.utf8",
				"%s.UTF8",
			}
			for _, template := range fmtTemplates {
				val := fmt.Sprintf(template, value)
				if _, err := PostgresqlDatabaseCollation(val, "collation"); err != nil {
					t.Errorf("Expected no error for %q but got %+v", val, err)
				}
			}
		})
	}
}

func TestPostgresqlDatabaseCollation_SpecialCases(t *testing.T) {
	cases := map[string]struct{}{
		// these are special-cases
		".utf8":                       {},
		"C":                           {},
		"POSIX":                       {},
		"English_United Kingdom.1252": {},
		"English_United States.1252":  {},
		"En-US":                       {},
		"ucs_basic":                   {},
		"default":                     {},
	}
	for value := range cases {
		t.Run(fmt.Sprintf("Value %q", value), func(t *testing.T) {
			if _, err := PostgresqlDatabaseCollation(value, "collation"); err != nil {
				t.Errorf("expected no error for value %q but got: %+v", value, err)
			}
		})
	}
}
