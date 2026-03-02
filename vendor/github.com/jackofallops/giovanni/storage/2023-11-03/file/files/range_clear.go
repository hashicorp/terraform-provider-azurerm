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
	HttpResponse *http.Response
}

// ClearByteRange clears the specified Byte Range from within the specified File
func (c Client) ClearByteRange(ctx context.Context, shareName, path, fileName string, input ClearByteRangeInput) (result ClearByteRangeResponse, err error) {

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

	if input.StartBytes < 0 {
		err = fmt.Errorf("`input.StartBytes` must be greater or equal to 0")
		return
	}

	if input.EndBytes <= 0 {
		err = fmt.Errorf("`input.EndBytes` must be greater than 0")
		return
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
