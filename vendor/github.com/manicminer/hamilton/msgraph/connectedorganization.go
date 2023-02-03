package msgraph

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/internal/utils"
	"github.com/manicminer/hamilton/odata"
)

type ConnectedOrganizationClient struct {
	BaseClient Client
}

func NewConnectedOrganizationClient(tenantId string) *ConnectedOrganizationClient {
	return &ConnectedOrganizationClient{
		BaseClient: NewClient(Version10, tenantId),
	}
}

// List returns a list of ConnectedOrganization
// https://docs.microsoft.com/graph/api/entitlementmanagement-list-connectedorganizations
func (c *ConnectedOrganizationClient) List(ctx context.Context, query odata.Query) (*[]ConnectedOrganization, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/connectedOrganizations",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		ConnectedOrganizations []ConnectedOrganization `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.ConnectedOrganizations, status, nil
}

// Create creates a new ConnectedOrganization.
// https://docs.microsoft.com/graph/api/entitlementmanagement-post-connectedorganizations
func (c *ConnectedOrganizationClient) Create(ctx context.Context, connectedOrganization ConnectedOrganization) (*ConnectedOrganization, int, error) {
	var status int
	body, err := json.Marshal(connectedOrganization)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      "/identityGovernance/entitlementManagement/connectedOrganizations",
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newConnectedOrganization ConnectedOrganization
	if err := json.Unmarshal(respBody, &newConnectedOrganization); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newConnectedOrganization, status, nil
}

// Get retrieves a ConnectedOrganization.
// https://docs.microsoft.com/graph/api/connectedorganization-get
func (c *ConnectedOrganizationClient) Get(ctx context.Context, id string, query odata.Query) (*ConnectedOrganization, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		OData:                  query,
		ValidStatusCodes:       []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var connectedOrganization ConnectedOrganization
	if err := json.Unmarshal(respBody, &connectedOrganization); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &connectedOrganization, status, nil
}

// Update amends an existing ConnectedOrganization.
// https://docs.microsoft.com/graph/api/connectedorganization-update
func (c *ConnectedOrganizationClient) Update(ctx context.Context, connectedOrganization ConnectedOrganization) (int, error) {
	var status int

	if connectedOrganization.ID == nil {
		return status, errors.New("cannot update ConnectedOrganization with nil ID")
	}

	// These are the only properties that can up updated.
	updatedOrg := ConnectedOrganization{
		DisplayName: connectedOrganization.DisplayName,
		Description: connectedOrganization.Description,
		State:       connectedOrganization.State,
	}

	body, err := json.Marshal(updatedOrg)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	_, status, _, err = c.BaseClient.Patch(ctx, PatchHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s", *connectedOrganization.ID),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Patch(): %v", err)
	}

	return status, nil
}

