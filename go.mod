module github.com/grantbow/bug/main // fork

// for go.mod fork https://github.com/golang/go/issues/28514

require github.com/driusan/bug/bugapp v0.0.0

replace github.com/driusan/bug/bugapp => ../../grantbow/bug/bugapp // fork

require github.com/driusan/bug/bugs v0.0.0

replace github.com/driusan/bug/bugs => ../../grantbow/bug/bugs // fork

require github.com/driusan/bug/scm v0.0.0

replace github.com/driusan/bug/scm => ../../grantbow/bug/scm // fork

require (
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-querystring v1.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)
