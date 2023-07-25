// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

func TestSortVersions_valid(t *testing.T) {
	testData := []struct {
		input    []compute.GalleryImageVersion
		expected []compute.GalleryImageVersion
	}{
		{
			input: []compute.GalleryImageVersion{
				{Name: utils.String("1.0.1")},
				{Name: utils.String("1.2.15.0")},
				{Name: utils.String("1.0.8")},
				{Name: utils.String("1.0.9")},
				{Name: utils.String("1.0.1.1")},
				{Name: utils.String("1.0.10")},
			},
			expected: []compute.GalleryImageVersion{
				{Name: utils.String("1.0.1")},
				{Name: utils.String("1.0.1.1")},
				{Name: utils.String("1.0.8")},
				{Name: utils.String("1.0.9")},
				{Name: utils.String("1.0.10")},
				{Name: utils.String("1.2.15.0")},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing sortSharedImageVersions..")

		actual, errs := sortSharedImageVersions(v.input)
		if len(errs) > 0 {
			t.Fatalf("Error parsing version: %v", errs)
		}
		if eq := reflect.DeepEqual(v.expected, actual); !eq {
			t.Fatalf("Expected %+v but got %+v", v.expected, actual)
		}
	}
}

func TestSortVersions_invalid(t *testing.T) {
	testData := []struct {
		input    []compute.GalleryImageVersion
		expected []compute.GalleryImageVersion
	}{
		{
			input: []compute.GalleryImageVersion{
				{Name: utils.String("1.0.1")},
				{Name: utils.String("1.2.15.0")},
				{Name: utils.String("1.0.8")},
				{Name: utils.String("1.0.9")},
				{Name: utils.String("1.0.1.1")},
				{Name: utils.String("1.0.10")},
				{Name: utils.String("latest")},
			},
			expected: []compute.GalleryImageVersion{
				{Name: utils.String("1.0.1")},
				{Name: utils.String("1.2.15.0")},
				{Name: utils.String("1.0.8")},
				{Name: utils.String("1.0.9")},
				{Name: utils.String("1.0.1.1")},
				{Name: utils.String("1.0.10")},
				{Name: utils.String("latest")},
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing sortSharedImageVersions with invalid input..")

		_, errs := sortSharedImageVersions(v.input)
		if len(errs) == 0 {
			t.Fatalf("Expected an error, got none")
		}
	}
}
