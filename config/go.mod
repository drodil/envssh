module github.com/drodil/envssh/config

go 1.14

replace github.com/drodil/envssh/util => ../util

require (
	github.com/drodil/envssh/util v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v2 v2.4.0
)
