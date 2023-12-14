package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ChangeLeaseInput struct {
	ExistingLeaseID string
	ProposedLeaseID string
}

type ChangeLeaseResponse struct {
	HttpResponse *client.Response
	Model        *ChangeLeaseModel
}

type ChangeLeaseModel struct {
	LeaseID string
}

// ChangeLease changes the lock from one Lease ID to another Lease ID
func (c Client) ChangeLease(ctx context.Context, containerName string, input ChangeLeaseInput) (resp ChangeLeaseResponse, err error) {
	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}
	if input.ExistingLeaseID == "" {
		return resp, fmt.Errorf("`input.ExistingLeaseID` cannot be an empty string")
	}
	if input.ProposedLeaseID == "" {
		return resp, fmt.Errorf("`input.ProposedLeaseID` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: changeLeaseOptions{
			existingLeaseId: input.ExistingLeaseID,
			proposedLeaseId: input.ProposedLeaseID,
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
		resp.Model = &ChangeLeaseModel{
			LeaseID: resp.HttpResponse.Header.Get("x-ms-lease-id"),
		}
	}

	return
}

var _ client.Options = changeLeaseOptions{}

type changeLeaseOptions struct {
	existingLeaseId string
	proposedLeaseId string
}

func (o changeLeaseOptions) ToHeaders() *client.Headers {
	headers := containerOptions{}.ToHeaders()

	headers.Append("x-ms-lease-action", "change")
	headers.Append("x-ms-lease-id", o.existingLeaseId)
	headers.Append("x-ms-proposed-lease-id", o.proposedLeaseId)

	return headers
}

func (o changeLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (o changeLeaseOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "lease")
	return query
}
