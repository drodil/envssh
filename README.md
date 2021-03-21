# envssh - Take your environment with you

Environment (friendly) SSH client.

Brings your environment with you to the remote machine including 
environment variables, configuration files and other.

**Note: This is very POC implementation and needs a lot of refactoring and more testing!**

## Installing

```bash
go get -u github.com/drodil/envssh
```

## Configuration

Configuration file is created automatically to $HOME/.ssh/envssh.yml. 
The file contains the following sections:

* global - Global level configuration which can be overriden per hostname
	* env - Environment variables to move to the remote 
		* static - Key-value pairs that are moved to the remote
		* moved - Keys from current environment that are evaluated and moved
	* files - List of files that are moved from local to the remote, each containing:
		* local - Local path of the file to move
		* remote - Destination path on the remote
	* commands - List of commands that are run on remote before starting the session
* servers - Server specific configurations that override the global ones based on hostname
	* host - Hostname of the remote
	* port - Port that is used to connect to remote
	* env - See above
	* files - See above
	* commands - See above

See example configuration file [envssh.yml](envssh.yml)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).
