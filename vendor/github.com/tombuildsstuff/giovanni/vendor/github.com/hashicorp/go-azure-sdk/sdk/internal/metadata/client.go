// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package metadata

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"
	"time"
)

// NOTE: this Client cannot use the base client since it'd cause a circular reference

type Client struct {
	endpoint string
}

func NewClientWithEndpoint(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

// GetMetaData connects to the ARM metadata service at the configured endpoint, to retrieve information about the
// current environment. Sometimes an endpoint will not support the latest schema, in such cases it will not be
// possible to configure all services but a best effort will be made to request and parse an earlier schema version.
// `name` is used when falling back to an earlier schema version where multiple environments are returned and the
// desired one must be matched by name.
func (c *Client) GetMetaData(ctx context.Context, name string) (*MetaData, error) {
	metadata, err := c.getMetaDataFrom2022API(ctx, name)
	if err != nil {
		log.Printf("[DEBUG] Falling back to ARM Metadata version 2019-05-01 for %s", c.endpoint)
		metadata, err = c.getMetaDataFrom2019API(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("retrieving metadata from the 2022-09-01 and 2019-05-01 APIs: %+v", err)
		}
	}

	return &MetaData{
		Authentication: Authentication{
			Audiences:        metadata.Authentication.Audiences,
			LoginEndpoint:    metadata.Authentication.LoginEndpoint,
			IdentityProvider: metadata.Authentication.IdentityProvider,
			Tenant:           metadata.Authentication.Tenant,
		},
		DnsSuffixes: DnsSuffixes{
			Attestation: metadata.Suffixes.AttestationEndpoint,
			FrontDoor:   metadata.Suffixes.AzureFrontDoorEndpointSuffix,
			KeyVault:    metadata.Suffixes.KeyVaultDns,
			ManagedHSM:  metadata.Suffixes.MhsmDns,
			MariaDB:     metadata.Suffixes.MariadbServerEndpoint,
			MySql:       metadata.Suffixes.MysqlServerEndpoint,
			Postgresql:  metadata.Suffixes.PostgresqlServerEndpoint,
			SqlServer:   metadata.Suffixes.SqlServerHostname,
			Storage:     metadata.Suffixes.Storage,
			StorageSync: metadata.Suffixes.StorageSyncEndpointSuffix,
			Synapse:     metadata.Suffixes.SynapseAnalytics,
		},
		Name: metadata.Name,
		ResourceIdentifiers: ResourceIdentifiers{
			Attestation:    normalizeResourceId(metadata.AttestationResourceId),
			Batch:          normalizeResourceId(metadata.Batch),
			LogAnalytics:   normalizeResourceId(metadata.LogAnalyticsResourceId),
			Media:          normalizeResourceId(metadata.Media),
			MicrosoftGraph: normalizeResourceId(metadata.MicrosoftGraphResourceId),
			OSSRDBMS:       normalizeResourceId(metadata.OssrDbmsResourceId),
			Synapse:        normalizeResourceId(metadata.SynapseAnalyticsResourceId),
		},
		ResourceManagerEndpoint: metadata.ResourceManager,
	}, nil
}

func (c *Client) getMetaDataFrom2022API(ctx context.Context, name string) (*metaDataResponse, error) {
	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := &net.Dialer{Resolver: &net.Resolver{}}
				return d.DialContext(ctx, network, addr)
			},
			TLSClientConfig:       &tlsConfig,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
	uri := fmt.Sprintf("%s/metadata/endpoints?api-version=2022-09-01", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %+v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("performing request: expected 200 OK but got %d %s", resp.StatusCode, resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %+v", err)
	}
	resp.Body.Close()

	// Trim away a BOM if present
	respBody = bytes.TrimPrefix(respBody, []byte("\xef\xbb\xbf"))

	var model *metaDataResponse
	if err := json.Unmarshal(respBody, &model); err != nil {
		log.Printf("[DEBUG] Unrecognised metadata response for %s: %s", uri, respBody)
		return nil, fmt.Errorf("unmarshaling response: %+v", err)
	}

	return model, nil
}

func (c *Client) getMetaDataFrom2019API(ctx context.Context, name string) (*metaDataResponse, error) {
	tlsConfig := tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				d := &net.Dialer{Resolver: &net.Resolver{}}
				return d.DialContext(ctx, network, addr)
			},
			TLSClientConfig:       &tlsConfig,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			ForceAttemptHTTP2:     true,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		},
	}
	uri := fmt.Sprintf("%s/metadata/endpoints?api-version=2019-05-01", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %+v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing request: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("performing request: expected 200 OK but got %d %s", resp.StatusCode, resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing response body: %+v", err)
	}
	resp.Body.Close()

	// Trim away a BOM if present
	respBody = bytes.TrimPrefix(respBody, []byte("\xef\xbb\xbf"))

	var model *[]metaDataResponse
	if err := json.Unmarshal(respBody, &model); err != nil {
		log.Printf("[DEBUG] Unrecognised metadata response for %s: %s", uri, respBody)
		return nil, fmt.Errorf("unmarshaling response: %+v", err)
	}

	if model == nil {
		return nil, fmt.Errorf("unmarshaling response: no environments returned")
	}

	// This version returns an array of environments, we are only interested in one
	var env metaDataResponse
	for _, e := range *model {
		if name == "" || e.Name == name {
			env = e
			break
		}
	}
	return &env, nil
}

