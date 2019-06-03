FROM golang:1.12.5-alpine3.9

ENV GO111MODULE on

WORKDIR /go/src/app
COPY . .

RUN go install -mod=vendor ./cmd/payments
RUN go install -mod=vendor ./tools/migrate
EXPOSE 80
