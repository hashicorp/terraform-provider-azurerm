// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type accountEndpoints struct {
	primaryBlobEndpoint            string
	primaryBlobHostName            string
	primaryBlobInternetEndpoint    string
	primaryBlobInternetHostName    string
	primaryBlobMicrosoftEndpoint   string
	primaryBlobMicrosoftHostName   string
	secondaryBlobEndpoint          string
	secondaryBlobHostName          string
	secondaryBlobInternetEndpoint  string
	secondaryBlobInternetHostName  string
	secondaryBlobMicrosoftEndpoint string
	secondaryBlobMicrosoftHostName string

	primaryDfsEndpoint            string
	primaryDfsHostName            string
	primaryDfsInternetEndpoint    string
	primaryDfsInternetHostName    string
	primaryDfsMicrosoftEndpoint   string
	primaryDfsMicrosoftHostName   string
	secondaryDfsInternetEndpoint  string
	secondaryDfsInternetHostName  string
	secondaryDfsEndpoint          string
	secondaryDfsHostName          string
	secondaryDfsMicrosoftEndpoint string
	secondaryDfsMicrosoftHostName string

	primaryFileEndpoint            string
	primaryFileHostName            string
	primaryFileInternetEndpoint    string
	primaryFileInternetHostName    string
	primaryFileMicrosoftEndpoint   string
	primaryFileMicrosoftHostName   string
	secondaryFileInternetEndpoint  string
	secondaryFileInternetHostName  string
	secondaryFileEndpoint          string
	secondaryFileHostName          string
	secondaryFileMicrosoftEndpoint string
	secondaryFileMicrosoftHostName string

	primaryQueueEndpoint            string
	primaryQueueHostName            string
	primaryQueueMicrosoftEndpoint   string
	primaryQueueMicrosoftHostName   string
	secondaryQueueEndpoint          string
	secondaryQueueHostName          string
	secondaryQueueMicrosoftEndpoint string
	secondaryQueueMicrosoftHostName string

	primaryTableEndpoint            string
	primaryTableHostName            string
	primaryTableMicrosoftEndpoint   string
	primaryTableMicrosoftHostName   string
	secondaryTableEndpoint          string
	secondaryTableHostName          string
	secondaryTableMicrosoftEndpoint string
	secondaryTableMicrosoftHostName string

	primaryWebEndpoint            string
	primaryWebHostName            string
	primaryWebInternetEndpoint    string
	primaryWebInternetHostName    string
	primaryWebMicrosoftEndpoint   string
	primaryWebMicrosoftHostName   string
	secondaryWebInternetEndpoint  string
	secondaryWebInternetHostName  string
	secondaryWebEndpoint          string
	secondaryWebHostName          string
	secondaryWebMicrosoftEndpoint string
	secondaryWebMicrosoftHostName string
}

