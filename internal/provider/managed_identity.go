// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// The Managed Identity source detection and header logic in this file references the Go MSAL
// library apps/managedidentity/managedidentity.go
// (github.com/AzureAD/microsoft-authentication-library-for-go).

package provider

import (
	"os"
	"path/filepath"
	"runtime"
)

// managedIdentitySource mirrors the Managed Identity sources detected by the Azure Identity
// libraries. The hosting environment is determined by inspecting a well-known set of environment
// variables so the correct token request can be issued.
type managedIdentitySource string

const (
	managedIdentitySourceDefaultToIMDS managedIdentitySource = "DefaultToIMDS"
	managedIdentitySourceAzureArc      managedIdentitySource = "AzureArc"
	managedIdentitySourceServiceFabric managedIdentitySource = "ServiceFabric"
	managedIdentitySourceCloudShell    managedIdentitySource = "CloudShell"
	managedIdentitySourceAzureML       managedIdentitySource = "AzureML"
	managedIdentitySourceAppService    managedIdentitySource = "AppService"
)

const (
	identityEndpointEnvVar         = "IDENTITY_ENDPOINT"
	identityHeaderEnvVar           = "IDENTITY_HEADER"
	identityServerThumbprintEnvVar = "IDENTITY_SERVER_THUMBPRINT"
	msiEndpointEnvVar              = "MSI_ENDPOINT"
	msiSecretEnvVar                = "MSI_SECRET"
	imdsEndpointEnvVar             = "IMDS_ENDPOINT"
)

// managedIdentitySourceFromEnvironment detects the Managed Identity source available in the current
// environment, replicating the detection logic used by the Azure Identity libraries.
func managedIdentitySourceFromEnvironment() managedIdentitySource {
	identityEndpoint := os.Getenv(identityEndpointEnvVar)
	identityHeader := os.Getenv(identityHeaderEnvVar)
	identityServerThumbprint := os.Getenv(identityServerThumbprintEnvVar)
	msiEndpoint := os.Getenv(msiEndpointEnvVar)
	msiSecret := os.Getenv(msiSecretEnvVar)
	imdsEndpoint := os.Getenv(imdsEndpointEnvVar)

	switch {
	case identityEndpoint != "" && identityHeader != "":
		if identityServerThumbprint != "" {
			return managedIdentitySourceServiceFabric
		}
		return managedIdentitySourceAppService
	case msiEndpoint != "":
		if msiSecret != "" {
			return managedIdentitySourceAzureML
		}
		return managedIdentitySourceCloudShell
	case isAzureArcEnvironment(identityEndpoint, imdsEndpoint):
		return managedIdentitySourceAzureArc
	}

	return managedIdentitySourceDefaultToIMDS
}

func isAzureArcEnvironment(identityEndpoint, imdsEndpoint string) bool {
	if identityEndpoint != "" && imdsEndpoint != "" {
		return true
	}

	himdsFilePath := azureArcHimdsFilePath(runtime.GOOS)
	if himdsFilePath != "" {
		if _, err := os.Stat(himdsFilePath); err == nil {
			return true
		}
	}

	return false
}

func azureArcHimdsFilePath(platform string) string {
	switch platform {
	case "windows":
		return filepath.Join(os.Getenv("ProgramData"), "AzureConnectedMachineAgent", "himds.exe")
	case "linux":
		return "/opt/azcmagent/bin/himds"
	default:
		return ""
	}
}

// ManagedIdentityCustomHeaders returns the HTTP headers required to request a Managed Identity token
// for the current hosting environment.
func ManagedIdentityCustomHeaders() map[string][]string {
	switch managedIdentitySourceFromEnvironment() {
	case managedIdentitySourceAppService:
		return map[string][]string{
			"X-IDENTITY-HEADER": {os.Getenv(identityHeaderEnvVar)},
		}
	case managedIdentitySourceServiceFabric:
		return map[string][]string{
			"Secret": {os.Getenv(identityHeaderEnvVar)},
		}
	}

	return nil
}
