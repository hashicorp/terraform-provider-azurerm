// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "./testdata/oidc_test_input.txt")

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
	t.Setenv("AZURE_FEDERATED_TOKEN_FILE", "./testdata/oidc_test_input.txt")

	p := &ProviderModel{
		OIDCToken:              basetypes.NewStringValue(configuredString),
		UseAKSWorkloadIdentity: basetypes.NewBoolValue(true),
	}

	_, err := getOidcToken(p)
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
	t.Setenv("ARM_CLIENT_SECRET", "")
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
	t.Setenv("ARM_CLIENT_SECRET", "foo")
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
	t.Setenv("ARM_CLIENT_ID", expectedString)

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
	t.Setenv("ARM_CLIENT_ID", "testClientID")

	p := &ProviderModel{
		ClientId:               basetypes.NewStringValue(configuredString),
		UseAKSWorkloadIdentity: basetypes.NewBoolValue(true),
	}

	_, err := getClientId(p)
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

func Test_extractDefaultTags_Null(t *testing.T) {
	frameworkTags := basetypes.NewMapNull(basetypes.StringType{})
	var diags diag.Diagnostics

	result := extractDefaultTags(context.Background(), frameworkTags, &diags)

	if len(result) != 0 {
		t.Fatalf("expected empty map for null, got %d items", len(result))
	}
	if diags.HasError() {
		t.Fatalf("expected no diagnostics for null, got %v", diags.Errors())
	}
}

func Test_extractDefaultTags_Unknown(t *testing.T) {
	frameworkTags := basetypes.NewMapUnknown(basetypes.StringType{})
	var diags diag.Diagnostics

	result := extractDefaultTags(context.Background(), frameworkTags, &diags)

	if len(result) != 0 {
		t.Fatalf("expected empty map for unknown, got %d items", len(result))
	}
	if diags.HasError() {
		t.Fatalf("expected no diagnostics for unknown, got %v", diags.Errors())
	}
}

func Test_extractDefaultTags_SingleTag(t *testing.T) {
	expectedKey := "env"
	expectedValue := "dev"

	elements := map[string]basetypes.StringValue{
		expectedKey: basetypes.NewStringValue(expectedValue),
	}
	frameworkTags, _ := basetypes.NewMapValueFrom(context.Background(), basetypes.StringType{}, elements)
	var diags diag.Diagnostics

	result := extractDefaultTags(context.Background(), frameworkTags, &diags)

	if len(result) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(result))
	}
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got %v", diags.Errors())
	}
	if result[expectedKey] == nil {
		t.Fatalf("expected tag %s to exist", expectedKey)
	}
	if *result[expectedKey] != expectedValue {
		t.Fatalf("expected %s=%s, got %s=%s", expectedKey, expectedValue, expectedKey, *result[expectedKey])
	}
}

func Test_extractDefaultTags_MultipleTags(t *testing.T) {
	elements := map[string]basetypes.StringValue{
		"managed_by": basetypes.NewStringValue("terraform"),
		"owner":      basetypes.NewStringValue("platform"),
		"env":        basetypes.NewStringValue("prod"),
	}
	frameworkTags, _ := basetypes.NewMapValueFrom(context.Background(), basetypes.StringType{}, elements)
	var diags diag.Diagnostics

	result := extractDefaultTags(context.Background(), frameworkTags, &diags)

	if len(result) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(result))
	}
	if diags.HasError() {
		t.Fatalf("expected no diagnostics, got %v", diags.Errors())
	}

	expectedTags := map[string]string{
		"managed_by": "terraform",
		"owner":      "platform",
		"env":        "prod",
	}

	for key, expectedValue := range expectedTags {
		if result[key] == nil {
			t.Fatalf("expected tag %s to exist", key)
		}
		if *result[key] != expectedValue {
			t.Fatalf("expected %s=%s, got %s", key, expectedValue, *result[key])
		}
	}
}

func Test_extractDefaultTags_PointerCapture(t *testing.T) {
	// This test verifies that pointers correctly capture unique values
	// and don't all point to the same address (a common Go gotcha)
	elements := map[string]basetypes.StringValue{
		"tag1": basetypes.NewStringValue("value1"),
		"tag2": basetypes.NewStringValue("value2"),
	}
	frameworkTags, _ := basetypes.NewMapValueFrom(context.Background(), basetypes.StringType{}, elements)
	var diags diag.Diagnostics

	result := extractDefaultTags(context.Background(), frameworkTags, &diags)

	tag1Ptr := result["tag1"]
	tag2Ptr := result["tag2"]

	if tag1Ptr == tag2Ptr {
		t.Fatal("pointers should not be equal (indicates loop variable capture bug)")
	}
	if *tag1Ptr == *tag2Ptr {
		t.Fatal("values should not be equal (indicates loop variable capture bug)")
	}
}
