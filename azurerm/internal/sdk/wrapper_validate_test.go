package sdk

import "testing"

func TestValidateTopLevelObjectValid(t *testing.T) {
	type Person struct {
		Name string `tfschema:"name"`
		Age  int    `tfschema:"int"`
	}
	if err := ValidateModelObject(&Person{}); err != nil {
		t.Fatalf("error: %+v", err)
	}
}

func TestValidateTopLevelObjectInvalid(t *testing.T) {
	t.Log("Person1")
	type Person1 struct {
		Age int `json:"int"`
	}
	if err := ValidateModelObject(&Person1{}); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}

	t.Log("Person2")
	type Person2 struct {
		Name string
	}
	if err := ValidateModelObject(&Person2{}); err == nil {
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
	if err := ValidateModelObject(&Person{}); err != nil {
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
	if err := ValidateModelObject(&Person{}); err == nil {
		t.Fatalf("expected an error but didn't get one")
	}
}
