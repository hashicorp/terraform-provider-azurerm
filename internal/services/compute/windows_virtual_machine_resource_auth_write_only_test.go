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
	"maps"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
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

func TestAccWindowsVirtualMachine_authPasswordWriteOnlyImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_windows_virtual_machine", "test")
	r := WindowsVirtualMachineResource{}
	password := "P@ss-" + rand.Text()
	var resourceID, machineID string

	step := func(label string, action plancheck.ResourceActionType) acceptance.TestStep {
		resourceName := data.ResourceType + "." + label
		return acceptance.TestStep{
			Config:          r.authPasswordSourceNamed(data, label, true, 1),
			ConfigVariables: config.Variables{"password": config.StringVariable(password)},
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(resourceName, action),
				},
			},
			Check: acceptance.ComposeTestCheckFunc(
				check.That(resourceName).ExistsInAzure(r),
				resource.TestCheckResourceAttrWith(resourceName, "id", func(value string) error {
					if value == "" {
						return errors.New("VM resource ID is empty")
					}
					if resourceID != "" && resourceID != value {
						return errors.New("replacement changed the test-owned VM resource address")
					}
					resourceID = value
					return nil
				}),
				resource.TestCheckResourceAttr(resourceName, "admin_password_wo_version", "1"),
				windowsVMPasswordSourceState(resourceName, ""),
				windowsVMPasswordMachineID(resourceName, &machineID, action == plancheck.ResourceActionReplace),
			),
		}
	}
	importStep := data.ImportStep("admin_password", "admin_password_wo_version")
	importStep.ConfigVariables = config.Variables{"password": config.StringVariable(password)}

	persistedImport := importStep
	persistedImport.Config = r.authPasswordSource(data, true, 1)
	persistedImport.ImportStatePersist = true
	persistedImport.ImportStateIdFunc = func(_ *terraform.State) (string, error) {
		if resourceID == "" {
			return "", errors.New("test-owned VM resource ID was not recorded")
		}
		return resourceID, nil
	}
	persistedImport.ImportStateCheck = windowsVMPasswordImportedState(&resourceID, &machineID)

	replaceStep := step("test", plancheck.ResourceActionReplace)
	// Retain the setup address until import succeeds so a failed import can still destroy the VM.
	replaceStep.Config += `
removed {
  from = azurerm_windows_virtual_machine.before_import
  lifecycle {
    destroy = false
  }
}
`
	replaceStep.Check = acceptance.ComposeTestCheckFunc(replaceStep.Check, func(state *terraform.State) error {
		if _, ok := state.RootModule().Resources[data.ResourceType+".before_import"]; ok {
			return errors.New("setup VM address was not removed from managed state")
		}
		return nil
	})

	runWindowsVMWriteOnlyPasswordTest(t, data, []string{password}, []acceptance.TestStep{
		step("before_import", plancheck.ResourceActionCreate),
		persistedImport,
		// Reapply the original password and version: only import reset the version.
		replaceStep,
		importStep,
	})
}

func runWindowsVMWriteOnlyPasswordTest(t *testing.T, data acceptance.TestData, passwords []string, steps []acceptance.TestStep) {
	t.Helper()
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skip("acceptance test requires TF_ACC")
	}

	workingDir := windowsVMPasswordTestDirectory(t)

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
			steps[i].ImportStateCheck = windowsVMPasswordImportStateCheck(steps[i].ImportStateCheck, steps[i].ImportStatePersist, checkStateFiles)
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
	// Process-wide temporary-directory overrides require a serial test.
	resource.Test(t, resource.TestCase{
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

func windowsVMPasswordImportStateCheck(importCheck resource.ImportStateCheckFunc, persisted bool, checkStateFiles func(bool) error) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		var importErr error
		if importCheck != nil {
			importErr = importCheck(states)
		}
		return errors.Join(importErr, checkStateFiles(!persisted))
	}
}

