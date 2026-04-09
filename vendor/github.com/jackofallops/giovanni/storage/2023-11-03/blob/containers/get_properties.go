package containers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type GetPropertiesInput struct {
	LeaseId string
}

type GetPropertiesResponse struct {
	ContainerProperties
	HttpResponse *http.Response
}

// GetProperties returns the properties for this Container without a Lease
func (c Client) GetProperties(ctx context.Context, containerName string, input GetPropertiesInput) (result GetPropertiesResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: getPropertiesOptions{
			leaseId: input.LeaseId,
		},
		Path: fmt.Sprintf("/%s", containerName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.DefaultEncryptionScope = resp.Header.Get("x-ms-default-encryption-scope")
				result.LeaseStatus = LeaseStatus(resp.Header.Get("x-ms-lease-status"))
				result.LeaseState = LeaseState(resp.Header.Get("x-ms-lease-state"))
				if result.LeaseStatus == Locked {
					duration := LeaseDuration(resp.Header.Get("x-ms-lease-duration"))
					result.LeaseDuration = &duration
				}

				// If this header is not returned in the response, the container is private to the account owner.
				accessLevel := resp.Header.Get("x-ms-blob-public-access")
				if accessLevel != "" {
					result.AccessLevel = AccessLevel(accessLevel)
				} else {
					result.AccessLevel = Private
				}

				// we can't necessarily use strconv.ParseBool here since this could be nil (only in some API versions)
				result.EncryptionScopeOverrideDisabled = strings.EqualFold(resp.Header.Get("x-ms-deny-encryption-scope-override"), "true")
				result.HasImmutabilityPolicy = strings.EqualFold(resp.Header.Get("x-ms-has-immutability-policy"), "true")
				result.HasLegalHold = strings.EqualFold(resp.Header.Get("x-ms-has-legal-hold"), "true")

				result.MetaData = metadata.ParseFromHeaders(resp.Header)
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

var _ client.Options = getPropertiesOptions{}

type getPropertiesOptions struct {
	leaseId string
}

func (o getPropertiesOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	// If specified, Get Container Properties only succeeds if the container’s lease is active and matches this ID.
	// If there is no active lease or the ID does not match, 412 (Precondition Failed) is returned.
	if o.leaseId != "" {
		headers.Append("x-ms-lease-id", o.leaseId)
	}

	return headers
}

func (getPropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (getPropertiesOptions) ToQuery() *client.QueryParams {
	return containerOptions{}.ToQuery()
}
