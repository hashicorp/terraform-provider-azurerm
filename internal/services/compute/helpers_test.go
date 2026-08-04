// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesets"
)

func TestExpandVirtualMachineScaleSetAutomaticUpgradePolicy(t *testing.T) {
	input := []interface{}{
		map[string]interface{}{
			"automatic_rollback_enabled":   false,
			"automatic_os_upgrade_enabled": true,
		},
	}

	for _, useRollingUpgradePolicy := range []bool{false, true} {
		t.Run(fmt.Sprintf("UseRollingUpgradePolicy-%t", useRollingUpgradePolicy), func(t *testing.T) {
			actual := ExpandVirtualMachineScaleSetAutomaticUpgradePolicy(input, useRollingUpgradePolicy)
			if pointer.From(actual.UseRollingUpgradePolicy) != useRollingUpgradePolicy {
				t.Fatalf("expected useRollingUpgradePolicy to be %t", useRollingUpgradePolicy)
			}
		})
	}
}

func TestFlattenVirtualMachineScaleSetRollingUpgradePolicy(t *testing.T) {
	input := &virtualmachinescalesets.RollingUpgradePolicy{
		MaxBatchInstancePercent: pointer.To(int64(20)),
	}

	testCases := []struct {
		name                    string
		useRollingUpgradePolicy *bool
		expectedBlockCount      int
	}{
		{
			name:                    "explicitly disabled",
			useRollingUpgradePolicy: pointer.To(false),
			expectedBlockCount:      0,
		},
		{
			name:                    "explicitly enabled",
			useRollingUpgradePolicy: pointer.To(true),
			expectedBlockCount:      1,
		},
		{
			name:               "not specified",
			expectedBlockCount: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := FlattenVirtualMachineScaleSetRollingUpgradePolicy(input, testCase.useRollingUpgradePolicy)
			if len(actual) != testCase.expectedBlockCount {
				t.Fatalf("expected %d rolling upgrade policy blocks, got %d", testCase.expectedBlockCount, len(actual))
			}
		})
	}
}

func TestSortVersions_valid(t *testing.T) {
	testData := []struct {
		input    []galleryimageversions.GalleryImageVersion
		expected []galleryimageversions.GalleryImageVersion
	}{
		{
			input: []galleryimageversions.GalleryImageVersion{
				{Name: pointer.To("1.0.1")},
				{Name: pointer.To("1.2.15.0")},
				{Name: pointer.To("1.0.8")},
				{Name: pointer.To("1.0.9")},
				{Name: pointer.To("1.0.1.1")},
				{Name: pointer.To("1.0.10")},
			},
			expected: []galleryimageversions.GalleryImageVersion{
				{Name: pointer.To("1.0.1")},
				{Name: pointer.To("1.0.1.1")},
				{Name: pointer.To("1.0.8")},
				{Name: pointer.To("1.0.9")},
				{Name: pointer.To("1.0.10")},
				{Name: pointer.To("1.2.15.0")},
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
		input    []galleryimageversions.GalleryImageVersion
		expected []galleryimageversions.GalleryImageVersion
	}{
		{
			input: []galleryimageversions.GalleryImageVersion{
				{Name: pointer.To("1.0.1")},
				{Name: pointer.To("1.2.15.0")},
				{Name: pointer.To("1.0.8")},
				{Name: pointer.To("1.0.9")},
				{Name: pointer.To("1.0.1.1")},
				{Name: pointer.To("1.0.10")},
				{Name: pointer.To("latest")},
			},
			expected: []galleryimageversions.GalleryImageVersion{
				{Name: pointer.To("1.0.1")},
				{Name: pointer.To("1.2.15.0")},
				{Name: pointer.To("1.0.8")},
				{Name: pointer.To("1.0.9")},
				{Name: pointer.To("1.0.1.1")},
				{Name: pointer.To("1.0.10")},
				{Name: pointer.To("latest")},
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
