default: tools

tools:
	mvn -U dependency:resolve-sources

test: validate
	mvn test -DtestPackage=tests

validate:
	mvn teamcity-configs:generate
