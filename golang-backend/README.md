# golang backend

- receives merge intervals request from the JavaScript/Vue3 frontend application
- runs the algorithm for merging intervals
- returns the merged intervals back to the frontend

## starting the golang backend (without Docker)

prerequsites: existing golang installation

```sh
cd /into/golang-backend
go run main.go
```

## running unit tests for the backend

### running test with basic output

```sh
go test
```

### running test with verbose output

```sh
go test -v
```

### running test with coverage information

```sh
go test -v -cover
```

### running test with visual coverage information displayed in a web browser

```sh
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## generate 1000 random intervals for testing (change number of generated intervals as desired)
```sh
go run generator/generate_intervals.go -n 1000
```
