module github.com/drodil/envssh

go 1.14

replace github.com/drodil/envssh/ssh => ./ssh

replace github.com/drodil/envssh/util => ./util

require (
	github.com/drodil/envssh/ssh v0.0.0-00010101000000-000000000000
	github.com/drodil/envssh/util v0.0.0-00010101000000-000000000000 // indirect
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670
)
