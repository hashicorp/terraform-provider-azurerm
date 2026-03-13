// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2021-07-01/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/testclient"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/resource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceFeatureRegistrationResource struct{}

func TestAccResourceFeatureRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_feature_registration", "test")
	r := ResourceFeatureRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			PreConfig: func() {
				// Last error may cause resource provider still in `Registered` status.Need to unregister it before a new test.
				if err := r.unRegisterFeature("EncryptionAtHost", "Microsoft.Compute"); err != nil {
					t.Fatalf("Failed to reset feature registration with error: %+v", err)
				}
			},
			Config: r.basic("EncryptionAtHost", "Microsoft.Compute"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceFeatureRegistration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_feature_registration", "test")
	r := ResourceFeatureRegistrationResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic("EncryptionAtHost", "Microsoft.Compute"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(func(data acceptance.TestData) string {
			return r.requiresImport("EncryptionAtHost", "Microsoft.Compute")
		}),
	})
}

func (ResourceFeatureRegistrationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := features.ParseFeatureID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.FeaturesClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	isRegistered := false
	if model := resp.Model; model != nil && model.Properties != nil && model.Properties.State != nil {
		isRegistered = strings.EqualFold(*model.Properties.State, "Registered")
	}

	return pointer.To(isRegistered), nil
}

func (ResourceFeatureRegistrationResource) basic(name, providerName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}

  resource_provider_registrations = "none"
  resource_providers_to_register  = [%[2]q]
}

resource "azurerm_resource_feature_registration" "test" {
  name          = %[1]q
  provider_name = %[2]q
}
`, name, providerName)
}

func (r ResourceFeatureRegistrationResource) requiresImport(name, providerName string) string {
	template := r.basic(name, providerName)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_feature_registration" "import" {
  name          = azurerm_resource_feature_registration.test.name
  provider_name = azurerm_resource_feature_registration.test.provider_name
}
`, template)
}

func (r ResourceFeatureRegistrationResource) unRegisterFeature(name, providerName string) error {
	client, err := testclient.Build()
	if err != nil {
		return fmt.Errorf("building client: %+v", err)
	}

	ctx, cancel := context.WithDeadline(client.StopContext, time.Now().Add(30*time.Minute))
	defer cancel()

	featuresClient := client.Resource.FeaturesClient
	id := features.NewFeatureID(client.Account.SubscriptionId, providerName, name)
	feature, err := featuresClient.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(feature.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	registrationState := ""
	featureName := ""
	if model := feature.Model; model != nil && feature.Model != nil && feature.Model.Properties != nil && feature.Model.Properties.State != nil {
		registrationState = *feature.Model.Properties.State
		featureName = *feature.Model.Name
	}

	if registrationState == "" {
		return fmt.Errorf("retrieving Feature %q: `state` was nil", featureName)
	}

	if !strings.EqualFold(registrationState, resource.Registered) {
		return nil
	}

	if _, err := featuresClient.Unregister(ctx, id); err != nil {
		return fmt.Errorf("unregistering %s: %+v", id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{resource.Unregistering},
		Target:     []string{resource.NotRegistered, resource.Unregistered},
		Refresh:    r.featureRegisteringStateRefreshFunc(ctx, featuresClient, id),
		MinTimeout: 3 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be complete unregistering: %+v", id, err)
	}

	return nil
}

func (r ResourceFeatureRegistrationResource) featureRegisteringStateRefreshFunc(ctx context.Context, client *features.FeaturesClient, id features.FeatureId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}
		if res.Model == nil || res.Model.Properties == nil || res.Model.Properties.State == nil {
			return nil, "", fmt.Errorf("error reading %s registering status: %+v", id, err)
		}

		return res, *res.Model.Properties.State, nil
	}
}
