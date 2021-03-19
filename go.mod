module github.com/drodil/envssh

go 1.14

replace github.com/drodil/envssh/ssh => ./ssh

require (
	github.com/drodil/envssh/ssh v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670
)
