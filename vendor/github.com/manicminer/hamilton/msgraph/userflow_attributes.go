package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

// UserFlowAttributesClient performs operations on UserFlowAttributes.
type UserFlowAttributesClient struct {
	BaseClient Client
}

// NewUserFlowAttributesClient returns a new UserFlowAttributesClient.
func NewUserFlowAttributesClient(tenantId string) *UserFlowAttributesClient {
	return &UserFlowAttributesClient{
		BaseClient: NewClient(Version10, tenantId),
	}
}

// List returns a list of UserFlowAttributes, optionally queried using OData.
func (c *UserFlowAttributesClient) List(ctx context.Context, query odata.Query) (*[]UserFlowAttribute, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/identity/userFlowAttributes",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("UserFlowAttributesClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		UserFlowAttributes []UserFlowAttribute `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.UserFlowAttributes, status, nil
}

// Create creates a new UserFlowAttribute.
func (c *UserFlowAttributesClient) Create(ctx context.Context, userFlowAttribute UserFlowAttribute) (*UserFlowAttribute, int, error) {
	var status int

	body, err := json.Marshal(userFlowAttribute)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body: body,
		OData: odata.Query{
			Metadata: odata.MetadataFull,
		},
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/identity/userFlowAttributes",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("UserFlowAttributesClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newUserFlowAttribute UserFlowAttribute
	if err := json.Unmarshal(respBody, &newUserFlowAttribute); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newUserFlowAttribute, status, nil
}

// Delete returns a UserFlowAttribute.
func (c *UserFlowAttributesClient) Get(ctx context.Context, id string, query odata.Query) (*UserFlowAttribute, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identity/userFlowAttributes/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("UserFlowAttributesClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var userflowAttribute UserFlowAttribute
	if err := json.Unmarshal(respBody, &userflowAttribute); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &userflowAttribute, status, nil
}

// Update amends an existing UserFlowAttribute.
func (c *UserFlowAttributesClient) Update(ctx context.Context, userflowAttribute UserFlowAttribute) (int, error) {
	var status int
	if userflowAttribute.ID == nil {
		return status, fmt.Errorf("cannot update userflowAttribute with nil ID")
	}

	userflowID := *userflowAttribute.ID
	userflowAttribute.ID = nil

	body, err := json.Marshal(userflowAttribute)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes: []int{
			http.StatusOK,
			http.StatusNoContent,
		},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identity/userFlowAttributes//%s", userflowID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("UserFlowAttributesClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}

// Delete removes a UserFlowAttribute.
func (c *UserFlowAttributesClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identity/userFlowAttributes/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("UserFlowAttributesClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}
