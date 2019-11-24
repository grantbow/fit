module scm

replace github.com/driusan/bug/scm => ../../../grantbow/bug/scm // fork

require github.com/driusan/bug/bugs v0.0.0

replace github.com/driusan/bug/bugs => ../../../grantbow/bug/issues // fork

require (
	github.com/driusan/bug/bugapp v0.0.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/driusan/bug/bugapp => ../../../grantbow/bug/fitapp // fork

go 1.13
