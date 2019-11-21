module fitapp

replace github.com/driusan/bug/bugapp => ../../../grantbow/bug/fitapp // fork

require github.com/driusan/bug/bugs v0.0.0

replace github.com/driusan/bug/bugs => ../../../grantbow/bug/issues // fork

require github.com/driusan/bug/scm v0.0.0

replace github.com/driusan/bug/scm => ../../../grantbow/bug/scm // fork

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
)

go 1.13
