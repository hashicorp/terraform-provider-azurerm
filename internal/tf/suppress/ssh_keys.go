// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package suppress

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func SSHKey(_, old, new string, _ *pluginsdk.ResourceData) bool {
	oldNormalized, err := NormalizeSSHKey(old)
	if err != nil {
		log.Printf("[DEBUG] error normalising ssh key %q: %+v", old, err)
		return false
	}

	newNormalized, err := NormalizeSSHKey(new)
	if err != nil {
		log.Printf("[DEBUG] error normalising ssh key %q: %+v", new, err)
		return false
	}

	if *oldNormalized == *newNormalized {
		return true
	}

	return false
}

// NormalizeSSHKey attempts to remove invalid formatting and line breaks that can be present in some cases
// when querying the Azure APIs
func NormalizeSSHKey(input string) (*string, error) {
	if input == "" {
		return nil, fmt.Errorf("empty string supplied")
	}

	output := input
	output = strings.ReplaceAll(output, "<<~EOT", "")
	output = strings.ReplaceAll(output, "EOT", "")
	output = strings.ReplaceAll(output, "\r", "")

	lines := make([]string, 0)
	for _, line := range strings.Split(output, "\n") {
		lines = append(lines, strings.TrimSpace(line))
	}

	normalised := strings.Join(lines, "")

	return pointer.To(normalised), nil
}
