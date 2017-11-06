package proxy

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Stdlib:
	"fmt"
	"net/http"

	// Community:
	log "github.com/Sirupsen/logrus"
)

//-----------------------------------------------------------------------------
// func: return200
//-----------------------------------------------------------------------------

func return200(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

//-----------------------------------------------------------------------------
// func: Start
//-----------------------------------------------------------------------------

// Start starts a proxy.
func Start() {

	// Request routing:
	http.HandleFunc("/", return200)

	// Start the web server:
	log.Info("Starting the web server...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
