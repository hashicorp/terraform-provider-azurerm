package authentication

import (
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest/azure/cli"
)

func TestAzureFindValidAccessTokenForTenant_InvalidDate(t *testing.T) {
	tenantId := "c056adac-c6a6-4ddf-ab20-0f26d47f7eea"
	expectedToken := cli.Token{
		ExpiresOn:    "invalid date",
		AccessToken:  "7cabcf30-8dca-43f9-91e6-fd56dfb8632f",
		TokenType:    "9b10b986-7a61-4542-8d5a-9fcd96112585",
		RefreshToken: "4ec3874d-ee2e-4980-ba47-b5bac11ddb94",
		Resource:     "https://management.core.windows.net/",
		Authority:    tenantId,
	}
	tokens := []cli.Token{expectedToken}
	token, err := findValidAccessTokenForTenant(tokens, tenantId)

	if err == nil {
		t.Fatalf("Expected an error to be returned but got nil")
	}

	if token != nil {
		t.Fatalf("Expected Token to be nil but got: %+v", token)
	}
}

func TestAzureFindValidAccessTokenForTenant_Expired(t *testing.T) {
	expirationDate := time.Now().Add(time.Minute * -1)
	tenantId := "c056adac-c6a6-4ddf-ab20-0f26d47f7eea"
	expectedToken := cli.Token{
		ExpiresOn:    expirationDate.Format("2006-01-02 15:04:05.999999"),
		AccessToken:  "7cabcf30-8dca-43f9-91e6-fd56dfb8632f",
		TokenType:    "9b10b986-7a61-4542-8d5a-9fcd96112585",
		RefreshToken: "4ec3874d-ee2e-4980-ba47-b5bac11ddb94",
		Resource:     "https://management.core.windows.net/",
		Authority:    tenantId,
	}
	tokens := []cli.Token{expectedToken}
	token, err := findValidAccessTokenForTenant(tokens, tenantId)

	if err == nil {
		t.Fatalf("Expected an error to be returned but got nil")
	}

	if token != nil {
		t.Fatalf("Expected Token to be nil but got: %+v", token)
	}
}

func TestAzureFindValidAccessTokenForTenant_ExpiringIn(t *testing.T) {
	minutesToVerify := []int{1, 30, 60}

	for _, minute := range minutesToVerify {
		expirationDate := time.Now().Add(time.Minute * time.Duration(minute))
		tenantId := "c056adac-c6a6-4ddf-ab20-0f26d47f7eea"
		expectedToken := cli.Token{
			ExpiresOn:    expirationDate.Format("2006-01-02 15:04:05.999999"),
			AccessToken:  "7cabcf30-8dca-43f9-91e6-fd56dfb8632f",
			TokenType:    "9b10b986-7a61-4542-8d5a-9fcd96112585",
			RefreshToken: "4ec3874d-ee2e-4980-ba47-b5bac11ddb94",
			Resource:     "https://management.core.windows.net/",
			Authority:    tenantId,
		}
		tokens := []cli.Token{expectedToken}
		token, err := findValidAccessTokenForTenant(tokens, tenantId)

		if err != nil {
			t.Fatalf("Expected no error to be returned for minute %d but got %+v", minute, err)
		}

		if token == nil {
			t.Fatalf("Expected Token to have a value for minute %d but it was nil", minute)
		}

		if token.AccessToken.AccessToken != expectedToken.AccessToken {
			t.Fatalf("Expected the Access Token to be %q for minute %d but got %q", expectedToken.AccessToken, minute, token.AccessToken.AccessToken)
		}

		if token.ClientID != expectedToken.ClientID {
			t.Fatalf("Expected the Client ID to be %q for minute %d but got %q", expectedToken.ClientID, minute, token.ClientID)
		}

		if token.IsCloudShell != false {
			t.Fatalf("Expected `IsCloudShell` to be false for minute %d but got true", minute)
		}
	}
}

func TestAzureFindValidAccessTokenForTenant_InvalidManagementDomain(t *testing.T) {
	expirationDate := time.Now().Add(1 * time.Hour)
	tenantId := "c056adac-c6a6-4ddf-ab20-0f26d47f7eea"
	expectedToken := cli.Token{
		ExpiresOn:   expirationDate.Format("2006-01-02 15:04:05.999999"),
		AccessToken: "7cabcf30-8dca-43f9-91e6-fd56dfb8632f",
		TokenType:   "9b10b986-7a61-4542-8d5a-9fcd96112585",
		Resource:    "https://portal.azure.com/",
		Authority:   tenantId,
	}
	tokens := []cli.Token{expectedToken}
	token, err := findValidAccessTokenForTenant(tokens, tenantId)

	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}

	if token != nil {
		t.Fatalf("Expected Token to be nil but got: %+v", token)
	}
}

