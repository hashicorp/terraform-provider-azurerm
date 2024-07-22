// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func Test_getOidcToken(t *testing.T) {
	expectedString := "testOidc"
	p := &ProviderModel{
		OIDCToken: basetypes.NewStringValue(expectedString),
	}

	result, err := getOidcToken(p)
	if err != nil {
		t.Fatalf("getOidcToken returned unexpected error %v", err)
	}

	if result == nil {
		t.Fatalf("getOidcToken returned nil result without an error")
	}

	if *result != expectedString {
		t.Fatalf("getOidcToken did not return expected string (%s), got %+v", expectedString, *result)
	}
}

func Test_getOidcTokenFromFile(t *testing.T) {
	expectedString := "testOidcFromFile"
	p := &ProviderModel{
		OIDCTokenFilePath: basetypes.NewStringValue("./testdata/oidc_test_input.txt"),
	}

	result, err := getOidcToken(p)
	if err != nil {
		t.Fatalf("getOidcToken returned unexpected error %v", err)
	}

	if result == nil {
		t.Fatalf("getOidcToken returned nil result without an error")
	}

	if *result != expectedString {
		t.Fatalf("getOidcToken did not return expected string `%s`, got `%s`", expectedString, *result)
	}
}

func Test_getOidcTokenExpectMismatch(t *testing.T) {
	configuredString := "testOidc"
	p := &ProviderModel{
		OIDCToken:         basetypes.NewStringValue(configuredString),
		OIDCTokenFilePath: basetypes.NewStringValue("./testdata/oidc_test_input.txt"),
	}

	_, err := getOidcToken(p)
	if err == nil {
		t.Fatal("expected an error but did not get one")
	}

	if !strings.EqualFold(err.Error(), "mismatch between supplied OIDC token and supplied OIDC token file contents - please either remove one or ensure they match") {
		t.Fatal("did not get expected error")
	}
}

func Test_getOidcTokenAKSWorkload(t *testing.T) {
	expectedString := "testOidcFromFile"
	err := os.Setenv("AZURE_FEDERATED_TOKEN_FILE", "./testdata/oidc_test_input.txt")
	if err != nil {
		t.Fatalf("could not set env var (`AZURE_FEDERATED_TOKEN_FILE`) for test: %+v", err)
	}

	p := &ProviderModel{
		UseAKSWorkloadIdentity: basetypes.NewBoolValue(true),
	}

	result, err := getOidcToken(p)
	if err != nil {
		t.Fatalf("getOidcToken returned unexpected error %v", err)
	}

	if *result != expectedString {
		t.Fatalf("getOidcToken did not return expected string `%s`, got `%s`", expectedString, *result)
	}
}

func Test_getOidcTokenAKSWorkloadExpectMismatch(t *testing.T) {
	configuredString := "testOidc"
	err := os.Setenv("AZURE_FEDERATED_TOKEN_FILE", "./testdata/oidc_test_input.txt")
	if err != nil {
		t.Fatalf("could not set env var (`AZURE_FEDERATED_TOKEN_FILE`) for test: %+v", err)
	}

	p := &ProviderModel{
		OIDCToken:              basetypes.NewStringValue(configuredString),
		UseAKSWorkloadIdentity: basetypes.NewBoolValue(true),
	}

	_, err = getOidcToken(p)
	if err == nil {
		t.Fatal("expected an error but did not get one")
	}

	if !strings.EqualFold(err.Error(), "mismatch between supplied OIDC token and OIDC token file contents provided by AKS Workload Identity - please either remove one, ensure they match, or disable use_aks_workload_identity") {
		t.Fatal("did not get expected error")
	}
}

func Test_getClientSecret(t *testing.T) {
	expectedString := "testClientSecret"

	p := &ProviderModel{
		ClientSecret: basetypes.NewStringValue(expectedString),
	}

	result, err := getClientSecret(p)
	if err != nil {
		t.Fatalf("getClientSecret returned unexpected error %v", err)
	}
	if result == nil {
		t.Fatalf("getClientSecret returned nil result without an error")
	}
	if *result != expectedString {
		t.Fatalf("getCLientSecret did not return expected string `%s`, got `%s`", expectedString, *result)
	}
}

func Test_getClientSecretExpectMismatch(t *testing.T) {
	configuredString := "testClientSecret"
	p := &ProviderModel{
		ClientSecret:         basetypes.NewStringValue(configuredString),
		ClientSecretFilePath: basetypes.NewStringValue("./testdata/client_secret_test_input.txt"),
	}

	_, err := getClientSecret(p)
	if err == nil {
		t.Fatalf("expected an error but did not get one")
	}
	if !strings.EqualFold(err.Error(), "mismatch between supplied Client Secret and supplied Client Secret file contents - please either remove one or ensure they match") {
		t.Fatal("did not get expected error")
	}
}

