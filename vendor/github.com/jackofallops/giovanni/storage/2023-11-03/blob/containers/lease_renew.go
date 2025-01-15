package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type RenewLeaseInput struct {
	LeaseId string
}

type RenewLeaseResponse struct {
	HttpResponse *http.Response
}

// RenewLease renews the lock based on the Lease ID
func (c Client) RenewLease(ctx context.Context, containerName string, input RenewLeaseInput) (result RenewLeaseResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}
	if input.LeaseId == "" {
		err = fmt.Errorf("`input.LeaseId` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: renewLeaseOptions{
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
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

var _ client.Options = renewLeaseOptions{}

type renewLeaseOptions struct {
	leaseId string
}

func (o renewLeaseOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	headers.Append("x-ms-lease-action", "renew")
	headers.Append("x-ms-lease-id", o.leaseId)

	return headers
}

func (o renewLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (o renewLeaseOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "lease")
	return query
}
