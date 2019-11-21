module github.com/grantbow/bug/main // fork

// for go.mod fork https://github.com/golang/go/issues/28514

require github.com/driusan/bug/bugapp v0.0.0

replace github.com/driusan/bug/bugapp => ../../grantbow/bug/fitapp // fork

require github.com/driusan/bug/bugs v0.0.0

replace github.com/driusan/bug/bugs => ../../grantbow/bug/issues // fork

require github.com/driusan/bug/scm v0.0.0

replace github.com/driusan/bug/scm => ../../grantbow/bug/scm // fork

require (
	github.com/FabianWe/etherpadlite-golang v0.0.0-20190415145731-46b2da95f3b7 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/matttproud/gochecklist v0.0.0-20150912192500-26fd8564d1e9 // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	golang.org/x/review v0.0.0-20190422220318-83908358f3a5 // indirect
)

go 1.13
