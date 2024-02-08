// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"fmt"
	"strings"
)

// mergedConfig prepends any necessary terraform configuration blocks to the
// TestStep Config.
//
// If there are ExternalProviders configurations in either the TestCase or
// TestStep, the terraform configuration block should be included with the
// step configuration to prevent errors with providers outside the
// registry.terraform.io hostname or outside the hashicorp namespace.
// This is only necessary when using TestStep.Config.
//
// When TestStep.ConfigDirectory is used, the expectation is that the
// Terraform configuration files will specify a terraform configuration
// block and/or provider blocks as necessary.
func (s TestStep) mergedConfig(ctx context.Context, testCase TestCase, configHasTerraformBlock, configHasProviderBlock bool) string {
	var config strings.Builder

	// Prevent issues with existing configurations containing the terraform
	// configuration block.
	if configHasTerraformBlock {
		config.WriteString(s.Config)

		return config.String()
	}

	if testCase.hasProviders(ctx) {
		config.WriteString(testCase.providerConfig(ctx, configHasProviderBlock))
	} else {
		config.WriteString(s.providerConfig(ctx, configHasProviderBlock))
	}

	config.WriteString(s.Config)

	return config.String()
}

// providerConfig takes the list of providers in a TestStep and returns a
// config with only empty provider blocks. This is useful for Import, where no
// config is provided, but the providers must be defined.
func (s TestStep) providerConfig(_ context.Context, skipProviderBlock bool) string {
	var providerBlocks, requiredProviderBlocks strings.Builder

	for name, externalProvider := range s.ExternalProviders {
		if !skipProviderBlock {
			providerBlocks.WriteString(fmt.Sprintf("provider %q {}\n", name))
		}

		if externalProvider.Source == "" && externalProvider.VersionConstraint == "" {
			continue
		}

		requiredProviderBlocks.WriteString(fmt.Sprintf("    %s = {\n", name))

		if externalProvider.Source != "" {
			requiredProviderBlocks.WriteString(fmt.Sprintf("      source = %q\n", externalProvider.Source))
		}

		if externalProvider.VersionConstraint != "" {
			requiredProviderBlocks.WriteString(fmt.Sprintf("      version = %q\n", externalProvider.VersionConstraint))
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
