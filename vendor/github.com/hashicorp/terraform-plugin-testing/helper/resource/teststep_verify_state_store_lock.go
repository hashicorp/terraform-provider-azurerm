// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sync"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testprovider"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/providerserver"
	"github.com/hashicorp/terraform-plugin-testing/internal/testing/testsdk/resource"
	"github.com/hashicorp/terraform-plugin-testing/internal/teststep"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mitchellh/go-testing-interface"
)

// testStepVerifyStateStoreLock will run a series of Terraform commands with the goal of ensuring that the state store (defined in config):
//   - Supports locking, acquired during `terraform apply`
//   - Prevents clients from acquiring a lock for an already locked state by returning an error message.
//   - Supports unlocking, by releasing a previously locked state after an operation is complete.
//
// This method also indirectly tests that workspaces and reading/writing state work properly, but testStepNewStateStore should be run prior
// to verify that logic.
func testStepVerifyStateStoreLock(ctx context.Context, t testing.T, step TestStep, providers *providerFactories, stateStoreCfg teststep.Config, helper *plugintest.Helper) error {
	t.Helper()

	// ----- Initialize TF working directory with a single resource that will pause during the apply operation until we indicate it can complete.
	pauseProvider, pauseChan := pauseProvider()
	providers.protov6 = providers.protov6.merge(protov6ProviderFactories{"tfplugintesting": pauseProvider})

	pauseCfg := stateStoreCfg.Append(`resource "tfplugintesting_pause" "resource" {}`)

	pauseWorkingDir := helper.RequireNewWorkingDir(ctx, t, "")
	defer pauseWorkingDir.Close()

	err := pauseWorkingDir.SetConfig(ctx, pauseCfg, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	err = runProviderCommand(ctx, t, pauseWorkingDir, providers, func() error {
		return pauseWorkingDir.Init(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error running init: %w", err)
	}

	// ----- Run terraform apply on working directory with pause resource (locking the state)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		err = runProviderCommand(ctx, t, pauseWorkingDir, providers, func() error {
			return pauseWorkingDir.Apply(ctx)
		})
		if err != nil {
			t.Fatalf("Unexpected error when running apply on working directory with pause resource: %s", err.Error())
		}
	}()

	// Attempt to acquire lock with new client
	err = assertClientCannotAcquireLock(ctx, t, step, providers, stateStoreCfg, helper)
	if err != nil {
		return fmt.Errorf("Failed client lock assertion: %w", err)
	}

	// ----- Send message to tfplugintesting resource to let pause working directory successfully complete the apply operation
	expectedID := "id-1234"
	pauseChan <- expectedID

	// Wait for the paused working directory apply to finish
	wg.Wait()

	// ----- Verify the pause resource is in state
	var populatedState *tfjson.State
	err = runProviderCommand(ctx, t, pauseWorkingDir, providers, func() error {
		populatedState, err = pauseWorkingDir.State(ctx)
		return err
	})
	if err != nil {
		return fmt.Errorf("Error retrieving \"default\" state: %w", err)
	}

	checkOutput := statecheck.ExpectKnownValue(
		"tfplugintesting_pause.resource",
		tfjsonpath.New("id"),
		knownvalue.StringExact(expectedID),
	)

	checkResp := statecheck.CheckStateResponse{}
	checkOutput.CheckState(ctx, statecheck.CheckStateRequest{State: populatedState}, &checkResp)

	if checkResp.Error != nil {
		return fmt.Errorf("After writing a test resource instance object to \"default\" and re-reading it, the object has vanished: %w", err)
	}

	// ----- Run terraform apply on original working directory now that state has been unlocked
	err = pauseWorkingDir.SetConfig(ctx, stateStoreCfg, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	err = runProviderCommand(ctx, t, pauseWorkingDir, providers, func() error {
		return pauseWorkingDir.Apply(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error removing fake resource in \"default\" workspace: %w", err)
	}

	// ----- Ensure the resources are deleted from pause working directory state
	err = assertEmptyWorkspace(ctx, t, pauseWorkingDir, providers, "default")
	if err != nil {
		return fmt.Errorf("After applying empty config to \"default\" workspace, the state should be empty: %w", err)
	}

	return nil
}

// assertClientCannotAcquireLock will create a new client working directory, then attempt to apply
// to the "default" workspace (i.e. indicating it could successfully acquire a lock). If the client is able to
// successfully apply, or receives an error message that is not related to acquiring the lock, the assertion will fail.
func assertClientCannotAcquireLock(ctx context.Context, t testing.T, step TestStep, providers *providerFactories, stateStoreCfg teststep.Config, helper *plugintest.Helper) error {
	clientWorkingDir := helper.RequireNewWorkingDir(ctx, t, "")
	defer clientWorkingDir.Close()
	err := clientWorkingDir.SetConfig(ctx, stateStoreCfg, step.ConfigVariables)
	if err != nil {
		return fmt.Errorf("Error setting config: %w", err)
	}

	err = runProviderCommand(ctx, t, clientWorkingDir, providers, func() error {
		return clientWorkingDir.Init(ctx)
	})
	if err != nil {
		return fmt.Errorf("Error running init: %w", err)
	}

	// ----- Get state of client working directory (does not require a lock)
	err = assertEmptyWorkspace(ctx, t, clientWorkingDir, providers, "default")
	if err != nil {
		return fmt.Errorf("After creating a new workspace, the state should be empty: %w", err)
	}

	// ----- Run terraform apply on client working directory, assert that an error occurs (indicating the state is locked)
	err = runProviderCommand(ctx, t, clientWorkingDir, providers, func() error {
		return clientWorkingDir.Apply(ctx)
	})
	if err == nil {
		return errors.New("Expected an error when attempting to apply to locked \"default\" state, but received none")
	}

	// This is the only part of the error diagnostic that is predictable as the rest is controlled by the provider
	// state store implementation. We assume any error to acquire the state lock was because the lock already exists.
	tfCoreLockErr := regexp.MustCompile(`Error acquiring the state lock`)
	if !tfCoreLockErr.MatchString(err.Error()) {
		return fmt.Errorf("Expected lock error when attempting to apply to locked \"default\" state, received different error: %w", err)
	}

	return nil
}

// pauseProvider returns a test provider server with a single resource "tfplugintesting_pause" and
// a channel waiting for a string. The "tfplugintesting_pause" resource will pause the terraform apply operation until a
// string is sent to the provided channel, which will then be stored in the state.
func pauseProvider() (func() (tfprotov6.ProviderServer, error), chan string) {
	waitForID := make(chan string)
	providerServer := providerserver.NewProviderServer(testprovider.Provider{
		Resources: map[string]testprovider.Resource{
			"tfplugintesting_pause": {
				CreateFunc: func(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
					id := <-waitForID

					resp.NewState = tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"id": tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"id": tftypes.NewValue(tftypes.String, id),
						},
					)
				},
				SchemaResponse: &resource.SchemaResponse{
					Schema: &tfprotov6.Schema{
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:     "id",
									Type:     tftypes.String,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	})

	return providerServer, waitForID
}
