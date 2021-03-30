// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ManagedIdentityCredentialOptions contains parameters that can be used to configure the pipeline used with Managed Identity Credential.
// All zero-value fields will be initialized with their default values.
type ManagedIdentityCredentialOptions struct {
	// HTTPClient sets the transport for making HTTP requests.
	// Leave this as nil to use the default HTTP transport.
	HTTPClient azcore.Transport

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry azcore.TelemetryOptions

	// Logging configures the built-in logging policy behavior.
	Logging azcore.LogOptions
}

// ManagedIdentityCredential attempts authentication using a managed identity that has been assigned to the deployment environment. This authentication type works in several
// managed identity environments such as Azure VMs, App Service, Azure Functions, Azure CloudShell, among others. More information about configuring managed identities can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview
type ManagedIdentityCredential struct {
	clientID string
	client   *managedIdentityClient
}

// NewManagedIdentityCredential creates an instance of the ManagedIdentityCredential capable of authenticating a resource that has a managed identity.
// clientID: The client ID to authenticate for a user assigned managed identity.
// options: ManagedIdentityCredentialOptions that configure the pipeline for requests sent to Azure Active Directory.
// More information on user assigned managed identities cam be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview#how-a-user-assigned-managed-identity-works-with-an-azure-vm
func NewManagedIdentityCredential(clientID string, options *ManagedIdentityCredentialOptions) (*ManagedIdentityCredential, error) {
	// Create a new Managed Identity Client with default options
	if options == nil {
		options = &ManagedIdentityCredentialOptions{}
	}
	client := newManagedIdentityClient(options)
	msiType, err := client.getMSIType()
	// If there is an error that means that the code is not running in a Managed Identity environment
	if err != nil {
		credErr := &CredentialUnavailableError{credentialType: "Managed Identity Credential", message: "Please make sure you are running in a managed identity environment, such as a VM, Azure Functions, Cloud Shell, etc..."}
		logCredentialError(credErr.credentialType, credErr)
		return nil, credErr
	}
	// Assign the msiType discovered onto the client
	client.msiType = msiType
	// check if no clientID is specified then check if it exists in an environment variable
	if len(clientID) == 0 {
		clientID = os.Getenv("AZURE_CLIENT_ID")
	}
	return &ManagedIdentityCredential{clientID: clientID, client: client}, nil
}

// GetToken obtains an AccessToken from the Managed Identity service if available.
// scopes: The list of scopes for which the token will have access.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *ManagedIdentityCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticate(ctx, c.clientID, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Managed Identity Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	logMSIEnv(c.client.msiType)
	return tk, err
}

// AuthenticationPolicy implements the azcore.Credential interface on ManagedIdentityCredential.
// NOTE: The TokenRequestOptions included in AuthenticationPolicyOptions must be a slice of resources in this case and not scopes.
func (c *ManagedIdentityCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	// The following code will remove the /.default suffix from any scopes passed into the method since ManagedIdentityCredentials expect a resource string instead of a scope string
	for i := range options.Options.Scopes {
		options.Options.Scopes[i] = strings.TrimSuffix(options.Options.Scopes[i], defaultSuffix)
	}
	return newBearerTokenPolicy(c, options)
}
