I. Installations instructions

a. Update .env file config

b. To build services and notifiers 
```.env
    go build main/listener.go
    go build main/notifier.go
```

c. To run services and notifiers 
```.env
    go run main/listener.go
    go run main/notifier.go
```