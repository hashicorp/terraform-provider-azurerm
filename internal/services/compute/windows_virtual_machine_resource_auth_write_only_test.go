// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccWindowsVirtualMachine_authPasswordWriteOnlyLifecycle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}
	passwords := []string{"P@ss-" + rand.Text(), "P@ss-" + rand.Text()}
	var machineID string

	step := func(password string, version int, action plancheck.ResourceActionType) acceptance.TestStep {
		return acceptance.TestStep{
			Config:          r.authPasswordSource(data, true, version),
			ConfigVariables: config.Variables{"password": config.StringVariable(password)},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, action),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttr(data.ResourceName, "admin_password_wo_version", fmt.Sprint(version)),
				windowsVMPasswordSourceState(data.ResourceName, ""),
				windowsVMPasswordMachineID(data.ResourceName, &machineID, action == plancheck.ResourceActionReplace),
			),
		}
	}
	importStep := func(password string) acceptance.TestStep {
		step := data.ImportStep("admin_password", "admin_password_wo_version")
		step.ConfigVariables = config.Variables{"password": config.StringVariable(password)}
		return step
	}

	runWindowsVMWriteOnlyPasswordTest(t, data, passwords, []acceptance.TestStep{
		step(passwords[0], 1, plancheck.ResourceActionCreate),
		importStep(passwords[0]),
		// A write-only value alone cannot produce a diff or change the VM.
		step(passwords[1], 1, plancheck.ResourceActionNoop),
		step(passwords[1], 2, plancheck.ResourceActionReplace),
		importStep(passwords[1]),
		// The version itself must trigger replacement, even with the same password.
		step(passwords[1], 3, plancheck.ResourceActionReplace),
		importStep(passwords[1]),
	})
}

func TestAccWindowsVirtualMachine_authPasswordWriteOnlyTransitions(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}
	ordinaryPassword := "P@ss-" + rand.Text()
	writeOnlyPassword := "P@ss-" + rand.Text()
	var machineID string

	step := func(writeOnly bool, password string, action plancheck.ResourceActionType) acceptance.TestStep {
		expectedPassword := ordinaryPassword
		if writeOnly {
			expectedPassword = ""
		}
		return acceptance.TestStep{
			Config:          r.authPasswordSource(data, writeOnly, 1),
			ConfigVariables: config.Variables{"password": config.StringVariable(password)},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(data.ResourceName, action),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				windowsVMPasswordMachineID(data.ResourceName, &machineID, action == plancheck.ResourceActionReplace),
				windowsVMPasswordSourceState(data.ResourceName, expectedPassword),
			),
		}
	}
	importStep := func(password string) acceptance.TestStep {
		step := data.ImportStep("admin_password", "admin_password_wo_version")
		step.ConfigVariables = config.Variables{"password": config.StringVariable(password)}
		return step
	}

	runWindowsVMWriteOnlyPasswordTest(t, data, []string{writeOnlyPassword}, []acceptance.TestStep{
		step(false, ordinaryPassword, plancheck.ResourceActionCreate),
		importStep(ordinaryPassword),
		step(true, writeOnlyPassword, plancheck.ResourceActionReplace),
		importStep(writeOnlyPassword),
		step(false, ordinaryPassword, plancheck.ResourceActionReplace),
		importStep(ordinaryPassword),
	})
}

func runWindowsVMWriteOnlyPasswordTest(t *testing.T, data acceptance.TestData, passwords []string, steps []acceptance.TestStep) {
	t.Helper()
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skip("acceptance test requires TF_ACC")
	}

	workingDir, err := filepath.Abs(".acctest-password-" + data.RandomString)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(workingDir, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(workingDir); err != nil {
			t.Error(err)
		}
	})

	checkStateFiles := func(importing bool) error {
		paths, err := filepath.Glob(filepath.Join(workingDir, "work*", "terraform.tfstate"))
		if err != nil {
			return err
		}
		minimum := 1
		if importing {
			minimum = 2 // Both the managed state and the separate import state must be checked.
		}
		if len(paths) < minimum {
			return errors.New("expected complete Terraform state files were not found")
		}
		for _, path := range paths {
			for _, filename := range []string{path, path + ".backup"} {
				state, err := os.ReadFile(filename)
				if errors.Is(err, os.ErrNotExist) && strings.HasSuffix(filename, ".backup") {
					continue
				}
				if err != nil {
					return err
				}
				if err := windowsVMPasswordAbsentFromState(state, passwords); err != nil {
					return err
				}
			}
		}
		return nil
	}

	var verifiedSteps []acceptance.TestStep
	for i := range steps {
		if steps[i].ImportState {
			steps[i].ImportStateCheck = func(_ []*terraform.InstanceState) error {
				return checkStateFiles(true)
			}
		} else {
			steps[i].Check = acceptance.ComposeTestCheckFunc(steps[i].Check, func(_ *terraform.State) error {
				return checkStateFiles(false)
			})
			steps[i].PostApplyFunc = func() {
				if err := checkStateFiles(false); err != nil {
					t.Error(err)
				}
			}
		}
		verifiedSteps = append(verifiedSteps, steps[i])
		if !steps[i].ImportState {
			// A refresh plan does not persist refreshed state; check an explicit refresh too.
			verifiedSteps = append(verifiedSteps, acceptance.TestStep{
				RefreshState: true,
				Check: func(_ *terraform.State) error {
					return checkStateFiles(false)
				},
			})
		}
	}

	r := WindowsVirtualMachineResource{}
	testclient.RegisterTestT(t)
	defer testclient.UnregisterTestT()

	// The acceptance wrapper replaces PreApply checks and does not expose WorkingDir.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_11_0)},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInitWithTestName(context.Background(), t.Name(), "azurerm"),
		WorkingDir:               workingDir,
		CheckDestroy: func(state *terraform.State) error {
			client, err := testclient.BuildWithTestName(t.Name())
			if err != nil {
				return err
			}
			return errors.Join(checkStateFiles(false), helpers.CheckDestroyedFunc(client, r, data.ResourceType, data.ResourceName)(state))
		},
		Steps: verifiedSteps,
	})
}

