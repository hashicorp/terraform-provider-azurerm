package parsers

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/azure"
)

func TestParseResourceID(t *testing.T) {
	testData := []struct {
		Environment azure.Environment
		Input       string
	}{
		{
			Environment: azure.ChinaCloud,
			Input:       "https://account1.blob.core.chinacloudapi.cn/container1",
		},
		{
			Environment: azure.GermanCloud,
			Input:       "https://account1.blob.core.cloudapi.de/container1",
		},
		{
			Environment: azure.PublicCloud,
			Input:       "https://account1.blob.core.windows.net/container1",
		},
		{
			Environment: azure.USGovernmentCloud,
			Input:       "https://account1.blob.core.usgovcloudapi.net/container1",
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing Environment %q", v.Environment.Name)
		actual, err := ParseContainerID(v.Input)
		if err != nil {
			t.Fatal(err)
		}

		if actual.AccountName != "account1" {
			t.Fatalf("Expected the account name to be `account1` but got %q", actual.AccountName)
		}

		if actual.ContainerName != "container1" {
			t.Fatalf("Expected the container name to be `container1` but got %q", actual.ContainerName)
		}
	}
}
