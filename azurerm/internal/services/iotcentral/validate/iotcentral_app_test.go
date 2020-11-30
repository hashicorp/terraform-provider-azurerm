package validate

import (
	"testing"
)

func TestIotCentralAppName(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a1",
			Error: false,
		},
		{
			Value: "11",
			Error: false,
		},
		{
			Value: "1a",
			Error: false,
		},
		{
			Value: "aa",
			Error: false,
		},
		{
			Value: "1-1",
			Error: false,
		},
		{
			Value: "aaa-aa",
			Error: false,
		},
		{
			Value: "a--a-aa",
			Error: false,
		},
		{
			Value: "a1-1",
			Error: false,
		},
		{
			Value: "a1-a",
			Error: false,
		},
		{
			Value: "1a-1",
			Error: false,
		},
		{
			Value: "1a-a-1-2",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde123",
			Error: false,
		},
		{
			Value: "a",
			Error: true,
		},
		{
			Value: "1",
			Error: true,
		},
		{
			Value: "1-",
			Error: true,
		},
		{
			Value: "a-",
			Error: true,
		},
		{
			Value: "a1-",
			Error: true,
		},
		{
			Value: "1a--1-1-a-",
			Error: true,
		},
		{
			Value: "aa-",
			Error: true,
		},
		{
			Value: "a1-",
			Error: true,
		},
		{
			Value: "1a--a1-",
			Error: true,
		},
		{
			Value: "aa-",
			Error: true,
		},
		{
			Value: "-",
			Error: true,
		},
		{
			Value: "-1",
			Error: true,
		},
		{
			Value: "-a",
			Error: true,
		},
		{
			Value: "AA",
			Error: true,
		},
		{
			Value: "AA-1",
			Error: true,
		},
		{
			Value: "AA-a",
			Error: true,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := IotCentralAppName(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}

func TestIotCentralAppSubdomain(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a1",
			Error: false,
		},
		{
			Value: "11",
			Error: false,
		},
		{
			Value: "1a",
			Error: false,
		},
		{
			Value: "aa",
			Error: false,
		},
		{
			Value: "1-1",
			Error: false,
		},
		{
			Value: "a-a",
			Error: false,
		},
		{
			Value: "a1-1",
			Error: false,
		},
		{
			Value: "a1-a",
			Error: false,
		},
		{
			Value: "1a-1",
			Error: false,
		},
		{
			Value: "1a-a",
			Error: false,
		},
		{
			Value: "a1-11",
			Error: false,
		},
		{
			Value: "aa-11",
			Error: false,
		},
		{
			Value: "11-1a",
			Error: false,
		},
		{
			Value: "11-a1",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde123",
			Error: false,
		},
		{
			Value: "a",
			Error: true,
		},
		{
			Value: "1",
			Error: true,
		},
		{
			Value: "1-",
			Error: true,
		},
		{
			Value: "a-",
			Error: true,
		},
		{
			Value: "a1-",
			Error: true,
		},
		{
			Value: "1a-",
			Error: true,
		},
		{
			Value: "aa-",
			Error: true,
		},
		{
			Value: "-",
			Error: true,
		},
		{
			Value: "-1",
			Error: true,
		},
		{
			Value: "-a",
			Error: true,
		},
		{
			Value: "AA",
			Error: true,
		},
		{
			Value: "AA-1",
			Error: true,
		},
		{
			Value: "AA-a",
			Error: true,
		},
		{
			Value: "A1-",
			Error: true,
		},
		{
			Value: "AA-A",
			Error: true,
		},
		{
			Value: "AA-aA",
			Error: true,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := IotCentralAppSubdomain(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}

func TestIotCentralAppDisplayName(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a",
			Error: false,
		},
		{
			Value: "A",
			Error: false,
		},
		{
			Value: "1",
			Error: false,
		},
		{
			Value: "1-",
			Error: false,
		},
		{
			Value: "a-",
			Error: false,
		},
		{
			Value: "A-",
			Error: false,
		},
		{
			Value: "a1-",
			Error: false,
		},
		{
			Value: "1a-",
			Error: false,
		},
		{
			Value: "aA-",
			Error: false,
		},
		{
			Value: "Aa-",
			Error: false,
		},
		{
			Value: "-",
			Error: false,
		},
		{
			Value: "-1",
			Error: false,
		},
		{
			Value: "_-a",
			Error: false,
		},
		{
			Value: "#$%$#!",
			Error: false,
		},
		{
			Value: "AA",
			Error: false,
		},
		{
			Value: "AA-1",
			Error: false,
		},
		{
			Value: "AA-a",
			Error: false,
		},
		{
			Value: "A1-",
			Error: false,
		},
		{
			Value: "AA-A",
			Error: false,
		},
		{
			Value: "AA-aA",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
			Error: false,
		},

		{
			Value: "",
			Error: true,
		},
		{
			Value: "adcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdssdavcadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdssdavcc",
			Error: true,
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := IotCentralAppDisplayName(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}

func TestIotCentralAppTemplateName(t *testing.T) {
	testData := []struct {
		Value string
		Error bool
	}{
		{
			Value: "a",
			Error: false,
		},
		{
			Value: "A",
			Error: false,
		},
		{
			Value: "1",
			Error: false,
		},
		{
			Value: "1-",
			Error: false,
		},
		{
			Value: "a-",
			Error: false,
		},
		{
			Value: "A-",
			Error: false,
		},
		{
			Value: "-",
			Error: false,
		},
		{
			Value: "-1",
			Error: false,
		},
		{
			Value: "_-a",
			Error: false,
		},
		{
			Value: "#$%$#!",
			Error: false,
		},
		{
			Value: "abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
			Error: false,
		},
		{
			Value: "",
			Error: true,
		},
		{
			Value: "abcdeabcdeabcdeabcde@$#%abcdeabcdeadeabcdeabcdeabcdeabcde-1a",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Value)

		_, err := IotCentralAppTemplateName(v.Value, "unit test")
		if err != nil && !v.Error {
			t.Fatalf("Expected pass but got an error: %s", err)
		}
	}
}
