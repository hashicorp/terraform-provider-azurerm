default: tools

tools:
	mvn -U dependency:sources

test: validate
	mvn test -DtestPackage=tests

validate:
	mvn teamcity-configs:generate
