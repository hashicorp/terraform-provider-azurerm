package cdn

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importCdnFrontDoorCustomDomainAssociation() pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		return []*pluginsdk.ResourceData{d}, nil
	}
}
