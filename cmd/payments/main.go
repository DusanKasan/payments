package main

import (
	"flag"
	"fmt"
	"github.com/DusanKasan/payments/internal/app/payments"
	"log"
	"math"
	"net/http"
	"os"
)
var DSN string
var Port int

func init() {
	flag.StringVar(&DSN, "dsn", "", "dsn to a postgres database")
	flag.IntVar(&Port, "port", 80, "port at which the server will listen")
	flag.Parse()

	if DSN == "" {
		log.Println("dsn cannot be empty")
		os.Exit(1)
	}

	if Port < 1 || Port > math.MaxUint16 {
		log.Println("port value out of bounds")
		os.Exit(1)
	}
}

func main() {
	server, err := payments.Handler(DSN)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), server); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
