<h1 align="center">
<br>
Auth
<br>
<br>
</h1>

Auth is a simple cookie/session based authentication middleware for HTTP handlers.

* Easy to use with minimal API surface
* 100% compatible with `http.Handler`
* Bring your own user type (uses Generics)
* Uses [github.com/alexedwards/scs](https://github.com/alexedwards/scs) under the hood

## Installation

```
go get -u github.com/philippta/auth
```

## Usage

View [example/main.go](example/main.go) for a complete example.

```go
type user struct {
    name string
}

func lookupUser(id string) (*user, ok bool) {
    return &user{name: "John Doe"}, true
}

auth := auth.New(lookupUser)

// Start and stop session
auth.Login(r.Context(), userID)
auth.Logout(r.Context())

// Get user or login state from session
user, ok := auth.User(r.Context())
      ok := auth.Ok(r.Context())

// Wrap handler or router
auth.Handler(router)
```

## License

Copyright (c) 2022 [Philipp Tanlak](https://github.com/philippta)

Licensed under [MIT License](LICENSE)
