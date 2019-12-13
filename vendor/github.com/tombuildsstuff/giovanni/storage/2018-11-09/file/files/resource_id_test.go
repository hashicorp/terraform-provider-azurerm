package files

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
)

func TestGetResourceID(t *testing.T) {
	testData := []struct {
		Environment azure.Environment
		Expected    string
	}{
		{
			Environment: azure.ChinaCloud,
			Expected:    "https://account1.file.core.chinacloudapi.cn/share1/directory1/file1.txt",
		},
		{
			Environment: azure.GermanCloud,
			Expected:    "https://account1.file.core.cloudapi.de/share1/directory1/file1.txt",
		},
		{
			Environment: azure.PublicCloud,
			Expected:    "https://account1.file.core.windows.net/share1/directory1/file1.txt",
		},
		{
			Environment: azure.USGovernmentCloud,
			Expected:    "https://account1.file.core.usgovcloudapi.net/share1/directory1/file1.txt",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		c := NewWithEnvironment(v.Environment)
		actual := c.GetResourceID("account1", "share1", "directory1", "file1.txt")
		if actual != v.Expected {
			t.Fatalf("Expected the Resource ID to be %q but got %q", v.Expected, actual)
		}
	}
}

func TestParseResourceID(t *testing.T) {
	testData := []struct {
		Environment azure.Environment
		Input       string
	}{
		{
			Environment: azure.ChinaCloud,
			Input:       "https://account1.file.core.chinacloudapi.cn/share1/directory1/file1.txt",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.file.core.cloudapi.de/share1/directory1/file1.txt",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.file.core.windows.net/share1/directory1/file1.txt",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.file.core.usgovcloudapi.net/share1/directory1/file1.txt",
		},
	}

	t.Logf("[DEBUG] Top Level Files")
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		actual, err := ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected Account Name to be `account1` but got %q", actual.AccountName)
		}
		if actual.ShareName != "share1" {
			t.Fatalf("Expected Share Name to be `share1` but got %q", actual.ShareName)
		}
		if actual.DirectoryName != "directory1" {
			t.Fatalf("Expected Directory Name to be `directory1` but got %q", actual.DirectoryName)
		}
		if actual.FileName != "file1.txt" {
			t.Fatalf("Expected File Name to be `file1.txt` but got %q", actual.FileName)
		}
	}

	testData = []struct {
		Environment azure.Environment
		Input       string
	}{
		{
			Environment: azure.ChinaCloud,
			Input:       "https://account1.file.core.chinacloudapi.cn/share1/directory1/directory2/file1.txt",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.file.core.cloudapi.de/share1/directory1/directory2/file1.txt",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.file.core.windows.net/share1/directory1/directory2/file1.txt",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.file.core.usgovcloudapi.net/share1/directory1/directory2/file1.txt",
		},
	}

	t.Logf("[DEBUG] Nested Files")
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		actual, err := ParseResourceID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected Account Name to be `account1` but got %q", actual.AccountName)
		}
		if actual.ShareName != "share1" {
			t.Fatalf("Expected Share Name to be `share1` but got %q", actual.ShareName)
		}
		if actual.DirectoryName != "directory1/directory2" {
			t.Fatalf("Expected Directory Name to be `directory1/directory2` but got %q", actual.DirectoryName)
		}
		if actual.FileName != "file1.txt" {
			t.Fatalf("Expected File Name to be `file1.txt` but got %q", actual.FileName)
		}
	}
}
