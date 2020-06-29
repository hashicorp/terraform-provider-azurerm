package examples

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/tombuildsstuff/terraform-configuration-tester/locator"
	"github.com/tombuildsstuff/terraform-configuration-tester/runner"
)

func TestRunExamples(t *testing.T) {
	if os.Getenv("TF_EXAMPLE_TEST") == "" {
		log.Printf("`TF_EXAMPLE_TEST` is not set - skipping")
		t.Skip()
	}

	examplesDirectory := "./"
	directories := locator.DiscoverExamples(examplesDirectory)

	input := runner.TestRunInput{
		// TODO: from Environment Variables in the future
		ProviderVersion:  "1.23.0",
		ProviderName:     "azurerm",
		TerraformVersion: "0.11.13",
		AvailableVariables: []runner.AvailableVariable{
			{
				Name:     "prefix",
				Generate: true,
			},
			{
				Name:       "location",
				EnvKeyName: "ARM_LOCATION",
			},
			{
				Name:       "alt_location",
				EnvKeyName: "ARM_LOCATION_ALT",
			},
			{
				Name:       "kubernetes_client_id",
				EnvKeyName: "ARM_CLIENT_ID",
			},
			{
				Name:       "kubernetes_client_secret",
				EnvKeyName: "ARM_CLIENT_SECRET",
			},
		},
	}

	for _, directoryPath := range directories {
		shortDirName := strings.Replace(directoryPath, examplesDirectory, "", -1)
		testName := fmt.Sprintf("examples/%s", shortDirName)
		t.Run(testName, func(t *testing.T) {
			if err := input.Run(directoryPath); err != nil {
				t.Fatalf("Error running %q: %s", shortDirName, err)
			}
		})
	}
}
