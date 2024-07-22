// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func decodeCertificate(clientCertificate string) ([]byte, error) {
	var pfx []byte
	if clientCertificate != "" {
		out := make([]byte, base64.StdEncoding.DecodedLen(len(clientCertificate)))
		n, err := base64.StdEncoding.Decode(out, []byte(clientCertificate))
		if err != nil {
			return pfx, fmt.Errorf("could not decode client certificate data: %v", err)
		}
		pfx = out[:n]
	}
	return pfx, nil
}

func getClientSecret(d *ProviderModel) (*string, error) {
	clientSecret := strings.TrimSpace(getEnvStringIfValueAbsent(d.ClientSecret, "ARM_CLIENT_SECRET"))

	if path := d.ClientSecretFilePath.ValueString(); path != "" {
		fileSecretRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client Secret from file %q: %v", path, err)
		}

		fileSecret := strings.TrimSpace(string(fileSecretRaw))

		if clientSecret != "" && clientSecret != fileSecret {
			return nil, fmt.Errorf("mismatch between supplied Client Secret and supplied Client Secret file contents - please either remove one or ensure they match")
		}

		clientSecret = fileSecret
	}

	return &clientSecret, nil
}

func getOidcToken(d *ProviderModel) (*string, error) {
	idToken := getEnvStringOrDefault(d.OIDCToken, "ARM_OIDC_TOKEN", "")

	if path := getEnvStringOrDefault(d.OIDCTokenFilePath, "ARM_OIDC_TOKEN_FILE_PATH", ""); path != "" {
		fileTokenRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading OIDC Token from file %q: %v", path, err)
		}

		fileToken := strings.TrimSpace(string(fileTokenRaw))

		if idToken != "" && idToken != fileToken {
			return nil, fmt.Errorf("mismatch between supplied OIDC token and supplied OIDC token file contents - please either remove one or ensure they match")
		}

		idToken = fileToken
	}

	if getEnvBoolIfValueAbsent(d.UseAKSWorkloadIdentity, "ARM_USE_AKS_WORKLOAD_IDENTITY") && os.Getenv("AZURE_FEDERATED_TOKEN_FILE") != "" {
		path := os.Getenv("AZURE_FEDERATED_TOKEN_FILE")
		fileTokenRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading OIDC Token from file %q provided by AKS Workload Identity: %v", path, err)
		}

		fileToken := strings.TrimSpace(string(fileTokenRaw))

		if idToken != "" && idToken != fileToken {
			return nil, fmt.Errorf("mismatch between supplied OIDC token and OIDC token file contents provided by AKS Workload Identity - please either remove one, ensure they match, or disable use_aks_workload_identity")
		}

		idToken = fileToken
	}

	return &idToken, nil
}

func getClientId(d *ProviderModel) (*string, error) {
	clientId := getEnvStringOrDefault(d.ClientId, "ARM_CLIENT_ID", "")

	if path := getEnvStringIfValueAbsent(d.ClientIdFilePath, ""); path != "" {
		fileClientIdRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client ID from file %q: %v", path, err)
		}

		fileClientId := strings.TrimSpace(string(fileClientIdRaw))

		if clientId != "" && clientId != fileClientId {
			return nil, fmt.Errorf("mismatch between supplied Client ID and supplied Client ID file contents - please either remove one or ensure they match")
		}

		clientId = fileClientId
	}

	if d.UseAKSWorkloadIdentity.ValueBool() && clientId != "" {
		aksClientId := os.Getenv("ARM_CLIENT_ID")
		if clientId != "" && clientId != aksClientId {
			return nil, fmt.Errorf("mismatch between supplied Client ID and that provided by AKS Workload Identity - please remove, ensure they match, or disable use_aks_workload_identity")
		}
		clientId = aksClientId
	}

	return &clientId, nil
}

// getEnvStringIfValueAbsent takes a Framework StringValue and a corresponding Environment Variable name and returns
// either the string value set in the StringValue if not Null / Unknown _or_ the os.GetEnv() value of the Environment
// Variable provided. If both of these are empty, an empty string "" is returned.
func getEnvStringIfValueAbsent(val types.String, envVar string) string {
	if val.IsNull() || val.IsUnknown() {
		return os.Getenv(envVar)
	}

	return val.ValueString()
}

// getEnvStringIfValueAbsent takes a Framework StringValue and a corresponding Environment Variable name and returns
// either the string value set in the StringValue if not Null / Unknown _or_ the os.GetEnv() value of the Environment
// Variable provided. If both of these are empty, an empty string "" is returned.
func getEnvStringOrDefault(val types.String, envVar string, defaultValue string) string {
	if val.IsNull() || val.IsUnknown() {
		if v := os.Getenv(envVar); v != "" {
			return os.Getenv(envVar)
		}
		return defaultValue
	}

	return val.ValueString()
}

// getEnvBoolIfValueAbsent takes a Framework BoolValue and a corresponding Environment Variable name and returns
// one of the following in priority order:
// 1 - the Boolean value set in the BoolValue if this is not Null / Unknown.
// 2 - the boolean representation of the os.GetEnv() value of the Environment Variable provided (where anything but
// 'true' or '1' is 'false').
// 3 - `false` in all other cases.
func getEnvBoolIfValueAbsent(val types.Bool, envVar string) bool {
	if val.IsNull() || val.IsUnknown() {
		v := os.Getenv(envVar)
		if strings.EqualFold(v, "true") || strings.EqualFold(v, "1") || v == "" {
			return true
		}
	}

	return val.ValueBool()
}

func getEnvBoolOrDefault(val types.Bool, envVar string, def bool) bool {
	if val.IsNull() || val.IsUnknown() {
		v := os.Getenv(envVar)
		if strings.EqualFold(v, "true") || strings.EqualFold(v, "1") || v == "" {
			return true
		} else {
			return def
		}
	}

	return val.ValueBool()
}

// getEnvListOfStringsIfAbsent returns a []string for the types.List, or the contents of the supplied Environment
// Variable `envVar` if set. If the separator is an empty string, then "," will be used as a default.
func getEnvListOfStringsIfAbsent(val types.List, envVar string, separator string) []string {
	result := make([]string, 0)
	if separator == "" {
		separator = ","
	}
	if val.IsNull() || val.IsUnknown() {
		if v := os.Getenv(envVar); v != "" {
			return strings.Split(v, separator)
		}
		return result
	}

	// we can skip the diags here as failing to decode into the result will return an empty list anyway
	val.ElementsAs(context.Background(), &result, false)

	return result
}
