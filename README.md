# GoLogic Micro-Components

A collection of micro-components that can be used individually or together that provide a clean api around many features used in developing web applications. These are a work in progress. Comments and documentation are currently sparse or non-existent.

Inspiration for the components api has been taken from the Laravel Framework for PHP. If you are familiar with that framework, you will probably find these components quite handy.

## Config

The Config component reads a file of key-value pairs. Keys may contain alpha-numeric characters, spaces, and underscores. Here is an example:

```
# This is a comment line. It will be ignored.

DB_CONNECTION=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_DATABASE=my_db
DB_USERNAME=my_user
DB_PASSWORD=password

```

### Usage:

```
config, _ := config.Load(".env")
dbName := config.Get("DB_NAME")
```

## Crypto

The Crypto component provides a wrapper around encryption, decryption, and key generation.

### Usage:

```
key := crypto.GenerateKey()
crypto := crypto.New(key)

message := "Here is a secret message"
encrypted, _ := crypto.Encrypt(message)
decrypted, _ := crypto.Decrypt(encrypted)

fmt.Println("Original: " + message)
fmt.Println("Encrypted: " + encrypted)
fmt.Println("Decrypted: " + decrypted)
```

## Hash

The Hash component provides a wrapper around password hashing and verification.

### Usage:

```
password := "secret"
hashed, _ := hash.Make(password)
valid := hash.Check(password, hashed)

fmt.Println(valid)
```

## Input

The Input component provides a wrapper around retrieving json or form input from a web api route.

### Usage:

```
func Handler(w http.ResponseWriter, r *http.Request) {
    input := input.Parse(r)
    if input.Has("name") {
        fmt.Println("Input has name with value: " + input.Get("name"))
    }
}
```

## Log

Coming soon.

## Mail

Coming soon.

## Validator

The Validator component provides a variety of validation functions. The validation rules and syntax follow that of the Laravel Framework as closely as possible.

### Usage:

```
// setup your rules
rules := make(map[string]string)
rules["name"] = "required|alpha"
rules["email"] = "required|email|hello"

// extend the validator by adding new functions
validator.AddValidator("hello", func(name string, value string, inputs map[string]string, params []string) bool {
    return false
}, "Hello input: %s")

// run the validator
input := input.Parse(r)
success, messages := validator.Validate(input.All(), rules)

fmt.Println(messages)

if success {
    fmt.Println("yay!")
} else {
    fmt.Println("oops!")
}
```

## License

MIT
