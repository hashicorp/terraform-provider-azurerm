package azure

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

	// Strip the Azure generated multiline wrapper
	Normalised := strings.ReplaceAll(input, "<<~EOT\r\n", "")
	Normalised = strings.ReplaceAll(Normalised, "EOT", "")

	// strip Windows flavour new-lines
	Normalised = strings.ReplaceAll(Normalised, "\r\n", "")

	// strip Linux flavour nee lines
	Normalised = strings.ReplaceAll(Normalised, "\n", "")

	return utils.String(Normalised), nil
}
