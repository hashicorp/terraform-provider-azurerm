package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SupportedLocations struct {
	// Locations is a list of Locations which are supported on this Azure Endpoint.
	// This could be nil when the user is offline, or the Azure MetaData Service does not have this
	// information and as such this should be used as best-effort, rather than guaranteed
	Locations *[]string
}

type cloudEndpoint struct {
	Endpoint  string    `json:"endpoint"`
	Locations *[]string `json:"locations"`
}

type metaDataResponse struct {
	CloudEndpoint map[string]cloudEndpoint `json:"cloudEndpoint"`
}

// AvailableAzureLocations returns a list of the Azure Locations which are available on the specified endpoint
func AvailableAzureLocations(ctx context.Context, endpoint string) (*SupportedLocations, error) {
	uri := fmt.Sprintf("https://%s//metadata/endpoints?api-version=2018-01-01", endpoint)
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("retrieving supported locations from Azure MetaData service: %+v", err)
	}
	var out metaDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("deserializing JSON from Azure MetaData service: %+v", err)
	}

	var locations *[]string
	for _, v := range out.CloudEndpoint {
		// one of the endpoints on this endpoint should reference itself
		// however this is best-effort, so if it doesn't, it's not the end of the world
		if strings.EqualFold(v.Endpoint, endpoint) {
			locations = v.Locations
		}
	}

	return &SupportedLocations{
		Locations: locations,
	}, nil
}
