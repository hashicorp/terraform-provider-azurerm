package compute

import (
	"bytes"
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SSHKeysSchema(isVirtualMachine bool) *pluginsdk.Schema {
	// the SSH Keys for a Virtual Machine cannot be changed once provisioned:
	// Code="PropertyChangeNotAllowed" Message="Changing property 'linuxConfiguration.ssh.publicKeys' is not allowed."

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		ForceNew: isVirtualMachine,
		Set:      SSHKeySchemaHash,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"public_key": {
					Type:             pluginsdk.TypeString,
					Required:         true,
					ForceNew:         isVirtualMachine,
					ValidateFunc:     validate.SSHKey,
					DiffSuppressFunc: SSHKeyDiffSuppress,
				},

				"username": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     isVirtualMachine,
					ValidateFunc: validation.StringIsNotEmpty,
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
	r := regexp.MustCompile("(/home/)+(?P<username>.*?)(/.ssh/authorized_keys)+")

	keys := r.SubexpNames()
	values := r.FindStringSubmatch(input)

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

func SSHKeyDiffSuppress(_, old, new string, _ *pluginsdk.ResourceData) bool {
	oldNormalised, err := NormaliseSSHKey(old)
	if err != nil {
		log.Printf("[DEBUG] error normalising ssh key %q: %+v", old, err)
		return false
	}

	newNormalised, err := NormaliseSSHKey(new)
	if err != nil {
		log.Printf("[DEBUG] error normalising ssh key %q: %+v", new, err)
		return false
	}

	if *oldNormalised == *newNormalised {
		return true
	}

	return false
}

func SSHKeySchemaHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		normalisedKey, err := NormaliseSSHKey(m["public_key"].(string))
		if err != nil {
			log.Printf("[DEBUG] error normalising ssh key %q: %+v", m["public_key"].(string), err)
		}
		buf.WriteString(fmt.Sprintf("%s-", *normalisedKey))
		buf.WriteString(fmt.Sprintf("%s", m["username"]))
	}

	return pluginsdk.HashString(buf.String())
}
