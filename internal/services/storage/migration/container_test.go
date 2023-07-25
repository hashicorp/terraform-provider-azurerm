// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func TestContainerV0ToV1(t *testing.T) {
	clouds := []*environments.Environment{
		environments.AzurePublic(),
		environments.AzureChina(),
		environments.AzureUSGovernment(),
	}

	for _, cloud := range clouds {
		t.Logf("[DEBUG] Testing with Cloud %q", cloud.Name)

		input := map[string]interface{}{
			"id":                   "old-id",
			"name":                 "some-name",
			"storage_account_name": "some-account",
		}

		meta := &clients.Client{
			Account: &clients.ResourceManagerAccount{
				Environment: *cloud,
			},
		}

		suffix, ok := meta.Account.Environment.Storage.DomainSuffix()
		if !ok {
			t.Fatalf("could not determine Storage domain suffix for environment %q", meta.Account.Environment.Name)
		}

		expected := map[string]interface{}{
			"id":                   fmt.Sprintf("https://some-account.blob.%s/some-name", *suffix),
			"name":                 "some-name",
			"storage_account_name": "some-account",
		}

		actual, err := ContainerV0ToV1{}.UpgradeFunc()(context.TODO(), input, meta)
		if err != nil {
			t.Fatalf("Expected no error but got: %s", err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("Expected %+v. Got %+v. But expected them to be the same", expected, actual)
		}

		t.Logf("[DEBUG] Ok!")
	}
}