func TestWindowsVMPasswordImportStateCheck(t *testing.T) {
	importErr := errors.New("import check failed")
	stateErr := errors.New("state check failed")
	for _, persisted := range []bool{false, true} {
		for name, test := range map[string]struct {
			customCheck bool
			importError error
			stateError  error
		}{
			"privacy only": {stateError: stateErr},
			"valid":        {customCheck: true},
			"import error": {customCheck: true, importError: importErr},
			"state error":  {customCheck: true, stateError: stateErr},
			"both errors":  {customCheck: true, importError: importErr, stateError: stateErr},
		} {
			t.Run(fmt.Sprintf("persisted=%t/%s", persisted, name), func(t *testing.T) {
				states := []*terraform.InstanceState{{ID: "test-owned-resource"}}
				importCalled, stateCalled := false, false
				var importCheck resource.ImportStateCheckFunc
				if test.customCheck {
					importCheck = func(actual []*terraform.InstanceState) error {
						importCalled = true
						if len(actual) != 1 || actual[0] != states[0] {
							t.Error("import check did not receive the imported state")
						}
						return test.importError
					}
				}
				err := windowsVMPasswordImportStateCheck(importCheck, persisted, func(separate bool) error {
					stateCalled = true
					if separate == persisted {
						t.Error("incorrect separate import state requirement")
					}
					return test.stateError
				})(states)
				if importCalled != test.customCheck || !stateCalled {
					t.Error("import or privacy check was not executed")
				}
				if errors.Is(err, importErr) != (test.importError != nil) || errors.Is(err, stateErr) != (test.stateError != nil) {
					t.Error("import or privacy check failure was lost")
				}
				if (err == nil) != (test.importError == nil && test.stateError == nil) {
					t.Error("unexpected import check error status")
				}
			})
		}
	}
}

func windowsVMPasswordTestDirectory(t *testing.T) string {
	t.Helper()
	checkout, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	workingDir, err := os.MkdirTemp(checkout, ".acctest-password-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(workingDir); err != nil {
			t.Error(err)
		}
	})
	for _, name := range []string{"TMPDIR", "TMP", "TEMP", "GOTMPDIR", "TF_ACC_TEMP_DIR"} {
		t.Setenv(name, workingDir)
	}
	if runtime.GOOS != "windows" {
		// Unix socket paths are length-limited; keep only sockets outside the checkout.
		socketDir, err := os.MkdirTemp("/tmp", "acctest-plugin-")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.RemoveAll(socketDir); err != nil {
				t.Error(err)
			}
		})
		t.Setenv("PLUGIN_UNIX_SOCKET_DIR", socketDir)
	}
	return workingDir
}

func TestWindowsVMPasswordTestDirectory(t *testing.T) {
	t.Setenv("PLUGIN_UNIX_SOCKET_DIR", "inherited-socket-directory")
	previous := make(map[string]string)
	wasSet := make(map[string]bool)
	for _, name := range []string{"TMPDIR", "TMP", "TEMP", "GOTMPDIR", "TF_ACC_TEMP_DIR", "PLUGIN_UNIX_SOCKET_DIR"} {
		previous[name], wasSet[name] = os.LookupEnv(name)
	}
	var directory, socketDirectory string
	t.Run("contained", func(t *testing.T) {
		directory = windowsVMPasswordTestDirectory(t)
		checkout, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Dir(directory) != checkout {
			t.Fatal("password test directory must be inside the checkout")
		}
		info, err := os.Stat(directory)
		if err != nil {
			t.Fatal(err)
		}
		if runtime.GOOS != "windows" && info.Mode().Perm() != 0o700 {
			t.Fatal("password test directory must be private")
		}
		for _, name := range []string{"TMPDIR", "TMP", "TEMP", "GOTMPDIR", "TF_ACC_TEMP_DIR"} {
			if os.Getenv(name) != directory {
				t.Errorf("%s must point to the private password test directory", name)
			}
		}
		socketDirectory = os.Getenv("PLUGIN_UNIX_SOCKET_DIR")
		if runtime.GOOS == "windows" {
			if socketDirectory != previous["PLUGIN_UNIX_SOCKET_DIR"] {
				t.Fatal("Windows TCP plugins do not need a socket directory override")
			}
		} else {
			info, err := os.Stat(socketDirectory)
			if err != nil {
				t.Fatal(err)
			}
			if !info.IsDir() || info.Mode().Perm() != 0o700 {
				t.Fatal("plugin socket directory must be private")
			}
		}
		file, err := os.CreateTemp("", "containment-check-")
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := os.Remove(file.Name()); err != nil {
				t.Error(err)
			}
		})
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		relative, err := filepath.Rel(directory, file.Name())
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			t.Fatal("temporary file escaped the private password test directory")
		}
	})
	if _, err := os.Stat(directory); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("password test directory was not removed")
	}
	if runtime.GOOS != "windows" {
		if _, err := os.Stat(socketDirectory); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("plugin socket directory was not removed")
		}
	}
	for name, expected := range previous {
		value, set := os.LookupEnv(name)
		if value != expected || set != wasSet[name] {
			t.Errorf("%s was not restored after the password test", name)
		}
	}
}

