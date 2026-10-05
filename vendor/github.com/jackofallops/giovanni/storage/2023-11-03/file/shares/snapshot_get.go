package shares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type GetSnapshotPropertiesResponse struct {
	HttpResponse *http.Response

	MetaData map[string]string
}

type GetSnapshotPropertiesInput struct {
	snapshotShare string
}

// GetSnapshot gets information about the specified Snapshot of the specified Storage Share
func (c Client) GetSnapshot(ctx context.Context, shareName string, input GetSnapshotPropertiesInput) (result GetSnapshotPropertiesResponse, err error) {
	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if input.snapshotShare == "" {
		err = fmt.Errorf("`snapshotShare` cannot be an empty string")
		return
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.MetaData = metadata.ParseFromHeaders(resp.Header)
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
