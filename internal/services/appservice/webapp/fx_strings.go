package webapp

import (
	"strings"
)

func decodeApplicationStackLinux(fxString string) ApplicationStackLinux {
	parts := strings.Split(fxString, "|")
	result := ApplicationStackLinux{}
	if len(parts) != 2 {
		return result
	}

	switch parts[0] {
	case "DOTNETCORE", "DOTNET":
		result.NetFrameworkVersion = parts[1]

	case "NODE":
		result.NodeVersion = parts[1]

	case "JAVA", "TOMCAT", "JBOSSEAP":
		result.JavaServer = parts[0]
		javaParts := strings.Split(parts[1], "-")
		if len(javaParts) == 2 {
			// e.g. 8-jre8
			result.JavaServerVersion = javaParts[0]
			result.JavaVersion = javaParts[1]
		} else {
			// e.g. 8u242 or 11.0.9
			result.JavaVersion = parts[1]
		}

	case "PHP":
		result.PhpVersion = parts[1]

	case "PYTHON":
		result.PythonVersion = parts[1]

	case "RUBY":
		result.RubyVersion = parts[1]

	default: // DOCKER is the expected default here as "custom" images require it
		if dockerParts := strings.Split(parts[1], ":"); len(dockerParts) == 2 {
			result.DockerImage = dockerParts[0]
			result.DockerImageTag = dockerParts[1]
		}
	}

	return result
}