func TestWindowsVMPasswordTestDirectoryProviderStartup(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("go-plugin uses TCP rather than Unix sockets on Windows")
	}

	checkout := filepath.Join(t.TempDir(), strings.Repeat("checkout-", 16))
	if err := os.Mkdir(checkout, 0o700); err != nil {
		t.Fatal(err)
	}
	t.Chdir(checkout)
	t.Setenv("PLUGIN_PROTOCOL_VERSIONS", "5")
	t.Setenv("PLUGIN_CLIENT_CERT", "")
	t.Setenv("PLUGIN_UNIX_SOCKET_GROUP", "")
	t.Setenv("PLUGIN_UNIX_SOCKET_DIR", "")

	workingDir := windowsVMPasswordTestDirectory(t)
	longSocketPath := filepath.Join(workingDir, "plugin1234567890")
	if len(longSocketPath) <= 108 {
		t.Fatal("control must exceed Linux and macOS Unix socket path limits")
	}
	listenConfig := net.ListenConfig{}
	if listener, err := listenConfig.Listen(t.Context(), "unix", longSocketPath); err == nil {
		if err := listener.Close(); err != nil {
			t.Error(err)
		}
		t.Fatal("control unexpectedly accepted an overlong Unix socket path")
	}

	// Exercise the harness's startup path without Terraform, Azure or passwords.
	provider := schema.NewGRPCProviderServer(&schema.Provider{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	reattach, closed, serveErr := plugin.DebugServe(ctx, &plugin.ServeOpts{
		GRPCProviderFunc:    func() tfprotov5.ProviderServer { return provider },
		NoLogOutputOverride: true,
		Logger:              hclog.NewNullLogger(),
	})
	t.Cleanup(func() {
		cancel()
		select {
		case <-closed:
			if serveErr == nil {
				if _, err := os.Stat(reattach.Addr.String); !errors.Is(err, os.ErrNotExist) {
					t.Error("plugin socket was not removed after shutdown")
				}
			}
		case <-time.After(5 * time.Second):
			t.Error("provider did not shut down")
		}
	})
	if serveErr != nil {
		t.Fatalf("provider startup failed: %s", serveErr)
	}
	if reattach.Addr.Network != "unix" || reattach.ProtocolVersion != 5 {
		t.Fatal("provider did not start a protocol 5 Unix socket server")
	}
	if filepath.Dir(reattach.Addr.String) != os.Getenv("PLUGIN_UNIX_SOCKET_DIR") {
		t.Fatal("provider socket escaped the private socket directory")
	}
	dialer := net.Dialer{Timeout: time.Second}
	connection, err := dialer.DialContext(ctx, "unix", reattach.Addr.String)
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.Close(); err != nil {
		t.Fatal(err)
	}
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

func windowsVMPasswordImportedState(resourceID, machineID *string) resource.ImportStateCheckFunc {
	return func(states []*terraform.InstanceState) error {
		if *resourceID == "" || *machineID == "" {
			return errors.New("test-owned VM identity was not recorded")
		}
		for _, state := range states {
			if state == nil || state.ID != *resourceID || state.Attributes["admin_password_wo_version"] != "0" {
				continue
			}
			if state.Attributes["virtual_machine_id"] != *machineID {
				return errors.New("import changed the Azure VM unique ID")
			}
			if state.Attributes["admin_password"] != "" || state.Attributes["admin_password_wo"] != "" {
				return errors.New("imported VM state contains a password")
			}
			return nil
		}
		return errors.New("test-owned VM with a reset write-only password version was not found in imported state")
	}
}

func TestWindowsVMPasswordImportedState(t *testing.T) {
	const password = "sentinel-<secret>&"
	for name, change := range map[string]map[string]string{
		"imported":          {},
		"imported first":    {},
		"managed version":   {"admin_password_wo_version": "1"},
		"missing version":   {"admin_password_wo_version": ""},
		"different VM":      {"virtual_machine_id": "different-vm"},
		"ordinary password": {"admin_password": password},
		"write-only value":  {"admin_password_wo": password},
		"missing VM":        {"id": "different-resource"},
		"missing identity":  {},
	} {
		t.Run(name, func(t *testing.T) {
			resourceID, machineID := "test-owned-resource", "test-owned-vm"
			attributes := map[string]string{
				"id":                        resourceID,
				"virtual_machine_id":        machineID,
				"admin_password_wo_version": "0",
			}
			maps.Copy(attributes, change)
			if name == "missing identity" {
				machineID = ""
			}
			states := []*terraform.InstanceState{
				nil,
				{ID: "test-owned-sibling"},
				{ID: resourceID, Attributes: map[string]string{"virtual_machine_id": machineID, "admin_password_wo_version": "1"}},
				{ID: attributes["id"], Attributes: attributes},
			}
			if name == "imported first" {
				slices.Reverse(states)
			}
			err := windowsVMPasswordImportedState(&resourceID, &machineID)(states)
			if (err == nil) != (name == "imported" || name == "imported first") {
				t.Fatal("imported password state check returned an unexpected result")
			}
			if err != nil && strings.Contains(err.Error(), password) {
				t.Fatal("import check exposed the sentinel in its diagnostic")
			}
		})
	}
}

// Read the raw state, not terraform show's projection, to include outputs, other
// resources, deposed instances and provider-private data. Never include values in errors.
func windowsVMPasswordAbsentFromState(state []byte, passwords []string) error {
	var value any
	if err := json.Unmarshal(state, &value); err != nil {
		return errors.New("could not decode complete Terraform state")
	}
	var containsPassword func(any) bool
	containsPassword = func(value any) bool {
		switch value := value.(type) {
		case map[string]any:
			for key, nested := range value {
				if containsPassword(key) || containsPassword(nested) {
					return true
				}
			}
		case []any:
			return slices.ContainsFunc(value, containsPassword)
		case string:
			decoded, _ := base64.StdEncoding.DecodeString(value)
			for _, password := range passwords {
				if strings.Contains(value, password) || bytes.Contains(decoded, []byte(password)) {
					return true
				}
			}
			var private any
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
	for name, state := range map[string]any{
		"ordinary attribute": map[string]any{"resources": []any{map[string]any{"instances": []any{map[string]any{"attributes": map[string]any{"admin_password": password}}}}}},
		"other resource":     map[string]any{"resources": []any{map[string]any{"type": "azurerm_key_vault_secret", "instances": []any{map[string]any{"attributes": map[string]any{"value": password}}}}}},
		"output":             map[string]any{"outputs": map[string]any{"secret": map[string]any{"value": password, "sensitive": true}}},
		"deposed":            map[string]any{"resources": []any{map[string]any{"instances": []any{map[string]any{"deposed": "old", "attributes": map[string]any{"secret": password}}}}}},
		"private":            map[string]any{"private": base64.StdEncoding.EncodeToString([]byte(password))},
		"private JSON":       map[string]any{"private": base64.StdEncoding.EncodeToString(private)},
		"map key":            map[string]any{password: true},
		"clean":              map[string]any{"resources": []any{map[string]any{"attributes": map[string]any{"admin_password_wo": nil, "admin_password_wo_version": 1}}}},
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
	return r.authPasswordSourceNamed(data, "test", writeOnly, version)
}

func (r WindowsVirtualMachineResource) authPasswordSourceNamed(data acceptance.TestData, label string, writeOnly bool, version int) string {
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

resource "azurerm_windows_virtual_machine" %q {
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
`, r.template(data), writeOnly, label, password)
}
