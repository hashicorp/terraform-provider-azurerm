package compute

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SSHKeysSchema(isVirtualMachine bool) *schema.Schema {
	// the SSH Keys for a Virtual Machine cannot be changed once provisioned:
	// Code="PropertyChangeNotAllowed" Message="Changing property 'linuxConfiguration.ssh.publicKeys' is not allowed."

	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: isVirtualMachine,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"public_key": {
					Type:         schema.TypeString,
					Required:     true,
					ForceNew:     isVirtualMachine,
					ValidateFunc: ValidateSSHKey,
				},

				"username": {
					Type:         schema.TypeString,
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

// TODO - leverage "golang.org/x/crypto/ssh" for key validation?
// ValidateSSHKey performs some basic validation on supplied SSH Keys - Signature and basic key length are evaluated
func ValidateSSHKey(i interface{}, k string) (warnings []string, errors []error) {
	permittedKeyTypes := []string{
		"ssh-rsa",
	}

	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string or whitespace", k)}
	}

	validSig := ""
	for _, t := range permittedKeyTypes {
		if strings.HasPrefix(v, t) {
			validSig = t
		}
	}

	switch validSig {
	case "ssh-rsa":
		strLen := len(strings.TrimSpace(v))
		if strLen < 379 {
			return nil, []error{fmt.Errorf("expected %q to have a key size of at least 2048 bits", k)}
		}
	case "":
		return nil, []error{fmt.Errorf("Bad: Azure currently supports SSH protocol 2 (SSH-2) RSA public-private key pairs with a minimum length of 2048 bits")}
	}

	return warnings, errors
}
