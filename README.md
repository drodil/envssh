# envssh

Environment (friendly) SSH client.

Brings your environment with you to the remote machine.

## TODO

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
* Support for ~/.ssh/config
	* Use possible SSH configs from there in connections
* More CLI options
	* Configuration file to use
	* Log file
	* Identity file
	* Port
	* Command to run instead starting session
	* Other useful options from plain ssh
* SSH client
	* Print banner
	* Private key authentication
	* Password prompt
	* Moving files to remote
	* Setting environment variables on remote
	* Session starting
	* Running commands
* Session clean up
	* When connection is closed, the user configuration files can be cleaned
	* Requires making backups of existing config files if any on remote and
	  restoring the backups
* Other things
	* Use colors
	* Add more tests
