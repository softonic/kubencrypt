package proxy

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Stdlib:
	"fmt"
	"net/http"
	"strconv"
	"time"

	// Community:
	log "github.com/Sirupsen/logrus"
)

//-----------------------------------------------------------------------------
// Typedefs:
//-----------------------------------------------------------------------------

// Data holds proxy data.
type Data struct {

	// Flags:
	ServicePort int
	Domain      string
}

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
func (d *Data) Start() {

	// Request routing:
	http.HandleFunc("/", return200)

	// Start the web server:
	log.Info("Starting the proxy server...")
	err := http.ListenAndServe(":"+strconv.Itoa(d.ServicePort), nil)
	if err != nil {
		log.Fatal(err)
	}
}

//-----------------------------------------------------------------------------
// func: Reachable
//-----------------------------------------------------------------------------

// Reachable returns true if the proxy is reachable.
func (d *Data) Reachable() bool {

	// Sensible timeout:
	var netClient = &http.Client{
		Timeout: time.Second * 5,
	}

	// Log:
	log.Info("Checking whether I am reachable...")

	for {

		// Send an HTTP/GET request:
		resp, err := netClient.Get(d.Domain)
		if (err == nil) && (resp.StatusCode == 200) {
			return true
		}

		// Sleep before next attempt:
		log.Info(err.Error)
		time.Sleep(time.Second * 10)
	}
}
