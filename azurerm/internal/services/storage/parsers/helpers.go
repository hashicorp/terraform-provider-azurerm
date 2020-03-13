package parsers

import (
	"fmt"
	"strings"
)

func getAccountNameFromEndpoint(endpoint string) (*string, error) {
	segments := strings.Split(endpoint, ".")
	if len(segments) == 0 {
		return nil, fmt.Errorf("The Endpoint contained no segments")
	}
	return &segments[0], nil
}

// getBlobEndpoint returns the endpoint for Blob API Operations on this storage account
func getBlobEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.blob.%s", accountName, baseUri)
}

// getDataLakeStoreEndpoint returns the endpoint for Data Lake Store API Operations on this storage account
func getDataLakeStoreEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.dfs.%s", accountName, baseUri)
}

// getFileEndpoint returns the endpoint for File Share API Operations on this storage account
func getFileEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.file.%s", accountName, baseUri)
}

// getQueueEndpoint returns the endpoint for Queue API Operations on this storage account
func getQueueEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.queue.%s", accountName, baseUri)
}

// getTableEndpoint returns the endpoint for Table API Operations on this storage account
func getTableEndpoint(baseUri string, accountName string) string {
	return fmt.Sprintf("https://%s.table.%s", accountName, baseUri)
}
