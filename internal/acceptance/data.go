// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

const (
	// charSetAlphaNum is the alphanumeric character set for use with randStringFromCharSet
	charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"
)

func init() {
	// unit testing
	if os.Getenv("TF_ACC") == "" {
		return
	}
}

type TestData struct {
	// Subscriptions is a set of AAD Subscriptions which should be used for this Test
	Subscriptions Subscriptions

	// Locations is a set of Azure Regions which should be used for this Test
	Locations Regions

	// RandomInteger is a random integer which is unique to this test case
	RandomInteger int

	// RandomString is a random 5 character string is unique to this test case
	RandomString string

	// ResourceName is the fully qualified resource name, comprising the
	// resource type and then the resource label
	// e.g. `azurerm_resource_group.test`
	ResourceName string

	// ResourceType is the Terraform Resource Type - `azurerm_resource_group`
	ResourceType string

	// EnvironmentName is the name of the Azure Environment where we're running
	EnvironmentName string

	// MetadataURL is the url of the endpoint where the environment is obtained
	MetadataURL string

	// resourceLabel is the local used for the resource - generally "test""
	resourceLabel string
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	testData := TestData{
		RandomInteger:   RandTimeInt(),
		RandomString:    randString(5),
		ResourceName:    fmt.Sprintf("%s.%s", resourceType, resourceLabel),
		EnvironmentName: EnvironmentName(),
		MetadataURL:     os.Getenv("ARM_METADATA_HOSTNAME"),

		ResourceType:  resourceType,
		resourceLabel: resourceLabel,
	}

	if features.UseDynamicTestLocations() {
		testData.Locations = availableLocations()
	} else {
		testData.Locations = Regions{
			Primary:   os.Getenv("ARM_TEST_LOCATION"),
			Secondary: os.Getenv("ARM_TEST_LOCATION_ALT"),
			Ternary:   os.Getenv("ARM_TEST_LOCATION_ALT2"),
		}
	}

	testData.Subscriptions = Subscriptions{
		Primary:   os.Getenv("ARM_SUBSCRIPTION_ID"),
		Secondary: os.Getenv("ARM_TEST_SUBSCRIPTION_ID_ALT"),
	}

	return testData
}

// RandomIntOfLength is a random 8 to 18 digit integer which is unique to this test case
func (td *TestData) RandomIntOfLength(len int) int {
	// len should not be
	//  - greater then 18, longest a int can represent
	//  - less then 8, as that gives us YYMMDDRR
	if 8 > len || len > 18 {
		panic("Invalid Test: RandomIntOfLength: len is not between 8 or 18 inclusive")
	}

	// 18 - just return the int
	if len >= 18 {
		return td.RandomInteger
	}

	// 16-17 just strip off the last 1-2 digits
	if len >= 16 {
		return td.RandomInteger / int(math.Pow10(18-len))
	}

	// 8-15 keep len - 2 digits and add 2 characters of randomness on
	s := strconv.Itoa(td.RandomInteger)
	r := s[16:18]
	v := s[0 : len-2]
	i, _ := strconv.Atoi(v + r)

	return i
}

// RandomStringOfLength is a random 1 to 1024 character string which is unique to this test case
func (td *TestData) RandomStringOfLength(len int) string {
	// len should not be less then 1 or greater than 1024
	if 1 > len || len > 1024 {
		panic("Invalid Test: RandomStringOfLength: length argument must be between 1 and 1024 characters")
	}

	return randString(len)
}

// randString generates a random alphanumeric string of the length specified
func randString(strlen int) string {
	return randStringFromCharSet(strlen, charSetAlphaNum)
}

// randStringFromCharSet generates a random string by selecting characters from
// the charset provided
func randStringFromCharSet(strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}
	return string(result)
}
