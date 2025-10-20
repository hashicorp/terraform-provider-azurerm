package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type CopyAbortInput struct {
	copyID string
}

type CopyAbortResponse struct {
	HttpResponse *http.Response
}

// AbortCopy aborts a pending Copy File operation, and leaves a destination file with zero length and full metadata
func (c Client) AbortCopy(ctx context.Context, shareName, path, fileName string, input CopyAbortInput) (result CopyAbortResponse, err error) {

	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if fileName == "" {
		err = fmt.Errorf("`fileName` cannot be an empty string")
		return
	}

	if input.copyID == "" {
		err = fmt.Errorf("`copyID` cannot be an empty string")
		return
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: CopyAbortOptions{
			copyId: input.copyID,
		},
		Path: fmt.Sprintf("/%s/%s%s", shareName, path, fileName),
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

type CopyAbortOptions struct {
	copyId string
}

func (c CopyAbortOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-copy-action", "abort")
	return headers
}

func (c CopyAbortOptions) ToOData() *odata.Query {
	return nil
}

func (c CopyAbortOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "copy")
	out.Append("copyid", c.copyId)
	return out
}
