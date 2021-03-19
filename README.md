# envssh

Environment (friendly) SSH client.

Brings your environment with you to the remote machine.

## Installing

TODO

## Configuration

TODO

## Contributing

TODO

## TODO

*More todos in the code!*

*Should be moved to be github issues instead*

* Configuration file support
    * Preferable YAML, maybe JSON support
	* Global config for all hosts
	* Host specific config
	* Environment variables to move
		* Static environment variables to set on connect f.eg. FOO=BAR
		* List of local environment variables to move f.eg. [EDITOR, VISUAL,
		  HTTP_PROXY, LANG]
		* Wildcard for all env variables f.eg. [*] though it should skip
		  defaults such as PATH, PWTD...
	* Configuration files to move
		* For example vim/bash/...
	* Configuration location perhaps ~/.ssh/envssh.yaml
	* Automatic creation with sane defaults if file does not exist
	* Config handling to own module
	* PTY term
	* Clean up / not to clean up after disconnect
	* Automatic installation of packages on connect (have to think about this..
	  requires sudo and checking package manager type etc.)
* Support for ~/.ssh/config
	* Use possible SSH configs from there in connections
	* Tunneling/proxy/etc.?
* More CLI options
	* Configuration file to use
	* Log file
	* Identity file
	* Port
	* Command to run instead starting session
	* To clean up/not to clean up after disconnect
	* Other useful options from plain ssh
	* Self-extraction on/off
* SSH client
	* envssh self-extraction to remote
		* Including config file
		* This to allow jumping from server to server and getting all stuff with
		  you all the time
	* Print banner before auth
	* Private key authentication
	* SSHAgent support
	* Improve and test moving files to/from remote
	* Running commands
	* Fix bugs with the interactive session
	* Automatic resizing of interactive session based on local terminal
	* Getting env variable(s) from remote
* Session clean up
	* When connection is closed, the user configuration files can be cleaned
	* Requires making backups of existing config files if any on remote and
	  restoring the backups
* Other things
	* Use colors
	* Add more tests
	* Interactive shell does not support up/down arrows or other special keys

