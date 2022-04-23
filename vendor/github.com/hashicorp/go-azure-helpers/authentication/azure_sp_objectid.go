package authentication

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/go-autorest/autorest"
	authWrapper "github.com/manicminer/hamilton-autorest/auth"
	envWrapper "github.com/manicminer/hamilton-autorest/environments"
	"github.com/manicminer/hamilton/auth"
	"github.com/manicminer/hamilton/environments"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/manicminer/hamilton/odata"

	"github.com/hashicorp/go-azure-helpers/sender"
)

func buildServicePrincipalObjectIDFunc(c *Config) func(ctx context.Context) (*string, error) {
	return func(ctx context.Context) (*string, error) {
		if c.UseMicrosoftGraph {
			objectId, err := objectIdFromMSALTokenClaims(ctx, c)
			if err != nil {
				log.Printf("could not parse objectId from claims, retrying via Microsoft Graph: %v", err)
				return objectIdFromMsGraph(ctx, c)
			}

			return objectId, err
		} else {
			objectId, err := objectIdFromADALTokenClaims(ctx, c)
			if err != nil {
				log.Printf("could not parse objectId from claims, retrying via Azure Active Directory Graph: %v", err)
				return objectIdFromAadGraph(ctx, c)
			}

			return objectId, err
		}
	}
}

func claimsFromAutorestAuthorizer(authorizer autorest.Authorizer) (*auth.Claims, error) {
	wrapper, err := authWrapper.NewAuthorizerWrapper(authorizer)
	if err != nil {
		return nil, fmt.Errorf("wrapping autorest.Authorizer: %v", err)
	}

	token, err := wrapper.Token()
	if err != nil {
		return nil, fmt.Errorf("acquiring access token: %v", err)
	}

	claims, err := auth.ParseClaims(token)
	if err != nil {
		return nil, fmt.Errorf("parsing claims from access token: %v", err)
	}

	return &claims, nil
}

func objectIdFromADALTokenClaims(ctx context.Context, c *Config) (*string, error) {
	env, err := AzureEnvironmentByNameFromEndpoint(ctx, c.MetadataHost, c.Environment)
	if err != nil {
		return nil, fmt.Errorf("determining environment: %v", err)
	}

	s := sender.BuildSender("GoAzureHelpers")

	oauthConfig, err := c.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("building oauthConfig: %v", err)
	}

	authorizer, err := c.GetADALToken(ctx, s, oauthConfig, env.GraphEndpoint)
	if err != nil {
		return nil, fmt.Errorf("configuring Authorizer: %v", err)
	}

	claims, err := claimsFromAutorestAuthorizer(authorizer)
	if err != nil {
		return nil, err
	}

	return &claims.ObjectId, nil
}