func (a accountEndpoints) set(d *pluginsdk.ResourceData) error {
	d.Set("primary_blob_endpoint", a.primaryBlobEndpoint)
	d.Set("primary_blob_host", a.primaryBlobHostName)
	d.Set("primary_blob_internet_endpoint", a.primaryBlobInternetEndpoint)
	d.Set("primary_blob_internet_host", a.primaryBlobInternetHostName)
	d.Set("primary_blob_microsoft_endpoint", a.primaryBlobMicrosoftEndpoint)
	d.Set("primary_blob_microsoft_host", a.primaryBlobMicrosoftHostName)
	d.Set("secondary_blob_endpoint", a.secondaryBlobEndpoint)
	d.Set("secondary_blob_host", a.secondaryBlobHostName)
	d.Set("secondary_blob_internet_endpoint", a.secondaryBlobInternetEndpoint)
	d.Set("secondary_blob_internet_host", a.secondaryBlobInternetHostName)
	d.Set("secondary_blob_microsoft_endpoint", a.secondaryBlobMicrosoftEndpoint)
	d.Set("secondary_blob_microsoft_host", a.secondaryBlobMicrosoftHostName)

	d.Set("primary_dfs_endpoint", a.primaryDfsEndpoint)
	d.Set("primary_dfs_host", a.primaryDfsHostName)
	d.Set("primary_dfs_internet_endpoint", a.primaryDfsInternetEndpoint)
	d.Set("primary_dfs_internet_host", a.primaryDfsInternetHostName)
	d.Set("primary_dfs_microsoft_endpoint", a.primaryDfsMicrosoftEndpoint)
	d.Set("primary_dfs_microsoft_host", a.primaryDfsMicrosoftHostName)
	d.Set("secondary_dfs_endpoint", a.secondaryDfsEndpoint)
	d.Set("secondary_dfs_host", a.secondaryDfsHostName)
	d.Set("secondary_dfs_internet_endpoint", a.secondaryDfsInternetEndpoint)
	d.Set("secondary_dfs_internet_host", a.secondaryDfsInternetHostName)
	d.Set("secondary_dfs_microsoft_endpoint", a.secondaryDfsMicrosoftEndpoint)
	d.Set("secondary_dfs_microsoft_host", a.secondaryDfsMicrosoftHostName)

	d.Set("primary_file_endpoint", a.primaryFileEndpoint)
	d.Set("primary_file_host", a.primaryFileHostName)
	d.Set("primary_file_internet_endpoint", a.primaryFileInternetEndpoint)
	d.Set("primary_file_internet_host", a.primaryFileInternetHostName)
	d.Set("primary_file_microsoft_endpoint", a.primaryFileMicrosoftEndpoint)
	d.Set("primary_file_microsoft_host", a.primaryFileMicrosoftHostName)
	d.Set("secondary_file_endpoint", a.secondaryFileEndpoint)
	d.Set("secondary_file_host", a.secondaryFileHostName)
	d.Set("secondary_file_internet_endpoint", a.secondaryFileInternetEndpoint)
	d.Set("secondary_file_internet_host", a.secondaryFileInternetHostName)
	d.Set("secondary_file_microsoft_endpoint", a.secondaryFileMicrosoftEndpoint)
	d.Set("secondary_file_microsoft_host", a.secondaryFileMicrosoftHostName)

	d.Set("primary_queue_endpoint", a.primaryQueueEndpoint)
	d.Set("primary_queue_host", a.primaryQueueHostName)
	d.Set("primary_queue_microsoft_endpoint", a.primaryQueueMicrosoftEndpoint)
	d.Set("primary_queue_microsoft_host", a.primaryQueueMicrosoftHostName)
	d.Set("secondary_queue_endpoint", a.secondaryQueueEndpoint)
	d.Set("secondary_queue_host", a.secondaryQueueHostName)
	d.Set("secondary_queue_microsoft_endpoint", a.secondaryQueueMicrosoftEndpoint)
	d.Set("secondary_queue_microsoft_host", a.secondaryQueueMicrosoftHostName)

	d.Set("primary_table_endpoint", a.primaryTableEndpoint)
	d.Set("primary_table_host", a.primaryTableHostName)
	d.Set("primary_table_microsoft_endpoint", a.primaryTableMicrosoftEndpoint)
	d.Set("primary_table_microsoft_host", a.primaryTableMicrosoftHostName)
	d.Set("secondary_table_endpoint", a.secondaryTableEndpoint)
	d.Set("secondary_table_host", a.secondaryTableHostName)
	d.Set("secondary_table_microsoft_endpoint", a.secondaryTableMicrosoftEndpoint)
	d.Set("secondary_table_microsoft_host", a.secondaryTableMicrosoftHostName)

	d.Set("primary_web_endpoint", a.primaryWebEndpoint)
	d.Set("primary_web_host", a.primaryWebHostName)
	d.Set("secondary_web_endpoint", a.secondaryWebEndpoint)
	d.Set("secondary_web_host", a.secondaryWebHostName)
	d.Set("primary_web_microsoft_endpoint", a.primaryWebMicrosoftEndpoint)
	d.Set("primary_web_microsoft_host", a.primaryWebMicrosoftHostName)
	d.Set("primary_web_internet_endpoint", a.primaryWebInternetEndpoint)
	d.Set("primary_web_internet_host", a.primaryWebInternetHostName)
	d.Set("secondary_web_internet_endpoint", a.secondaryWebInternetEndpoint)
	d.Set("secondary_web_internet_host", a.secondaryWebInternetHostName)
	d.Set("secondary_web_microsoft_endpoint", a.secondaryWebMicrosoftEndpoint)
	d.Set("secondary_web_microsoft_host", a.secondaryWebMicrosoftHostName)

	return nil
}

