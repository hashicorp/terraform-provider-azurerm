package compute

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NormaliseSSHKey attempts to remove invalid formatting and line breaks that can be present in some cases
// when querying the Azure APIs
func NormaliseSSHKey(input string) (*string, error) {
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

	return utils.String(normalised), nil
}
