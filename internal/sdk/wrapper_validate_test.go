package sdk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"testing"
)

func TestValidateTopLevelObjectValid(t *testing.T) {
	type Person struct {
		Name string `tfschema:"name"`
		Age  int    `tfschema:"int"`
	}
	// fields := make(map[string]*schema.Schema, 0)
	fields := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
		"int": {
			Type: pluginsdk.TypeInt,
		},
	}

	if err := ValidateModelObject(&Person{}, fields); err != nil {
		t.Fatalf("error: %+v", err)
	}
}

func TestValidateTopLevelObjectInvalid(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		Age int `json:"int"`
	}
	fields := map[string]*schema.Schema{
		"int": {
			Type: pluginsdk.TypeInt,
		},
	}
	if err := ValidateModelObject(&Person1{}, fields); err != nil {
		t.Fatalf("expected an error but didn't get one")
	}

	t.Log("Person2")
	type Person2 struct {
		Name string
	}

	fields2 := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
	}
	if err := ValidateModelObject(&Person2{}, fields2); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}

func TestValidateTopLevelObjectInvalidInterface(t *testing.T) {
	type Person struct {
		Name string `tfschema:"name"`
	}
	var p interface{} = Person{}
	if err := ValidateModelObject(&p, nil); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}

func TestValidateNestedObjectValid(t *testing.T) {
	type Pet struct {
		Name string `tfschema:"name"`
	}
	type Person struct {
		Name string `tfschema:"name"`
		Pets []Pet  `tfschema:"pets"`
	}
	if err := ValidateModelObject(&Person{}, nil); err != nil {
		t.Fatalf("error: %+v", err)
	}
}

func TestValidateNestedObjectInvalid(t *testing.T) {
	type Pet struct {
		Name string `tfschema:"name"`
		Age  int
	}
	type Person struct {
		Name string `tfschema:"name"`
		Pets []Pet  `tfschema:"pets"`
	}
	if err := ValidateModelObject(&Person{}, nil); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}
