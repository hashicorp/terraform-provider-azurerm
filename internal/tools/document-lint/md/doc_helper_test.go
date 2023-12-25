// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package md_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/md"
)

func TestMDPathFor(t *testing.T) {
	cases := [][2]string{
		{
			"azurerm_api_management_api_policy",
			"api_management_api_policy.html.markdown",
		},
		{
			"not_exists",
			"",
		},
	}
	for _, c := range cases {
		got := md.MDPathFor(c[0])
		if !strings.Contains(got, c[1]) {
			t.Fatalf("%s: \nwant: %s,\ngot:  %s", c[0], c[1], got)
		}
	}
}

func TestResourceNameReg(t *testing.T) {
	var titleReg = regexp.MustCompile(`\npage_title:[^\n]*(azurerm_[a-zA-Z0-9_]+)"`)

	subs := titleReg.FindStringSubmatch(`---
subcategory: "AAD B2C"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_aadb2c_directory"
description: |-
  Manages an AAD B2C Directory.
---

# azurerm_aadb2c_directory

Manages an AAD B2C Directory.

## Example Usage`)
	t.Logf("%v", subs)
}
