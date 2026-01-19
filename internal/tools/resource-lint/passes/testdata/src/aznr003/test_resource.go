package aznr003

import (
	"context"

	"testdata/src/mockpkg/pluginsdk"
	"testdata/src/mockpkg/sdk"
)

// Mock types for testing
type CustomerManagedKey struct {
	KeyVaultKeyID string
}

type Encryption struct {
	KeySource string
}

type NetworkACLs struct {
	DefaultAction string
}

type NetworkRuleSet struct {
	DefaultAction string
}

// Model for the typed resource
type AIServicesModel struct {
	Name string `tfschema:"name"`
}

// Mock resource type
type AIServices struct{}

// Mark this file as containing a typed resource
var _ sdk.Resource = AIServices{}

func (AIServices) ResourceType() string {
	return "azurerm_ai_services"
}

func (AIServices) ModelObject() interface{} {
	return &AIServicesModel{}
}

func (AIServices) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
	}
}

func (AIServices) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (AIServices) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (AIServices) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

func (AIServices) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return nil
		},
	}
}

// CORRECT: expand function as receiver method
func (r AIServices) expandCustomerManagedKeyCorrect(input []CustomerManagedKey) (*Encryption, error) {
	return &Encryption{}, nil
}

// CORRECT: flatten function as receiver method
func (r AIServices) flattenNetworkACLsCorrect(input *NetworkRuleSet) []NetworkACLs {
	return []NetworkACLs{}
}

// VIOLATION: expand function as global function
func expandCustomerManagedKey(input []CustomerManagedKey) (*Encryption, error) { // want `AZNR003`
	return &Encryption{}, nil
}

// VIOLATION: flatten function as global function
func flattenNetworkACLs(input *NetworkRuleSet) []NetworkACLs { // want `AZNR003`
	return []NetworkACLs{}
}

// VIOLATION: ExpandXxx (capital E) as global function
func ExpandSomething(input string) string { // want `AZNR003`
	return input
}

// VIOLATION: FlattenXxx (capital F) as global function
func FlattenSomething(input string) string { // want `AZNR003`
	return input
}

// NOT a violation: regular function without expand/flatten prefix
func processData(input string) string {
	return input
}

// NOT a violation: function with "expand" in the middle
func doExpandWork(input string) string {
	return input
}
