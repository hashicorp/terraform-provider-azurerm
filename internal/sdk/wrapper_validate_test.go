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
	schemaFields := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
		"int": {
			Type: pluginsdk.TypeInt,
		},
	}

	if err := ValidateModelObject(&Person{}, schemaFields); err != nil {
		t.Fatalf("error: %+v", err)
	}
}

func TestValidateTopLevelObjectInvalid(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		Age int `json:"int"`
	}
	schemaFields := map[string]*schema.Schema{
		"int": {
			Type: pluginsdk.TypeInt,
		},
	}
	if err := ValidateModelObject(&Person1{}, schemaFields); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}

	t.Log("Person2")
	type Person2 struct {
		Name string
	}

	schemaFields2 := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
	}
	if err := ValidateModelObject(&Person2{}, schemaFields2); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}

func TestValidateTopLevelObjectInvalidInterface(t *testing.T) {
	type Person struct {
		Name string `tfschema:"name"`
	}
	schemaFields := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
	}
	var p interface{} = Person{}
	if err := ValidateModelObject(&p, schemaFields); err == nil {
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
	schemaFields := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
		"pets": {
			Type: schema.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
	if err := ValidateModelObject(&Person{}, schemaFields); err != nil {
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
	schemaFields := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
		"pets": {
			Type: schema.TypeList,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
	if err := ValidateModelObject(&Person{}, schemaFields); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}

func TestValidateSchemaMissingTopLevelModel(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		Age int `tfschema:"int"`
	}
	schemaFields := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
	}
	if err := ValidateModelObject(&Person1{}, schemaFields); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}

func TestValidateSchemaWrongType(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		Age int `tfschema:"int"`
	}
	schemaFields := map[string]*schema.Schema{
		"int": {
			Type: pluginsdk.TypeString,
		},
	}
	if err := ValidateModelObject(&Person1{}, schemaFields); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}

func TestValidateSchemaInt64(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		Age int64 `tfschema:"int"`
	}
	schemaFields := map[string]*schema.Schema{
		"int": {
			Type: pluginsdk.TypeInt,
		},
	}
	if err := ValidateModelObject(&Person1{}, schemaFields); err != nil {
		t.Fatalf("expected no error but got one: %+v", err)
	}
}

func TestValidateSchemaString(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		Name string `tfschema:"name"`
	}
	schemaFields := map[string]*schema.Schema{
		"name": {
			Type: pluginsdk.TypeString,
		},
	}
	if err := ValidateModelObject(&Person1{}, schemaFields); err != nil {
		t.Fatalf("expected no error but got one: %+v", err)
	}
}

func TestValidateSchemaBool(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		IsHuman bool `tfschema:"is_human"`
	}
	schemaFields := map[string]*schema.Schema{
		"is_human": {
			Type: pluginsdk.TypeBool,
		},
	}
	if err := ValidateModelObject(&Person1{}, schemaFields); err != nil {
		t.Fatalf("expected no error but got one: %+v", err)
	}
}
