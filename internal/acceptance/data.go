// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package acceptance

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/vcr"
)

const (
	// charSetAlphaNum is the alphanumeric character set for use with randStringFromCharSet
	charSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"
)

func init() {
	// unit testing via go-vcr
	if os.Getenv("TF_ACC") == "" {
		return
	}

	if os.Getenv("TC_TEST_VIA_VCR") == "replay" {
		// Override real subscription IDs with placeholders so we natively use the placeholders during replay. This
		// is required for the ImportStep to work
		os.Setenv("ARM_SUBSCRIPTION_ID", vcr.SubscriptionPlaceholder)
		os.Setenv("ARM_SUBSCRIPTION_ID_ALT", vcr.SubscriptionPlaceholderAlt)
		os.Setenv("ARM_SUBSCRIPTION_ID_ALT2", vcr.SubscriptionPlaceholderAlt2)
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

// vcrRandTimeInt produces a stable 18-digit integer from a hash of the test name.
// It mimics the YYMMddHHmmsshhRRRR shape of RandTimeInt but is deterministic.
func vcrRandTimeInt(testName string) int {
	h := fnv.New64a()
	h.Write([]byte(testName))
	u := h.Sum64()

	// Use a fixed date prefix so the result is always 18 digits and never time-dependent.
	// We use 20450101 (8 digits) followed by 10 digits from the hash.
	const fixedPrefix = "20450101"
	postfix := fmt.Sprintf("%010d", u%10000000000)
	i, _ := strconv.Atoi(fixedPrefix + postfix)
	return i
}

// vcrRandString produces a stable random string from the provided rng.
func vcrRandString(rng *rand.Rand, strlen int) string {
	result := make([]byte, strlen)
	for i := range result {
		result[i] = charSetAlphaNum[rng.Intn(len(charSetAlphaNum))]
	}
	return string(result)
}

// BuildTestData generates some test data for the given resource
func BuildTestData(t *testing.T, resourceType string, resourceLabel string) TestData {
	var randomInt int
	var randomString string
	if os.Getenv("TC_TEST_VIA_VCR") != "" {
		// In VCR mode, seed from the test name so all random values are
		// stable across runs. Both values share the same rng so they are
		// deterministic relative to each other as well.
		h := fnv.New64a()
		h.Write([]byte(t.Name()))
		rng := rand.New(rand.NewSource(int64(h.Sum64())))
		randomInt = vcrRandTimeInt(t.Name())
		randomString = vcrRandString(rng, 5)
	} else {
		randomInt = RandTimeInt()
		randomString = randString(5)
	}

	testData := TestData{
		RandomInteger:   randomInt,
		RandomString:    randomString,
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
		Secondary: os.Getenv("ARM_SUBSCRIPTION_ID_ALT"),
	}

	return testData
}

// RandomIntOfLength is a random 8 to 18 digit integer which is unique to this test case
func (td *TestData) RandomIntOfLength(length int) int {
	// length should not be
	//  - greater then 18, longest a int can represent
	//  - less then 8, as that gives us YYMMDDRR
	if 8 > length || length > 18 {
		panic("Invalid Test: RandomIntOfLength: length is not between 8 or 18 inclusive")
	}

	// 18 - just return the int
	if length >= 18 {
		return td.RandomInteger
	}

	// 16-17 just strip off the last 1-2 digits
	if length >= 16 {
		return td.RandomInteger / int(math.Pow10(18-length))
	}

	// 8-15 keep length - 2 digits and add 2 characters of randomness on
	s := strconv.Itoa(td.RandomInteger)
	r := s[16:18]
	v := s[0 : length-2]
	i, _ := strconv.Atoi(v + r)

	return i
}

// RandomStringOfLength is a random 1 to 1024 character string which is unique to this test case
func (td *TestData) RandomStringOfLength(length int) string {
	// len should not be less then 1 or greater than 1024
	if 1 > length || length > 1024 {
		panic("Invalid Test: RandomStringOfLength: length argument must be between 1 and 1024 characters")
	}

	return randString(length)
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
