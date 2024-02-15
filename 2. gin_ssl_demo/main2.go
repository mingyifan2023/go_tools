package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	environment := flag.String("environment", "development", "Environment")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", greet)

	if *environment == "production" {
		serveProduction(mux)
	} else {
		serveDevelopment(mux)
	}
}

func serveDevelopment(h http.Handler) {
	err := http.ListenAndServe(":4000", h)
	log.Fatal(err)
}

func serveProduction(h http.Handler) {
	// Configure autocert settings
	autocertManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("mojoru.com"),
		// 设置路由
		// localDir, err := os.Getwd()
		// if err != nil {
		// 	fmt.Println("获取当前目录失败:", err)
		// 	return
		// }

		Cache: autocert.DirCache("/home/gin-ssl-demo"),
	}

	// Listen for HTTP requests on port 80 in a new goroutine. Use
	// autocertManager.HTTPHandler(nil) as the handler. This will send ACME
	// "http-01" challenge responses as necessary, and 302 redirect all other
	// requests to HTTPS.
	go func() {
		srv := &http.Server{
			Addr:         ":80",
			Handler:      autocertManager.HTTPHandler(nil),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		err := srv.ListenAndServe()
		log.Fatal(err)
	}()

	// Configure the TLS config to use the autocertManager.GetCertificate function.
	tlsConfig := &tls.Config{
		GetCertificate:           autocertManager.GetCertificate,
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         ":443",
		Handler:      h,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServeTLS("", "") // Key and cert provided automatically by autocert.
	log.Fatal(err)
}

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}
