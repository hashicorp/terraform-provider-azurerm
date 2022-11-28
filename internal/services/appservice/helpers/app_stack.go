package helpers

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApplicationStackWindows struct {
	NetFrameworkVersion     string `tfschema:"dotnet_version"`
	PhpVersion              string `tfschema:"php_version"`
	JavaVersion             string `tfschema:"java_version"`
	PythonVersion           string `tfschema:"python_version"`
	NodeVersion             string `tfschema:"node_version"`
	JavaContainer           string `tfschema:"java_container"`
	JavaContainerVersion    string `tfschema:"java_container_version"`
	DockerContainerName     string `tfschema:"docker_container_name"`
	DockerContainerRegistry string `tfschema:"docker_container_registry"`
	DockerContainerTag      string `tfschema:"docker_container_tag"`
	CurrentStack            string `tfschema:"current_stack"`
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
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
					},
				},

				"dotnet_core_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.dotnet_core_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
					},
					Description: fmt.Sprintf(`The version of DotNet to use.`),
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"7.4", // Deprecated
						"Off", // Really?!?! Should be `AutoUpdate` or `Latest` ?
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.python",
					},
				},

				"python_version": {
					Type:       pluginsdk.TypeString,
					Optional:   true,
					Computed:   true,
					Deprecated: "This property is deprecated. Values set are not used by the service.",
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
					},
					ConflictsWith: []string{
						"site_config.0.application_stack.0.python",
					},
				},

				"python": {
					Type:       pluginsdk.TypeBool,
					Optional:   true,
					Default:    false,
					Deprecated: "This property is deprecated. Values set are not used by the service.",
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.python",
					},
					ConflictsWith: []string{
						"site_config.0.application_stack.0.python_version",
					},
				},

				"node_version": {
					// Not used directly - Set via app_settings.0.WEBSITE_NODE_DEFAULT_VERSION
					// deprecate? how?
					Type:       pluginsdk.TypeString,
					Optional:   true,
					Computed:   true,
					Deprecated: "This property is no longer configurable, please set `WEBSITE_NODE_DEFAULT_VERSION` in `app_settings`.",
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
					},
					ConflictsWith: []string{
						"site_config.0.application_stack.0.node",
					},
				},

				"java_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
					},
				},

				"java_embedded_server_enabled": {
					Type:     pluginsdk.TypeBool,
					Optional: true,
					Default:  true,
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
				},

				"java_container_version": {
					Type:       pluginsdk.TypeString,
					Optional:   true,
					Deprecated: "This property has been deprecated in favour of `tomcat_version` and `java_embedded_server_enabled`",
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_container",
					},
				},

				"docker_container_name": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_container_name",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
					},
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_container_tag",
					},
				},

				"docker_container_registry": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"docker_container_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_container_name",
					},
				},

				"current_stack": {
					Type:       pluginsdk.TypeString,
					Optional:   true,
					Computed:   true,
					Deprecated: "This value has been deprecated. Values set here are ignored and are configured automatically based on the choice of application software.",
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

				"docker_container_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_container_registry": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_container_tag": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"current_stack": {
					Type:     pluginsdk.TypeString,
					Computed: true,
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
	DockerImageTag      string `tfschema:"docker_image_tag"`
	DockerImage         string `tfschema:"docker_image"`
	RubyVersion         string `tfschema:"ruby_version"`
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
					ValidateFunc: validation.StringInSlice([]string{ // TODO replace with major.minor regex?
						"3.1",
						"5.0", // deprecated
						"6.0",
						"7.0",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
					},
				},

				"go_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{ // TODO replace with major.minor regex?
						"1.19",
						"1.18",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
					},
				},

				"php_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{ // TODO replace with major.minor regex?
						"7.4",
						"8.0",
						"8.1",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
					},
				},

				"python_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{ // TODO replace with major.minor regex?
						"3.7",
						"3.8",
						"3.9",
						"3.10",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
					},
				},

				"node_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"12-lts",
						"14-lts",
						"16-lts",
						"18-lts",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
					},
				},

				"ruby_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{ // TODO replace with major.minor regex?
						"2.6", // Deprecated - accepted but not offered in the portal. Remove in 4.0
						"2.7",
					}, false),
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
					},
				},

				"java_version": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty, // There a significant number of variables here, and the versions are not uniformly formatted.
					// TODO - Needs notes in the docs for this to help users navigate the inconsistencies in the service. e.g. jre8 va java8 etc
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
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
						"site_config.0.application_stack.0.java_version",
					},
				},

				"java_server_version": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					RequiredWith: []string{
						"site_config.0.application_stack.0.java_server",
					},
				},

				"docker_image": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					ExactlyOneOf: []string{
						"site_config.0.application_stack.0.docker_image",
						"site_config.0.application_stack.0.dotnet_version",
						"site_config.0.application_stack.0.java_version",
						"site_config.0.application_stack.0.node_version",
						"site_config.0.application_stack.0.php_version",
						"site_config.0.application_stack.0.python_version",
						"site_config.0.application_stack.0.ruby_version",
						"site_config.0.application_stack.0.go_version",
					},
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_image_tag",
					},
				},

				"docker_image_tag": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{
						"site_config.0.application_stack.0.docker_image",
					},
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

				"docker_image": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"docker_image_tag": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
