module main

replace github.com/driusan/bug => ../../grantbow/bug // fork

require github.com/driusan/bug/bugapp v0.0.0

replace github.com/driusan/bug/bugapp => ../../grantbow/bug/bugapp // fork

require github.com/driusan/bug/bugs v0.0.0

replace github.com/driusan/bug/bugs => ../../grantbow/bug/bugs // fork

require github.com/driusan/bug/scm v0.0.0

replace github.com/driusan/bug/scm => ../../grantbow/bug/scm // fork

require github.com/ghodss/yaml v1.0.0