func flattenAccountEndpoints(primaryEndpoints, secondaryEndpoints *storageaccounts.Endpoints, routingPreference *storageaccounts.RoutingPreference) accountEndpoints {
	output := accountEndpoints{}

	if primaryEndpoints != nil {
		output.primaryBlobEndpoint, output.primaryBlobHostName = flattenAccountEndpointAndHost(primaryEndpoints.Blob)
		output.primaryDfsEndpoint, output.primaryDfsHostName = flattenAccountEndpointAndHost(primaryEndpoints.Dfs)
		output.primaryFileEndpoint, output.primaryFileHostName = flattenAccountEndpointAndHost(primaryEndpoints.File)
		output.primaryQueueEndpoint, output.primaryQueueHostName = flattenAccountEndpointAndHost(primaryEndpoints.Queue)
		output.primaryTableEndpoint, output.primaryTableHostName = flattenAccountEndpointAndHost(primaryEndpoints.Table)
		output.primaryWebEndpoint, output.primaryWebHostName = flattenAccountEndpointAndHost(primaryEndpoints.Web)

		if routingPreference != nil {
			if primaryEndpoints.InternetEndpoints != nil && pointer.From(routingPreference.PublishInternetEndpoints) {
				output.primaryBlobInternetEndpoint, output.primaryBlobInternetHostName = flattenAccountEndpointAndHost(primaryEndpoints.InternetEndpoints.Blob)
				output.primaryDfsInternetEndpoint, output.primaryDfsInternetHostName = flattenAccountEndpointAndHost(primaryEndpoints.InternetEndpoints.Dfs)
				output.primaryFileInternetEndpoint, output.primaryFileInternetHostName = flattenAccountEndpointAndHost(primaryEndpoints.InternetEndpoints.File)
				output.primaryWebInternetEndpoint, output.primaryWebInternetHostName = flattenAccountEndpointAndHost(primaryEndpoints.InternetEndpoints.Web)
			}

			if primaryEndpoints.MicrosoftEndpoints != nil && pointer.From(routingPreference.PublishMicrosoftEndpoints) {
				output.primaryBlobMicrosoftEndpoint, output.primaryBlobMicrosoftHostName = flattenAccountEndpointAndHost(primaryEndpoints.MicrosoftEndpoints.Blob)
				output.primaryDfsMicrosoftEndpoint, output.primaryDfsMicrosoftHostName = flattenAccountEndpointAndHost(primaryEndpoints.MicrosoftEndpoints.Dfs)
				output.primaryFileMicrosoftEndpoint, output.primaryFileMicrosoftHostName = flattenAccountEndpointAndHost(primaryEndpoints.MicrosoftEndpoints.File)
				output.primaryQueueMicrosoftEndpoint, output.primaryQueueMicrosoftHostName = flattenAccountEndpointAndHost(primaryEndpoints.MicrosoftEndpoints.Queue)
				output.primaryTableMicrosoftEndpoint, output.primaryTableMicrosoftHostName = flattenAccountEndpointAndHost(primaryEndpoints.MicrosoftEndpoints.Table)
				output.primaryWebMicrosoftEndpoint, output.primaryWebMicrosoftHostName = flattenAccountEndpointAndHost(primaryEndpoints.MicrosoftEndpoints.Web)
			}
		}
	}

	if secondaryEndpoints != nil {
		output.secondaryBlobEndpoint, output.secondaryBlobHostName = flattenAccountEndpointAndHost(secondaryEndpoints.Blob)
		output.secondaryDfsEndpoint, output.secondaryDfsHostName = flattenAccountEndpointAndHost(secondaryEndpoints.Dfs)
		output.secondaryFileEndpoint, output.secondaryFileHostName = flattenAccountEndpointAndHost(secondaryEndpoints.File)
		output.secondaryQueueEndpoint, output.secondaryQueueHostName = flattenAccountEndpointAndHost(secondaryEndpoints.Queue)
		output.secondaryTableEndpoint, output.secondaryTableHostName = flattenAccountEndpointAndHost(secondaryEndpoints.Table)
		output.secondaryWebEndpoint, output.secondaryWebHostName = flattenAccountEndpointAndHost(secondaryEndpoints.Web)

		if routingPreference != nil {
			if secondaryEndpoints.InternetEndpoints != nil && pointer.From(routingPreference.PublishInternetEndpoints) {
				output.secondaryBlobInternetEndpoint, output.secondaryBlobInternetHostName = flattenAccountEndpointAndHost(secondaryEndpoints.InternetEndpoints.Blob)
				output.secondaryDfsInternetEndpoint, output.secondaryDfsInternetHostName = flattenAccountEndpointAndHost(secondaryEndpoints.InternetEndpoints.Dfs)
				output.secondaryFileInternetEndpoint, output.secondaryFileInternetHostName = flattenAccountEndpointAndHost(secondaryEndpoints.InternetEndpoints.File)
				output.secondaryWebInternetEndpoint, output.secondaryWebInternetHostName = flattenAccountEndpointAndHost(secondaryEndpoints.InternetEndpoints.Web)
			}

			if secondaryEndpoints.MicrosoftEndpoints != nil && pointer.From(routingPreference.PublishMicrosoftEndpoints) {
				output.secondaryBlobMicrosoftEndpoint, output.secondaryBlobMicrosoftHostName = flattenAccountEndpointAndHost(secondaryEndpoints.MicrosoftEndpoints.Blob)
				output.secondaryDfsMicrosoftEndpoint, output.secondaryDfsMicrosoftHostName = flattenAccountEndpointAndHost(secondaryEndpoints.MicrosoftEndpoints.Dfs)
				output.secondaryFileMicrosoftEndpoint, output.secondaryFileMicrosoftHostName = flattenAccountEndpointAndHost(secondaryEndpoints.MicrosoftEndpoints.File)
				output.secondaryQueueMicrosoftEndpoint, output.secondaryQueueMicrosoftHostName = flattenAccountEndpointAndHost(secondaryEndpoints.MicrosoftEndpoints.Queue)
				output.secondaryTableMicrosoftEndpoint, output.secondaryTableMicrosoftHostName = flattenAccountEndpointAndHost(secondaryEndpoints.MicrosoftEndpoints.Table)
				output.secondaryWebMicrosoftEndpoint, output.secondaryWebMicrosoftHostName = flattenAccountEndpointAndHost(secondaryEndpoints.MicrosoftEndpoints.Web)
			}
		}
	}

	return output
}

