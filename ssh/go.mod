module github.com/drodil/envssh/ssh

go 1.14

replace github.com/drodil/envssh/util => ../util

require (
	github.com/drodil/envssh/util v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1
)
