module bug-import

//replace github.com/driusan/bug => ../../../grantbow/bug // fork
replace github.com/driusan/bug/bug-import => ../../../grantbow/bug/bug-import // fork

//require github.com/driusan/bug/bugapp v0.0.0

//replace github.com/driusan/bug/bugapp => ../../../grantbow/bug/bugapp // fork

//require github.com/driusan/bug/bugs v0.0.0

//replace github.com/driusan/bug/bugs => ../../../grantbow/bug/bugs // fork

require (
	github.com/driusan/bug v0.3.1
	//github.com/driusan/bug/scm v0.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
)

//replace github.com/driusan/bug/scm => ../../../grantbow/bug/scm // fork
