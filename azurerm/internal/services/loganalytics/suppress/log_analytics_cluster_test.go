package suppress

import "testing"

func TestCaseClusterUrl(t *testing.T) {
	cases := []struct {
		Name        string
		ClusterURL  string
		KeyVaultURL string
		Suppress    bool
	}{
		{
			Name:        "empty URL",
			ClusterURL:  "",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    false,
		},
		{
			Name:        "URL with port and wrong scheme",
			ClusterURL:  "http://flynns.arcade.com:443",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    false,
		},
		{
			Name:        "invalid URL scheme",
			ClusterURL:  "https//flynns.arcade.com",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    false,
		},
		{
			Name:        "invalid URL character",
			ClusterURL:  "https://flynns^arcade.com/",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    false,
		},
		{
			Name:        "invalid URL missing scheme",
			ClusterURL:  "//flynns.arcade.com/",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    false,
		},
		{
			Name:        "URL with wrong scheme no port",
			ClusterURL:  "http://flynns.arcade.com",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    false,
		},
		{
			Name:        "same URL different case",
			ClusterURL:  "https://Flynns.Arcade.com/",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    false,
		},
		{
			Name:        "full URL with username@host/path?query#fragment",
			ClusterURL:  "https://Creator4983@flynns.arcade.com/ENCOM?games=MatrixBlaster#MCP",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    true,
		},
		{
			Name:        "full URL with username:password@host/path?query#fragment",
			ClusterURL:  "https://Creator4983:7898@flynns.arcade.com/ENCOM?games=SpaceParanoids&developer=KevinFlynn#MCP",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    true,
		},
		{
			Name:        "URL missing path separator",
			ClusterURL:  "https://flynns.arcade.com",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    true,
		},
		{
			Name:        "same URL",
			ClusterURL:  "https://flynns.arcade.com/",
			KeyVaultURL: "https://flynns.arcade.com/",
			Suppress:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if LogAnalyticsClusterUrl("test", tc.ClusterURL, tc.KeyVaultURL, nil) != tc.Suppress {
				t.Fatalf("Expected LogAnalyticsClusterUrl to return %t for '%q' == '%q'", tc.Suppress, tc.ClusterURL, tc.KeyVaultURL)
			}
		})
	}
}
