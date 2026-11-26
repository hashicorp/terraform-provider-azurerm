// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers_test

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/linux_web_app_slot"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_web_app"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/resources/windows_web_app_slot"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/helpers"
)

var dockerFxStringCases = []struct {
	name      string
	fxVersion string
	expected  string
}{
	{
		name:      "uppercase prefix",
		fxVersion: "DOCKER|mcr.microsoft.com/appsvc/wordpress-debian-php:8.4",
		expected:  "appsvc/wordpress-debian-php:8.4",
	},
	{
		name:      "lowercase prefix",
		fxVersion: "docker|mcr.microsoft.com/appsvc/wordpress-debian-php:8.4",
		expected:  "appsvc/wordpress-debian-php:8.4",
	},
	{
		name:      "mixed case prefix",
		fxVersion: "Docker|mcr.microsoft.com/appsvc/wordpress-debian-php:8.4",
		expected:  "appsvc/wordpress-debian-php:8.4",
	},
}

var dockerAppSettings = map[string]string{
	"DOCKER_REGISTRY_SERVER_URL": "https://mcr.microsoft.com",
}

func TestDecodeDockerAppStackLinux(t *testing.T) {
	for _, tc := range dockerFxStringCases {
		t.Run(tc.name, func(t *testing.T) {
			webApp := helpers.SiteConfigLinux{LinuxFxVersion: tc.fxVersion}
			webApp.DecodeDockerAppStack(dockerAppSettings)
			if got := webApp.ApplicationStack[0].DockerImageName; got != tc.expected {
				t.Fatalf("web app: expected docker_image_name %q, got %q", tc.expected, got)
			}

			slot := linux_web_app_slot.SiteConfigLinuxWebAppSlot{LinuxFxVersion: tc.fxVersion}
			slot.DecodeDockerAppStack(dockerAppSettings)
			if got := slot.ApplicationStack[0].DockerImageName; got != tc.expected {
				t.Fatalf("web app slot: expected docker_image_name %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestDecodeDockerAppStackWindows(t *testing.T) {
	for _, tc := range dockerFxStringCases {
		t.Run(tc.name, func(t *testing.T) {
			webApp := windows_web_app.SiteConfigWindows{WindowsFxVersion: tc.fxVersion}
			webApp.DecodeDockerAppStack(dockerAppSettings)
			if got := webApp.ApplicationStack[0].DockerImageName; got != tc.expected {
				t.Fatalf("web app: expected docker_image_name %q, got %q", tc.expected, got)
			}

			slot := windows_web_app_slot.SiteConfigWindowsWebAppSlot{WindowsFxVersion: tc.fxVersion}
			slot.DecodeDockerAppStack(dockerAppSettings)
			if got := slot.ApplicationStack[0].DockerImageName; got != tc.expected {
				t.Fatalf("web app slot: expected docker_image_name %q, got %q", tc.expected, got)
			}
		})
	}
}
