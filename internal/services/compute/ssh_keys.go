// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"bytes"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
					DiffSuppressFunc: suppress.SSHKey,
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

func expandSSHKeys(input []interface{}) []virtualmachines.SshPublicKey {
	output := make([]virtualmachines.SshPublicKey, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		username := raw["username"].(string)
		output = append(output, virtualmachines.SshPublicKey{
			KeyData: pointer.To(raw["public_key"].(string)),
			Path:    pointer.To(formatUsernameForAuthorizedKeysPath(username)),
		})
	}

	return output
}

func expandSSHKeysVMSS(input []interface{}) []virtualmachinescalesets.SshPublicKey {
	output := make([]virtualmachinescalesets.SshPublicKey, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		username := raw["username"].(string)
		output = append(output, virtualmachinescalesets.SshPublicKey{
			KeyData: pointer.To(raw["public_key"].(string)),
			Path:    pointer.To(formatUsernameForAuthorizedKeysPath(username)),
		})
	}

	return output
}

func flattenSSHKeys(input *virtualmachines.SshConfiguration) (*[]interface{}, error) {
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
			return nil, fmt.Errorf("parsing username from %q", *v.Path)
		}

		output = append(output, map[string]interface{}{
			"public_key": *v.KeyData,
			"username":   *username,
		})
	}

	return &output, nil
}

func flattenSSHKeysVMSS(input *virtualmachinescalesets.SshConfiguration) (*[]interface{}, error) {
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
			return nil, fmt.Errorf("parsing username from %q", *v.Path)
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

func SSHKeySchemaHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		normalisedKey, err := suppress.NormalizeSSHKey(m["public_key"].(string))
		if err != nil {
			log.Printf("[DEBUG] error normalising ssh key %q: %+v", m["public_key"].(string), err)
		}
		buf.WriteString(fmt.Sprintf("%s-", *normalisedKey))
		buf.WriteString(fmt.Sprintf("%s", m["username"]))
	}

	return pluginsdk.HashString(buf.String())
}
