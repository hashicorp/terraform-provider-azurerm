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

type azAccount struct {
	EnvironmentName *string `json:"environmentName"`
	HomeTenantId    *string `json:"homeTenantId"`
	Id              *string `json:"id"`
	Default         *bool   `json:"isDefault"`
	Name            *string `json:"name"`
	State           *string `json:"state"`
	TenantId        *string `json:"tenantId"`

	ManagedByTenants *[]struct {
		TenantId *string `json:"tenantId"`
	} `json:"managedByTenants"`

	User *struct {
		AssignedIdentityInfo *string `json:"assignedIdentityInfo"`
		Name                 *string `json:"name"`
		Type                 *string `json:"type"`
	}
}

type azVersion struct {
	AzureCli          *string      `json:"azure-cli,omitempty"`
	AzureCliCore      *string      `json:"azure-cli-core,omitempty"`
	AzureCliTelemetry *string      `json:"azure-cli-telemetry,omitempty"`
	Extensions        *interface{} `json:"extensions,omitempty"`
}

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
	account, err := getAzAccount()
	if err != nil {
		return nil, fmt.Errorf("obtaining tenant ID: %s", err)
	}
	return account.TenantId, nil
}

// GetDefaultSubscriptionID tries to determine the default subscription
func GetDefaultSubscriptionID() (*string, error) {
	account, err := getAzAccount()
	if err != nil {
		return nil, fmt.Errorf("obtaining subscription ID: %s", err)
	}
	return account.Id, nil
}

// ListAvailableSubscriptionIDs lists the available subscriptions
func ListAvailableSubscriptionIDs() (*[]string, error) {
	accounts, err := listAzAccounts()
	if err != nil {
		return nil, fmt.Errorf("obtaining subscription ID: %s", err)
	}
	subscriptionIds := make([]string, 0)
	if accounts != nil {
		for _, account := range *accounts {
			if account.Id != nil {
				subscriptionIds = append(subscriptionIds, *account.Id)
			}
		}
	}
	return &subscriptionIds, nil
}

// GetAccountName returns the name of the authenticated principal
func GetAccountName() (*string, error) {
	account, err := getAzAccount()
	if err != nil {
		return nil, fmt.Errorf("obtaining account name: %s", err)
	}
	if account.User == nil {
		return nil, fmt.Errorf("account details were nil: %s", err)
	}
	return account.User.Name, nil
}

// GetAccountType returns the account type of the authenticated principal
func GetAccountType() (*string, error) {
	account, err := getAzAccount()
	if err != nil {
		return nil, fmt.Errorf("obtaining account type: %s", err)
	}
	if account.User == nil {
		return nil, fmt.Errorf("account details were nil: %s", err)
	}
	return account.User.Type, nil
}

// getAzAccount returns the output of `az account show`
func getAzAccount() (*azAccount, error) {
	var account azAccount
	if err := JSONUnmarshalAzCmd(true, &account, "account", "show"); err != nil {
		return nil, fmt.Errorf("obtaining account details: %s", err)
	}
	return &account, nil
}

// listAzAccounts returns the output of `az account list`
func listAzAccounts() (*[]azAccount, error) {
	var account []azAccount
	if err := JSONUnmarshalAzCmd(true, &account, "account", "list"); err != nil {
		return nil, fmt.Errorf("obtaining account details: %s", err)
	}
	return &account, nil
}

// getAzVersion tries to determine the version of Azure CLI in the path.
func getAzVersion() (*string, error) {
	var cliVersion azVersion

	if err := JSONUnmarshalAzCmd(true, &cliVersion, "version"); err != nil {
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
