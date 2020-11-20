package check

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/helpers"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/types"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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

// ExistsInAzure validates that the specified resource exists within Azure
func (t thatType) ExistsInAzure(testResource types.TestResource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client)
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

// DoesNotExist returns a TestCheckFunc which validates that the specific key
// does not exist on the resource
func (t thatWithKeyType) DoesNotExist() resource.TestCheckFunc {
	return resource.TestCheckNoResourceAttr(t.resourceName, t.key)
}

// Exists returns a TestCheckFunc which validates that the specific key exists on the resource
func (t thatWithKeyType) Exists() resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(t.resourceName, t.key)
}

// IsEmpty returns a TestCheckFunc which validates that the specific key is empty on the resource
func (t thatWithKeyType) IsEmpty() resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(t.resourceName, t.key, "")
}

// HasValue returns a TestCheckFunc which validates that the specific key has the
// specified value on the resource
func (t thatWithKeyType) HasValue(value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(t.resourceName, t.key, value)
}

// MatchesOtherKey returns a TestCheckFunc which validates that the key on this resource
// matches another other key on another resource
func (t thatWithKeyType) MatchesOtherKey(other thatWithKeyType) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrPair(t.resourceName, t.key, other.resourceName, other.key)
}

// MatchesRegex returns a TestCheckFunc which validates that the key on this resource matches
// the given regular expression
func (t thatWithKeyType) MatchesRegex(r *regexp.Regexp) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(t.resourceName, t.key, r)
}
