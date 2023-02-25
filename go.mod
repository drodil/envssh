module github.com/drodil/envssh

go 1.14

replace github.com/drodil/envssh/ssh => ./ssh

replace github.com/drodil/envssh/util => ./util

replace github.com/drodil/envssh/config => ./config

require (
	github.com/drodil/envssh/config v0.0.0-00010101000000-000000000000
	github.com/drodil/envssh/ssh v0.0.0-00010101000000-000000000000
	github.com/drodil/envssh/util v0.0.0-00010101000000-000000000000
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.0.0-20210317153231-de623e64d2a6 // indirect
)
