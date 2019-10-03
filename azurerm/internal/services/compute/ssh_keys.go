package compute

import (
	"fmt"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SSHKeysSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"public_key": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.NoZeroValues,
				},

				"username": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.NoZeroValues,
				},
			},
		},
	}
}

func ExpandSSHKeys(input []interface{}) []compute.SSHPublicKey {
	output := make([]compute.SSHPublicKey, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		username := raw["username"].(string)
		output = append(output, compute.SSHPublicKey{
			KeyData: utils.String(raw["public_key"].(string)),
			Path:    utils.String(formatUsernameForAuthorizedKeysPath(username)),
		})
	}

	return output
}

func FlattenSSHKeys(input *compute.SSHConfiguration) (*[]interface{}, error) {
	if input == nil || input.PublicKeys == nil {
		return &[]interface{}{}, nil
	}

	output := make([]interface{}, 0)
	for _, v := range *input.PublicKeys {
		if v.KeyData == nil || v.Path == nil {
			continue
		}

		username := parseUsernameFromAuthorizedKeysPath(*v.Path)
		if username == nil {
			return nil, fmt.Errorf("Error parsing username from %q", *v.Path)
		}

		output = append(output, map[string]interface{}{
			"public_key": *v.KeyData,
			"username":   *username,
		})
	}

	return &output, nil
}

// formatUsernameForAuthorizedKeysPath returns the path to the authorized keys file
// for the specified username
func formatUsernameForAuthorizedKeysPath(username string) string {
	return fmt.Sprintf("/home/%s/.ssh/authorized_keys", username)
}

// parseUsernameFromAuthorizedKeysPath parses the username out of the authorized keys
// path returned from the Azure API
func parseUsernameFromAuthorizedKeysPath(input string) *string {
	// the Azure VM agent hard-codes this to `/home/username/.ssh/authorized_keys`
	// as such we can hard-code this for a better UX
	compiled, err := regexp.Compile("(/home/)+(?P<username>.*?)(/.ssh/authorized_keys)+")
	if err != nil {
		return nil
	}

	keys := compiled.SubexpNames()
	values := compiled.FindStringSubmatch(input)

	if values == nil {
		return nil
	}

	for i, k := range keys {
		if k == "username" {
			value := values[i]
			return &value
		}
	}

	return nil
}
