// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package providerserver

import (
	"context"
	"fmt"
	"strings"
)

// ServeOpts are options for serving the provider.
type ServeOpts struct {
	// Address is the full address of the provider. Full address form has three
	// parts separated by forward slashes (/): Hostname, namespace, and
	// provider type ("name").
	//
	// For example: registry.terraform.io/hashicorp/random.
	Address string

	// Debug runs the provider in a mode acceptable for debugging and testing
	// processes, such as delve, by managing the process lifecycle. Information
	// needed for Terraform CLI to connect to the provider is output to stdout.
	// os.Interrupt (Ctrl-c) can be used to stop the provider.
	Debug bool

	// ProtocolVersion is the protocol version that should be used when serving
	// the provider. Either protocol version 5 or protocol version 6 can be
	// used. Defaults to protocol version 6.
	//
	// Protocol version 5 has the following functionality limitations, which
	// will raise an error during the GetProviderSchema or other RPCs:
	//
	//     - tfsdk.Attribute cannot use Attributes field (nested attributes).
	//
	ProtocolVersion int
}

// Validate a given provider address. This is only used for the Address field
// to preserve backwards compatibility for the Name field.
//
// This logic is manually implemented over importing
// github.com/hashicorp/terraform-registry-address as its functionality such as
// ParseAndInferProviderSourceString and ParseRawProviderSourceString allow
// shorter address formats, which would then require post-validation anyways.
func (opts ServeOpts) validateAddress(_ context.Context) error {
	addressParts := strings.Split(opts.Address, "/")
	formatErr := fmt.Errorf("expected hostname/namespace/type format, got: %s", opts.Address)

	if len(addressParts) != 3 {
		return formatErr
	}

	if addressParts[0] == "" || addressParts[1] == "" || addressParts[2] == "" {
		return formatErr
	}

	return nil
}

// Validation checks for provider defined ServeOpts.
//
// Current checks which return errors:
//
//   - If Address is not set
//   - Address is a valid full provider address
//   - ProtocolVersion, if set, is 5 or 6
func (opts ServeOpts) validate(ctx context.Context) error {
	if opts.Address == "" {
		return fmt.Errorf("Address must be provided")
	}

	err := opts.validateAddress(ctx)

	if err != nil {
		return fmt.Errorf("unable to validate Address: %w", err)
	}

	switch opts.ProtocolVersion {
	// 0 represents unset, which Serve will use default.
	case 0, 5, 6:
	default:
		return fmt.Errorf("ProtocolVersion, if set, must be 5 or 6")
	}

	return nil
}
