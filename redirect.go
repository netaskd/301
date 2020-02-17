package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/kabukky/httpscerts"
)

func handleRedirect(res http.ResponseWriter, req *http.Request, proto, domain string, port int) {
	log.Println("Got request for", req.URL.Path)
	http.Redirect(res, req, proto+"://"+domain+"/"+req.URL.Path[1:], 301)
}

func redirServer(proto, domain string, port int) {

	// Go, WHY U NO INTERPOLATE??
	portstring := fmt.Sprintf(":%d", port)
	fmt.Println("Listening on port", port, "Redirecting to", domain)

	redirServer := http.NewServeMux()

	// This feels super dirty
	redirServer.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		handleRedirect(res, req, proto, domain, port)
	})

	if proto == "https" {
		err := httpscerts.Check("cert.pem", "key.pem")
		if err != nil {
			err = httpscerts.Generate("cert.pem", "key.pem", fmt.Sprintf("%s:%s", domain, port))
			if err != nil {
				log.Fatal("Can not generate certs!")
			}
		}
		http.ListenAndServeTLS(portstring, "cert.pem", "key.pem", redirServer)
	} else {
		http.ListenAndServe(portstring, redirServer)
	}
}