func Test_getClientSecretFromFile(t *testing.T) {
	os.Setenv("ARM_CLIENT_SECRET", "")
	expectedString := "testClientSecretFromFile"
	p := &ProviderModel{
		ClientSecretFilePath: basetypes.NewStringValue("./testdata/client_secret_test_input.txt"),
	}

	result, err := getClientSecret(p)
	if err != nil {
		t.Fatalf("getClientSecretFromFile returned unexpected error %v", err)
	}
	if result == nil {
		t.Fatalf("getClientSecretFromFile returned nil result without an error")
	}
	if *result != expectedString {
		t.Fatalf("getClientSecret did not return expected string `%s`, got `%s`", expectedString, *result)
	}
}

func Test_getClientSecretFromFileMismatch(t *testing.T) {
	os.Setenv("ARM_CLIENT_SECRET", "foo")
	p := &ProviderModel{
		ClientSecretFilePath: basetypes.NewStringValue("./testdata/client_secret_test_input.txt"),
	}

	result, err := getClientSecret(p)
	if err == nil {
		t.Fatalf("expected an error but did not get one")
	}
	if result != nil {
		t.Fatalf("getClientSecretFromFile returned a result with an error")
	}
}

func Test_getClientID(t *testing.T) {
	expectedString := "testClientID"
	p := &ProviderModel{
		ClientId: basetypes.NewStringValue(expectedString),
	}

	result, err := getClientId(p)
	if err != nil {
		t.Fatalf("getClientID returned unexpected error %v", err)
	}
	if result == nil {
		t.Fatalf("getClientID returned nil result without an error")
	}
	if *result != expectedString {
		t.Fatalf("getClientID did not return expected string `%s`, got `%s`", expectedString, *result)
	}
}

func Test_getClientIDFromFileExpectMismatch(t *testing.T) {
	configuredString := "testClientID"
	p := &ProviderModel{
		ClientId:         basetypes.NewStringValue(configuredString),
		ClientIdFilePath: basetypes.NewStringValue("./testdata/client_id_test_input.txt"),
	}
	_, err := getClientId(p)
	if err == nil {
		t.Fatalf("expected an error but did not get one")
	}
	if !strings.EqualFold(err.Error(), "mismatch between supplied Client ID and supplied Client ID file contents - please either remove one or ensure they match") {
		t.Fatal("did not get expected error")
	}
}

func Test_getClientIDFromFile(t *testing.T) {
	expectedString := "testClientIDFromFile"
	p := &ProviderModel{
		ClientIdFilePath: basetypes.NewStringValue("./testdata/client_id_test_input.txt"),
	}

	result, err := getClientId(p)
	if err != nil {
		t.Fatalf("getClientID returned unexpected error %v", err)
	}
	if result == nil {
		t.Fatalf("getClientID returned nil result without an error")
	}
	if *result != expectedString {
		t.Fatalf("getClientID did not return expected string `%s`, got `%s`", expectedString, *result)
	}
}

func Test_getClientIDAKSWorkload(t *testing.T) {
	expectedString := "testClientIDAKSWorkload"
	err := os.Setenv("ARM_CLIENT_ID", expectedString)
	if err != nil {
		t.Fatalf("failed to set environment variable ARM_CLIENT_ID: %v", err)
	}

	p := &ProviderModel{
		UseAKSWorkloadIdentity: basetypes.NewBoolValue(true),
	}
	result, err := getClientId(p)
	if err != nil {
		t.Fatalf("getClientID returned unexpected error %v", err)
	}
	if result == nil {
		t.Fatalf("getClientID returned nil result without an error")
	}
	if *result != expectedString {
		t.Fatalf("getClientID did not return expected string `%s`, got `%s`", expectedString, *result)
	}
}

func Test_getClientIDAKSWorkloadExpectMismatch(t *testing.T) {
	configuredString := "testClientIDAKSWorkload"
	err := os.Setenv("ARM_CLIENT_ID", "testClientID")
	if err != nil {
		t.Fatalf("failed to set environment variable ARM_CLIENT_ID: %v", err)
	}
	p := &ProviderModel{
		ClientId:               basetypes.NewStringValue(configuredString),
		UseAKSWorkloadIdentity: basetypes.NewBoolValue(true),
	}

	_, err = getClientId(p)
	if err == nil {
		t.Fatalf("expected an error but did not get one")
	}
	if !strings.EqualFold(err.Error(), "mismatch between supplied Client ID and that provided by AKS Workload Identity - please remove, ensure they match, or disable use_aks_workload_identity") {
		t.Fatal("did not get expected error")
	}
}

func Test_decodeCertificate(t *testing.T) {
	input := "dGVzdERhdGE="
	expected := "testData"

	result, err := decodeCertificate(input)
	if err != nil {
		t.Fatalf("decodeCertificate returned unexpected error %v", err)
	}
	if result == nil {
		t.Fatalf("decodeCertificate returned nil result without an error")
	}
	if string(result) != expected {
		t.Fatalf("decodeCertificate did not return expected result `%s`, got `%s`", expected, result)
	}
}

func Test_decodeCertificateExpectError(t *testing.T) {
	input := "NotValidInput"

	_, err := decodeCertificate(input)
	if err == nil {
		t.Fatalf("expected an error but did not get one")
	}
	if !strings.HasPrefix(err.Error(), "could not decode client certificate data:") {
		t.Fatalf("did not get expected error, got '%v'", err)
	}
}
