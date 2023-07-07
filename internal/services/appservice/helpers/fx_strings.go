// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	FxStringPrefixDocker         FxStringPrefix = "DOCKER"
	FxStringPrefixDotNet         FxStringPrefix = "DOTNET"
	FxStringPrefixDotNetCore     FxStringPrefix = "DOTNETCORE"
	FxStringPrefixDotNetIsolated FxStringPrefix = "DOTNET-ISOLATED"
	FxStringPrefixGo             FxStringPrefix = "GO"
	FxStringPrefixJava           FxStringPrefix = "JAVA"
	FxStringPrefixJBoss          FxStringPrefix = "JBOSSEAP"
	FxStringPrefixNode           FxStringPrefix = "NODE"
	FxStringPrefixPhp            FxStringPrefix = "PHP"
	FxStringPrefixPowerShell     FxStringPrefix = "POWERSHELL"
	FxStringPrefixPython         FxStringPrefix = "PYTHON"
	FxStringPrefixRuby           FxStringPrefix = "RUBY"
	FxStringPrefixTomcat         FxStringPrefix = "TOMCAT"
)

type FxStringPrefix string

var urlSchemes = []string{
	"https://",
	"http://",
}

func decodeApplicationStackLinux(fxString string) ApplicationStackLinux {
	parts := strings.Split(fxString, "|")
	result := ApplicationStackLinux{}
	if len(parts) != 2 {
		return result
	}

	switch FxStringPrefix(strings.ToUpper(parts[0])) {
	case FxStringPrefixDotNetIsolated, FxStringPrefixDotNet, FxStringPrefixDotNetCore:
		result.NetFrameworkVersion = parts[1]

	case FxStringPrefixGo:
		result.GoVersion = parts[1]

	case FxStringPrefixNode:
		result.NodeVersion = parts[1]

	case FxStringPrefixJava:
		result.JavaServer = LinuxJavaServerJava
		javaParts := strings.Split(parts[1], "-")
		if strings.HasPrefix(parts[1], "8") {
			result.JavaVersion = "8"
		}
		if strings.HasPrefix(javaParts[0], "11") {
			result.JavaVersion = "11"
		}
		if strings.HasPrefix(javaParts[0], "17") {
			result.JavaVersion = "17"
		}
		result.JavaServerVersion = javaParts[0]

	case FxStringPrefixTomcat:
		result.JavaServer = LinuxJavaServerTomcat
		javaParts := strings.Split(parts[1], "-")
		if len(javaParts) == 2 {
			result.JavaServerVersion = javaParts[0]
			javaVersion := strings.TrimPrefix(javaParts[1], "jre")
			javaVersion = strings.TrimPrefix(javaVersion, "java")
			result.JavaVersion = javaVersion
		}

	case FxStringPrefixJBoss:
		result.JavaServer = LinuxJavaServerJboss
		javaParts := strings.Split(parts[1], "-")
		if len(javaParts) == 2 {
			result.JavaServerVersion = javaParts[0]
			javaVersion := strings.TrimPrefix(javaParts[1], "jre")
			javaVersion = strings.TrimPrefix(javaVersion, "java")
			result.JavaVersion = javaVersion
		}

	case FxStringPrefixPhp:
		result.PhpVersion = parts[1]

	case FxStringPrefixPython:
		result.PythonVersion = parts[1]

	case FxStringPrefixRuby:
		result.RubyVersion = parts[1]
	}

	return result
}

