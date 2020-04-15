I. Installations instructions

a. Update .env file config

b. Update dependencies by running commands
```.env
    go get ./...
```

c. To build services and notifiers 
```.env
    go build main/listener.go
    go build main/notifier.go
```

d. To run services and notifiers 
```.env
    go run main/listener.go
    go run main/notifier.go
```