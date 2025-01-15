package containers

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"net/http"
	"strconv"
)

type BreakLeaseInput struct {
	//  For a break operation, proposed duration the lease should continue
	//  before it is broken, in seconds, between 0 and 60.
	//  This break period is only used if it is shorter than the time remaining on the lease.
	//  If longer, the time remaining on the lease is used.
	//  A new lease will not be available before the break period has expired,
	//  but the lease may be held for longer than the break period.
	//  If this header does not appear with a break operation, a fixed-duration lease breaks
	//  after the remaining lease period elapses, and an infinite lease breaks immediately.
	BreakPeriod *int

	LeaseID string
}

type BreakLeaseResponse struct {
	BreakLeaseModel
	HttpResponse *http.Response
}

type BreakLeaseModel struct {
	// Approximate time remaining in the lease period, in seconds.
	// If the break is immediate, 0 is returned.
	LeaseTime int
}

// BreakLease breaks a lock based on it's Lease ID
func (c Client) BreakLease(ctx context.Context, containerName string, input BreakLeaseInput) (result BreakLeaseResponse, err error) {
	if containerName == "" {
		return result, fmt.Errorf("`containerName` cannot be an empty string")
	}
	if input.LeaseID == "" {
		return result, fmt.Errorf("`input.LeaseID` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: breakLeaseOptions{
			breakPeriod: input.BreakPeriod,
			leaseId:     input.LeaseID,
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
				if leaseTimeRaw := resp.Header.Get("x-ms-lease-time"); leaseTimeRaw != "" {
					if leaseTime, err := strconv.Atoi(leaseTimeRaw); err == nil {
						result.LeaseTime = leaseTime
					}
				}
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

var _ client.Options = breakLeaseOptions{}

type breakLeaseOptions struct {
	breakPeriod *int
	leaseId     string
}

func (o breakLeaseOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	headers.Append("x-ms-lease-action", "break")
	headers.Append("x-ms-lease-id", o.leaseId)

	if o.breakPeriod != nil {
		headers.Append("x-ms-lease-break-period", fmt.Sprintf("%d", *o.breakPeriod))
	}

	return headers
}

func (o breakLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (o breakLeaseOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "lease")
	return query
}
