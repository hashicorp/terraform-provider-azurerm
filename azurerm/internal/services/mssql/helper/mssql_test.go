package helper

import "testing"

// Your server name can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.
func TestValidateMsSqlServerName(t *testing.T) {
	cases := []struct {
		Value  string
		Errors bool
	}{
		{
			Value:  "",
			Errors: true,
		},
		{
			Value:  "k",
			Errors: false,
		},
		{
			Value:  "K",
			Errors: true,
		},
		{
			Value:  "k-",
			Errors: true,
		},
		{
			Value:  "k-t",
			Errors: false,
		},
		{
			Value:  "K-T",
			Errors: true,
		},
		{
			Value:  "validname",
			Errors: false,
		},
		{
			Value:  "invalid_name",
			Errors: true,
		},
		{
			Value:  "123456789112345678921234567893123456789412345678951234567896123",
			Errors: false,
		},
		{
			Value:  "01234567891123456789212345678931234567894123456789512345678961234",
			Errors: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateMsSqlServerName(tc.Value, "name")

		if len(errors) > 0 != tc.Errors {
			if tc.Errors {
				t.Fatalf("Expected ValidateMsSqlServerName to have errors for '%s', got %d ", tc.Value, len(errors))
			} else {
				t.Fatalf("Expected ValidateMsSqlServerName to not have errors for '%s', got %d ", tc.Value, len(errors))
			}
		}
	}
}

// Your database name can't end with '.' or ' ', can't contain '<,>,*,%,&,:,\,/,?' or control characters, and can't have more than 128 characters.
func TestValidateMsSqlDatabaseName(t *testing.T) {
	cases := []struct {
		Value  string
		Errors bool
	}{
		{
			Value:  "",
			Errors: true,
		},
		{
			Value:  "k",
			Errors: false,
		},
		{
			Value:  "K",
			Errors: false,
		},
		{
			Value:  "space ",
			Errors: true,
		},
		{
			Value:  "dot.",
			Errors: true,
		},
		{
			Value:  "data_base-name",
			Errors: false,
		},
		{
			Value:  "ends.with.dash-",
			Errors: false,
		},
		{
			Value:  "ends.with.underscore_",
			Errors: false,
		},
		{
			Value:  "fail:semicolon",
			Errors: true,
		},
		{
			Value:  "fail?question",
			Errors: true,
		},
		{
			Value:  "fail&ampersand",
			Errors: true,
		},
		{
			Value:  "fail%percent",
			Errors: true,
		},
		{
			Value:  "12345678911234567892123456789312345678941234567895123456789612345678971234567898123456789912345678901234567891123456789212345678",
			Errors: false,
		},
		{
			Value:  "123456789112345678921234567893123456789412345678951234567896123456789712345678981234567899123456789012345678911234567892123456789",
			Errors: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateMsSqlDatabaseName(tc.Value, "name")

		if len(errors) > 0 != tc.Errors {
			if tc.Errors {
				t.Fatalf("Expected ValidateMsSqlDatabaseName to have errors for '%s', got %d ", tc.Value, len(errors))
			} else {
				t.Fatalf("Expected ValidateMsSqlDatabaseName to not have errors for '%s', got %d ", tc.Value, len(errors))
			}
		}
	}
}

// Following characters and any control characters are not allowed for resource name '%,&,\\\\,?,/'.\"
// The name can not end with characters: '. '
func TestValidateMsSqlElasticPoolName(t *testing.T) {
	cases := []struct {
		Value  string
		Errors bool
	}{
		{
			Value:  "",
			Errors: true,
		},
		{
			Value:  "k",
			Errors: false,
		},
		{
			Value:  "K",
			Errors: false,
		},
		{
			Value:  "space ",
			Errors: true,
		},
		{
			Value:  "dot.",
			Errors: true,
		},
		{
			Value:  "data_base-name",
			Errors: false,
		},
		{
			Value:  "ends.with.dash-",
			Errors: false,
		},
		{
			Value:  "ends.with.underscore_",
			Errors: false,
		},
		{
			Value:  "fail?question",
			Errors: true,
		},
		{
			Value:  "fail&ampersand",
			Errors: true,
		},
		{
			Value:  "fail%percent",
			Errors: true,
		},
		{
			Value:  "12345678911234567892123456789312345678941234567895123456789612345678971234567898123456789912345678901234567891123456789212345678",
			Errors: false,
		},
		{
			Value:  "123456789112345678921234567893123456789412345678951234567896123456789712345678981234567899123456789012345678911234567892123456789",
			Errors: true,
		},
	}

	for _, tc := range cases {
		_, errors := ValidateMsSqlElasticPoolName(tc.Value, "name")

		if len(errors) > 0 != tc.Errors {
			if tc.Errors {
				t.Fatalf("Expected ValidateMsSqlServerName to have errors for '%s', got %d ", tc.Value, len(errors))
			} else {
				t.Fatalf("Expected ValidateMsSqlServerName to not have errors for '%s', got %d ", tc.Value, len(errors))
			}
		}
	}
}
