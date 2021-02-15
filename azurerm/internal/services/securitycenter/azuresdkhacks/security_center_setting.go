package azuresdkhacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

func GetSecurityCenterSetting(client *security.SettingsClient, ctx context.Context, settingName string) (setting security.DataExportSettings, err error) {
	// NOTE: client.Get() returns security.Setting, which doesn't contain the "Enabled" property
	// https://github.com/Azure/azure-sdk-for-go/issues/12724
	req, err := client.GetPreparer(ctx, settingName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SettingsClient", "Get", nil, "Failure preparing request")
		return setting, fmt.Errorf("Reading Security Center setting: %+v", err)
	}
	resp, err := client.GetSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "security.SettingsClient", "Get", resp, "Failure sending request")
		return setting, fmt.Errorf("Reading Security Center setting: %+v", err)
	}

	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&setting),
		autorest.ByClosing())
	if err != nil {
		return setting, fmt.Errorf("Reading Security Center setting: %+v", err)
	}

	return setting, nil
}
