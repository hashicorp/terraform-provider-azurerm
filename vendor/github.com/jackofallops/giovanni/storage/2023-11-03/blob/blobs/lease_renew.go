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
	HttpResponse *http.Response
}

type RenewLeaseInput struct {
	LeaseID string
}

func (c Client) RenewLease(ctx context.Context, containerName, blobName string, input RenewLeaseInput) (result RenewLeaseResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}

	if strings.ToLower(containerName) != containerName {
		err = fmt.Errorf("`containerName` must be a lower-cased string")
		return
	}

	if blobName == "" {
		err = fmt.Errorf("`blobName` cannot be an empty string")
		return
	}

	if input.LeaseID == "" {
		err = fmt.Errorf("`input.LeaseID` cannot be an empty string")
		return
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