func EncodeFunctionAppLinuxFxVersion(input []ApplicationStackLinuxFunctionApp) *string {
	if len(input) == 0 || input[0].CustomHandler {
		return utils.String("")
	}

	appStack := input[0]
	var appType, appString string
	switch {
	case appStack.NodeVersion != "":
		appType = string(FxStringPrefixNode)
		appString = appStack.NodeVersion

	case appStack.DotNetVersion != "":
		if appStack.DotNetIsolated {
			appType = string(FxStringPrefixDotNetIsolated)
		} else {
			appType = string(FxStringPrefixDotNet)
		}
		appString = appStack.DotNetVersion

	case appStack.PythonVersion != "":
		appType = string(FxStringPrefixPython)
		appString = appStack.PythonVersion

	case appStack.JavaVersion != "":
		appType = string(FxStringPrefixJava)
		appString = appStack.JavaVersion

	case appStack.PowerShellCoreVersion != "":
		appType = string(FxStringPrefixPowerShell)
		appString = appStack.PowerShellCoreVersion

	case len(appStack.Docker) > 0 && appStack.Docker[0].ImageName != "":
		appType = string(FxStringPrefixDocker)
		dockerCfg := appStack.Docker[0]
		if dockerCfg.RegistryURL != "" {
			dockerUrl := dockerCfg.RegistryURL
			httpPrefixes := []string{"https://", "http://"}
			for _, prefix := range httpPrefixes {
				dockerUrl = strings.TrimPrefix(dockerUrl, prefix)
			}
			appString = fmt.Sprintf("%s/%s:%s", dockerUrl, dockerCfg.ImageName, dockerCfg.ImageTag)
		} else {
			appString = fmt.Sprintf("%s:%s", dockerCfg.ImageName, dockerCfg.ImageTag)
		}
	}

	return utils.String(fmt.Sprintf("%s|%s", appType, appString))
}

func DecodeFunctionAppLinuxFxVersion(input string) ([]ApplicationStackLinuxFunctionApp, error) {
	if input == "" {
		// This is a valid string for "Custom" stack which we picked up earlier, so we can skip here
		return nil, nil
	}

	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unrecognised LinuxFxVersion format received, got %s", input)
	}

	result := make([]ApplicationStackLinuxFunctionApp, 0)

	switch FxStringPrefix(strings.ToUpper(parts[0])) {
	case FxStringPrefixDotNet:
		appStack := ApplicationStackLinuxFunctionApp{DotNetVersion: parts[1]}
		result = append(result, appStack)

	case FxStringPrefixDotNetIsolated:
		appStack := ApplicationStackLinuxFunctionApp{DotNetVersion: parts[1], DotNetIsolated: true}
		result = append(result, appStack)

	case FxStringPrefixNode:
		appStack := ApplicationStackLinuxFunctionApp{NodeVersion: parts[1]}
		result = append(result, appStack)

	case FxStringPrefixPython:
		appStack := ApplicationStackLinuxFunctionApp{PythonVersion: parts[1]}
		result = append(result, appStack)

	case FxStringPrefixJava:
		appStack := ApplicationStackLinuxFunctionApp{JavaVersion: parts[1]}
		result = append(result, appStack)

	case FxStringPrefixPowerShell:
		appStack := ApplicationStackLinuxFunctionApp{PowerShellCoreVersion: parts[1]}
		result = append(result, appStack)

	case "docker":
		// This is handled as part of unpacking the app_settings using DecodeFunctionAppDockerFxString but included here for signposting as this is not intuitive.
	}

	return result, nil
}

func DecodeFunctionAppDockerFxString(input string, partial ApplicationStackDocker) ([]ApplicationStackDocker, error) {
	if input == "" {
		// This is a valid string for "Custom" stack which we picked up earlier, so we can skip here
		return nil, nil
	}

	parts := strings.Split(input, "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unrecognised LinuxFxVersion format received, got %s", input)
	}

	if !strings.EqualFold(parts[0], "docker") {
		return nil, fmt.Errorf("expected a docker FX version, got %q", parts[0])
	}

	dockerUrl := partial.RegistryURL
	for _, prefix := range urlSchemes {
		if strings.HasPrefix(dockerUrl, prefix) {
			dockerUrl = strings.TrimPrefix(dockerUrl, prefix)
			continue
		}
	}
	dockerParts := strings.Split(strings.TrimPrefix(parts[1], dockerUrl), ":")
	if len(dockerParts) != 2 {
		return nil, fmt.Errorf("invalid docker image reference %q", parts[1])
	}

	partial.ImageName = strings.TrimPrefix(dockerParts[0], "/")
	partial.ImageTag = dockerParts[1]

	return []ApplicationStackDocker{partial}, nil
}

