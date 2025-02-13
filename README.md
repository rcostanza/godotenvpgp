# godotenvpgp - PGP-encrypted .env files

[![test status](https://github.com/rcostanza/godotenvpgp/actions/workflows/go.yml/badge.svg)](https://github.com/rcostanza/godotenvpgp/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/rcostanza/godotenvpgp)](https://goreportcard.com/report/github.com/rcostanza/godotenvpgp)
[![Go Coverage](https://github.com/rcostanza/godotenvpgp/wiki/coverage.svg)](https://raw.githack.com/wiki/rcostanza/godotenvpgp/coverage.html)

Utility package to use & load encrypted .env files.

# Table of contents

1. [Why](#why)
2. [Features](#features)
3. [Basic usage](#basic-usage)
    - [Library](#library)
    - [Cli tool](#cli-tool)
4. [Use cases](#use-cases)
    - [Specific values per environment on the same file](#use-cases-1)
    - [Separate .env files per environment](#use-cases-2)

## <a id="why"></a>Why?

- In some scenarios, managing (some) environment variables directly on a file can be an easier way, instead of being managed over the infrastructure domain (aka. you don't have to bug your devops teams as often);
- An encrypted .env file can be safely versioned along with your code;
- It makes easier to distribute development credentials among a team of developers.

## <a id="features"></a>Features

- Automatically loads any encrypted env files found with the patterns `.env.encrypted` or `.env.<environment>.encrypted`
- A bundled cli tool to encrypt, decrypt, and show contents of encrypted files
- Automatically preloads a regular `.env` file if present
- Supports environment-specific variables
    - Used to store variables for multiple environments on a single file
- Supports multiple environment-specific .env files
    - Used to isolate environments variables by environment, each on a separate file
    - Each file uses its own password

## <a id="basic-usage"></a>Basic usage

### <a id="library"></a>Library
Install the package:

```bash
go install github.com/rcostanza/godotenvpgp@latest
```

Create an environment variable to store the file password, either using a regular `.env` file or otherwise:

```ini
ENVFILE_PASSWORD=<password>
```

Create an unencrypted file named `.env.unencrypted` containing the variables that need to be encrypted.<br />
⚠️ **Warning: `.unencrypted` files should not be versioned** ⚠️

```ini
MY_VAR=123
```

Use the cli tool to encrypt the file. A `.env.encrypted` file will be generated:

```bash
godotenvpgp encrypt
```

Finally, either manually invoke it:

```go
import (
    "fmt"
    "os"
    "github.com/rcostanza/godotenvpgp/envfile"
)

func main() {
    envfile.Load() // Loads both .env and .env.encrypted
    fmt.Println(os.Getenv("MY_VAR")) // 123
}
```

Or autoload it:

```go
import (
    "fmt"
    "os"
    _ "github.com/rcostanza/godotenvpgp/autoload" // Autoloads both .env and .env.encrypted
)

func main() {
    fmt.Println(os.Getenv("MY_VAR")) // 123
}
```

### <a name="cli-tool"></a>Cli tool

Used to encrypt, decrypt, and read contents of an encrypted file.

```bash
godotenvpgp <command>

  encrypt       Encrypt all unencrypted .env files (.env*.unencrypted) to .env.*.encrypted files
  decrypt       Decrypt all encrypted .env files (.env*.encrypted) to .env.*.unencrypted files
  show <file>   Show the content of an encrypted .env file
```

## <a id="use-cases"></a>Use cases

### <a id="use-cases-1"></a>Specific values per environment on the same file:

When having multiple environments, it's possible to have different values for a variable depending on the environment on the same encrypted file. You can use optional `[tags]` to define values that are specific to certain environments.

Using the example keys below:

```ini
V1=123
[dev]
V2=456dev
V3=789dev
[production]
V2=456production
V3=789production
```

If you have an environment variable (declared outside the encrypted file) named `env/ENV` or `environment/ENVIRONMENT` with the value of `dev` (case sensitive), then V2 equals `123dev` and V3 equals `456dev`.

Variables defined below a tag will _only_ apply if the environment matches. For example, if your environment is not set, or set to `hml`, V2 and V3 will not exist.

Variables defined before any tags will apply to any environment. In this example, the variable V1 will always be present.


### <a id="use-cases-2"></a>Separate .env files per environment:

When having multiple environments and using the solution above to keep all variables in a single file, if allowing access to variables from all environments to anyone with the decryption password is undesirable, multiple encripted .env files, each with its own password, can be used for more granular access. That way it's possible, for example, to more freely distribute the decryption password of development variables to devs, while keeping the decription password of production variables more restricted and thus keeping the production encrypted file inaccessible.

Create an unencrypted file for each environment, providing variables for each:

```
.env.dev.unencrypted
.env.testing.unencrypted
.env.production.unencrypted
```

Define in your environment variables, either using a regular .env file or otherwise, variables with the password for each file, using the naming pattern `ENVFILE_PASSWORD_<environment name uppercase>`:

```ini
env=dev # case sensitive
ENVFILE_PASSWORD_DEV=123
ENVFILE_PASSWORD_TESTING=456
ENVFILE_PASSWORD_PRODUCTION=789
```

Encrypt the files using the cli tool; a separate encrypted file will be generated for each (`.env.dev.encrypted`, `.env.testing.encrypted` and `.env.production.encrypted`):

```bash
godotenvpgp encrypt
```

When using separate files, only the file matching the environment defined in the environment variable `env/ENV` or `environment/ENVIRONMENT` will be loaded. If no environment is specified, only the default `.env.encrypted` file will be loaded, if present.

For these unused, environment-specific files, the respective environment variable containing their password is not required.
