module bugapp

replace github.com/driusan/bug/bugapp => ../../../grantbow/bug/bugapp // fork

require github.com/driusan/bug/bugs v0.0.0

replace github.com/driusan/bug/bugs => ../../../grantbow/bug/bugs // fork

require github.com/driusan/bug/scm v0.0.0

replace github.com/driusan/bug/scm => ../../../grantbow/bug/scm // fork

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/google/go-github v17.0.0+incompatible
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
)
