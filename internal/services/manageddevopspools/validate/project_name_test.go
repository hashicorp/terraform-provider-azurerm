// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestProjectName(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "Valid project name",
			input:     "MyProject",
			shouldErr: false,
		},
		{
			name:      "Valid project name with numbers",
			input:     "Project123",
			shouldErr: false,
		},
		{
			name:      "Valid project name with hyphens",
			input:     "My-Project",
			shouldErr: false,
		},
		{
			name:      "Valid project name with spaces",
			input:     "My Project",
			shouldErr: false,
		},
		{
			name:      "Empty string",
			input:     "",
			shouldErr: true,
		},
		{
			name:      "Too long (65 characters)",
			input:     "ThisProjectNameIsWayTooLongAndExceedsTheMaximumAllowedLength65",
			shouldErr: true,
		},
		{
			name:      "Starts with underscore",
			input:     "_MyProject",
			shouldErr: true,
		},
		{
			name:      "Starts with period",
			input:     ".MyProject",
			shouldErr: true,
		},
		{
			name:      "Ends with period",
			input:     "MyProject.",
			shouldErr: true,
		},
		{
			name:      "Contains backslash",
			input:     "My\\Project",
			shouldErr: true,
		},
		{
			name:      "Contains forward slash",
			input:     "My/Project",
			shouldErr: true,
		},
		{
			name:      "Contains colon",
			input:     "My:Project",
			shouldErr: true,
		},
		{
			name:      "Contains asterisk",
			input:     "My*Project",
			shouldErr: true,
		},
		{
			name:      "Contains question mark",
			input:     "My?Project",
			shouldErr: true,
		},
		{
			name:      "Contains double quote",
			input:     `My"Project`,
			shouldErr: true,
		},
		{
			name:      "Contains single quote",
			input:     "My'Project",
			shouldErr: true,
		},
		{
			name:      "Contains less than",
			input:     "My<Project",
			shouldErr: true,
		},
		{
			name:      "Contains greater than",
			input:     "My>Project",
			shouldErr: true,
		},
		{
			name:      "Contains semicolon",
			input:     "My;Project",
			shouldErr: true,
		},
		{
			name:      "Contains hash",
			input:     "My#Project",
			shouldErr: true,
		},
		{
			name:      "Contains dollar sign",
			input:     "My$Project",
			shouldErr: true,
		},
		{
			name:      "Contains curly braces",
			input:     "My{Project}",
			shouldErr: true,
		},
		{
			name:      "Contains plus",
			input:     "My+Project",
			shouldErr: true,
		},
		{
			name:      "Contains equals",
			input:     "My=Project",
			shouldErr: true,
		},
		{
			name:      "Contains square brackets",
			input:     "My[Project]",
			shouldErr: true,
		},
		{
			name:      "Contains pipe",
			input:     "My|Project",
			shouldErr: true,
		},
		{
			name:      "Contains comma",
			input:     "My,Project",
			shouldErr: true,
		},
		{
			name:      "Contains control character",
			input:     "My\x00Project",
			shouldErr: true,
		},
		{
			name:      "Reserved name - App_Browsers",
			input:     "App_Browsers",
			shouldErr: true,
		},
		{
			name:      "Reserved name - App_Code",
			input:     "App_Code",
			shouldErr: true,
		},
		{
			name:      "Reserved name - App_Data",
			input:     "App_Data",
			shouldErr: true,
		},
		{
			name:      "Reserved name - App_GlobalResources",
			input:     "App_GlobalResources",
			shouldErr: true,
		},
		{
			name:      "Reserved name - App_LocalResources",
			input:     "App_LocalResources",
			shouldErr: true,
		},
		{
			name:      "Reserved name - App_Themes",
			input:     "App_Themes",
			shouldErr: true,
		},
		{
			name:      "Reserved name - App_WebResources",
			input:     "App_WebResources",
			shouldErr: true,
		},
		{
			name:      "Reserved name - bin",
			input:     "bin",
			shouldErr: true,
		},
		{
			name:      "Reserved name - web.config",
			input:     "web.config",
			shouldErr: true,
		},
		{
			name:      "Reserved name case insensitive - BIN",
			input:     "BIN",
			shouldErr: true,
		},
		{
			name:      "Valid 64 character name",
			input:     "ThisProjectNameIsExactly64CharactersLongAndShouldBeAccepted12",
			shouldErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := ProjectName(tc.input, "test")

			if tc.shouldErr && len(errors) == 0 {
				t.Errorf("Expected validation to fail for input %q, but it passed", tc.input)
			}

			if !tc.shouldErr && len(errors) > 0 {
				t.Errorf("Expected validation to pass for input %q, but it failed with errors: %v", tc.input, errors)
			}
		})
	}
}
