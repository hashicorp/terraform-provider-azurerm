// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package linux_function_app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	appserviceClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_function_app/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -test-params "B1"

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate         = Resource{}
	_ sdk.ResourceWithCustomImporter = Resource{}
	_ sdk.ResourceWithCustomizeDiff  = Resource{}
	_ sdk.ResourceWithStateMigration = Resource{}
	_ sdk.ResourceWithIdentity       = Resource{}
)

func (r Resource) ResourceType() string {
	return "azurerm_linux_function_app"
}

func (r Resource) Identity() resourceids.ResourceId {
	return &commonids.FunctionAppId{}
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateFunctionAppID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := appserviceClient.From(metadata.Client).WebAppsClient
		servicePlanClient := appserviceClient.From(metadata.Client).ServicePlanClient

		id, err := commonids.ParseFunctionAppID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}
		site, err := client.Get(ctx, *id)
		if err != nil || site.Model == nil || site.Model.Properties == nil {
			return fmt.Errorf("reading Linux %s: %+v", id, err)
		}
		props := *site.Model.Properties
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
		if !strings.Contains(strings.ToLower(*sp.Model.Kind), "linux") && !strings.Contains(strings.ToLower(*sp.Model.Kind), "elastic") && !strings.Contains(strings.ToLower(*sp.Model.Kind), "functionapp") {
			return errors.New("specified Service Plan is not a Linux Functionapp plan")
		}

		return nil
	}
}

func (r Resource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.LinuxFunctionAppV0toV1{},
		},
	}
}
