package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type RenewLeaseResponse struct {
	HttpResponse *client.Response
}

type RenewLeaseInput struct {
	LeaseID string
}

func (c Client) RenewLease(ctx context.Context, containerName, blobName string, input RenewLeaseInput) (resp RenewLeaseResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.LeaseID == "" {
		return resp, fmt.Errorf("`input.LeaseID` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: renewLeaseOptions{
			leaseID: input.LeaseID,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
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

	return
}

type renewLeaseOptions struct {
	leaseID string
}

func (r renewLeaseOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-lease-action", "renew")
	headers.Append("x-ms-lease-id", r.leaseID)
	return headers
}

func (r renewLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (r renewLeaseOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "lease")
	return out
}
