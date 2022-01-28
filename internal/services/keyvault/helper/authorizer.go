package helper

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/hashicorp/go-azure-helpers/authentication"
)

func CustomMultiTenantBearerAuthorizer(ctx context.Context, sender autorest.Sender, oauthConfig *authentication.OAuthConfig, clientID, clientSecret string) *autorest.BearerAuthorizerCallback {
	return autorest.NewBearerAuthorizerCallback(sender, func(tenantID, resource string) (*autorest.BearerAuthorizer, error) {
		oauthCfg := oauthConfig.OAuth
		oAuthTenantId := strings.TrimPrefix(oauthCfg.AuthorityEndpoint.Path, "/")
		if oAuthTenantId != tenantID && oauthConfig.MultiTenantOauth != nil {
			mto := *oauthConfig.MultiTenantOauth
			for _, cfg := range mto.AuxiliaryTenants() {
				if cfg != nil && strings.TrimPrefix(cfg.AuthorityEndpoint.Path, "/") == tenantID {
					oauthCfg = cfg
					break
				}
			}
		}
		spt, err := adal.NewServicePrincipalToken(*oauthCfg, clientID, clientSecret, resource)
		if err != nil {
			return nil, fmt.Errorf("while creating new service principal token: %+v", err)
		}
		spt.SetSender(sender)
		return autorest.NewBearerAuthorizer(spt), nil
	})
}
