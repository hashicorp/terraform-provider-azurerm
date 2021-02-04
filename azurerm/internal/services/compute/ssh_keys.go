package compute

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-06-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"golang.org/x/crypto/ssh"
)

func SSHKeysSchema(isVirtualMachine bool) *schema.Schema {
	// the SSH Keys for a Virtual Machine cannot be changed once provisioned:
	// Code="PropertyChangeNotAllowed" Message="Changing property 'linuxConfiguration.ssh.publicKeys' is not allowed."

	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: isVirtualMachine,
		Set:      SSHKeySchemaHash,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"public_key": {
					Type:             schema.TypeString,
					Required:         true,
					ForceNew:         isVirtualMachine,
					ValidateFunc:     ValidateSSHKey,
					DiffSuppressFunc: SSHKeyDiffSuppress,
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

// ValidateSSHKey performs some basic validation on supplied SSH Keys - Encoded Signature and Key Size are evaluated
// Will require rework if/when other Key Types are supported
func ValidateSSHKey(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.TrimSpace(v) == "" {
		return nil, []error{fmt.Errorf("expected %q to not be an empty string or whitespace", k)}
	}

	keyParts := strings.Fields(v)
	if len(keyParts) > 1 {
		byteStr, err := base64.StdEncoding.DecodeString(keyParts[1])
		if err != nil {
			return nil, []error{fmt.Errorf("Error decoding %q for public key data", k)}
		}
		pubKey, err := ssh.ParsePublicKey(byteStr)
		if err != nil {
			return nil, []error{fmt.Errorf("Error parsing %q as a public key object", k)}
		}

		if pubKey.Type() != ssh.KeyAlgoRSA {
			return nil, []error{fmt.Errorf("Error - the provided %s SSH key is not supported. Only RSA SSH keys are supported by Azure", pubKey.Type())}
		} else {
			rsaPubKey, ok := pubKey.(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
			if !ok {
				return nil, []error{fmt.Errorf("Error - could not retrieve the RSA public key from the SSH public key")}
			}
			rsaPubKeyBits := rsaPubKey.Size() * 8
			if rsaPubKeyBits < 2048 {
				return nil, []error{fmt.Errorf("Error - the provided RSA SSH key has %d bits. Only ssh-rsa keys with 2048 bits or higher are supported by Azure", rsaPubKeyBits)}
			}
		}
	} else {
		return nil, []error{fmt.Errorf("Error %q is not a complete SSH2 Public Key", k)}
	}

	return warnings, errors
}

func SSHKeyDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
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

	return schema.HashString(buf.String())
}
