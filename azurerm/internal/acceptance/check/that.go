package check

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/helpers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/testclient"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type thatType struct {
	// resourceName being the full resource name e.g. azurerm_foo.bar
	resourceName string
}

// Key returns a type which can be used for more fluent assertions for a given Resource
func That(resourceName string) thatType {
	return thatType{
		resourceName: resourceName,
	}
}

// DoesNotExistInAzure validates that the specified resource does not exist within Azure
func (t thatType) DoesNotExistInAzure(testResource types.TestResource) pluginsdk.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := testclient.Build()
		if err != nil {
			return fmt.Errorf("building client: %+v", err)
		}
		return helpers.DoesNotExistInAzure(client, testResource, t.resourceName)(s)
	}
}

// ExistsInAzure validates that the specified resource exists within Azure
func (t thatType) ExistsInAzure(testResource types.TestResource) pluginsdk.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := testclient.Build()
		if err != nil {
			return fmt.Errorf("building client: %+v", err)
		}
		return helpers.ExistsInAzure(client, testResource, t.resourceName)(s)
	}
}

// Key returns a type which can be used for more fluent assertions for a given Resource & Key combination
func (t thatType) Key(key string) thatWithKeyType {
	return thatWithKeyType{
		resourceName: t.resourceName,
		key:          key,
	}
}

type thatWithKeyType struct {
	// resourceName being the full resource name e.g. azurerm_foo.bar
	resourceName string

	// key being the specific field we're querying e.g. bar or a nested object ala foo.0.bar
	key string
}

// JsonAssertionFunc is a function which takes a deserialized JSON object and asserts on it
type JsonAssertionFunc func(input []interface{}) (*bool, error)

// ContainsKeyValue returns a TestCheckFunc which asserts upon a given JSON string set into
// the State by deserializing it and then asserting on it via the JsonAssertionFunc
func (t thatWithKeyType) ContainsJsonValue(assertion JsonAssertionFunc) pluginsdk.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, exists := s.RootModule().Resources[t.resourceName]
		if !exists {
			return fmt.Errorf("%q was not found in the state", t.resourceName)
		}

		value, exists := rs.Primary.Attributes[t.key]
		if !exists {
			return fmt.Errorf("the value %q does not exist within %q", t.key, t.resourceName)
		}

		if value == "" {
			return fmt.Errorf("the value for %q was empty", t.key)
		}

		var out []interface{}
		if err := json.Unmarshal([]byte(value), &out); err != nil {
			return fmt.Errorf("deserializing the value for %q (%q) to json: %+v", t.key, value, err)
		}

		ok, err := assertion(out)
		if err != nil {
			return fmt.Errorf("asserting value for %q: %+v", t.key, err)
		}

		if ok == nil || !*ok {
			return fmt.Errorf("assertion failed for %q: %+v", t.key, err)
		}

		return nil
	}
}

// DoesNotExist returns a TestCheckFunc which validates that the specific key
// does not exist on the resource
func (t thatWithKeyType) DoesNotExist() pluginsdk.TestCheckFunc {
	return resource.TestCheckNoResourceAttr(t.resourceName, t.key)
}

// Exists returns a TestCheckFunc which validates that the specific key exists on the resource
func (t thatWithKeyType) Exists() pluginsdk.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(t.resourceName, t.key)
}

// IsEmpty returns a TestCheckFunc which validates that the specific key is empty on the resource
func (t thatWithKeyType) IsEmpty() pluginsdk.TestCheckFunc {
	return resource.TestCheckResourceAttr(t.resourceName, t.key, "")
}

// HasValue returns a TestCheckFunc which validates that the specific key has the
// specified value on the resource
func (t thatWithKeyType) HasValue(value string) pluginsdk.TestCheckFunc {
	return resource.TestCheckResourceAttr(t.resourceName, t.key, value)
}

// MatchesOtherKey returns a TestCheckFunc which validates that the key on this resource
// matches another other key on another resource
func (t thatWithKeyType) MatchesOtherKey(other thatWithKeyType) pluginsdk.TestCheckFunc {
	return resource.TestCheckResourceAttrPair(t.resourceName, t.key, other.resourceName, other.key)
}

// MatchesRegex returns a TestCheckFunc which validates that the key on this resource matches
// the given regular expression
func (t thatWithKeyType) MatchesRegex(r *regexp.Regexp) pluginsdk.TestCheckFunc {
	return resource.TestMatchResourceAttr(t.resourceName, t.key, r)
}
