// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
	helperlogging "github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	testing "github.com/mitchellh/go-testing-interface"
)

// InitTestContext registers the terraform-plugin-log/tfsdklog test sink,
// configures the standard library log package, and creates SDK logger
// contexts. The incoming context is expected to be devoid of logging setup.
//
// The standard library log package handling is important as provider code
// under test may be using that package or another logging library outside of
// terraform-plugin-log.
func InitTestContext(ctx context.Context, t testing.T) context.Context {
	helperlogging.SetOutput(t)

	ctx = tfsdklog.ContextWithTestLogging(ctx, t.Name())
	ctx = tfsdklog.NewRootSDKLogger(ctx, tfsdklog.WithLevelFromEnv(EnvTfLogSdk))
	ctx = tfsdklog.NewSubsystem(ctx, SubsystemHelperResource,
		// All calls are through the HelperResource* helper functions
		tfsdklog.WithAdditionalLocationOffset(1),
		tfsdklog.WithLevelFromEnv(EnvTfLogSdkHelperResource),
	)
	ctx = TestNameContext(ctx, t.Name())

	return ctx
}

// TestNameContext adds the current test name to loggers.
func TestNameContext(ctx context.Context, testName string) context.Context {
	ctx = tfsdklog.SubsystemSetField(ctx, SubsystemHelperResource, KeyTestName, testName)

	return ctx
}

// TestStepNumberContext adds the current test step number to loggers.
func TestStepNumberContext(ctx context.Context, stepNumber int) context.Context {
	ctx = tfsdklog.SubsystemSetField(ctx, SubsystemHelperResource, KeyTestStepNumber, stepNumber)

	return ctx
}

// TestTerraformPathContext adds the current test Terraform CLI path to loggers.
func TestTerraformPathContext(ctx context.Context, terraformPath string) context.Context {
	ctx = tfsdklog.SubsystemSetField(ctx, SubsystemHelperResource, KeyTestTerraformPath, terraformPath)

	return ctx
}

// TestWorkingDirectoryContext adds the current test working directory to loggers.
func TestWorkingDirectoryContext(ctx context.Context, workingDirectory string) context.Context {
	ctx = tfsdklog.SubsystemSetField(ctx, SubsystemHelperResource, KeyTestWorkingDirectory, workingDirectory)

	return ctx
}
