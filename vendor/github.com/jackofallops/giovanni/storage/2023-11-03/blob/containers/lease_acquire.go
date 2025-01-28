package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type AcquireLeaseInput struct {
	// Specifies the duration of the lease, in seconds, or negative one (-1) for a lease that never expires.
	// A non-infinite lease can be between 15 and 60 seconds
	LeaseDuration int

	ProposedLeaseID string
}

type AcquireLeaseResponse struct {
	AcquireLeaseModel
	HttpResponse *http.Response
}

type AcquireLeaseModel struct {
	LeaseID string
}

// AcquireLease establishes and manages a lock on a container for delete operations.
func (c Client) AcquireLease(ctx context.Context, containerName string, input AcquireLeaseInput) (result AcquireLeaseResponse, err error) {
	if containerName == "" {
		return result, fmt.Errorf("`containerName` cannot be an empty string")
	}
	// An infinite lease duration is -1 seconds. A non-infinite lease can be between 15 and 60 seconds
	if input.LeaseDuration != -1 && (input.LeaseDuration <= 15 || input.LeaseDuration >= 60) {
		return result, fmt.Errorf("`input.LeaseDuration` must be -1 (infinite), or between 15 and 60 seconds")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: acquireLeaseOptions{
			leaseDuration:   input.LeaseDuration,
			proposedLeaseId: input.ProposedLeaseID,
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
				result.LeaseID = resp.Header.Get("x-ms-lease-id")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

var _ client.Options = acquireLeaseOptions{}

type acquireLeaseOptions struct {
	leaseDuration   int
	proposedLeaseId string
}

func (o acquireLeaseOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	headers.Append("x-ms-lease-action", "acquire")
	headers.Append("x-ms-lease-duration", fmt.Sprintf("%d", o.leaseDuration))

	if o.proposedLeaseId != "" {
		headers.Append("x-ms-proposed-lease-id", o.proposedLeaseId)
	}

	return headers
}

func (o acquireLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (o acquireLeaseOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "lease")
	return query
}
