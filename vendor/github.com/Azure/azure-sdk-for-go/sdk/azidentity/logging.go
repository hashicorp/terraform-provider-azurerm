// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

// LogCredential entries contain information about authentication.
// This includes information like the names of environment variables
// used when obtaining credentials and the type of credential used.
const LogCredential azcore.LogClassification = "Credential"

// log environment variables that can be used for credential types
func logEnvVars() {
	if !azcore.Log().Should(LogCredential) {
		return
	}
	// Log available environment variables
	envVars := []string{}
	if envCheck := os.Getenv("AZURE_TENANT_ID"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_TENANT_ID")
	}
	if envCheck := os.Getenv("AZURE_CLIENT_ID"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_CLIENT_ID")
	}
	if envCheck := os.Getenv("AZURE_CLIENT_SECRET"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_CLIENT_SECRET")
	}
	if envCheck := os.Getenv("AZURE_AUTHORITY_HOST"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_AUTHORITY_HOST")
	}
	if envCheck := os.Getenv("AZURE_CLI_PATH"); len(envCheck) > 0 {
		envVars = append(envVars, "AZURE_CLI_PATH")
	}
	if len(envVars) > 0 {
		azcore.Log().Writef(LogCredential, "Azure Identity => Found the following environment variables:\n\t%s", strings.Join(envVars, ", "))
	}
}

func logGetTokenSuccess(cred azcore.TokenCredential, opts azcore.TokenRequestOptions) {
	if !azcore.Log().Should(LogCredential) {
		return
	}
	msg := fmt.Sprintf("Azure Identity => GetToken() result for %T: SUCCESS\n", cred)
	msg += fmt.Sprintf("\tCredential Scopes: [%s]", strings.Join(opts.Scopes, ", "))
	azcore.Log().Write(LogCredential, msg)
}

func logCredentialError(credName string, err error) {
	azcore.Log().Writef(LogCredential, "Azure Identity => ERROR in %s: %s", credName, err.Error())
}

func logMSIEnv(msi msiType) {
	if !azcore.Log().Should(LogCredential) {
		return
	}
	var msg string
	switch msi {
	case msiTypeIMDS:
		msg = "Azure Identity => Managed Identity environment: IMDS"
	case msiTypeAppServiceV20170901, msiTypeCloudShell, msiTypeAppServiceV20190801:
		msg = "Azure Identity => Managed Identity environment: MSI_ENDPOINT"
	case msiTypeUnavailable:
		msg = "Azure Identity => Managed Identity environment: Unavailable"
	default:
		msg = "Azure Identity => Managed Identity environment: Unknown"
	}
	azcore.Log().Write(LogCredential, msg)
}

func addGetTokenFailureLogs(credName string, err error, includeStack bool) {
	if !azcore.Log().Should(LogCredential) {
		return
	}
	stack := ""
	if includeStack {
		// skip the stack trace frames and ourself
		stack = "\n" + runtime.StackTrace(3, azcore.StackFrameCount)
	}
	azcore.Log().Writef(LogCredential, "Azure Identity => ERROR in GetToken() call for %s: %s%s", credName, err.Error(), stack)
}