func JavaLinuxFxStringBuilder(javaMajorVersion, javaServer, javaServerVersion string) (*string, error) {
	switch javaMajorVersion {
	case "8":
		{
			switch javaServer {
			case LinuxJavaServerJava:
				if strings.Contains(javaServerVersion, "u") {
					return pointer.To(fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion)), nil // e.g. JAVA|8u302
				} else {
					return pointer.To(fmt.Sprintf("%s|%s-jre8", LinuxJavaServerJava, javaServerVersion)), nil // e.g. "JAVA|8-jre8"
				}
			case LinuxJavaServerTomcat:
				if len(strings.Split(javaServerVersion, ".")) == 3 {
					return pointer.To(fmt.Sprintf("%s|%s-java8", LinuxJavaServerTomcat, javaServerVersion)), nil // e.g. TOMCAT|10.0.20-java8
				} else {
					return pointer.To(fmt.Sprintf("%s|%s-jre8", LinuxJavaServerTomcat, javaServerVersion)), nil // e.g. TOMCAT|10.0-jre8
				}
			case LinuxJavaServerJboss:
				return pointer.To(fmt.Sprintf("%s|%s-java8", LinuxJavaServerJboss, javaServerVersion)), nil
			}
		}
	case "11":
		switch javaServer {
		case LinuxJavaServerJava:
			if len(strings.Split(javaServerVersion, ".")) == 3 {
				return pointer.To(fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion)), nil // e.g. JAVA|11.0.13
			} else {
				return pointer.To(fmt.Sprintf("%s|%s-java11", LinuxJavaServerJava, javaServerVersion)), nil // e.g.JAVA|11-java1
			}
		case LinuxJavaServerTomcat:
			return pointer.To(fmt.Sprintf("%s|%s-java11", LinuxJavaServerTomcat, javaServerVersion)), nil // e.g. TOMCAT|10.0-java11 and TOMCAT|10.0.20-java11

		case LinuxJavaServerJboss:
			return pointer.To(fmt.Sprintf("%s|%s-java11", LinuxJavaServerJboss, javaServerVersion)), nil // e.g. TOMCAT|10.0-java11 and TOMCAT|10.0.20-java11// e.g. JBOSSEAP|7-java11 / JBOSSEAP|7.4.2-java11
		}

	case "17":
		switch javaServer {
		case LinuxJavaServerJava:
			if len(strings.Split(javaServerVersion, ".")) == 3 {
				return pointer.To(fmt.Sprintf("%s|%s", LinuxJavaServerJava, javaServerVersion)), nil // "JAVA|17.0.2"
			} else {
				return pointer.To(fmt.Sprintf("%s|%s-java17", LinuxJavaServerJava, javaServerVersion)), nil // "JAVA|17-java17"
			}

		case LinuxJavaServerTomcat:
			return pointer.To(fmt.Sprintf("%s|%s-java17", LinuxJavaServerTomcat, javaServerVersion)), nil // e,g, TOMCAT|10.0-java17 / TOMCAT|10.0.20-java17
		case LinuxJavaServerJboss:
			return nil, fmt.Errorf("java 17 is not supported on %s", LinuxJavaServerJboss)
		default:
			return pointer.To(fmt.Sprintf("%s|%s-java17", javaServer, javaServerVersion)), nil
		}

	default:
		return pointer.To(fmt.Sprintf("%s|%s-%s", javaServer, javaServerVersion, javaMajorVersion)), nil

	}
	return nil, fmt.Errorf("unsupported combination of `java_version`, `java_server`, and `java_server_version`")
}

func EncodeDockerFxString(image string, registryUrl string) string {
	template := "DOCKER|%s/%s"

	registryUrl = trimURLScheme(registryUrl)

	return fmt.Sprintf(template, registryUrl, image)
}

func EncodeDockerFxStringWindows(image string, registryUrl string) string {
	template := "DOCKER|%s/%s"
	dockerHubTemplate := "DOCKER|%s"

	registryUrl = trimURLScheme(registryUrl)

	// Windows App Services fail to auth if the index portion of the image ref is `index.docker.io` so it must not be included.
	if strings.EqualFold(registryUrl, "index.docker.io") {
		return fmt.Sprintf(dockerHubTemplate, image)
	}

	return fmt.Sprintf(template, registryUrl, image)
}

func FxStringHasPrefix(input string, prefix FxStringPrefix) bool {
	return strings.HasPrefix(strings.ToUpper(input), string(prefix))
}
