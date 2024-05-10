# go-bankgiro
Go library and CLI tool to encode and seal files sent to the Swedish Bankgiro

## Installation

### Binaries
You can download pre-built binaries from the releases page.

### Debian/Ubuntu Package Repository
You can download the package from the hoglandet package repository. For more information, visit [deb.pkg.hoglan.dev](https://git-registry.hoglandet.se/debian/repo-usage/).

### Building from source
To build the CLI tool from source, you need to have Go installed. You can then run the following command to build the CLI tool:

```bash
git clone https://github.com/hoglandets-it/go-bankgiro.git
cd go-bankgiro
go build ./cmd/bankgiro -o go-bankgiro
chmod +x go-bankgiro

# Optional: Move binary to PATH
sudo mv go-bankgiro /usr/local/bin/
```

## Usage
```bash
$ go-bankgiro
NAME:
   go-bankgiro - A new cli application

USAGE:
   go-bankgiro [global options] command [command options] 

DESCRIPTION:
   A tool to seal and validate Bankgiro files with HMAC

COMMANDS:
   seal, s      seal a file with a given key
   validate, v  validate a file with a given key
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### Seal a file
```bash
$ go-bankgiro seal --help

NAME:
   go-bankgiro seal - seal a file with a given key

USAGE:
   go-bankgiro seal [command options]key file

OPTIONS:
   --key value, -k value  key to seal the file with
   --kvv value, -v value  kvv to check the seal with (optional)
   --help, -h             show help
```
