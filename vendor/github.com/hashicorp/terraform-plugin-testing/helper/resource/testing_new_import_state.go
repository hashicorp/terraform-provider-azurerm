// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-version"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/google/go-cmp/cmp"
	"github.com/mitchellh/go-testing-interface"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func testStepNewImportState(ctx context.Context, t testing.T, helper *plugintest.Helper, testCaseWorkingDir *plugintest.WorkingDir, step TestStep, priorStepCfg teststep.Config, providers *providerFactories, stepNumber int) error {
	t.Helper()

	// step.ImportStateKind implicitly defaults to the zero-value (ImportCommandWithID) for backward compatibility
	kind := step.ImportStateKind
	importStatePersist := step.ImportStatePersist

	if err := importStatePreconditions(t, helper, step); err != nil {
		return err
	}

	resourceName := step.ResourceName
	if resourceName == "" {
		t.Fatal("ResourceName is required for an import state test")
	}

	// get state from check sequence
	var state *terraform.State
	var stateJSON *tfjson.State
	var err error

	err = runProviderCommand(ctx, t, testCaseWorkingDir, providers, func() error {
		stateJSON, state, err = getState(ctx, t, testCaseWorkingDir)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	// Determine the ID to import
	var importId string
	switch {
	case step.ImportStateIdFunc != nil:
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateIdFunc for import identifier")

		var err error

		logging.HelperResourceDebug(ctx, "Calling TestStep ImportStateIdFunc")

		importId, err = step.ImportStateIdFunc(state)

		if err != nil {
			t.Fatal(err)
		}

		logging.HelperResourceDebug(ctx, "Called TestStep ImportStateIdFunc")
	case step.ImportStateId != "":
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateId for import identifier")

		importId = step.ImportStateId
	default:
		logging.HelperResourceTrace(ctx, "Using resource identifier for import identifier")

		resource, err := testResource(resourceName, state)
		if err != nil {
			t.Fatal(err)
		}
		importId = resource.Primary.ID
	}

	if step.ImportStateIdPrefix != "" {
		logging.HelperResourceTrace(ctx, "Prepending TestStep ImportStateIdPrefix for import identifier")

		importId = step.ImportStateIdPrefix + importId
	}

	logging.HelperResourceTrace(ctx, fmt.Sprintf("Using import identifier: %s", importId))

	var priorIdentityValues map[string]any

	if kind.plannable() && kind.resourceIdentity() {
		priorIdentityValues = identityValuesFromStateValues(stateJSON.Values, resourceName)
		if len(priorIdentityValues) == 0 {
			return fmt.Errorf("importing resource %s: expected prior state to have resource identity values, got none", resourceName)
		}
	}

	testStepConfigRequest := config.TestStepConfigRequest{
		StepNumber: stepNumber,
		TestName:   t.Name(),
	}
	testStepConfig := teststep.Configuration(teststep.PrepareConfigurationRequest{
		Directory:             step.ConfigDirectory,
		File:                  step.ConfigFile,
		Raw:                   step.Config,
		TestStepConfigRequest: testStepConfigRequest,
	}.Exec())

	// If the current import state test step doesn't have configuration, use the prior test step config
	if testStepConfig == nil {
		if priorStepCfg == nil {
			t.Fatal("Cannot import state with no specified config")
		}

		logging.HelperResourceTrace(ctx, "Using prior TestStep Config for import")

		testStepConfig = priorStepCfg
	}

	switch {
	case step.ImportStateConfigExact:
		break

	case kind.plannable() && kind.resourceIdentity():
		testStepConfig = appendImportBlockWithIdentity(testStepConfig, resourceName, priorIdentityValues)

	case kind.plannable():
		testStepConfig = appendImportBlock(testStepConfig, resourceName, importId)
	}

	var workingDir *plugintest.WorkingDir
	if importStatePersist {
		workingDir = testCaseWorkingDir
	} else {
		workingDir = helper.RequireNewWorkingDir(ctx, t, "")
		defer workingDir.Close()
	}

	err = workingDir.SetConfig(ctx, testStepConfig, step.ConfigVariables)
	if err != nil {
		t.Fatalf("Error setting test config: %s", err)
	}

	if kind.plannable() {
		if stepNumber > 1 {
			err = workingDir.CopyState(ctx, testCaseWorkingDir.StateFilePath())
			if err != nil {
				t.Fatalf("copying state: %s", err)
			}

			err = runProviderCommand(ctx, t, workingDir, providers, func() error {
				return workingDir.RemoveResource(ctx, resourceName)
			})
			if err != nil {
				t.Fatalf("removing resource %s from copied state: %s", resourceName, err)
			}
		}
	}

	if !importStatePersist {
		err = runProviderCommand(ctx, t, workingDir, providers, func() error {
			return workingDir.Init(ctx)
		})
		if err != nil {
			t.Fatalf("Error running init: %s", err)
		}
	}

	if kind.plannable() {
		return testImportBlock(ctx, t, workingDir, providers, resourceName, step, priorIdentityValues)
	} else {
		return testImportCommand(ctx, t, workingDir, providers, resourceName, importId, step, state)
	}
}

func testImportBlock(ctx context.Context, t testing.T, workingDir *plugintest.WorkingDir, providers *providerFactories, resourceName string, step TestStep, priorIdentityValues map[string]any) error {
	kind := step.ImportStateKind

	err := runProviderCommandCreatePlan(ctx, t, workingDir, providers)
	if err != nil {
		return fmt.Errorf("generating plan with import config: %s", err)
	}

	plan, err := runProviderCommandSavedPlan(ctx, t, workingDir, providers)
	if err != nil {
		return fmt.Errorf("reading generated plan with import config: %s", err)
	}

	logging.HelperResourceDebug(ctx, fmt.Sprintf("ImportBlockWithId: %d resource changes", len(plan.ResourceChanges)))

	// Verify reasonable things about the plan
	var resourceChangeUnderTest *tfjson.ResourceChange

	if len(plan.ResourceChanges) == 0 {
		return fmt.Errorf("importing resource %s: expected a resource change, got no changes", resourceName)
	}

	for _, change := range plan.ResourceChanges {
		if change.Address == resourceName {
			resourceChangeUnderTest = change
		}
	}

	if resourceChangeUnderTest == nil || resourceChangeUnderTest.Change == nil || resourceChangeUnderTest.Change.Actions == nil {
		return fmt.Errorf("importing resource %s: expected a resource change, got no changes", resourceName)
	}

	change := resourceChangeUnderTest.Change
	actions := change.Actions
	importing := change.Importing

	switch {
	case importing == nil:
		return fmt.Errorf("importing resource %s: expected an import operation, got %q action with plan \nstdout:\n\n%s", resourceChangeUnderTest.Address, actions, savedPlanRawStdout(ctx, t, workingDir, providers))
	// By default we want to ensure there isn't a proposed plan after importing, but for some resources this is unavoidable.
	// An example would be importing a resource that cannot read it's entire value back from the remote API.
	case !step.ExpectNonEmptyPlan && !actions.NoOp():
		return fmt.Errorf("importing resource %s: expected a no-op import operation, got %q action with plan \nstdout:\n\n%s", resourceChangeUnderTest.Address, actions, savedPlanRawStdout(ctx, t, workingDir, providers))
	}

	if err := runPlanChecks(ctx, t, plan, step.ImportPlanChecks.PreApply); err != nil {
		return err
	}

	if kind.resourceIdentity() {
		newIdentityValues := identityValuesFromStateValues(plan.PlannedValues, resourceName)
		if !cmp.Equal(priorIdentityValues, newIdentityValues) {
			return fmt.Errorf("importing resource %s: expected identity values %v, got %v", resourceName, priorIdentityValues, newIdentityValues)
		}
	}

	return nil
}

func testImportCommand(ctx context.Context, t testing.T, workingDir *plugintest.WorkingDir, providers *providerFactories, resourceName string, importId string, step TestStep, state *terraform.State) error {
	err := runProviderCommand(ctx, t, workingDir, providers, func() error {
		return workingDir.Import(ctx, resourceName, importId)
	})
	if err != nil {
		return err
	}

	var importState *terraform.State
	err = runProviderCommand(ctx, t, workingDir, providers, func() error {
		_, importState, err = getState(ctx, t, workingDir)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Error getting state: %s", err)
	}

	logging.HelperResourceDebug(ctx, fmt.Sprintf("State after import: %d resources in the root module", len(importState.RootModule().Resources)))

	// Go through the imported state and verify
	if step.ImportStateCheck != nil {
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateCheck")
		runImportStateCheckFunction(ctx, t, importState, step)
	}

	// Verify that all the states match
	if step.ImportStateVerify {
		logging.HelperResourceTrace(ctx, "Using TestStep ImportStateVerify")

		// Ensure that we do not match against data sources as they
		// cannot be imported and are not what we want to verify.
		// Mode is not present in ResourceState so we use the
		// stringified ResourceStateKey for comparison.
		newResources := make(map[string]*terraform.ResourceState)
		for k, v := range importState.RootModule().Resources {
			if !strings.HasPrefix(k, "data.") {
				newResources[k] = v
			}
		}
		oldResources := make(map[string]*terraform.ResourceState)
		for k, v := range state.RootModule().Resources {
			if !strings.HasPrefix(k, "data.") {
				oldResources[k] = v
			}
		}

		identifierAttribute := step.ImportStateVerifyIdentifierAttribute

		if identifierAttribute == "" {
			identifierAttribute = "id"
		}

		for _, r := range newResources {
			rIdentifier, ok := r.Primary.Attributes[identifierAttribute]

			if !ok {
				t.Fatalf("ImportStateVerify: New resource missing identifier attribute %q, ensure attribute value is properly set or use ImportStateVerifyIdentifierAttribute to choose different attribute", identifierAttribute)
			}

			// Find the existing resource
			var oldR *terraform.ResourceState
			for _, r2 := range oldResources {
				if r2.Primary == nil || r2.Type != r.Type || r2.Provider != r.Provider {
					continue
				}

				r2Identifier, ok := r2.Primary.Attributes[identifierAttribute]

				if !ok {
					t.Fatalf("ImportStateVerify: Old resource missing identifier attribute %q, ensure attribute value is properly set or use ImportStateVerifyIdentifierAttribute to choose different attribute", identifierAttribute)
				}

				if r2Identifier == rIdentifier {
					oldR = r2
					break
				}
			}
			if oldR == nil || oldR.Primary == nil {
				t.Fatalf(
					"Failed state verification, resource with ID %s not found",
					rIdentifier)
			}

			// don't add empty flatmapped containers, so we can more easily
			// compare the attributes
			skipEmpty := func(k, v string) bool {
				if strings.HasSuffix(k, ".#") || strings.HasSuffix(k, ".%") {
					if v == "0" {
						return true
					}
				}
				return false
			}

			// Compare their attributes
			actual := make(map[string]string)
			for k, v := range r.Primary.Attributes {
				if skipEmpty(k, v) {
					continue
				}
				actual[k] = v
			}

			expected := make(map[string]string)
			for k, v := range oldR.Primary.Attributes {
				if skipEmpty(k, v) {
					continue
				}
				expected[k] = v
			}

			// Remove fields we're ignoring
			for _, v := range step.ImportStateVerifyIgnore {
				for k := range actual {
					if strings.HasPrefix(k, v) {
						delete(actual, k)
					}
				}
				for k := range expected {
					if strings.HasPrefix(k, v) {
						delete(expected, k)
					}
				}
			}

			// timeouts are only _sometimes_ added to state. To
			// account for this, just don't compare timeouts at
			// all.
			for k := range actual {
				if strings.HasPrefix(k, "timeouts.") {
					delete(actual, k)
				}
				if k == "timeouts" {
					delete(actual, k)
				}
			}
			for k := range expected {
				if strings.HasPrefix(k, "timeouts.") {
					delete(expected, k)
				}
				if k == "timeouts" {
					delete(expected, k)
				}
			}

			if !reflect.DeepEqual(actual, expected) {
				// Determine only the different attributes
				// go-cmp tries to show surrounding identical map key/value for
				// context of differences, which may be confusing.
				for k, v := range expected {
					if av, ok := actual[k]; ok && v == av {
						delete(expected, k)
						delete(actual, k)
					}
				}

				if diff := cmp.Diff(expected, actual); diff != "" {
					return fmt.Errorf("ImportStateVerify attributes not equivalent. Difference is shown below. The - symbol indicates attributes missing after import.\n\n%s", diff)
				}
			}
		}
	}

	return nil
}

func appendImportBlock(config teststep.Config, resourceName string, importID string) teststep.Config {
	return config.Append(
		fmt.Sprintf(``+"\n"+
			`import {`+"\n"+
			`	to = %s`+"\n"+
			`	id = %q`+"\n"+
			`}`,
			resourceName, importID))
}

func appendImportBlockWithIdentity(config teststep.Config, resourceName string, identityValues map[string]any) teststep.Config {
	configBuilder := strings.Builder{}
	configBuilder.WriteString(fmt.Sprintf(``+"\n"+
		`import {`+"\n"+
		`	to = %s`+"\n"+
		`	identity = {`+"\n",
		resourceName))

	for k, v := range identityValues {
		// It's valid for identity attributes to be null, we can just omit it from config
		if v == nil {
			continue
		}

		switch v := v.(type) {
		case bool:
			configBuilder.WriteString(fmt.Sprintf(`		%q = %t`+"\n", k, v))

		case []any:
			var quotedV []string
			for _, v := range v {
				quotedV = append(quotedV, fmt.Sprintf(`%q`, v))
			}
			configBuilder.WriteString(fmt.Sprintf(`		%q = [%s]`+"\n", k, strings.Join(quotedV, ", ")))

		case json.Number:
			configBuilder.WriteString(fmt.Sprintf(`		%q = %s`+"\n", k, v))

		case string:
			configBuilder.WriteString(fmt.Sprintf(`		%q = %q`+"\n", k, v))

		default:
			panic(fmt.Sprintf("unexpected type %T for identity value %q", v, k))
		}
	}

	configBuilder.WriteString(`	}` + "\n")
	configBuilder.WriteString(`}` + "\n")

	return config.Append(configBuilder.String())
}

func importStatePreconditions(t testing.T, helper *plugintest.Helper, step TestStep) error {
	t.Helper()

	kind := step.ImportStateKind
	versionUnderTest := *helper.TerraformVersion().Core()
	resourceIdentityMinimumVersion := version.Must(version.NewVersion("1.12.0"))

	// Instead of calling [t.Fatal], we return an error. This package's unit tests can use [TestStep.ExpectError] to match
	// on the error message. An alternative, [plugintest.TestExpectTFatal], does not have access to logged error messages,
	// so it is open to false positives on this complex code path.
	//
	// Multiple cases may match, so check the most specific cases first
	switch {
	case kind.resourceIdentity() && versionUnderTest.LessThan(resourceIdentityMinimumVersion):
		return fmt.Errorf(
			`ImportState steps using resource identity require Terraform 1.12.0 or later. Either ` +
				`upgrade the Terraform version running the test or add a ` + "`TerraformVersionChecks`" + ` to ` +
				`the test case to skip this test.` + "\n\n" +
				`https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/tfversion-checks#skip-version-checks`)

	case kind.plannable() && versionUnderTest.LessThan(tfversion.Version1_5_0):
		return fmt.Errorf(
			`ImportState steps using plannable import blocks require Terraform 1.5.0 or later. Either ` +
				`upgrade the Terraform version running the test or add a ` + "`TerraformVersionChecks`" + ` to ` +
				`the test case to skip this test.` + "\n\n" +
				`https://developer.hashicorp.com/terraform/plugin/testing/acceptance-tests/tfversion-checks#skip-version-checks`)

	case kind.plannable() && step.ImportStatePersist:
		return fmt.Errorf(`ImportStatePersist is not supported with plannable import blocks`)

	case kind.plannable() && step.ImportStateVerify:
		return fmt.Errorf(`ImportStateVerify is not supported with plannable import blocks`)
	}

	return nil
}

func resourcesFromState(stateValues *tfjson.StateValues) []*tfjson.StateResource {
	if stateValues == nil || stateValues.RootModule == nil {
		return []*tfjson.StateResource{}
	}

	return stateValues.RootModule.Resources
}

func identityValuesFromStateValues(stateValues *tfjson.StateValues, resourceName string) map[string]any {
	var resource *tfjson.StateResource
	resources := resourcesFromState(stateValues)

	for _, r := range resources {
		if r.Address == resourceName {
			resource = r
			break
		}
	}

	if resource == nil || len(resource.IdentityValues) == 0 {
		return map[string]any{}
	}

	return resource.IdentityValues
}

func runImportStateCheckFunction(ctx context.Context, t testing.T, importState *terraform.State, step TestStep) {
	t.Helper()

	var states []*terraform.InstanceState
	for address, r := range importState.RootModule().Resources {
		if strings.HasPrefix(address, "data.") {
			continue
		}

		if r.Primary == nil {
			continue
		}

		is := r.Primary.DeepCopy() //nolint:staticcheck // legacy usage
		is.Ephemeral.Type = r.Type // otherwise the check function cannot see the type
		states = append(states, is)
	}

	logging.HelperResourceTrace(ctx, "Calling TestStep ImportStateCheck")

	if err := step.ImportStateCheck(states); err != nil {
		t.Fatal(err)
	}

	logging.HelperResourceTrace(ctx, "Called TestStep ImportStateCheck")
}

func savedPlanRawStdout(ctx context.Context, t testing.T, wd *plugintest.WorkingDir, providers *providerFactories) string {
	t.Helper()

	var stdout string

	err := runProviderCommand(ctx, t, wd, providers, func() error {
		var err error
		stdout, err = wd.SavedPlanRawStdout(ctx)
		return err
	})

	if err != nil {
		return fmt.Sprintf("error retrieving formatted plan output: %s", err)
	}
	return stdout
}
