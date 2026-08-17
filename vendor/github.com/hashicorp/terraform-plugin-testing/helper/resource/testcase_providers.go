// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"strings"
)

// providerConfig takes the list of providers in a TestCase and returns a
// config with only empty provider blocks. This is useful for Import, where no
// config is provided, but the providers must be defined.
func (c TestCase) providerConfig(_ context.Context, skipProviderBlock bool) string {
	var providerBlocks, requiredProviderBlocks strings.Builder

	// [BF] The Providers field handling predates the logic being moved to this
	//      method. It's not entirely clear to me at this time why this field
	//      is being used and not the others, but leaving it here just in case
	//      it does have a special purpose that wasn't being unit tested prior.
	for name := range c.Providers {
		fmt.Fprintf(&providerBlocks, "provider %q {}\n", name)
	}

	for name, externalProvider := range c.ExternalProviders {
		if !skipProviderBlock {
			fmt.Fprintf(&providerBlocks, "provider %q {}\n", name)
		}

		if externalProvider.Source == "" && externalProvider.VersionConstraint == "" {
			continue
		}

		fmt.Fprintf(&requiredProviderBlocks, "    %s = {\n", name)

		if externalProvider.Source != "" {
			fmt.Fprintf(&requiredProviderBlocks, "      source = %q\n", externalProvider.Source)
		}

		if externalProvider.VersionConstraint != "" {
			fmt.Fprintf(&requiredProviderBlocks, "      version = %q\n", externalProvider.VersionConstraint)
		}

		requiredProviderBlocks.WriteString("    }\n")
	}

	if requiredProviderBlocks.Len() > 0 {
		return fmt.Sprintf(`
terraform {
  required_providers {
%[1]s
  }
}

%[2]s
`, strings.TrimSuffix(requiredProviderBlocks.String(), "\n"), providerBlocks.String())
	}

	return providerBlocks.String()
}
