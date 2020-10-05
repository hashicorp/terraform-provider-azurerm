package suppress

import (
	"encoding/xml"
	"io"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	`github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress`
)

// Deprecated: moved to internal and will be removed in 3.0
func XmlDiff(k, old, new string, d *schema.ResourceData) bool {
	return suppress.XmlDiff(k, old, new, d)
}

// This function will extract all XML tokens from a string, but ignoring all white-space tokens
func expandXmlTokensFromString(input string) ([]xml.Token, error) {
	decoder := xml.NewDecoder(strings.NewReader(input))
	tokens := make([]xml.Token, 0)
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if chars, ok := token.(xml.CharData); ok {
			text := string(chars)
			if strings.TrimSpace(text) == "" {
				continue
			}
		}
		tokens = append(tokens, xml.CopyToken(token))
	}
	return tokens, nil
}
