package queues

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageQueue = Client{}

func TestQueuesLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	queueName := fmt.Sprintf("queue-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}
	queuesClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.%s.%s", accountName, "queue", *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(queuesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	// first let's test an empty container
	_, err = queuesClient.Create(ctx, queueName, CreateInput{MetaData: map[string]string{}})
	if err != nil {
		t.Fatal(fmt.Errorf("error creating: %s", err))
	}

	// then let's retrieve it to ensure there's no metadata..
	resp, err := queuesClient.GetMetaData(ctx, queueName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(resp.MetaData) != 0 {
		t.Fatalf("Expected no MetaData but got: %s", err)
	}

	// then let's add some..
	updatedMetaData := map[string]string{
		"band":  "panic",
		"boots": "the-overpass",
	}
	_, err = queuesClient.SetMetaData(ctx, queueName, SetMetaDataInput{MetaData: updatedMetaData})
	if err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	resp, err = queuesClient.GetMetaData(ctx, queueName)
	if err != nil {
		t.Fatalf("Error re-retrieving MetaData: %s", err)
	}

	if len(resp.MetaData) != 2 {
		t.Fatalf("Expected metadata to have 2 items but got: %s", resp.MetaData)
	}
	if resp.MetaData["band"] != "panic" {
		t.Fatalf("Expected `band` to be `panic` but got: %s", resp.MetaData["band"])
	}
	if resp.MetaData["boots"] != "the-overpass" {
		t.Fatalf("Expected `boots` to be `the-overpass` but got: %s", resp.MetaData["boots"])
	}

	// and woo let's remove it again
	_, err = queuesClient.SetMetaData(ctx, queueName, SetMetaDataInput{MetaData: map[string]string{}})
	if err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	resp, err = queuesClient.GetMetaData(ctx, queueName)
	if err != nil {
		t.Fatalf("Error retrieving MetaData: %s", err)
	}
	if len(resp.MetaData) != 0 {
		t.Fatalf("Expected no MetaData but got: %s", err)
	}

	// set some properties
	props := StorageServiceProperties{
		Logging: &LoggingConfig{
			Version: "1.0",
			Delete:  true,
			Read:    true,
			Write:   true,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
		Cors: &Cors{
			CorsRule: []CorsRule{
				CorsRule{
					AllowedMethods:  "GET,PUT",
					AllowedOrigins:  "http://www.example.com",
					ExposedHeaders:  "x-tempo-*",
					AllowedHeaders:  "x-tempo-*",
					MaxAgeInSeconds: 500,
				},
				CorsRule{
					AllowedMethods:  "POST",
					AllowedOrigins:  "http://www.test.com",
					ExposedHeaders:  "*",
					AllowedHeaders:  "x-method-*",
					MaxAgeInSeconds: 200,
				},
			},
		},
		HourMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
		MinuteMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
	}
	_, err = queuesClient.SetServiceProperties(ctx, SetStorageServicePropertiesInput{Properties: props})
	if err != nil {
		t.Fatalf("SetServiceProperties failed: %s", err)
	}

	properties, err := queuesClient.GetServiceProperties(ctx)
	if err != nil {
		t.Fatalf("GetServiceProperties failed: %s", err)
	}

	if len(properties.Cors.CorsRule) > 1 {
		if properties.Cors.CorsRule[0].AllowedMethods != "GET,PUT" {
			t.Fatalf("CORS Methods weren't set!")
		}
		if properties.Cors.CorsRule[1].AllowedMethods != "POST" {
			t.Fatalf("CORS Methods weren't set!")
		}
	} else {
		t.Fatalf("CORS Methods weren't set!")
	}

	if properties.HourMetrics.Enabled {
		t.Fatalf("HourMetrics were enabled when they shouldn't be!")
	}

	if properties.MinuteMetrics.Enabled {
		t.Fatalf("MinuteMetrics were enabled when they shouldn't be!")
	}

	if !properties.Logging.Write {
		t.Fatalf("Logging Write's was not enabled when they should be!")
	}

	includeAPIS := true
	// set some properties
	props2 := StorageServiceProperties{
		Logging: &LoggingConfig{
			Version: "1.0",
			Delete:  true,
			Read:    true,
			Write:   true,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
		Cors: &Cors{
			CorsRule: []CorsRule{
				CorsRule{
					AllowedMethods:  "PUT",
					AllowedOrigins:  "http://www.example.com",
					ExposedHeaders:  "x-tempo-*",
					AllowedHeaders:  "x-tempo-*",
					MaxAgeInSeconds: 500,
				},
			},
		},
		HourMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: true,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
			IncludeAPIs: &includeAPIS,
		},
		MinuteMetrics: &MetricsConfig{
			Version: "1.0",
			Enabled: false,
			RetentionPolicy: RetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
	}

	_, err = queuesClient.SetServiceProperties(ctx, SetStorageServicePropertiesInput{Properties: props2})
	if err != nil {
		t.Fatalf("SetServiceProperties failed: %s", err)
	}

	properties, err = queuesClient.GetServiceProperties(ctx)
	if err != nil {
		t.Fatalf("GetServiceProperties failed: %s", err)
	}

	if len(properties.Cors.CorsRule) == 1 {
		if properties.Cors.CorsRule[0].AllowedMethods != "PUT" {
			t.Fatalf("CORS Methods weren't set!")
		}
	} else {
		t.Fatalf("CORS Methods weren't set!")
	}

	if !properties.HourMetrics.Enabled {
		t.Fatalf("HourMetrics were enabled when they shouldn't be!")
	}

	if properties.MinuteMetrics.Enabled {
		t.Fatalf("MinuteMetrics were enabled when they shouldn't be!")
	}

	if !properties.Logging.Write {
		t.Fatalf("Logging Write's was not enabled when they should be!")
	}

	log.Printf("[DEBUG] Deleting..")
	_, err = queuesClient.Delete(ctx, queueName)
	if err != nil {
		t.Fatal(fmt.Errorf("error deleting: %s", err))
	}
}