func windowsVMPasswordSourceState(resourceName, expected string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		vm, ok := state.RootModule().Resources[resourceName]
		if !ok || vm.Primary == nil {
			return errors.New("VM password state was not found")
		}
		if vm.Primary.Attributes["admin_password"] != expected {
			return errors.New("ordinary password state does not match the configured password source")
		}
		return nil
	}
}

func windowsVMPasswordMachineID(resourceName string, previous *string, replacement bool) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrWith(resourceName, "virtual_machine_id", func(value string) error {
		if value == "" {
			return errors.New("Azure VM unique ID is empty")
		}
		if *previous != "" && (*previous != value) != replacement {
			return errors.New("Azure VM unique ID did not match the expected replacement behaviour")
		}
		*previous = value
		return nil
	})
}

// Read the raw state, not terraform show's projection, to include outputs, other
// resources, deposed instances and provider-private data. Never include values in errors.
func windowsVMPasswordAbsentFromState(state []byte, passwords []string) error {
	var value interface{}
	if err := json.Unmarshal(state, &value); err != nil {
		return errors.New("could not decode complete Terraform state")
	}
	var containsPassword func(interface{}) bool
	containsPassword = func(value interface{}) bool {
		switch value := value.(type) {
		case map[string]interface{}:
			for key, nested := range value {
				if containsPassword(key) || containsPassword(nested) {
					return true
				}
			}
		case []interface{}:
			return slices.ContainsFunc(value, containsPassword)
		case string:
			decoded, _ := base64.StdEncoding.DecodeString(value)
			for _, password := range passwords {
				if strings.Contains(value, password) || bytes.Contains(decoded, []byte(password)) {
					return true
				}
			}
			var private interface{}
			if json.Unmarshal(decoded, &private) == nil && containsPassword(private) {
				return true
			}
		}
		return false
	}
	if containsPassword(value) {
		return errors.New("write-only password found in serialized Terraform state")
	}
	return nil
}

func TestWindowsVMPasswordAbsentFromState(t *testing.T) {
	const password = "sentinel-<secret>&"
	private, err := json.Marshal(map[string]string{"password": password})
	if err != nil {
		t.Fatal(err)
	}
	for name, state := range map[string]interface{}{
		"ordinary attribute": map[string]interface{}{"resources": []interface{}{map[string]interface{}{"instances": []interface{}{map[string]interface{}{"attributes": map[string]interface{}{"admin_password": password}}}}}},
		"other resource":     map[string]interface{}{"resources": []interface{}{map[string]interface{}{"type": "azurerm_key_vault_secret", "instances": []interface{}{map[string]interface{}{"attributes": map[string]interface{}{"value": password}}}}}},
		"output":             map[string]interface{}{"outputs": map[string]interface{}{"secret": map[string]interface{}{"value": password, "sensitive": true}}},
		"deposed":            map[string]interface{}{"resources": []interface{}{map[string]interface{}{"instances": []interface{}{map[string]interface{}{"deposed": "old", "attributes": map[string]interface{}{"secret": password}}}}}},
		"private":            map[string]interface{}{"private": base64.StdEncoding.EncodeToString([]byte(password))},
		"private JSON":       map[string]interface{}{"private": base64.StdEncoding.EncodeToString(private)},
		"map key":            map[string]interface{}{password: true},
		"clean":              map[string]interface{}{"resources": []interface{}{map[string]interface{}{"attributes": map[string]interface{}{"admin_password_wo": nil, "admin_password_wo_version": 1}}}},
	} {
		t.Run(name, func(t *testing.T) {
			raw, err := json.Marshal(state)
			if err != nil {
				t.Fatal(err)
			}
			err = windowsVMPasswordAbsentFromState(raw, []string{password})
			if (err == nil) != (name == "clean") {
				t.Fatal("whole-state password absence check returned an unexpected result")
			}
			if err != nil && strings.Contains(err.Error(), password) {
				t.Fatal("state check exposed the sentinel in its diagnostic")
			}
		})
	}
	t.Run("malformed state", func(t *testing.T) {
		if err := windowsVMPasswordAbsentFromState([]byte("{"), []string{password}); err == nil {
			t.Fatal("malformed state must fail closed")
		}
	})
}

func (r WindowsVirtualMachineResource) authPasswordSource(data acceptance.TestData, writeOnly bool, version int) string {
	password := "admin_password = var.password"
	if writeOnly {
		password = fmt.Sprintf("admin_password_wo = var.password\n  admin_password_wo_version = %d", version)
	}
	return fmt.Sprintf(`
%s

variable "password" {
  type      = string
  sensitive = true
  ephemeral = %t
}

resource "azurerm_windows_virtual_machine" "test" {
  name                = local.vm_name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  %s
  network_interface_ids = [
    azurerm_network_interface.test.id,
  ]

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter"
    version   = "latest"
  }
}
`, r.template(data), writeOnly, password)
}
