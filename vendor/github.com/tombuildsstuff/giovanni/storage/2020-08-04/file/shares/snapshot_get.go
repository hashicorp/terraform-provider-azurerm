package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type GetSnapshotPropertiesResponse struct {
	HttpResponse *client.Response

	MetaData map[string]string
}

type GetSnapshotPropertiesInput struct {
	snapshotShare string
}

// GetSnapshot gets information about the specified Snapshot of the specified Storage Share
func (c Client) GetSnapshot(ctx context.Context, shareName string, input GetSnapshotPropertiesInput) (resp GetSnapshotPropertiesResponse, err error) {
	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if input.snapshotShare == "" {
		return resp, fmt.Errorf("`snapshotShare` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: snapShotGetOptions{
			snapshotShare: input.snapshotShare,
		},
		Path: fmt.Sprintf("/%s", shareName),
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
		if resp.HttpResponse.Header != nil {
			resp.MetaData = metadata.ParseFromHeaders(resp.HttpResponse.Header)
		}
	}

	return
}

type snapShotGetOptions struct {
	snapshotShare string
}

func (s snapShotGetOptions) ToHeaders() *client.Headers {
	return nil
}

func (s snapShotGetOptions) ToOData() *odata.Query {
	return nil
}

func (s snapShotGetOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("snapshot", s.snapshotShare)
	return out
}