func TestAzureFindValidAccessTokenForTenant_DifferentTenant(t *testing.T) {
	expirationDate := time.Now().Add(1 * time.Hour)
	expectedToken := cli.Token{
		ExpiresOn:   expirationDate.Format("2006-01-02 15:04:05.999999"),
		AccessToken: "7cabcf30-8dca-43f9-91e6-fd56dfb8632f",
		TokenType:   "9b10b986-7a61-4542-8d5a-9fcd96112585",
		Resource:    "https://management.core.windows.net/",
		Authority:   "9b5095de-5496-4b5e-9bc6-ef2c017b9d35",
	}
	tokens := []cli.Token{expectedToken}
	token, err := findValidAccessTokenForTenant(tokens, "c056adac-c6a6-4ddf-ab20-0f26d47f7eea")

	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}

	if token != nil {
		t.Fatalf("Expected Token to be nil but got: %+v", token)
	}
}

func TestAzureFindValidAccessTokenForTenant_ValidFromCloudShell(t *testing.T) {
	expirationDate := time.Now().Add(1 * time.Hour)
	tenantId := "c056adac-c6a6-4ddf-ab20-0f26d47f7eea"
	expectedToken := cli.Token{
		ExpiresOn:   expirationDate.Format(time.RFC3339),
		AccessToken: "7cabcf30-8dca-43f9-91e6-fd56dfb8632f",
		TokenType:   "9b10b986-7a61-4542-8d5a-9fcd96112585",
		Resource:    "https://management.core.windows.net/",
		Authority:   tenantId,
	}
	tokens := []cli.Token{expectedToken}
	token, err := findValidAccessTokenForTenant(tokens, tenantId)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %+v", err)
	}

	if token == nil {
		t.Fatalf("Expected Token to have a value but it was nil")
	}

	if token.AccessToken.AccessToken != expectedToken.AccessToken {
		t.Fatalf("Expected the Access Token to be %q but got %q", expectedToken.AccessToken, token.AccessToken.AccessToken)
	}

	if token.ClientID != expectedToken.ClientID {
		t.Fatalf("Expected the Client ID to be %q but got %q", expectedToken.ClientID, token.ClientID)
	}

	if token.IsCloudShell != true {
		t.Fatalf("Expected `IsCloudShell` to be true but got false")
	}
}

func TestAzureFindValidAccessTokenForTenant_ValidFromAzureCLI(t *testing.T) {
	expirationDate := time.Now().Add(1 * time.Hour)
	tenantId := "c056adac-c6a6-4ddf-ab20-0f26d47f7eea"
	expectedToken := cli.Token{
		ExpiresOn:    expirationDate.Format("2006-01-02 15:04:05.999999"),
		AccessToken:  "7cabcf30-8dca-43f9-91e6-fd56dfb8632f",
		TokenType:    "9b10b986-7a61-4542-8d5a-9fcd96112585",
		RefreshToken: "4ec3874d-ee2e-4980-ba47-b5bac11ddb94",
		Resource:     "https://management.core.windows.net/",
		Authority:    tenantId,
	}
	tokens := []cli.Token{expectedToken}
	token, err := findValidAccessTokenForTenant(tokens, tenantId)

	if err != nil {
		t.Fatalf("Expected no error to be returned but got %+v", err)
	}

	if token == nil {
		t.Fatalf("Expected Token to have a value but it was nil")
	}

	if token.AccessToken.AccessToken != expectedToken.AccessToken {
		t.Fatalf("Expected the Access Token to be %q but got %q", expectedToken.AccessToken, token.AccessToken.AccessToken)
	}

	if token.ClientID != expectedToken.ClientID {
		t.Fatalf("Expected the Client ID to be %q but got %q", expectedToken.ClientID, token.ClientID)
	}

	if token.IsCloudShell != false {
		t.Fatalf("Expected `IsCloudShell` to be false but got true")
	}
}

func TestAzureFindValidAccessTokenForTenant_NoTokens(t *testing.T) {
	tokens := make([]cli.Token, 0)
	token, err := findValidAccessTokenForTenant(tokens, "abc123")

	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}

	if token != nil {
		t.Fatalf("Expected a null token to be returned but got: %+v", token)
	}
}
