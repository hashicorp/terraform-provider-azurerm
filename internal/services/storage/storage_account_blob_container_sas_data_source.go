// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/hashicorp/go-azure-helpers/storage"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataSourceStorageAccountBlobContainerSharedAccessSignature() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceStorageContainerSasRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"connection_string": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ip_address": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: storageValidate.SharedAccessSignatureIP,
			},

			"start": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

			"expiry": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ISO8601DateTime,
			},

			"permissions": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"read": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"add": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"create": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"write": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"delete": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"delete_version": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"list": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"tags": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"find": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"move": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"execute": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"ownership": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"permissions": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"set_immutability_policy": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},
					},
				},
			},

			"cache_control": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"content_encoding": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"content_language": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"content_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"sas": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceStorageContainerSasRead(d *pluginsdk.ResourceData, _ interface{}) error {
	connString := d.Get("connection_string").(string)
	containerName := d.Get("container_name").(string)
	httpsOnly := d.Get("https_only").(bool)
	ip := d.Get("ip_address").(string)
	start := d.Get("start").(string)
	expiry := d.Get("expiry").(string)
	permissionsIface := d.Get("permissions").([]interface{})

	// response headers
	cacheControl := d.Get("cache_control").(string)
	contentDisposition := d.Get("content_disposition").(string)
	contentEncoding := d.Get("content_encoding").(string)
	contentLanguage := d.Get("content_language").(string)
	contentType := d.Get("content_type").(string)

	permissions := BuildContainerPermissionsString(permissionsIface[0].(map[string]interface{}))

	// Parse the connection string
	kvp, err := storage.ParseAccountSASConnectionString(connString)
	if err != nil {
		return err
	}

	// Create the string to sign with the key...
	accountName := kvp[connStringAccountNameKey]
	accountKey := kvp[connStringAccountKeyKey]
	signedProtocol := "https,http"
	if httpsOnly {
		signedProtocol = "https"
	}
	signedIp := ip
	signedIdentifier := ""
	signedSnapshotTime := ""

	sasToken, err := storage.ComputeContainerSASToken(permissions, start, expiry, accountName, accountKey,
		containerName, signedIdentifier, signedIp, signedProtocol, signedSnapshotTime, cacheControl,
		contentDisposition, contentEncoding, contentLanguage, contentType)
	if err != nil {
		return err
	}

	d.Set("sas", sasToken)
	tokenHash := sha256.Sum256([]byte(sasToken))
	d.SetId(hex.EncodeToString(tokenHash[:]))

	return nil
}

func BuildContainerPermissionsString(perms map[string]interface{}) string {
	orderedPermissions := []struct {
		name   string
		letter string
	}{
		{"read", "r"},
		{"add", "a"},
		{"create", "c"},
		{"write", "w"},
		{"delete", "d"},
		{"delete_version", "x"},
		{"list", "l"},
		{"tags", "t"},
		{"find", "f"},
		{"move", "m"},
		{"execute", "e"},
		{"ownership", "o"},
		{"permissions", "p"},
		{"set_immutability_policy", "i"},
	}

	retVal := ""

	for _, perm := range orderedPermissions {
		if val, pres := perms[perm.name].(bool); pres && val {
			retVal += perm.letter
		}
	}

	return retVal
}
