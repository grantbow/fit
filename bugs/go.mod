module bugs

replace github.com/driusan/bug/bugs => ../../../grantbow/bug/bugs // fork

require github.com/driusan/bug/scm v0.0.0

replace github.com/driusan/bug/scm => ../../../grantbow/bug/scm // fork

require github.com/driusan/bug/bugapp v0.0.0

replace github.com/driusan/bug/bugapp => ../../../grantbow/bug/bugapp // fork

require (
	github.com/ghodss/yaml v1.0.0
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
