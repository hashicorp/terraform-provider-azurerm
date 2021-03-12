package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/go-autorest/autorest/azure"
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

// availableAzureLocations returns a list of the Azure Locations which are available on the specified endpoint
func availableAzureLocations(ctx context.Context, env *azure.Environment) (*SupportedLocations, error) {
	// e.g. https://management.azure.com/ but we need management.azure.com
	endpoint := strings.TrimPrefix(env.ResourceManagerEndpoint, "https://")
	endpoint = strings.TrimSuffix(endpoint, "/")

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

	// TODO: remove this once Microsoft fixes the API
	// the Azure API returns the india locations the wrong way around
	// e.g. 'southindia' is returned as 'indiasouth'
	// so we need to conditionally switch these out until Microsoft fixes the API
	// $ az account list-locations -o table | grep india
	//  Central India             centralindia         (Asia Pacific) Central India
	//  South India               southindia           (Asia Pacific) South India
	//  West India                westindia            (Asia Pacific) West India
	if env.Name == azure.PublicCloud.Name && locations != nil {
		out := *locations
		out = switchLocationIfExists("indiacentral", "centralindia", out)
		out = switchLocationIfExists("indiasouth", "southindia", out)
		out = switchLocationIfExists("indiawest", "westindia", out)
		locations = &out
	}

	return &SupportedLocations{
		Locations: locations,
	}, nil
}

func switchLocationIfExists(find, replace string, locations []string) []string {
	out := locations

	for i, v := range out {
		if v == find {
			out[i] = replace
		}
	}

	return locations
}
