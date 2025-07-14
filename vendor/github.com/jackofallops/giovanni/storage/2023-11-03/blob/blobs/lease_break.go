package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
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
	HttpResponse *http.Response

	// Approximate time remaining in the lease period, in seconds.
	// If the break is immediate, 0 is returned.
	LeaseTime int
}

// BreakLease breaks an existing lock on a blob using the LeaseID.
func (c Client) BreakLease(ctx context.Context, containerName, blobName string, input BreakLeaseInput) (result BreakLeaseResponse, err error) {
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
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: breakLeaseOptions{
			input: input,
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

type breakLeaseOptions struct {
	input BreakLeaseInput
}

func (b breakLeaseOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-lease-action", "break")
	headers.Append("x-ms-lease-id", b.input.LeaseID)

	if b.input.BreakPeriod != nil {
		headers.Append("x-ms-lease-break-period", strconv.Itoa(*b.input.BreakPeriod))
	}

	return headers
}

func (b breakLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (b breakLeaseOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "lease")
	return out
}
