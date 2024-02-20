package containers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetPropertiesInput struct {
	LeaseId string
}

type GetPropertiesResponse struct {
	HttpResponse *client.Response
	Model        *ContainerProperties
}

// GetProperties returns the properties for this Container without a Lease
func (c Client) GetProperties(ctx context.Context, containerName string, input GetPropertiesInput) (resp GetPropertiesResponse, err error) {
	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
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
	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	if resp.HttpResponse != nil {
		resp.Model = &ContainerProperties{}
		resp.Model.LeaseStatus = LeaseStatus(resp.HttpResponse.Header.Get("x-ms-lease-status"))
		resp.Model.LeaseState = LeaseState(resp.HttpResponse.Header.Get("x-ms-lease-state"))
		if resp.Model.LeaseStatus == Locked {
			duration := LeaseDuration(resp.HttpResponse.Header.Get("x-ms-lease-duration"))
			resp.Model.LeaseDuration = &duration
		}

		// If this header is not returned in the response, the container is private to the account owner.
		accessLevel := resp.HttpResponse.Header.Get("x-ms-blob-public-access")
		if accessLevel != "" {
			resp.Model.AccessLevel = AccessLevel(accessLevel)
		} else {
			resp.Model.AccessLevel = Private
		}

		// we can't necessarily use strconv.ParseBool here since this could be nil (only in some API versions)
		resp.Model.HasImmutabilityPolicy = strings.EqualFold(resp.HttpResponse.Header.Get("x-ms-has-immutability-policy"), "true")
		resp.Model.HasLegalHold = strings.EqualFold(resp.HttpResponse.Header.Get("x-ms-has-legal-hold"), "true")
		resp.Model.MetaData = metadata.ParseFromHeaders(resp.HttpResponse.Header)
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
