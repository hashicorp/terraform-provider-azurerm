// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_web_app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_web_app/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var _ sdk.ResourceWithUpdate = Resource{}

var _ sdk.ResourceWithCustomImporter = Resource{}

var _ sdk.ResourceWithCustomizeDiff = Resource{}

var _ sdk.ResourceWithStateMigration = Resource{}

func (r Resource) ResourceType() string {
	return "azurerm_linux_web_app"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateAppServiceID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.AppService.WebAppsClient
		servicePlanClient := metadata.Client.AppService.ServicePlanClient

		id, err := commonids.ParseWebAppID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		site, err := client.Get(ctx, *id)
		if err != nil || site.Model == nil || site.Model.Properties == nil {
			return fmt.Errorf("reading Linux %s: %+v", id, err)
		}
		props := site.Model.Properties
		if props.ServerFarmId == nil {
			return fmt.Errorf("determining Service Plan ID for Linux %s: %+v", id, err)
		}
		servicePlanId, err := commonids.ParseAppServicePlanIDInsensitively(*props.ServerFarmId)
		if err != nil {
			return err
		}

		sp, err := servicePlanClient.Get(ctx, *servicePlanId)
		if err != nil || sp.Model == nil || sp.Model.Kind == nil {
			return fmt.Errorf("reading Service Plan for Linux %s: %+v", id, err)
		}
		if !strings.Contains(strings.ToLower(*sp.Model.Kind), "linux") {
			return errors.New("specified Service Plan is not a Linux plan")
		}

		return nil
	}
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.LinuxWebAppV0toV1{},
		},
	}
}
