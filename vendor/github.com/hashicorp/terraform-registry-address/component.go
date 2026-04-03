// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfaddr

import (
	"fmt"
	"strings"

	svchost "github.com/hashicorp/terraform-svchost"
)

// Component represents a component listed in a Terraform component registry.
// It combines a registry package address with an optional subdirectory path
// to specify the exact location of a component within a registry package.
//
// Components are resolved through a two-step process:
//  1. The registry package address is used to query the registry API
//  2. The registry returns a ComponentSourceRemote with the actual source location
//
// Example component addresses:
//   - registry.terraform.io/hashicorp/aws
//   - registry.terraform.io/hashicorp/aws//modules/vpc
//   - example.com/myorg/mycomponent//subdir/component
type Component struct {
	// Package is the registry package that the target component belongs to.
	// The component installer must translate this into a ComponentSourceRemote
	// using the registry API and then take that underlying address's
	// Package in order to find the actual package location.
	Package ComponentPackage

	// If Subdir is non-empty then it represents a sub-directory within the
	// remote package that the registry address eventually resolves to.
	// This will ultimately become the suffix of the Subdir of the
	// ComponentSourceRemote that the registry address translates to.
	//
	// Subdir uses a normalized forward-slash-based path syntax within the
	// virtual filesystem represented by the final package. It will never
	// include `../` or `./` sequences.
	Subdir string
}

// ParseComponentSource parses a string representation of a component registry
// address and returns a Component struct. This function only accepts component
// registry addresses and will reject any other address type.
//
// Supported address formats:
//   - "namespace/name" (uses default registry host)
//   - "hostname/namespace/name" (explicit registry host)
//   - "hostname/namespace/name//subdirectory" (with subdirectory)
//
// Examples:
//   - "hashicorp/aws" -> registry.terraform.io/hashicorp/aws
//   - "example.com/myorg/component"
//   - "registry.terraform.io/hashicorp/aws//modules/vpc"
//
// The function validates that:
//   - The address has 2 or 3 slash-separated segments
//   - Subdirectories don't escape the package with "../" paths
//   - The hostname (if provided) is not a reserved VCS host
//   - Namespace and component names follow registry naming conventions
func ParseComponentSource(raw string) (Component, error) {
	var err error

	// Extract subdirectory path if present (indicated by "//" separator)
	var subDir string
	raw, subDir = splitPackageSubdir(raw)
	if strings.HasPrefix(subDir, "../") {
		return Component{}, fmt.Errorf("subdirectory path %q leads outside of the component package", subDir)
	}

	// Split the main address into its components
	parts := strings.Split(raw, "/")
	// A valid registry component address has either two or three parts, because the leading hostname part is optional.
	if len(parts) != 2 && len(parts) != 3 {
		return Component{}, fmt.Errorf("a component registry source address must have either two or three slash-separated segments")
	}

	// Default to the standard Terraform registry if no hostname is specified
	host := DefaultModuleRegistryHost
	// Check if hostname segment is present
	if len(parts) == 3 {
		host, err = svchost.ForComparison(parts[0])
		if err != nil {
			// The svchost library doesn't produce very good error messages to
			// return to an end-user, so we'll use some custom ones here.
			switch {
			case strings.Contains(parts[0], "--"):
				// Looks like possibly punycode, which we don't allow here
				// to ensure that source addresses are written readably.
				return Component{}, fmt.Errorf("invalid component registry hostname %q; internationalized domain names must be given as direct unicode characters, not in punycode", parts[0])
			default:
				return Component{}, fmt.Errorf("invalid component registry hostname %q", parts[0])
			}
		}
		if !strings.Contains(host.String(), ".") {
			return Component{}, fmt.Errorf("invalid component registry hostname: must contain at least one dot")
		}
		// Discard the hostname prefix now that we've processed it
		parts = parts[1:]
	}

	ret := Component{
		Package: ComponentPackage{
			Host: host,
		},
		Subdir: subDir,
	}

	// Prevent potential parsing collisions with known VCS hosts
	// These hostnames are reserved for direct VCS installation
	if host == svchost.Hostname("github.com") || host == svchost.Hostname("bitbucket.org") || host == svchost.Hostname("gitlab.com") {
		// NOTE: This may change in the future if we allow VCS installations
		// 	from these hosts to be registered in the component registry.
		// 	For now, we disallow them to avoid confusion.
		return ret, fmt.Errorf("can't use %q as a component registry host, because it's reserved for installing directly from version control repositories", host)
	}

	// Validate and assign the namespace segment
	if ret.Package.Namespace, err = parseModuleRegistryName(parts[0]); err != nil {
		if strings.Contains(parts[0], ".") {
			// Seems like the user omitted one of the latter components in
			// an address with an explicit hostname.
			return ret, fmt.Errorf("source address must have two more components after the hostname: the namespace and the name")
		}
		return ret, fmt.Errorf("invalid namespace %q: %s", parts[0], err)
	}

	// Validate and assign the component name segment
	if ret.Package.Name, err = parseModuleRegistryName(parts[1]); err != nil {
		if strings.Contains(parts[1], "?") {
			// The user was trying to include a query string, probably?
			return ret, fmt.Errorf("component registry addresses may not include a query string portion")
		}
		return ret, fmt.Errorf("invalid component name %q: %s", parts[1], err)
	}

	return ret, nil
}

// String returns a full representation of the component address, including any
// additional components that are typically implied by omission in
// user-written addresses.
//
// We typically use this longer representation in error messages, in case
// the inclusion of normally-omitted components is helpful in debugging
// unexpected behavior.
func (c Component) String() string {
	if c.Subdir != "" {
		return c.Package.String() + "//" + c.Subdir
	}
	return c.Package.String()
}

// ForDisplay is similar to String but instead returns a representation of
// the idiomatic way to write the address in configuration, omitting
// components that are commonly just implied in addresses written by
// users.
//
// We typically use this shorter representation in informational messages,
// such as the note that we're about to start downloading a package.
func (c Component) ForDisplay() string {
	if c.Subdir != "" {
		return c.Package.ForDisplay() + "//" + c.Subdir
	}
	return c.Package.ForDisplay()
}
