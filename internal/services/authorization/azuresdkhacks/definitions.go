package azuresdkhacks

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/authorization/mgmt/2020-04-01-preview/authorization"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

type RoleDefinitionsWorkaroundClient struct {
	sdkClient *authorization.RoleDefinitionsClient
}

func NewRoleDefinitionsWorkaroundClient(client *authorization.RoleDefinitionsClient) RoleDefinitionsWorkaroundClient {
	return RoleDefinitionsWorkaroundClient{
		sdkClient: client,
	}
}

// CreateOrUpdate creates or updates a role definition.
// Parameters:
// scope - the scope of the role definition.
// roleDefinitionID - the ID of the role definition.
// roleDefinition - the values for the role definition.
func (client RoleDefinitionsWorkaroundClient) CreateOrUpdate(ctx context.Context, scope string, roleDefinitionID string, roleDefinition authorization.RoleDefinition) (result RoleDefinitionUpdateResponse, err error) {
	req, err := client.sdkClient.CreateOrUpdatePreparer(ctx, scope, roleDefinitionID, roleDefinition)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorization.RoleDefinitionsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	resp, err := client.sdkClient.CreateOrUpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authorization.RoleDefinitionsClient", "CreateOrUpdate", resp, "Failure sending request")
		return
	}

	result, err = client.CreateOrUpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorization.RoleDefinitionsClient", "CreateOrUpdate", resp, "Failure responding to request")
	}

	return
}

// CreateOrUpdateResponder handles the response to the CreateOrUpdate request. The method always
// closes the http.Response Body.
func (client RoleDefinitionsWorkaroundClient) CreateOrUpdateResponder(resp *http.Response) (result RoleDefinitionUpdateResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Get get role definition by name (GUID).
// Parameters:
// scope - the scope of the role definition.
// roleDefinitionID - the ID of the role definition.
func (client RoleDefinitionsWorkaroundClient) Get(ctx context.Context, scope string, roleDefinitionID string) (result RoleDefinitionGetResponse, err error) {
	req, err := client.sdkClient.GetPreparer(ctx, scope, roleDefinitionID)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorization.RoleDefinitionsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.sdkClient.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "authorization.RoleDefinitionsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "authorization.RoleDefinitionsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client RoleDefinitionsWorkaroundClient) GetResponder(resp *http.Response) (result RoleDefinitionGetResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// RoleDefinition role definition.
type RoleDefinitionGetResponse struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The role definition ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The role definition name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The role definition type.
	Type *string `json:"type,omitempty"`
	// RoleDefinitionProperties - Role definition properties.
	*RoleDefinitionProperties `json:"properties,omitempty"`
}

type RoleDefinitionUpdateResponse struct {
	autorest.Response `json:"-"`
	// ID - READ-ONLY; The role definition ID.
	ID *string `json:"id,omitempty"`
	// Name - READ-ONLY; The role definition name.
	Name *string `json:"name,omitempty"`
	// Type - READ-ONLY; The role definition type.
	Type *string `json:"type,omitempty"`
	// RoleDefinitionProperties - Role definition properties.
	*RoleDefinitionProperties `json:"properties,omitempty"`
}

// RoleDefinitionProperties role definition properties.
type RoleDefinitionProperties struct {
	// RoleName - The role name.
	RoleName *string `json:"roleName,omitempty"`
	// Description - The role definition description.
	Description *string `json:"description,omitempty"`
	// RoleType - The role type.
	RoleType *string `json:"type,omitempty"`
	// Permissions - Role definition permissions.
	Permissions *[]Permission `json:"permissions,omitempty"`
	// AssignableScopes - Role definition assignable scopes.
	AssignableScopes *[]string `json:"assignableScopes,omitempty"`

	// not exposed in the sdk
	CreatedOn *string `json:"createdOn,omitempty"`
	UpdatedOn *string `json:"updatedOn,omitempty"`
	CreatedBy *string `json:"createdBy,omitempty"`
	UpdatedBy *string `json:"updatedBy,omitempty"`
}

// Permission role definition permissions.
type Permission struct {
	// Actions - Allowed actions.
	Actions *[]string `json:"actions,omitempty"`
	// NotActions - Denied actions.
	NotActions *[]string `json:"notActions,omitempty"`
	// DataActions - Allowed Data actions.
	DataActions *[]string `json:"dataActions,omitempty"`
	// NotDataActions - Denied Data actions.
	NotDataActions *[]string `json:"notDataActions,omitempty"`
}