// Delete removes a ConnectedOrganization.
// https://docs.microsoft.com/graph/api/connectedorganization-delete
func (c *ConnectedOrganizationClient) Delete(ctx context.Context, id string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: RetryOn404ConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s", id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

// List the external sponsors for a connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-list-externalsponsors
func (c *ConnectedOrganizationClient) ListExternalSponsors(ctx context.Context, query odata.Query, id string) (*[]DirectoryObject, int, error) {
	return listSponsors(&c.BaseClient, ctx, query, id, true)
}

// List the internal sponsors for a connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-list-internalsponsors
func (c *ConnectedOrganizationClient) ListInternalSponsors(ctx context.Context, query odata.Query, id string) (*[]DirectoryObject, int, error) {
	return listSponsors(&c.BaseClient, ctx, query, id, false)
}

// Add a user as an external sponsor to the connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-post-externalsponsors
func (c *ConnectedOrganizationClient) AddExternalSponsorUser(ctx context.Context, orgId string, userId string) error {
	return addSponsor(&c.BaseClient, ctx, orgId, userId, true, false)
}

// Add a group as an external sponsor to the connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-post-externalsponsors
func (c *ConnectedOrganizationClient) AddExternalSponsorGroup(ctx context.Context, orgId string, grpId string) error {
	return addSponsor(&c.BaseClient, ctx, orgId, grpId, true, true)
}

// Add a user as an external sponsor to the connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-post-internalsponsors
func (c *ConnectedOrganizationClient) AddInternalSponsorUser(ctx context.Context, orgId string, userId string) error {
	return addSponsor(&c.BaseClient, ctx, orgId, userId, false, false)
}

// Add a group as an external sponsor to the connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-post-internalsponsors
func (c *ConnectedOrganizationClient) AddInternalSponsorGroup(ctx context.Context, orgId string, grpId string) error {
	return addSponsor(&c.BaseClient, ctx, orgId, grpId, false, true)
}

// Delete a user or group as an external sponsor to the connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-delete-externalsponsors
func (c *ConnectedOrganizationClient) DeleteExternalSponsor(ctx context.Context, orgId string, id string) error {
	return deleteSponsor(&c.BaseClient, ctx, orgId, id, true)
}

// Delete a user or group as an internal sponsor to the connected organization.
// https://docs.microsoft.com/graph/api/connectedorganization-delete-internalsponsors
func (c *ConnectedOrganizationClient) DeleteInternalSponsor(ctx context.Context, orgId string, id string) error {
	return deleteSponsor(&c.BaseClient, ctx, orgId, id, false)
}

func addSponsor(client *Client, ctx context.Context, orgId string, userOrGroupId string, external bool, group bool) error {
	if !ValidateId(&orgId) {
		return fmt.Errorf("the id %q is not a valid connected organization id", orgId)
	}
	if !ValidateId(&userOrGroupId) {
		return fmt.Errorf("the id %q is not a valid user/group id", userOrGroupId)
	}

	var userOrGroup string
	if group {
		userOrGroup = "/groups/"
	} else {
		userOrGroup = "/users/"
	}

	extUser := Ref{
		ObjectUri: utils.StringPtr(string(client.Endpoint) + "/" + string(client.ApiVersion) + userOrGroup + userOrGroupId),
	}

	body, err := json.Marshal(extUser)
	if err != nil {
		return fmt.Errorf("json.Marshal(): %v", err)
	}

	var internalOrExternal string
	if external {
		internalOrExternal = "externalSponsors"
	} else {
		internalOrExternal = "internalSponsors"
	}

	_, status, _, err := client.Post(ctx, PostHttpRequestInput{
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s/%s/$ref", orgId, internalOrExternal),
			HasTenantId: true,
		},
		ValidStatusCodes: []int{http.StatusNoContent},
		Body:             body,
	})

	if err != nil {
		return fmt.Errorf("AddExternalSponsorUser returned status code %d: %v", status, err)
	}

	return nil
}

// List the internal/external sponsors for a connected organization.
func listSponsors(c *Client, ctx context.Context, query odata.Query, id string, external bool) (*[]DirectoryObject, int, error) {
	if !ValidateId(&id) {
		return nil, 0, fmt.Errorf("the id %q is not a valid connected organization id", id)
	}

	var internalOrExternal string
	if external {
		internalOrExternal = "externalSponsors"
	} else {
		internalOrExternal = "internalSponsors"
	}

	resp, status, _, err := c.Get(ctx, GetHttpRequestInput{
		DisablePaging:    query.Top > 0,
		OData:            query,
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s/%s", id, internalOrExternal),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("ConnectedOrganizationClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		ExternalSponsors []DirectoryObject `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.ExternalSponsors, status, nil
}

func deleteSponsor(c *Client, ctx context.Context, orgId string, id string, external bool) error {
	if !ValidateId(&orgId) {
		return fmt.Errorf("the id %q is not a valid connected organization id", id)
	}
	if !ValidateId(&id) {
		return fmt.Errorf("the id %q is not a valid user/group id", id)
	}
	var internalOrExternal string
	if external {
		internalOrExternal = "externalSponsors"
	} else {
		internalOrExternal = "internalSponsors"
	}

	_, status, _, err := c.Delete(ctx, DeleteHttpRequestInput{
		Uri: Uri{
			Entity:      fmt.Sprintf("/identityGovernance/entitlementManagement/connectedOrganizations/%s/%s/%s/$ref", orgId, internalOrExternal, id),
			HasTenantId: true,
		},
		ValidStatusCodes: []int{http.StatusNoContent},
	})

	if err != nil {
		return fmt.Errorf("DeleteSponsor returned status code %d: %v", status, err)
	}

	return nil
}
