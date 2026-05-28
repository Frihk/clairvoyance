# backend

## backend structure

```
│   .env
│   .env.example
│   .gitignore
│   go.mod
│   go.sum
│   README.md
│   
├───cmd
│   └───server
│           main.go
│
└───internal
    ├───api
    │   ├───handlers
    │   │       auth.go
    │   │       
    │   ├───middleware
    │   │       auth.go
    │   │       cors.go
    │   │
    │   └───router
    │           router.go
    │
    ├───auth
    │       jwt.go
    │       password.go
    │
    ├───blockchain
    │       .keep
    │
    ├───config
    │       config.go
    │
    ├───db
    │       .keep
    │
    ├───models
    │       credentials.go
    │       users.go
    │
    ├───repositories
    │       auth.go
    │
    ├───services
    │       auth.go
    │
    └───utils
            errors.go
            logger.go
            responses.go
```