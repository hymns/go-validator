# go-validator

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/hymns/go-validator)](https://github.com/hymns/go-validator/releases) [![Go Version](https://img.shields.io/badge/go-1.24.0-blue.svg)](https://golang.org/dl/) [![Go Report Card](https://goreportcard.com/badge/github.com/hymns/go-validator)](https://goreportcard.com/report/github.com/hymns/go-validator) [![GoDoc](https://godoc.org/github.com/hymns/go-validator?status.svg)](https://pkg.go.dev/github.com/hymns/go-validator) [![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A Laravel-inspired validation package for Go. No struct tags required — pass your data as a plain map and declare rules as pipe-separated strings or use the fluent typed builder.

## Installation

```bash
go get github.com/hymns/go-validator
```

## Quick start

```go
import validator "github.com/hymns/go-validator"

v := validator.Make(
    validator.Input{
        "email":    "user@example.com",
        "password": "secret",
        "age":      "17",
    },
    validator.Rules{
        "email":    "required|email",
        "password": "required|min:8",
        "age":      "required|integer|min:18",
    },
)

if v.Fails() {
    fmt.Println(v.Errors()) // map[age:["The age must be at least 18."] password:["The password must be at least 8 characters."]]
}
```

Or use the **typed rule builder** for compile-time safety and IDE autocompletion:

```go
v := validator.Make(input, validator.Rules{
    "email": validator.R().Required().Email().Max(100).Build(),
    "age":   validator.R().Required().Integer().Min(18).Build(),
    "role":  validator.R().Required().In("admin", "user").Build(),
})
```

## API

### `Make(input, rules) *Validator`

Creates a validator. Validation is lazy — it runs on the first call to `Fails()`, `Passes()`, or `Errors()`.

### `(*Validator).Messages(msgs) *Validator`

Override error messages. Call before `Fails()` / `Passes()`.

```go
v := validator.Make(input, rules).Messages(validator.Messages{
    "email.required": "Email address is mandatory.",
    "email.email":    "That doesn't look like a valid email.",
})
```

### `(*Validator).Bail() *Validator`

Stop validation after the first field that produces an error (global bail). Useful when later fields depend on earlier ones passing.

```go
v := validator.Make(input, rules).Bail()
```

For per-field bail (stop checking remaining rules on that field on first failure), add `bail` to the rule string:

```go
validator.Rules{
    "email": "bail|required|email|unique:users,email",
    //         ↑ stops checking email rules on first failure
}
```

### `(*Validator).WithDB(db *sql.DB) *Validator`

Attach a database connection for `unique` and `exists` rules.

### `(*Validator).Fails() bool`

Returns `true` if validation failed.

### `(*Validator).Passes() bool`

Returns `true` if validation passed.

### `(*Validator).Errors() ErrorBag`

Returns all errors. `ErrorBag` is `map[string][]string`.

```go
errs := v.Errors()
errs.Has("email")      // bool
errs.First("email")    // string — first error message for the field
```

## Rules

Rules are pipe-separated strings. Parameters are colon-separated from the rule name, and multiple parameters are comma-separated.

```
"required|string|min:3|max:100"
"required|in:admin,user,editor"
"required|unique:users,email"
```

### Presence & required

| Rule | Description |
|------|-------------|
| `required` | Field must be present and non-empty |
| `nullable` | Field may be `nil`; skips all other rules when `nil` |
| `present` | Key must exist in input (value may be empty) |
| `filled` | If key is present, value must not be empty |
| `bail` | Stop checking remaining rules for this field on first failure |

### Conditional required

| Rule | Description |
|------|-------------|
| `required_if:field,value` | Required when another field equals a value |
| `required_unless:field,value` | Required unless another field equals a value |
| `required_with:field1,field2` | Required when any of the listed fields is present |
| `required_without:field1,field2` | Required when any of the listed fields is absent |
| `required_with_all:field1,field2` | Required when all listed fields are present |
| `required_without_all:field1,field2` | Required when all listed fields are absent |

### Prohibition

| Rule | Description |
|------|-------------|
| `prohibited` | Field must be absent or empty |
| `prohibited_if:field,value` | Prohibited when another field equals a value |
| `prohibited_unless:field,value` | Prohibited unless another field equals a value |

### Strings

| Rule | Description |
|------|-------------|
| `string` | Must be a string |
| `alpha` | Letters only (Unicode) |
| `alpha_num` | Letters and numbers only |
| `alpha_dash` | Letters, numbers, hyphens, and underscores |
| `email` | Valid email format |
| `url` | Valid URL (http/https) |
| `active_url` | URL with resolvable host |
| `uuid` | Valid UUID (any version) |
| `ulid` | Valid ULID |
| `hex_color` | Hex color code (`#RGB`, `#RGBA`, `#RRGGBB`, `#RRGGBBAA`) |
| `json` | Valid JSON string |
| `lowercase` | All lowercase |
| `uppercase` | All uppercase |
| `starts_with:a,b` | Must start with one of the given values |
| `ends_with:a,b` | Must end with one of the given values |
| `doesnt_start_with:a,b` | Must not start with any of the given values |
| `doesnt_end_with:a,b` | Must not end with any of the given values |
| `regex:pattern` | Must match the regex pattern |
| `not_regex:pattern` | Must not match the regex pattern |
| `same:field` | Must match another field's value |
| `different:field` | Must differ from another field's value |
| `confirmed` | Must have a matching `{field}_confirmation` field |

> **Note:** Regex patterns cannot contain `|` (the rule separator). Use `Extend()` for complex patterns.

### Numeric

| Rule | Description |
|------|-------------|
| `integer` | Must be an integer |
| `numeric` | Must be numeric (int or float) |
| `decimal:min,max` | Must have between `min` and `max` decimal places |
| `digits:n` | Exactly `n` digits |
| `digits_between:min,max` | Between `min` and `max` digits |
| `multiple_of:n` | Must be a multiple of `n` |
| `min:n` | Minimum value (numeric) or minimum length (string/array) |
| `max:n` | Maximum value (numeric) or maximum length (string/array) |
| `between:min,max` | Between `min` and `max` (numeric or string length) |
| `size:n` | Exact size — character count for strings, item count for arrays, numeric value |
| `gt:n` or `gt:field` | Greater than value or another field |
| `gte:n` or `gte:field` | Greater than or equal |
| `lt:n` or `lt:field` | Less than value or another field |
| `lte:n` or `lte:field` | Less than or equal |

### Boolean

| Rule | Description |
|------|-------------|
| `boolean` | Must be a boolean-like value (`true`, `false`, `1`, `0`, `"1"`, `"0"`) |
| `accepted` | Must be an accepted value (`yes`, `on`, `1`, `true`) |
| `declined` | Must be a declined value (`no`, `off`, `0`, `false`) |
| `accepted_if:field,value` | Must be accepted when another field equals a value |
| `declined_if:field,value` | Must be declined when another field equals a value |

### Dates

Supported date layouts: `2006-01-02`, `02/01/2006`, `01/02/2006`, `2006-01-02 15:04:05`, `2006-01-02T15:04:05Z07:00`.

| Rule | Description |
|------|-------------|
| `date` | Valid date string |
| `date_format:layout` | Matches the given Go time layout |
| `date_equals:date` | Equals the given date |
| `before:date` | Must be before the given date |
| `before_or_equal:date` | Must be before or equal to the given date |
| `after:date` | Must be after the given date |
| `after_or_equal:date` | Must be after or equal to the given date |
| `timezone` | Valid IANA timezone name (e.g. `Asia/Kuala_Lumpur`) |

### Network

| Rule | Description |
|------|-------------|
| `ip` | Valid IPv4 or IPv6 address |
| `ipv4` | Valid IPv4 address |
| `ipv6` | Valid IPv6 address |
| `mac_address` | Valid MAC address (colon or hyphen separated) |

### Arrays

| Rule | Description |
|------|-------------|
| `array` | Must be a `[]any` |
| `distinct` | Array items must be unique |
| `in:a,b,c` | Value must be one of the listed values |
| `not_in:a,b,c` | Value must not be one of the listed values |

### Database

Requires `.WithDB(db)`. Returns a `db_required` error if no DB is attached.

| Rule | Description |
|------|-------------|
| `unique:table,column` | Value must not exist in the given column |
| `unique:table,column,ignoreValue,ignoreColumn` | Unique but ignore a specific row (for updates) |
| `exists:table,column` | Value must exist in the given column |

```go
// On create
"email": "required|email|unique:users,email"

// On update — ignore the current user's row
"email": "required|email|unique:users,email,42,id"
```

`ignoreColumn` defaults to `id` if omitted.

## Custom messages

Override any message globally with `field.rule` keys, or just `rule` for field-agnostic overrides.

```go
v := validator.Make(input, rules).Messages(validator.Messages{
    "email.required":  "Please provide your email.",
    "email.email":     "That email address is not valid.",
    "password.min":    "Password needs at least :param characters.",
})
```

Available placeholders: `:field`, `:param`, `:other`, `:value`, `:min`, `:max`.

## Typed rule builder

`R()` returns a `*RuleBuilder`. Chain methods, then call `Build()` to produce the pipe-separated string. Both styles are fully compatible with `Make()`.

```go
validator.Rules{
    // string style
    "email": "required|email|max:100",

    // typed builder — compile-safe, IDE autocomplete
    "email": validator.R().Required().Email().Max(100).Build(),
    "age":   validator.R().Required().Integer().Min(18).Build(),
    "role":  validator.R().Required().In("admin", "user").Build(),
    "slug":  validator.R().Required().Regex(`^[a-z0-9-]+$`).Build(),
    "score": validator.R().Required().Numeric().Between(0, 100).Build(),
}
```

Every built-in rule has a corresponding method. A few naming notes:

| Rule string | Builder method |
| ----------- | -------------- |
| `string` | `.Str()` (reserved word in Go) |
| `required_if:field,val` | `.RequiredIf("field", "val")` |
| `starts_with:a,b` | `.StartsWith("a", "b")` |
| `in:a,b,c` | `.In("a", "b", "c")` |
| `unique:users,email` | `.Unique("users", "email")` |

## Nested validation

Use dot notation to validate fields inside nested maps:

```go
// Input
{
    "user": {
        "name": "hamizi",
        "address": {
            "postcode": "50000"
        }
    }
}

// Rules
validator.Make(input, validator.Rules{
    "user.name":             "required|min:3",
    "user.address.postcode": "required|digits:5",
})
```

## Array wildcard

Use `field.*` to validate every element of an array. Errors are keyed as `field.0`, `field.1`, etc.

```go
// Input
{
    "tags": ["laravel", "golang", ""]
}

// Rules
validator.Make(input, validator.Rules{
    "tags":   "required|array",
    "tags.*": "required|string|min:2",
})

// Errors
// tags.2 → ["The tags.2 field is required."]
```

## Custom rules

Register reusable rules with `Extend`:

```go
validator.Extend("strong_password", func(field string, value any, param string) error {
    s, ok := value.(string)
    if !ok || len(s) < 12 {
        return fmt.Errorf("The %s must be at least 12 characters.", field)
    }
    return nil
})

// Use it like any built-in rule
validator.Make(input, validator.Rules{
    "password": "required|strong_password",
})
```

## Usage with Fiber

```go
func CreateUser(c *fiber.Ctx) error {
    var body map[string]any
    if err := c.BodyParser(&body); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
    }

    v := validator.Make(validator.Input(body), validator.Rules{
        "name":     "required|string|min:2|max:100",
        "email":    "required|email|unique:users,email",
        "password": "required|min:8|confirmed",
        "role":     "required|in:admin,user,editor",
    }).WithDB(db)

    if v.Fails() {
        return c.Status(422).JSON(v.Errors())
    }

    // proceed with validated input ...
    return c.SendStatus(201)
}
```

## Testing

```bash
# Run all tests
go test ./tests/...

# Verbose
go test ./tests/... -v
```

## License

Distributed under MIT License, please see license file within the code for more details.
