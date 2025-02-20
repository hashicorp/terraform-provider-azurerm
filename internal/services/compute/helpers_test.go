// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
)

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
