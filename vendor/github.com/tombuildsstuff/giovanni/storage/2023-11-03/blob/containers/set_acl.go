package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type SetAccessControlInput struct {
	AccessLevel AccessLevel
	LeaseId     string
}

type SetAccessControlResponse struct {
	HttpResponse *http.Response
}

// SetAccessControl sets the Access Control for a Container without a Lease ID
// NOTE: The SetAccessControl operation only supports Shared Key authorization.
func (c Client) SetAccessControl(ctx context.Context, containerName string, input SetAccessControlInput) (result SetAccessControlResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: setAccessControlListOptions{
			accessLevel: input.AccessLevel,
			leaseId:     input.LeaseId,
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
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

var _ client.Options = setAccessControlListOptions{}

type setAccessControlListOptions struct {
	accessLevel AccessLevel
	leaseId     string
}

func (o setAccessControlListOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	// If this header is not included in the request, container data is private to the account owner.
	if o.accessLevel != Private {
		headers.Append("x-ms-blob-public-access", string(o.accessLevel))
	}

	// If specified, Get Container Properties only succeeds if the container’s lease is active and matches this ID.
	// If there is no active lease or the ID does not match, 412 (Precondition Failed) is returned.
	if o.leaseId != "" {
		headers.Append("x-ms-lease-id", o.leaseId)
	}

	return headers
}

func (o setAccessControlListOptions) ToOData() *odata.Query {
	return nil
}

func (o setAccessControlListOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "acl")
	return query
}
