// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package azurecli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/hashicorp/go-version"
)

// CheckAzVersion tries to determine the version of Azure CLI in the path and checks for a compatible version
func CheckAzVersion(minVersion string, nextMajorVersion *string) error {
	var cliVersion *struct {
		AzureCli          *string      `json:"azure-cli,omitempty"`
		AzureCliCore      *string      `json:"azure-cli-core,omitempty"`
		AzureCliTelemetry *string      `json:"azure-cli-telemetry,omitempty"`
		Extensions        *interface{} `json:"extensions,omitempty"`
	}
	err := JSONUnmarshalAzCmd(&cliVersion, "version")
	if err != nil {
		return fmt.Errorf("could not parse Azure CLI version: %v", err)
	}

	if cliVersion.AzureCli == nil {
		return fmt.Errorf("could not detect Azure CLI version. Please ensure you have installed Azure CLI version %s or newer", minVersion)
	}

	actual, err := version.NewVersion(*cliVersion.AzureCli)
	if err != nil {
		return fmt.Errorf("could not parse detected Azure CLI version %q: %+v", *cliVersion.AzureCli, err)
	}

	supported, err := version.NewVersion(minVersion)
	if err != nil {
		return fmt.Errorf("could not parse supported Azure CLI version: %+v", err)
	}

	if nextMajorVersion != nil {
		nextMajor, err := version.NewVersion(*nextMajorVersion)
		if err != nil {
			return fmt.Errorf("could not parse next major Azure CLI version: %+v", err)
		}

		if nextMajor.LessThanOrEqual(actual) {
			return fmt.Errorf("unsupported Azure CLI version %q detected, please install a version newer than %s but older than %s", actual, supported, nextMajor)
		}
	}

	if actual.LessThan(supported) {
		return fmt.Errorf("unsupported Azure CLI version %q detected, please install version %s or newer and ensure the `az` command is in your path", actual, supported)
	}

	return nil
}

// GetDefaultSubscriptionID tries to determine the default subscription
func GetDefaultSubscriptionID() (string, error) {
	var account struct {
		SubscriptionID string `json:"id"`
	}
	err := JSONUnmarshalAzCmd(&account, "account", "show")
	if err != nil {
		return "", fmt.Errorf("obtaining subscription ID: %s", err)
	}

	return account.SubscriptionID, nil
}

// CheckTenantID validates the supplied tenant ID, and tries to determine the default tenant if a valid one is not supplied.
func CheckTenantID(tenantId string) (string, error) {
	validTenantId, err := regexp.MatchString("^[a-zA-Z0-9._-]+$", tenantId)
	if err != nil {
		return "", fmt.Errorf("could not parse tenant ID %q: %s", tenantId, err)
	}

	if !validTenantId {
		var account struct {
			ID       string `json:"id"`
			TenantID string `json:"tenantId"`
		}
		err := JSONUnmarshalAzCmd(&account, "account", "show")
		if err != nil {
			return "", fmt.Errorf("obtaining tenant ID: %s", err)
		}
		tenantId = account.TenantID
	}

	return tenantId, nil
}

// JSONUnmarshalAzCmd executes an Azure CLI command and unmarshalls the JSON output.
func JSONUnmarshalAzCmd(i interface{}, arg ...string) error {
	var stderr bytes.Buffer
	var stdout bytes.Buffer

	arg = append(arg, "-o=json")
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

	if err := json.Unmarshal(stdout.Bytes(), &i); err != nil {
		return fmt.Errorf("unmarshaling the output of Azure CLI: %v", err)
	}

	return nil
}
