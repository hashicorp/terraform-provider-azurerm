package azurerm

import (
	"testing"

	"github.com/hashicorp/terraform/terraform"
)

func TestAzureRMVirtualMachineMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion       int
		ID                 string
		InputAttributes    map[string]string
		ExpectedAttributes map[string]string
		Meta               interface{}
	}{
		"v0_1": {
			StateVersion: 0,
			ID:           "azurevm-dummy-id",
			InputAttributes: map[string]string{
				"os_profile_windows_config.#":                                      "1",
				"os_profile_windows_config.429474957.additional_unattend_config.#": "0",
				"os_profile_windows_config.429474957.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.429474957.provision_vm_agent":           "false",
				"os_profile_windows_config.429474957.winrm.#":                      "0",
			},
			ExpectedAttributes: map[string]string{
				"os_profile_windows_config.#":                                       "1",
				"os_profile_windows_config.2229351482.additional_unattend_config.#": "0",
				"os_profile_windows_config.2229351482.enable_automatic_upgrades":    "false",
				"os_profile_windows_config.2229351482.provision_vm_agent":           "false",
				"os_profile_windows_config.2229351482.timezone":                     "UTC",
				"os_profile_windows_config.2229351482.winrm.#":                      "0",
			},
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.InputAttributes,
		}
		is, err := resourceAzureRMVirtualMachineMigrateState(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.ExpectedAttributes {
			actual := is.Attributes[k]
			if actual != v {
				t.Fatalf("Bad Virtual Machine Migrate for %q: %q\n\n expected: %q", k, actual, v)
			}
		}
	}
}
