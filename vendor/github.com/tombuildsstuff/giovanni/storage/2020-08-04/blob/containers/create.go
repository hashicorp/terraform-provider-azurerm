package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type CreateInput struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	AccessLevel AccessLevel

	// A name-value pair to associate with the container as metadata.
	MetaData map[string]string
}

type CreateResponse struct {
	HttpResponse *client.Response
	Error        *ErrorResponse `xml:"Error"`
}

// Create creates a new container under the specified account.
// If the container with the same name already exists, the operation fails.
func (c Client) Create(ctx context.Context, containerName string, input CreateInput) (resp CreateResponse, err error) {
	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}
	if err := metadata.Validate(input.MetaData); err != nil {
		return resp, fmt.Errorf("`input.MetaData` is not valid: %+v", err)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: createOptions{
			accessLevel: input.AccessLevel,
			metaData:    input.MetaData,
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

	return
}

var _ client.Options = createOptions{}

type createOptions struct {
	accessLevel AccessLevel
	metaData    map[string]string
}

func (o createOptions) ToHeaders() *client.Headers {
	headers := containerOptions{
		metaData: o.metaData,
	}.ToHeaders()

	// If this header is not included in the request, container data is private to the account owner.
	if o.accessLevel != Private {
		headers.Append("x-ms-blob-public-access", string(o.accessLevel))
	}

	return headers
}

func (createOptions) ToOData() *odata.Query {
	return nil
}

func (createOptions) ToQuery() *client.QueryParams {
	return containerOptions{}.ToQuery()
}
