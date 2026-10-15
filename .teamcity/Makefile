# .mvn/maven.config makes Maven order the classpath the way it did before 3.10. Without it the config
# generator picks up the Kotlin compiler's cut-down copy of jdom ahead of the real one and fails with
# `NoSuchMethodError: org.jdom.Element.<init>`. That file can't hold comments, so the reason lives here.
default: tools

tools:
	mvn -U dependency:resolve-sources

test: validate
	mvn test -DtestPackage=tests

validate:
	mvn teamcity-configs:generate
