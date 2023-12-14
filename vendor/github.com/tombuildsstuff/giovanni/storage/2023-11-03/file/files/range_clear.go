package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ClearByteRangeInput struct {
	StartBytes int64
	EndBytes   int64
}

type ClearByteRangeResponse struct {
	HttpResponse *client.Response
}

// ClearByteRange clears the specified Byte Range from within the specified File
func (c Client) ClearByteRange(ctx context.Context, shareName, path, fileName string, input ClearByteRangeInput) (resp ClearByteRangeResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if fileName == "" {
		return resp, fmt.Errorf("`fileName` cannot be an empty string")
	}

	if input.StartBytes < 0 {
		return resp, fmt.Errorf("`input.StartBytes` must be greater or equal to 0")
	}

	if input.EndBytes <= 0 {
		return resp, fmt.Errorf("`input.EndBytes` must be greater than 0")
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: ClearByteRangeOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s%s", shareName, path, fileName),
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

type ClearByteRangeOptions struct {
	input ClearByteRangeInput
}

func (c ClearByteRangeOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-write", "clear")
	headers.Append("x-ms-range", fmt.Sprintf("bytes=%d-%d", c.input.StartBytes, c.input.EndBytes))
	return headers
}

func (c ClearByteRangeOptions) ToOData() *odata.Query {
	return nil
}

func (c ClearByteRangeOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "range")
	return out
}
