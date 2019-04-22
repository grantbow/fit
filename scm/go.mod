module scm

replace github.com/driusan/bug/scm => ../../../grantbow/bug/scm // fork

require github.com/driusan/bug/bugs v0.0.0

replace github.com/driusan/bug/bugs => ../../../grantbow/bug/bugs // fork

require (
	github.com/driusan/bug/bugapp v0.0.0
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace github.com/driusan/bug/bugapp => ../../../grantbow/bug/bugapp // fork
