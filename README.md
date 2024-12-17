### Generate mocks

example
```
mockgen -source=api/users/store.go -destination=api/users/mocks/mock_store.go -package=mocks
```

### Run tests

```
go test -v ./...
```