default:
	go install .
	@echo "Done building."

build:
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o alerting .
run:
	@go run main.go