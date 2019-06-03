# Payments Server

This is a simple REST server allowing consumers to manipulate a collection of Payments.

## Usage

`go run $GOPATH/src/github.com/DusanKasan/payments/cmd/payments -dsn="postgres://user:pass@host:5432/schema"`

The `dsn` flag specifies a DSN to a postgres database.

The server uses [Gin Framework](https://github.com/gin-gonic/gin) for the HTTP layer and will by default start in a debug mode (you will be notified of this and how to switch to production mode in the CLI). The logs are outputted to stdout as is standard.

## Endpoints

The server exposes the following endpoints:
- GET /payments\[?afterID=:ID\]
- GET /payments/:ID
- POST /payments
- PUT /payments/:ID
- DELETE /payments/:ID

For more in-depth documentation and examples look at the OpenAPI docs in `api` directory and/or the [Pact](https://docs.pact.io/) tests in `tests/pact`

## Payment Schema & Validation

A Payment object schema is defined in the OpenAPI specification as following:

```
type: object
properties:
    id:
        type: string
        format: uuid
    type:
        type: string
        enum:
            - 'Payment'
    version:
        type: 'integer'
        format: 'int64'
        enum:
            - 0
    organisation_id:
          type: string
        format: uuid
    attributes:
        type: object
```

Since the structure of the Payment resource is domain specific and out of scope for this project, we chose to not validate the attributes part of the Payment object (see below) and handle it as an arbitrary JSON object. We validate the version field to be 0, which allows us to add any attributes validation in the future for next versions of the Payment resource.

As for the errors, it may be necessary to introduce error codes in the future, but currently it is not needed as we can easily cover the errors via HTTP codes. However we always return a JSON representation of an error which can easily be extended by adding an error code.

