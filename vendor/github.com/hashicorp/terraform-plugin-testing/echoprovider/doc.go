// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

// Package echoprovider contains a protocol v6 Terraform provider that can be used to transfer data from
// provider configuration to state via a managed resource. This is only meant for provider acceptance testing
// of data that cannot be stored in Terraform artifacts (plan/state), such as an ephemeral resource.
//
// Example Usage:
//
//	// Ephemeral resource that is under test
//	ephemeral "examplecloud_thing" "this" {
//		name = "thing-one"
//	}
//
//	provider "echo" {
//		data = ephemeral.examplecloud_thing.this
//	}
//
//	resource "echo" "test" {} // The `echo.test.data` attribute will contain the ephemeral data from `ephemeral.examplecloud_thing.this`
package echoprovider
