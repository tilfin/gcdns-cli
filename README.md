# gcdns

## Build

```
$ go build -o gcdns main.go
```

### Windows

```
$ GOOS=windows GOARCH=386 go build -o gcdns.exe main.go
```

## Usage

```
$ ./gcdns --help
usage: gcdns --project=PROJECT --mz=MZ [<flags>] <command> [<args> ...]

Google Cloud DNS CLI

Flags:
  --help             Show context-sensitive help (also try --help-long and
                     --help-man).
  --keyfile=KEYFILE  JSON key file
  --project=PROJECT  Project name of Google Cloud
  --mz=MZ            Target managed zone of project

Commands:
  help [<command>...]
    Show help.

  list
    list record

  set <host> <ip>
    set record

```

### References

* [
Google Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)
