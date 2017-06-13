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

## Examples

### Show list record sets

```
$ ./gcdns list --keyfile company-abcdef012345.json --project high-office-XXX --mz ddns-example-net
+-----------------------------+------+------------------------------------------------------------------+
| Name                        | Type | Value                                                            |
+-----------------------------+------+------------------------------------------------------------------+
| ddns.example.net.           | NS   | ns-cloud-c1.googledomains.com.                                   |
|                             |      | ns-cloud-c2.googledomains.com.                                   |
|                             |      | ns-cloud-c3.googledomains.com.                                   |
|                             |      | ns-cloud-c4.googledomains.com.                                   |
+-----------------------------+------+------------------------------------------------------------------+
| ddns.example.net.           | SOA  | ns.googledomains.com. admin.google.com. 0 21600 3600 1209600 300 |
+-----------------------------+------+------------------------------------------------------------------+
| local.ddns.example.net.     | A    | 127.0.0.1                                                        |
+-----------------------------+------+------------------------------------------------------------------+
```

### Set A record

Set IP address of the host where you logged in

```
$ ./gcdns set --keyfile company-abcdef012345.json --project high-office-XXX --mz ddns-example-net \
  mine.ddns.example.net `curl ifconfig.io 2> /dev/null`
A record of mine.ddns.example.net has changed (new IP address: 12.34.56.78).
```

## References

* [
Google Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials)