func flattenAccountEndpointAndHost(input *string) (string, string) {
	endpoint := ""
	host := ""
	if input != nil {
		endpoint = *input
		if u, _ := url.Parse(*input); u != nil {
			host = u.Host
		}
	}
	return endpoint, host
}

type accountAccessKeysAndConnectionStrings struct {
	primaryConnectionString       string
	secondaryConnectionString     string
	primaryBlobConnectionString   string
	secondaryBlobConnectionString string
	primaryAccessKey              string
	secondaryAccessKey            string
}

func (a accountAccessKeysAndConnectionStrings) set(d *pluginsdk.ResourceData) error {
	d.Set("primary_connection_string", a.primaryConnectionString)
	d.Set("secondary_connection_string", a.secondaryConnectionString)
	d.Set("primary_blob_connection_string", a.primaryBlobConnectionString)
	d.Set("secondary_blob_connection_string", a.secondaryBlobConnectionString)
	d.Set("primary_access_key", a.primaryAccessKey)
	d.Set("secondary_access_key", a.secondaryAccessKey)

	return nil
}

func flattenAccountAccessKeysAndConnectionStrings(accountName, domainSuffix string, keys []storageaccounts.StorageAccountKey, endpoints accountEndpoints) accountAccessKeysAndConnectionStrings {
	output := accountAccessKeysAndConnectionStrings{}

	// NOTE: users might not have access to list the keys, which is handled in the Data Source (optional) / Resource (required) respectively
	if len(keys) > 0 {
		output.primaryAccessKey = pointer.From(keys[0].Value)
		if len(keys) > 1 {
			output.secondaryAccessKey = pointer.From(keys[1].Value)
		}

		if output.primaryAccessKey != "" {
			output.primaryConnectionString = fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", accountName, output.primaryAccessKey, domainSuffix)

			if endpoints.primaryBlobEndpoint != "" {
				output.primaryBlobConnectionString = fmt.Sprintf("DefaultEndpointsProtocol=https;BlobEndpoint=%s;AccountName=%s;AccountKey=%s", endpoints.primaryBlobEndpoint, accountName, output.primaryAccessKey)
			}
		}

		if output.secondaryAccessKey != "" {
			output.secondaryConnectionString = fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=%s", accountName, output.secondaryAccessKey, domainSuffix)

			if endpoints.secondaryBlobEndpoint != "" {
				output.secondaryBlobConnectionString = fmt.Sprintf("DefaultEndpointsProtocol=https;BlobEndpoint=%s;AccountName=%s;AccountKey=%s", endpoints.secondaryBlobEndpoint, accountName, output.secondaryAccessKey)
			}
		}
	}

	return output
}
