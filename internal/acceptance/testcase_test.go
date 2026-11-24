// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
	"strings"
	"testing"
)

// unknownArgument and unknownBlock are what a regression test step reports when the previous version does not know a property, copied from real runs.
func unknownArgument(block, name string) string {
	return fmt.Sprintf("\nError: Unsupported argument\n\n  on terraform_plugin_test.tf line 47, in %s:\n  47:   %s = true\n\nAn argument named %q is not expected here.\n", block, name, name)
}

func unknownBlock(block, name string) string {
	return fmt.Sprintf("\nError: Unsupported block type\n\n  on terraform_plugin_test.tf line 48, in %s:\n  48:   %s {\n\nBlocks of type %q are not expected here.\n", block, name, name)
}

func TestOnlyNewProperties(t *testing.T) { // azignore:AZT004 - a unit test, it only reads the local schema and never runs Terraform
	const plan = "Error running pre-apply plan: exit status 1\n"
	const vnet = `resource "azurerm_virtual_network" "test"`

	cases := []struct {
		name string
		msg  string
		want bool
	}{
		{"new argument", plan + unknownArgument(vnet, "advertised_gateway_prefixes"), true},
		{"new block", plan + unknownBlock(vnet, "ip_address_pool"), true},
		{"new nested argument", plan + unknownArgument(`resource "azurerm_storage_blob_inventory_policy" "test"`, "created_within_days"), true},
		{"new data source argument", plan + unknownArgument(`data "azurerm_virtual_network" "test"`, "resource_group_name"), true},
		{"new provider block", plan + unknownBlock(`provider "azurerm"`, "features"), true},
		{"several new properties", plan + unknownArgument(vnet, "advertised_gateway_prefixes") + unknownBlock(vnet, "ip_address_pool"), true},
		{"new argument with the text wrapped", plan + strings.NewReplacer(", in ", ",\nin ", " named ", "\nnamed ").Replace(unknownArgument(vnet, "advertised_gateway_prefixes")), true},
		{"misspelt argument", plan + unknownArgument(vnet, "advertised_gateway_prefixez"), false},
		{"argument only reported, never set", plan + unknownArgument(vnet, "guid"), false},
		{"block given as an argument", plan + unknownArgument(vnet, "ip_address_pool"), false},
		{"new argument and a misspelt one", plan + unknownArgument(vnet, "advertised_gateway_prefixes") + unknownArgument(vnet, "bgp_communityy"), false},
		{"another provider's resource", plan + unknownArgument(`resource "random_string" "test"`, "bogus"), false},
		{"another error", plan + unknownArgument(vnet, "advertised_gateway_prefixes") + "\nError: Missing required argument\n\nThe argument \"location\" is required, but no definition was found.\n", false},
		{"no Terraform error", "exit status 1", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := onlyNewProperties(t, c.msg); got != c.want {
				t.Fatalf("onlyNewProperties() = %t, want %t for:\n%s", got, c.want, c.msg)
			}
		})
	}
}
