// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package azurecli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
)

// CheckAzVersion tries to determine the version of Azure CLI in the path and checks for a compatible version
func CheckAzVersion() error {
	currentVersion, err := getAzVersion()
	if err != nil {
		return err
	}

	actual, err := version.NewVersion(*currentVersion)
	if err != nil {
		return fmt.Errorf("could not parse detected Azure CLI version %q: %+v", *currentVersion, err)
	}

	supported, err := version.NewVersion(MinimumVersion)
	if err != nil {
		return fmt.Errorf("could not parse supported Azure CLI version: %+v", err)
	}

	nextMajor, err := version.NewVersion(NextMajorVersion)
	if err != nil {
		return fmt.Errorf("could not parse next major Azure CLI version: %+v", err)
	}

	if nextMajor.LessThanOrEqual(actual) {
		return fmt.Errorf("unsupported Azure CLI version %q detected, please install a version newer than %s but older than %s", actual, supported, nextMajor)
	}

	if actual.LessThan(supported) {
		return fmt.Errorf("unsupported Azure CLI version %q detected, please install version %s or newer and ensure the `az` command is in your path", actual, supported)
	}

	return nil
}

// ValidateTenantID validates the supplied tenant ID, and tries to determine the default tenant if a valid one is not supplied.
func ValidateTenantID(tenantId string) (bool, error) {
	validTenantId, err := regexp.MatchString("^[a-zA-Z0-9._-]+$", tenantId)
	if err != nil {
		return false, fmt.Errorf("could not parse tenant ID %q: %s", tenantId, err)
	}

	return validTenantId, nil
}

// GetDefaultTenantID tries to determine the default tenant
func GetDefaultTenantID() (*string, error) {
	var account struct {
		TenantID string `json:"tenantId"`
	}
	if err := JSONUnmarshalAzCmd(true, &account, "account", "show"); err != nil {
		return nil, fmt.Errorf("obtaining tenant ID: %s", err)
	}

	return &account.TenantID, nil
}

// GetDefaultSubscriptionID tries to determine the default subscription
func GetDefaultSubscriptionID() (*string, error) {
	var account struct {
		SubscriptionID string `json:"id"`
	}
	err := JSONUnmarshalAzCmd(true, &account, "account", "show")
	if err != nil {
		return nil, fmt.Errorf("obtaining subscription ID: %s", err)
	}

	return &account.SubscriptionID, nil
}

// getAzVersion tries to determine the version of Azure CLI in the path.
func getAzVersion() (*string, error) {
	var cliVersion *struct {
		AzureCli          *string      `json:"azure-cli,omitempty"`
		AzureCliCore      *string      `json:"azure-cli-core,omitempty"`
		AzureCliTelemetry *string      `json:"azure-cli-telemetry,omitempty"`
		Extensions        *interface{} `json:"extensions,omitempty"`
	}
	err := JSONUnmarshalAzCmd(true, &cliVersion, "version")
	if err != nil {
		return nil, fmt.Errorf("could not parse Azure CLI version: %v", err)
	}

	if cliVersion.AzureCli == nil {
		return nil, fmt.Errorf("could not detect Azure CLI version")
	}

	return cliVersion.AzureCli, nil
}

// JSONUnmarshalAzCmd executes an Azure CLI command and unmarshalls the JSON output, optionally retrieving from and
// populating the command result cache, to avoid unnecessary repeated invocations of Azure CLI.
func JSONUnmarshalAzCmd(cacheable bool, i interface{}, arg ...string) error {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	arg = append(arg, "-o=json")
	argstring := strings.Join(arg, " ")

	var result []byte
	if cacheable {
		if cachedResult, ok := cache.Get(argstring); ok {
			result = cachedResult
		}
	}

	if result == nil {
		log.Printf("[DEBUG] az-cli invocation: az %s", argstring)

		cmd := exec.Command("az", arg...)
		cmd.Stderr = &stderr
		cmd.Stdout = &stdout

		if err := cmd.Start(); err != nil {
			err := fmt.Errorf("launching Azure CLI: %+v", err)
			if stdErrStr := stderr.String(); stdErrStr != "" {
				err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
			}
			return err
		}

		if err := cmd.Wait(); err != nil {
			err := fmt.Errorf("running Azure CLI: %+v", err)
			if stdErrStr := stderr.String(); stdErrStr != "" {
				err = fmt.Errorf("%s: %s", err, strings.TrimSpace(stdErrStr))
			}
			return err
		}

		result = stdout.Bytes()

		if cacheable {
			cache.Set(argstring, result)
		}
	}

	if err := json.Unmarshal(result, &i); err != nil {
		return fmt.Errorf("unmarshaling the output of Azure CLI: %v", err)
	}

	return nil
}
