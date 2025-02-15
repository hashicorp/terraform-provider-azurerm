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

type CreateSnapshotInput struct {
	MetaData map[string]string
}

type CreateSnapshotResponse struct {
	HttpResponse *http.Response

	// This header is a DateTime value that uniquely identifies the share snapshot.
	// The value of this header may be used in subsequent requests to access the share snapshot.
	// This value is opaque.
	SnapshotDateTime string
}

// CreateSnapshot creates a read-only snapshot of the share
// A share can support creation of 200 share snapshots. Attempting to create more than 200 share snapshots fails with 409 (Conflict).
// Attempting to create a share snapshot while a previous Snapshot Share operation is in progress fails with 409 (Conflict).
func (c Client) CreateSnapshot(ctx context.Context, shareName string, input CreateSnapshotInput) (result CreateSnapshotResponse, err error) {

	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf("`input.MetaData` is not valid: %+v", err)
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: snapShotCreateOptions{
			metaData: input.MetaData,
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
				result.SnapshotDateTime = resp.Header.Get("x-ms-snapshot")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type snapShotCreateOptions struct {
	metaData map[string]string
}

func (s snapShotCreateOptions) ToHeaders() *client.Headers {
	headers := metadata.SetMetaDataHeaders(s.metaData)
	return &headers
}

func (s snapShotCreateOptions) ToOData() *odata.Query {
	return nil
}

func (s snapShotCreateOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "share")
	out.Append("comp", "snapshot")
	return out
}
