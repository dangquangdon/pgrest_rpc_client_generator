# Code generator for PostgREST Http Client

Generating Http Client for a RPC functions in a PostgREST API server.

# Usage

```
NAME:
   pgr-gen - Generate HTTP Client which calls the RPC functions defined in PostgREST server.

USAGE:
   pgr-gen [global options] command [command options]

COMMANDS:
   generate           Generate all the codes for types and requests
   generate-types     Generate only the types
   generate-requests  Generate only the functions for http requests
   help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --destination value, -d value  Destination directory where the code is going to be stored (default: ".")
   --url value, -u value          Base url of the PostgREST server
   --client-id value, -c value    Value for the User-Agent Header
   --help, -h                     show help
```

# Example

```
gpr-gen -u http://localhost:3000 -c 007-agent generate
```