type metaDataResponse struct {
	Portal         string `json:"portal"`
	Authentication struct {
		LoginEndpoint    string   `json:"loginEndpoint"`
		Audiences        []string `json:"audiences"`
		Tenant           string   `json:"tenant"`
		IdentityProvider string   `json:"identityProvider"`
	} `json:"authentication"`
	Media         string `json:"media"`
	GraphAudience string `json:"graphAudience"`
	Graph         string `json:"graph"`
	Name          string `json:"name"`
	Suffixes      struct {
		AzureDataLakeStoreFileSystem        string `json:"azureDataLakeStoreFileSystem"`
		AcrLoginServer                      string `json:"acrLoginServer"`
		SqlServerHostname                   string `json:"sqlServerHostname"`
		AzureDataLakeAnalyticsCatalogAndJob string `json:"azureDataLakeAnalyticsCatalogAndJob"`
		KeyVaultDns                         string `json:"keyVaultDns"`
		Storage                             string `json:"storage"`
		AzureFrontDoorEndpointSuffix        string `json:"azureFrontDoorEndpointSuffix"`
		StorageSyncEndpointSuffix           string `json:"storageSyncEndpointSuffix"`
		MhsmDns                             string `json:"mhsmDns"`
		MysqlServerEndpoint                 string `json:"mysqlServerEndpoint"`
		PostgresqlServerEndpoint            string `json:"postgresqlServerEndpoint"`
		MariadbServerEndpoint               string `json:"mariadbServerEndpoint"`
		SynapseAnalytics                    string `json:"synapseAnalytics"`
		AttestationEndpoint                 string `json:"attestationEndpoint"`
	} `json:"suffixes"`
	Batch                                 string `json:"batch"`
	ResourceManager                       string `json:"resourceManager"`
	VmImageAliasDoc                       string `json:"vmImageAliasDoc"`
	ActiveDirectoryDataLake               string `json:"activeDirectoryDataLake"`
	SqlManagement                         string `json:"sqlManagement"`
	MicrosoftGraphResourceId              string `json:"microsoftGraphResourceId"`
	AppInsightsResourceId                 string `json:"appInsightsResourceId"`
	AppInsightsTelemetryChannelResourceId string `json:"appInsightsTelemetryChannelResourceId"`
	AttestationResourceId                 string `json:"attestationResourceId"`
	SynapseAnalyticsResourceId            string `json:"synapseAnalyticsResourceId"`
	LogAnalyticsResourceId                string `json:"logAnalyticsResourceId"`
	OssrDbmsResourceId                    string `json:"ossrDbmsResourceId"`
	Gallery                               string `json:"gallery"`
}
