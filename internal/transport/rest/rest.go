package rest

import (
	"log"
	// "context"
	// "github.com/centrifugal/centrifuge"
	"net/http"
)

func ListenAndServe() {
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
