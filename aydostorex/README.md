# AYDO Store X

## Usage
```bash
NAME:
   AydoStoreX -

USAGE:
   aydostorex [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config, -c "config.toml"	Path to the configuration file
   --log, -l "4"		Level of logging (0 less verbose, to 5 most verbose) default to 4
   --debug			Launch web server in debug mode
   --help, -h			show help
   --version, -v		print the versi
```

## configuration file
```toml
listen_addr=":8080"
store_root="/tmp/aydostorex"
authentification="path/to/account.toml"
sources=["https://store2.com", "https://otherstore.org"]

[tls]
certificate="cert.pem"
key="key.rsa"
```
- **listen_addr** : address and port on which the service need to listen.
- **store_root** : root directory where the files will be stored
- **authentification** : path to a toml file that contains the account defining the ACL for every users. Leave empty if you don't want authentification/authrorization
- **sources** : The sources array contains the address of other stores. When a file can't be found on the local store, the request if forwarder to the stores defined in sources. This allow you to create some layers of stores.
- **TLS  section**:
if both certificate and key are specified, service will enable HTTPS
    - **certificate** : Path to the certificate
    - **certificate** : Path to the key

example of an ACL file:  
```toml
[[accounts]]
login="admin"
password="rooter"
read=true
write=true
delete=true
[[accounts]]
login="public"
password="public"
read=true
write=false
delete=false
```
