// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfaddr

import (
	"strings"

	svchost "github.com/hashicorp/terraform-svchost"
)

// ComponentPackage represents a symbolic address for a component in a Terraform
// component registry. It serves as an indirection layer that allows the component
// installer to query a registry and translate this symbolic address (along with
// version constraints provided separately) into a concrete physical source location.
//
// ComponentPackage is intentionally distinct from other package address types
// because they serve different purposes: registry package addresses are exclusively
// used for registry queries to discover the actual component package location.
// This separation helps maintainers understand the component installation workflow
// and enables the type system to enforce proper usage patterns.
//
// Example registry addresses:
//   - registry.terraform.io/hashicorp/aws
//   - example.com/myorg/mycomponent
type ComponentPackage struct {
	// Host is the hostname of the component registry (e.g., "registry.terraform.io")
	Host svchost.Hostname

	// Namespace identifies the organization or user that published the component
	Namespace string

	// Name is the component's name within the namespace
	Name string
}

// String returns the full registry address as a human-readable string.
// This method formats the address for display purposes, using the Unicode
// representation of hostnames rather than punycode encoding.
//
// The returned format is: "hostname/namespace/name"
// For example: "registry.terraform.io/hashicorp/aws"
func (s ComponentPackage) String() string {
	// Note: we're using the "display" form of the hostname here because
	// for our service hostnames "for display" means something different:
	// it means to render non-ASCII characters directly as Unicode
	// characters, rather than using the "punycode" representation we
	// use for internal processing, and so the "display" representation
	// is actually what users would write in their configurations.
	return s.Host.ForDisplay() + "/" + s.ForRegistryProtocol()
}

// ForDisplay returns a string representation suitable for display to users,
// omitting the registry hostname when it's the default registry host.
// This is the format users would typically write in their configurations.
//
// For the default registry, returns: "namespace/name"
// For custom registries, returns: "hostname/namespace/name"
func (s ComponentPackage) ForDisplay() string {
	if s.Host == DefaultModuleRegistryHost {
		return s.ForRegistryProtocol()
	}
	return s.Host.ForDisplay() + "/" + s.ForRegistryProtocol()
}

// ForRegistryProtocol returns a string representation of just the namespace,
// name, and target system portions of the address, always omitting the
// registry hostname and the subdirectory portion, if any.
//
// This is primarily intended for generating addresses to send to the
// registry in question via the registry protocol, since the protocol
// skips sending the registry its own hostname as part of identifiers.
//
// The returned format is: "namespace/name"
// For example: "hashicorp/aws"
func (s ComponentPackage) ForRegistryProtocol() string {
	var buf strings.Builder
	buf.WriteString(s.Namespace)
	buf.WriteByte('/')
	buf.WriteString(s.Name)
	return buf.String()
}
