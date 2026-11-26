// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package preflight_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight"
	preflightclient "github.com/hashicorp/terraform-provider-azurerm/internal/preflight/client"
	preflightsdk "github.com/hashicorp/terraform-provider-azurerm/internal/preflight/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/preflight/testdata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"golang.org/x/oauth2"
)

func TestAccPreflight_validateResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skipf("Acceptance tests skipped unless env 'TF_ACC' set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	p := provider.AzureProvider()
	d := p.Configure(ctx, terraform.NewResourceConfigRaw(nil))
	if d.HasError() {
		t.Fatalf("diags had error!")
	}
	client := p.Meta().(*clients.Client)

	cases := []struct {
		Name        string
		Location    string
		ID          resourceids.ResourceId
		APIVersion  string
		Properties  any
		ExpectError bool
	}{
		{
			Name:        "valid_webapp",
			Location:    "westeurope",
			ID:          pointer.To(commonids.NewAppServiceID(client.Account.SubscriptionId, "testResourceGroup", "testSiteName")),
			APIVersion:  "2024-11-01",
			Properties:  testdata.WebAppExample(client.Account.SubscriptionId),
			ExpectError: false,
		},
		{
			Name:        "invalid_vnet_missing_required_property",
			Location:    "westeurope",
			ID:          pointer.To(commonids.NewVirtualNetworkID(client.Account.SubscriptionId, "testResourceGroup", "testVirtualValid")),
			APIVersion:  "2024-01-01",
			Properties:  testdata.VirtualNetworkExampleMissingRequiredProperty(client.Account.SubscriptionId),
			ExpectError: true,
		},
	}

	metadata := sdk.ResourceMetaData{
		Client: client,
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			request, err := preflight.NewValidationRequest(pointer.To(tc.Location), tc.ID, tc.APIVersion, tc.Properties)
			if err != nil {
				t.Fatalf("building validation request: %s", err)
			}

			err = request.ValidateResource(ctx, metadata)
			if err != nil && !tc.ExpectError {
				t.Fatalf("expected no error, got: %s", err)
			}
			if err == nil && tc.ExpectError {
				t.Fatalf("expected an error but got none")
			}
		})
	}
}

type stubAuthorizer struct{}

func (stubAuthorizer) Token(_ context.Context, _ *http.Request) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}, nil
}

func (stubAuthorizer) AuxiliaryTokens(_ context.Context, _ *http.Request) ([]*oauth2.Token, error) {
	return nil, nil
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type trackingReadCloser struct {
	io.Reader
	closed *bool
}

func (t *trackingReadCloser) Close() error {
	*t.closed = true
	return nil
}

func TestValidateResource_extractsErrorBody(t *testing.T) {
	originalBodyClosed := false
	transport := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		bodyJSON := `{
			"error": {
				"code": "InvalidResourceLocation",
				"message": "The resource location 'invalidregion' is not valid.",
				"details": [
					{
						"code": "InvalidLocationDetail",
						"target": "location",
						"message": "Location 'invalidregion' is not supported for this resource type."
					}
				]
			}
		}`
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Status:     "400 Bad Request",
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: &trackingReadCloser{
				Reader: strings.NewReader(bodyJSON),
				closed: &originalBodyClosed,
			},
			Request: req,
		}, nil
	})

	preflightClient, err := preflightsdk.NewResourceValidationClientClientWithBaseURI(environments.ResourceManagerAPI("https://management.azure.com"))
	if err != nil {
		t.Fatalf("building preflight client: %+v", err)
	}
	preflightClient.Client.SetTransport(transport)
	preflightClient.Client.SetAuthorizer(stubAuthorizer{})

	metadata := sdk.ResourceMetaData{
		Client: &clients.Client{
			Preflight: &preflightclient.Client{
				PreflightClient: preflightClient,
			},
		},
	}

	id := pointer.To(commonids.NewVirtualNetworkID("00000000-0000-0000-0000-000000000000", "testResourceGroup", "testVnet"))
	request, err := preflight.NewValidationRequest(pointer.To("westeurope"), id, "2024-01-01", map[string]any{})
	if err != nil {
		t.Fatalf("building validation request: %+v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = request.ValidateResource(ctx, metadata)
	if err == nil {
		t.Fatalf("expected error from ValidateResource, got nil")
	}
	if !originalBodyClosed {
		t.Fatalf("expected underlying HTTP response body to be drained and closed by Execute()")
	}

	expectedSubstrings := []string{
		"InvalidResourceLocation",
		"The resource location 'invalidregion' is not valid.",
		"location: Location 'invalidregion' is not supported for this resource type.",
	}
	for _, want := range expectedSubstrings {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("expected error %q to contain %q", err.Error(), want)
		}
	}
}
