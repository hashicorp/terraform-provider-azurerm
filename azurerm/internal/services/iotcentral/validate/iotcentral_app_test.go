package validate

import (
	"testing"
)

func TestIotCentralAppName(t *testing.T) {
	validNames := []string{
		"a1", "11", "1a", "aa",
		"1-1", "aaa-aa", "a--a-aa",
		"a1-1", "a1-a", "1a-1", "1a-a-1-2",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde123",
	}

	for _, name := range validNames {
		_, err := IotCentralAppName(name, "unit test")
		if err != nil {
			t.Fatalf("%q should be a valid IoT Central Name: %q", name, err)
		}
	}

	invalidNames := []string{
		"a", "1",
		"1-", "a-",
		"a1-", "1a--1-1-a-", "aa-",
		"a1-", "1a--a1-", "aa-",
		"-", "-1", "-a",
		"AA", "AA-1", "AA-a",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
	}

	for _, name := range invalidNames {
		_, err := IotCentralAppName(name, "unit test")
		if err == nil {
			t.Fatalf("%q should be an invalid IoT Central Name: %q", name, err)
		}
	}
}

func TestIotCentralAppSubdomain(t *testing.T) {
	validNames := []string{
		"a1", "11", "1a", "aa",
		"1-1", "a-a",
		"a1-1", "a1-a", "1a-1", "1a-a",
		"a1-11", "aa-11", "11-1a", "11-a1",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde123",
	}

	for _, name := range validNames {
		_, err := IotCentralAppSubdomain(name, "unit test")
		if err != nil {
			t.Fatalf("%q should be a valid IoT Central Subdomain name: %q", name, err)
		}
	}

	invalidNames := []string{
		"a", "1",
		"1-", "a-",
		"a1-", "1a-", "aa-",
		"-", "-1", "-a",
		"AA", "AA-1", "AA-a",
		"A1-", "AA-A", "AA-aA",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
	}

	for _, name := range invalidNames {
		_, err := IotCentralAppSubdomain(name, "unit test")
		if err == nil {
			t.Fatalf("%q should be an invalid IoT Central Subdomain Name: %q", name, err)
		}
	}
}

func TestIotCentralAppDisplayName(t *testing.T) {
	validNames := []string{
		"a", "A", "1",
		"1-", "a-", "A-",
		"a1-", "1a-", "aA-", "Aa-",
		"-", "-1", "_-a", "#$%$#!",
		"AA", "AA-1", "AA-a",
		"A1-", "AA-A", "AA-aA",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde1234",
	}

	for _, name := range validNames {
		_, err := IotCentralAppDisplayName(name, "unit test")
		if err != nil {
			t.Fatalf("%q should be a valid IoT Central Display Name: %q", name, err)
		}
	}

	invalidNames := []string{
		"",
		"adcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdssdavcadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdsadcdssdavcc",
	}
	for _, name := range invalidNames {
		_, err := IotCentralAppDisplayName(name, "unit test")
		if err == nil {
			t.Fatalf("%q should be an invalid IoT Central Display Name: %q", name, err)
		}
	}
}

func TestIotCentralAppTemplateName(t *testing.T) {
	validNames := []string{
		"a", "A", "1",
		"1-", "a-", "A-",
		"-", "-1", "_-a", "#$%$#!",
		"abcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcdeabcde",
	}

	for _, name := range validNames {
		_, err := IotCentralAppTemplateName(name, "unit test")
		if err != nil {
			t.Fatalf("%q should be a valid IoT Central Template Name: %q", name, err)
		}
	}

	invalidNames := []string{
		"",
		"abcdeabcdeabcdeabcde@$#%abcdeabcdeadeabcdeabcdeabcdeabcde-1a",
	}
	for _, name := range invalidNames {
		_, err := IotCentralAppTemplateName(name, "unit test")
		if err == nil {
			t.Fatalf("%q should be an invalid IoT Central Template Name: %q", name, err)
		}
	}
}
