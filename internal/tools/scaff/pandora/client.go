// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pandora

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// DefaultBaseURL is the address the Pandora Data API serves on locally.
const DefaultBaseURL = "http://localhost:8080"

const resourceManagerBasePath = "/v1/resource-manager/services"

// Client is a read-only client for the Pandora Data API.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient returns a Client pointed at the given base URL. If baseURL is empty
// the DefaultBaseURL is used.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

// get fetches the given path (which may be an absolute path beginning with "/"
// as returned in the API's *Uri fields, or a relative path) and decodes the
// JSON body into out.
func (c *Client) get(path string, out interface{}) error {
	url := path
	if strings.HasPrefix(path, "/") {
		url = c.BaseURL + path
	} else if !strings.HasPrefix(path, "http") {
		url = c.BaseURL + "/" + path
	}

	resp, err := c.HTTP.Get(url)
	if err != nil {
		return fmt.Errorf("requesting %q: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found: %q", url)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d requesting %q", resp.StatusCode, url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response from %q: %w", url, err)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decoding response from %q: %w", url, err)
	}

	return nil
}

// ListServices returns the map of available services keyed by service name.
func (c *Client) ListServices() (map[string]ServiceSummary, error) {
	var out ServicesResponse
	if err := c.get(resourceManagerBasePath, &out); err != nil {
		return nil, err
	}
	return out.Services, nil
}

// GetService returns the details (versions, resource provider) for a single
// service by its Pandora service name (e.g. "RedHatOpenShift").
func (c *Client) GetService(service string) (*ServiceResponse, error) {
	var out ServiceResponse
	if err := c.get(fmt.Sprintf("%s/%s", resourceManagerBasePath, service), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetVersion returns the resources available for a service at a specific API
// version (e.g. "2025-07-25").
func (c *Client) GetVersion(service, version string) (*VersionResponse, error) {
	var out VersionResponse
	if err := c.get(fmt.Sprintf("%s/%s/%s", resourceManagerBasePath, service, version), &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetResourceSchema returns the schema document (constants, models, resourceIds)
// for a resource within a service/version.
func (c *Client) GetResourceSchema(service, version, resource string) (*SchemaResponse, error) {
	var out SchemaResponse
	path := fmt.Sprintf("%s/%s/%s/%s/schema", resourceManagerBasePath, service, version, resource)
	if err := c.get(path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetResourceOperations returns the operations document for a resource within a
// service/version.
func (c *Client) GetResourceOperations(service, version, resource string) (*OperationsResponse, error) {
	var out OperationsResponse
	path := fmt.Sprintf("%s/%s/%s/%s/operations", resourceManagerBasePath, service, version, resource)
	if err := c.get(path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetSchemaByURI fetches a schema document from an absolute *Uri value returned
// by the API (e.g. ResourceSummary.SchemaURI).
func (c *Client) GetSchemaByURI(uri string) (*SchemaResponse, error) {
	var out SchemaResponse
	if err := c.get(uri, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetOperationsByURI fetches an operations document from an absolute *Uri value
// returned by the API (e.g. ResourceSummary.OperationsURI).
func (c *Client) GetOperationsByURI(uri string) (*OperationsResponse, error) {
	var out OperationsResponse
	if err := c.get(uri, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
