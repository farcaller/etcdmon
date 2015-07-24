package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net/http"
)

var (
	target     = flag.String("target", "", "etcd target url.")
	listen     = flag.String("listen", "0.0.0.0:8000", "Listening port.")
	cert       = flag.String("cert", "", "etcd client certificate.")
	key        = flag.String("key", "", "etcd client key.")
	httpClient *http.Client
)

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := httpClient.Get(*target + "/metrics")
	if err != nil {
		log.Printf("failed to talk to etcd at %s: %v", *target, err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	flag.Parse()

	authCert, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		log.Panicf("failed to load cert/key: %v", err)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{authCert},
			InsecureSkipVerify: true,
		},
	}
	httpClient = &http.Client{Transport: tr}

	http.HandleFunc("/metrics", handler)
	http.ListenAndServe(*listen, nil)
}
