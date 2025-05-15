package workspace

import (
	"context"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type CreateOrUpdateOperationResponse struct {
	HttpResponse *http.Response
	Model        *WorkspaceContract
}

func (c WorkspaceClient) CreateOrUpdate(ctx context.Context, id WorkspaceId, input WorkspaceContract) (result CreateOrUpdateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		Path:       id.ID(),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	result.Model = &WorkspaceContract{}
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
