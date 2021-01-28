example queue reader
====================

This is an example queue-reader, it reads messages from an SQS Queue, and
logs them out through the zap logger included.

You can test this by:
1. Run a local `goaws` server
    - From the root of `service-framework`:`docker run -d --name goaws -p 4100:4100 -v "/$(pwd)/example/queue/goaws.yaml:/conf/goaws.yaml" pafortin/goaws`
    - From this directory: `docker run -d --name goaws -p 4100:4100 -v "/$(pwd)/goaws.yaml:/conf/goaws.yaml" pafortin/goaws`
2. Run the reader
    - From the root of `service-framework`: `go run example/queue/cmd/main.go read --aws-sqs-queue-url=http://localhost:4100/example-queue --aws-endpoint-url=http://localhost:4100`
    - From this directory: `go run cmd/main.go read --aws-sqs-queue-url=http://localhost:4100/example-queue --aws-endpoint-url=http://localhost:4100`
3. Send it a message: `aws --endpoint-url http://localhost:4100 sqs send-message --queue-url http://localhost:4100/example-queue --message-body "Hello, world!"`

## Notes
- This example server will shut down if it receives any errors from AWS
- It supports "graceful shutdown", so if it receives a `SIGINT`, it will shut down.