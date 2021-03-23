# envssh - Take your environment with you

Environment (friendly) SSH client.

Brings your environment with you to the remote machine including
environment variables, configuration files and other.

**Note: This is POC implementation and needs a lot of refactoring and more testing!**

## Installing

```bash
go get -u github.com/drodil/envssh
```

## Usage

Tool behaves much like the normal ssh client. To connect to remote, simply run:

```bash
envssh [remote]
```

The remote should be either hostname or IP address and can include username and
port which are used to connect f.eg. user@remote:1234.

On the first run, default configuration file will be created to
$HOME/.ssh/envssh.yml. See Configuration section for details how to configure
the tool.

Run `envssh --help` for more options.

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
	* aliases - List of aliases for this server, can be used instead hostname in destination
	* env - See above
	* files - See above
	* commands - See above

See example configuration file [envssh.yml](envssh.yml)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).
