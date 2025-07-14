// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const (
	JavaContainerEmbeddedServer        string = "JAVA"
	JavaContainerTomcat                string = "TOMCAT"
	JavaContainerEmbeddedServerVersion string = "SE"
	PhpVersionSevenPointOne            string = "7.1"
	PhpVersionSevenPointFour           string = "7.4"
	PhpVersionOff                      string = "Off"

	CurrentStackDotNet     string = "dotnet"
	CurrentStackDotNetCore string = "dotnetcore"
	CurrentStackJava       string = "java"
	CurrentStackNode       string = "node"
	CurrentStackPhp        string = "php"
	CurrentStackPython     string = "python"

	LinuxJavaServerJava   string = "JAVA"
	LinuxJavaServerTomcat string = "TOMCAT"
	LinuxJavaServerJboss  string = "JBOSSEAP"
)

type ApplicationStackWindows struct {
	CurrentStack         string `tfschema:"current_stack"`
	JavaContainer        string `tfschema:"java_container"`
	JavaContainerVersion string `tfschema:"java_container_version"`
	JavaEmbeddedServer   bool   `tfschema:"java_embedded_server_enabled"`
	JavaVersion          string `tfschema:"java_version"`
	NetFrameworkVersion  string `tfschema:"dotnet_version"`
	NetCoreVersion       string `tfschema:"dotnet_core_version"`
	NodeVersion          string `tfschema:"node_version"`
	PhpVersion           string `tfschema:"php_version"`
	Python               bool   `tfschema:"python"`
	TomcatVersion        string `tfschema:"tomcat_version"`

	DockerRegistryUrl      string `tfschema:"docker_registry_url"`
	DockerRegistryUsername string `tfschema:"docker_registry_username"`
	DockerRegistryPassword string `tfschema:"docker_registry_password"`
	DockerImageName        string `tfschema:"docker_image_name"`
}

var windowsApplicationStackConstraint = []string{
	"site_config.0.application_stack.0.docker_image_name",
	"site_config.0.application_stack.0.dotnet_version",
	"site_config.0.application_stack.0.dotnet_core_version",
	"site_config.0.application_stack.0.java_version",
	"site_config.0.application_stack.0.node_version",
	"site_config.0.application_stack.0.php_version",
	"site_config.0.application_stack.0.python",
}

func windowsApplicationStackSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{ // Note: DotNet versions are abstracted between API and Portal displayed values, so do not match 1:1. A table of the converted values is provided in the resource doc.
						"v2.0",
						"v3.0",
						"v4.0",
						"v5.0",
						"v6.0",
						"v7.0",
						"v8.0",
						"v9.0",
					}, false),
					AtLeastOneOf: windowsApplicationStackConstraint,
				},

				"dotnet_core_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"v4.0",
					}, false),
					AtLeastOneOf: windowsApplicationStackConstraint,
					Description:  "The version of DotNetCore to use.",
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						PhpVersionSevenPointOne,  // Deprecated
						PhpVersionSevenPointFour, // Deprecated
						PhpVersionOff,            // Portal displays `Off` for `""` meaning use latest available
					}, false),
					AtLeastOneOf: windowsApplicationStackConstraint,
				},

				"python": {
					Type:         pluginsdk.TypeBool,
					Optional:     true,
					Default:      false,
					AtLeastOneOf: windowsApplicationStackConstraint,
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"~14",
						"~16",
						"~18",
						"~20",
						"~22",
					}, false),
					AtLeastOneOf: windowsApplicationStackConstraint,
				},

				"java_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					AtLeastOneOf: windowsApplicationStackConstraint,
				},

				"java_embedded_server_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Computed: true,
					ConflictsWith: []string{
						"site_config.0.application_stack.0.tomcat_version",
					},
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_version",
					},
					Description: "Should the application use the embedded web server for the version of Java in use.",
				},

				"tomcat_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty, // This is a long list of regularly changing values, not all valid values of which are made known in the portal/docs
					ConflictsWith: []string{
						"site_config.0.application_stack.0.java_embedded_server_enabled",
					},
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_version",
					},
				},

				"java_container": {
					Type:       pluginsdk.TypeString,
					Optional:   true,
					Deprecated: "this property has been deprecated in favour of `tomcat_version` and `java_embedded_server_enabled`",
					ValidateFunc: validation.StringInSlice([]string{
						"JAVA",
						"JETTY", // No longer supported / offered - Java SE or Tomcat (10, 9.5, 8) only
						"TOMCAT",
					}, false),
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_container_version",
					},
					ConflictsWith: []string{
						"site_config.0.application_stack.0.tomcat_version",
					},
				},

				"java_container_version": {
					Type:       pluginsdk.TypeString,
					Optional:   true,
					Deprecated: "This property has been deprecated in favour of `tomcat_version` and `java_embedded_server_enabled`",
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_container",
					},
				},

				"docker_image_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					AtLeastOneOf: windowsApplicationStackConstraint,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"docker_registry_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					RequiredWith: []string{"site_config.0.application_stack.0.docker_image_name"},
				},

				"docker_registry_username": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"docker_registry_password": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
				},

				"current_stack": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true, // This will be set to the configured type from above if not explicitly set
					ValidateFunc: validation.StringInSlice([]string{
						"dotnet",
						"dotnetcore",
						"node",
						"python",
						"php",
						"java",
					}, false),
				},
			},
		},
	}
}

func windowsApplicationStackSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"dotnet_core_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"python": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"node_version": { // Discarded by service if JavaVersion is specified
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_embedded_server_enabled": {
					Type:     pluginsdk.TypeBool,
					Computed: true,
				},

				"tomcat_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_container": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_container_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"current_stack": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_image_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_registry_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_registry_username": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_registry_password": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
			},
		},
	}
}

type ApplicationStackLinux struct {
	NetFrameworkVersion string `tfschema:"dotnet_version"`
	GoVersion           string `tfschema:"go_version"`
	PhpVersion          string `tfschema:"php_version"`
	PythonVersion       string `tfschema:"python_version"`
	NodeVersion         string `tfschema:"node_version"`
	JavaVersion         string `tfschema:"java_version"`
	JavaServer          string `tfschema:"java_server"`
	JavaServerVersion   string `tfschema:"java_server_version"`
	RubyVersion         string `tfschema:"ruby_version"`

	DockerRegistryUrl      string `tfschema:"docker_registry_url"`
	DockerRegistryUsername string `tfschema:"docker_registry_username"`
	DockerRegistryPassword string `tfschema:"docker_registry_password"`
	DockerImageName        string `tfschema:"docker_image_name"`
}

var linuxApplicationStackConstraint = []string{
	"site_config.0.application_stack.0.docker_image_name",
	"site_config.0.application_stack.0.dotnet_version",
	"site_config.0.application_stack.0.java_version",
	"site_config.0.application_stack.0.node_version",
	"site_config.0.application_stack.0.php_version",
	"site_config.0.application_stack.0.python_version",
	"site_config.0.application_stack.0.ruby_version",
	"site_config.0.application_stack.0.go_version",
}

func linuxApplicationStackSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"3.1",
						"5.0", // deprecated
						"6.0",
						"7.0",
						"8.0",
						"9.0",
					}, false),
					ExactlyOneOf: linuxApplicationStackConstraint,
				},

				"go_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"1.19",
						"1.18",
					}, false),
					ExactlyOneOf: linuxApplicationStackConstraint,
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"7.4",
						"8.0",
						"8.1",
						"8.2",
						"8.3",
					}, false),
					ExactlyOneOf: linuxApplicationStackConstraint,
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"3.7",
						"3.8",
						"3.9",
						"3.10",
						"3.11",
						"3.12",
						"3.13",
					}, false),
					ExactlyOneOf: linuxApplicationStackConstraint,
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"12-lts",
						"14-lts",
						"16-lts",
						"18-lts",
						"20-lts",
						"22-lts",
					}, false),
					ExactlyOneOf: linuxApplicationStackConstraint,
				},

				"ruby_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"2.6", // TODO: 4.0 Deprecated - accepted but not offered in the portal. Remove in 4.0
						"2.7", // TODO: 4.0 EOL 31/03/2023 https://github.com/Azure/app-service-linux-docs/blob/master/Runtime_Support/ruby_support.md Remove Ruby support in 4.0?
					}, false),
					ExactlyOneOf: linuxApplicationStackConstraint,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"8",
						"11",
						"17",
						"21",
					}, false),
					ExactlyOneOf: linuxApplicationStackConstraint,
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_server_version", "site_config.0.application_stack.0.java_server",
					},
				},

				"java_server": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"JAVA",
						"TOMCAT",
						"JBOSSEAP",
					}, false),
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_version", "site_config.0.application_stack.0.java_server_version",
					},
				},

				"java_server_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_version", "site_config.0.application_stack.0.java_server",
					},
				},

				"docker_image_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ExactlyOneOf: linuxApplicationStackConstraint,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"docker_registry_url": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					RequiredWith: []string{"site_config.0.application_stack.0.docker_image_name"},
				},

				"docker_registry_username": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"docker_registry_password": {
					Type:      pluginsdk.TypeString,
					Optional:  true,
					Sensitive: true,
				},
			},
		},
	}
}

func linuxApplicationStackSchemaComputed() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"dotnet_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"go_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"node_version": { // Discarded by service if JavaVersion is specified
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ruby_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_server": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"java_server_version": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_image_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_registry_url": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_registry_username": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_registry_password": {
					Type:      pluginsdk.TypeString,
					Computed:  true,
					Sensitive: true,
				},
			},
		},
	}
}
