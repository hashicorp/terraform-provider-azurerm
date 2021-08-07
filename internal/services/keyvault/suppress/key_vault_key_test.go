package suppress

import "testing"

func TestSuppressKeyVaultKeyVersion(t *testing.T) {
	cases := []struct {
		Name     string
		OldID    string
		NewID    string
		Suppress bool
	}{
		{
			Name:     "same versionless ids",
			OldID:    "https://test-kv.vault.azure.net/keys/my-encryption-key",
			NewID:    "https://test-kv.vault.azure.net/keys/my-encryption-key",
			Suppress: true,
		},
		{
			Name:     "different versionless ids",
			OldID:    "https://test-kv.vault.azure.net/keys/my-encryption-key",
			NewID:    "https://test-kv.vault.azure.net/keys/your-encryption-key",
			Suppress: false,
		},
		{
			Name:     "same versioned ids",
			OldID:    "https://test-kv.vault.azure.net/keys/my-encryption-key/54f53756142643ecb420a8a5eaacad2d",
			NewID:    "https://test-kv.vault.azure.net/keys/my-encryption-key/54f53756142643ecb420a8a5eaacad2d",
			Suppress: true,
		},
		{
			Name:     "same key but different versions",
			OldID:    "https://test-kv.vault.azure.net/keys/my-encryption-key/54f53756142643ecb420a8a5eaacad2d",
			NewID:    "https://test-kv.vault.azure.net/keys/my-encryption-key/0ce6731616e54f0e82130d33172492d0",
			Suppress: true,
		},
		{
			Name:     "different key and additional version",
			OldID:    "https://test-kv.vault.azure.net/keys/my-encryption-key",
			NewID:    "https://test-kv.vault.azure.net/keys/your-encryption-key/54f53756142643ecb420a8a5eaacad2d",
			Suppress: false,
		},
		{
			Name:     "extra path separator",
			OldID:    "https://test-kv.vault.azure.net/keys/my-encryption-key/",
			NewID:    "https://test-kv.vault.azure.net/keys/my-encryption-key",
			Suppress: true,
		},
		{
			Name:     "invalid id",
			OldID:    "test-kv.vault.azure.net/keys/",
			NewID:    "https://test-kv.vault.azure.net/keys/my-encryption-key",
			Suppress: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if DiffSuppressIgnoreKeyVaultKeyVersion("test", tc.OldID, tc.NewID, nil) != tc.Suppress {
				t.Fatalf("Expected DiffSuppressIgnoreKeyVaultKeyVersion to return %t for '%q' == '%q'", tc.Suppress, tc.OldID, tc.NewID)
			}
		})
	}
}
