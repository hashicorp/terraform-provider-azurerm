package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ReleaseLeaseInput struct {
	LeaseId string
}

type ReleaseLeaseResponse struct {
	HttpResponse *client.Response
}

// ReleaseLease releases the lock based on the Lease ID
func (c Client) ReleaseLease(ctx context.Context, containerName string, input ReleaseLeaseInput) (resp ReleaseLeaseResponse, err error) {
	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}
	if input.LeaseId == "" {
		return resp, fmt.Errorf("`input.LeaseId` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: releaseLeaseOptions{
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

	return
}

var _ client.Options = releaseLeaseOptions{}

type releaseLeaseOptions struct {
	leaseId string
}

func (o releaseLeaseOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	headers.Append("x-ms-lease-action", "release")
	headers.Append("x-ms-lease-id", o.leaseId)

	return headers
}

func (o releaseLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (o releaseLeaseOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "lease")
	return query
}