func objectIdFromMSALTokenClaims(ctx context.Context, c *Config) (*string, error) {
	env, err := environments.EnvironmentFromString(c.Environment)
	if err != nil {
		// failed to find a suitable hamilton environment, so convert the provided autorest environment
		azureEnv, err := AzureEnvironmentByNameFromEndpoint(ctx, c.MetadataHost, c.Environment)
		if err != nil {
			return nil, fmt.Errorf("determining environment: %v", err)
		}

		env = envWrapper.EnvironmentFromAzureEnvironment(*azureEnv)
	}

	oauthConfig, err := c.BuildOAuthConfig(string(env.AzureADEndpoint))
	if err != nil {
		return nil, fmt.Errorf("building oauthConfig: %v", err)
	}

	authorizer, err := c.GetMSALToken(ctx, env.MsGraph, sender.BuildSender("GoAzureHelpers"), oauthConfig, string(env.MsGraph.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("configuring Authorizer: %v", err)
	}

	claims, err := claimsFromAutorestAuthorizer(authorizer)
	if err != nil {
		return nil, err
	}

	return &claims.ObjectId, nil
}

func objectIdFromAadGraph(ctx context.Context, c *Config) (*string, error) {
	env, err := AzureEnvironmentByNameFromEndpoint(ctx, c.MetadataHost, c.Environment)
	if err != nil {
		return nil, fmt.Errorf("determining environment: %v", err)
	}

	s := sender.BuildSender("GoAzureHelpers")

	oauthConfig, err := c.BuildOAuthConfig(env.ActiveDirectoryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("building oauthConfig: %v", err)
	}

	graphAuth, err := c.GetADALToken(ctx, s, oauthConfig, env.GraphEndpoint)
	if err != nil {
		return nil, fmt.Errorf("configuring Authorizer: %v", err)
	}

	client := graphrbac.NewServicePrincipalsClientWithBaseURI(env.GraphEndpoint, c.TenantID)
	client.Authorizer = graphAuth
	client.Sender = s

	filter := fmt.Sprintf("appId eq '%s'", c.ClientID)
	listResult, listErr := client.List(ctx, filter)

	if listErr != nil {
		return nil, fmt.Errorf("listing Service Principals: %#v", listErr)
	}

	if listResult.Values() == nil || len(listResult.Values()) != 1 || listResult.Values()[0].ObjectID == nil {
		return nil, fmt.Errorf("unexpected Service Principal query result: %#v", listResult.Values())
	}

	return listResult.Values()[0].ObjectID, nil
}

func objectIdFromMsGraph(ctx context.Context, c *Config) (*string, error) {
	env, err := environments.EnvironmentFromString(c.Environment)
	if err != nil {
		// failed to find a suitable hamilton environment, so convert the provided autorest environment
		azureEnv, err := AzureEnvironmentByNameFromEndpoint(ctx, c.MetadataHost, c.Environment)
		if err != nil {
			return nil, fmt.Errorf("determining environment: %v", err)
		}

		env = envWrapper.EnvironmentFromAzureEnvironment(*azureEnv)
	}

	oauthConfig, err := c.BuildOAuthConfig(string(env.AzureADEndpoint))
	if err != nil {
		return nil, fmt.Errorf("building oauthConfig: %v", err)
	}

	msGraphAuth, err := c.GetMSALToken(ctx, env.MsGraph, sender.BuildSender("GoAzureHelpers"), oauthConfig, string(env.MsGraph.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("configuring Authorizer: %v", err)
	}

	authorizerWrapper, err := authWrapper.NewAuthorizerWrapper(msGraphAuth)
	if err != nil {
		return nil, fmt.Errorf("configuring Authorizer wrapper: %v", err)
	}

	client := msgraph.NewServicePrincipalsClient(c.TenantID)
	client.BaseClient.ApiVersion = msgraph.Version10
	client.BaseClient.Authorizer = authorizerWrapper
	client.BaseClient.DisableRetries = true
	client.BaseClient.Endpoint = env.MsGraph.Endpoint
	client.BaseClient.RequestMiddlewares = &[]msgraph.RequestMiddleware{hamiltonRequestLogger}
	client.BaseClient.ResponseMiddlewares = &[]msgraph.ResponseMiddleware{hamiltonResponseLogger}

	result, status, err := client.List(ctx, odata.Query{Filter: fmt.Sprintf("appId eq '%s'", c.ClientID)})
	if err != nil {
		if status == http.StatusUnauthorized || status == http.StatusForbidden {
			return nil, fmt.Errorf("access denied when listing Service Principals: %+v", err)
		}
		return nil, fmt.Errorf("listing Service Principals: %+v", err)
	}

	if result == nil {
		return nil, fmt.Errorf("unexpected Service Principal query result, was nil")
	}

	if len(*result) != 1 || (*result)[0].ID == nil {
		return nil, fmt.Errorf("unexpected Service Principal query result: %+v", *result)
	}

	return (*result)[0].ID, nil
}

func hamiltonRequestLogger(req *http.Request) (*http.Request, error) {
	if req == nil {
		return nil, nil
	}

	if dump, err := httputil.DumpRequestOut(req, true); err == nil {
		log.Printf("[DEBUG] GoAzureHelpers Request: \n%s\n", dump)
	} else {
		log.Printf("[DEBUG] GoAzureHelpers Request: %s to %s\n", req.Method, req.URL)
	}

	return req, nil
}

func hamiltonResponseLogger(req *http.Request, resp *http.Response) (*http.Response, error) {
	if resp == nil {
		log.Printf("[DEBUG] GoAzureHelpers Request for %s %s completed with no response", req.Method, req.URL)
		return nil, nil
	}

	if dump, err := httputil.DumpResponse(resp, true); err == nil {
		log.Printf("[DEBUG] GoAzureHelpers Response: \n%s\n", dump)
	} else {
		log.Printf("[DEBUG] GoAzureHelpers Response: %s for %s %s\n", resp.Status, req.Method, req.URL)
	}

	return resp, nil
}
