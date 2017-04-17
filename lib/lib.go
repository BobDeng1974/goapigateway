/***********

 ▄▄▄██▀▀▀▓█████   █████▒ █████▒▒█████  ▄▄▄█████▓ ▒█████   ███▄    █  ██▓
   ▒██   ▓█   ▀ ▓██   ▒▓██   ▒▒██▒  ██▒▓  ██▒ ▓▒▒██▒  ██▒ ██ ▀█   █ ▓██▒
   ░██   ▒███   ▒████ ░▒████ ░▒██░  ██▒▒ ▓██░ ▒░▒██░  ██▒▓██  ▀█ ██▒▒██▒
▓██▄██▓  ▒▓█  ▄ ░▓█▒  ░░▓█▒  ░▒██   ██░░ ▓██▓ ░ ▒██   ██░▓██▒  ▐▌██▒░██░
 ▓███▒   ░▒████▒░▒█░   ░▒█░   ░ ████▓▒░  ▒██▒ ░ ░ ████▓▒░▒██░   ▓██░░██░
 ▒▓▒▒░   ░░ ▒░ ░ ▒ ░    ▒ ░   ░ ▒░▒░▒░   ▒ ░░   ░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░▓
 ▒ ░▒░    ░ ░  ░ ░      ░       ░ ▒ ▒░     ░      ░ ▒ ▒░ ░ ░░   ░ ▒░ ▒ ░
 ░ ░ ░      ░    ░ ░    ░ ░   ░ ░ ░ ▒    ░      ░ ░ ░ ▒     ░   ░ ░  ▒ ░
 ░   ░      ░  ░                  ░ ░               ░ ░           ░  ░

*
*
* project test restfull Api Gateway Aws
*
* @package     main
* @author      jeffotoni
* @copyright   Copyright (c) 2017
* @license     --
* @link        --
* @since       Version 0.1
*/

// This program is a restful server, its purpose is to receive multiple POST requests,
// GET so that we can test the various ways to send files to a restful server.
// Our goal is to discover the different ways upload and receive upload so that
// we can implement our file server.
//
// The Amazon is allowed to send binaries to its Api Gateway,
// and it is possible to use lambda functions so that the entire upload process is done by the
// Api Gateway without necessarily needing to send direct to a restful server,
// but our goal is to send direct to our server Restful
//
// We will use curl as our client to test our submissions and we will also use
// the Amazon Api Gateway and see if it is possible to send a binary to our restful
// server directly without using lambda functions.
// Each test can be automated, but at the beginning for didactic reasons we will do the
// entire process manual so that we can fully understand its operation.
package lib

import (
	"crypto/tls" //https
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeffotoni/goapigateway/conf"
)

// Our global variables
var (
	err                error
	returns            string
	confServerHttp     *http.Server
	confServerHttps    *http.Server
	confServerHttpTest *http.Server
	AUTHORIZATION      = `bc8c154ebabc6f3da724e9x5fef79238`

	// A directory is created by
	// the key we are simulating
	acessekey = "123456"
)

// To return the messages
type JsonPostTest1 struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// This method StartTestServer that will start our server,
// and mount our handler so we can work everything
// that arrives and everything that can come out.
func StartTestServer() {

	// set Config
	cfg := conf.Config()

	// Displaying server
	// screen message
	showMsg()

	// Our redirect to https: will not work, it will not overwrite the url,
	// so we will disable ..
	//
	// He listens on port 80 but we will not actually let him run,
	// we omit to close the connection.
	if cfg.Schema == "https" {

		// This method will only serve to redirect everything you get
		// on port 80 to port 443, we will ensure that all access
		// will come from https
		go http.ListenAndServe(":"+cfg.PortRedirect, http.HandlerFunc(redirect))
	}

	///create route
	router := mux.NewRouter().StrictSlash(true)

	// Opening an escape port for Homologation
	rTest := mux.NewRouter().StrictSlash(true)

	// Homologation
	// This handler is that we will test all the possibilities
	// that it can receive when the method is post coming from the api gateway of aws
	rTest.
		HandleFunc("/postest", func(w http.ResponseWriter, r *http.Request) {

			fmt.Println("Homologation Environment")
			fmt.Println("Subdomain " + "http://" + r.Host + r.URL.Path)

			// Showing the objects of the
			// http.request method
			showMsgHandler(r)

			// Executes all the functions
			// of the request,
			// get, post etc.
			allBodyExec(w, r)

			msgjson := conf.JsonMsg(200, "Homologation")
			fmt.Fprintln(w, msgjson)
		})

	// Production
	// Every time trying to access our api without a
	// method it fires to the root and sends a welcome message
	router.Handle("/", http.FileServer(http.Dir("msg")))

	// This handler is that we will test all the possibilities
	// that it can receive when the method is post coming from the api gateway of aws
	router.
		HandleFunc("/postest", func(w http.ResponseWriter, r *http.Request) {

			w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")

			// Showing the objects of the
			// http.request method
			showMsgHandler(r)

			// Executes all the functions
			// of the request,
			// get, post etc.
			allBodyExec(w, r)

		})

	// Now it is possible to get on port 443 using https,
	// we did a check when it is https, the system
	// parameterizes config as needed
	if cfg.Schema == "https" {

		// .crt — Alternate synonymous most common among *nix systems .pem (pubkey).
		// .csr — Certficate Signing Requests (synonymous most common among *nix systems).
		// .cer — Microsoft alternate form of .crt, you can use MS to convert .crt to .cer (DER encoded .cer, or base64[PEM] encoded .cer).
		// .pem = The PEM extension is used for different types of X.509v3 files which contain ASCII (Base64) armored data prefixed with a «—– BEGIN …» line. These files may also bear the cer or the crt extension.
		// .der — The DER extension is used for binary DER encoded certificates.
		//
		// Setting our tls to accept
		// our .key and .crt keys.
		cfgs := &tls.Config{

			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		// Config to upload our server
		confServerHttps = &http.Server{

			Handler:      router,
			Addr:         cfg.Host + ":" + cfg.ServerPort,
			TLSConfig:    cfgs,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),

			// Good idea, good live!!!
			//WriteTimeout: 10 * time.Second,
			//ReadTimeout:  10 * time.Second,
		}

	} else {

		// Config to upload our server
		confServerHttp = &http.Server{

			Handler: router,
			Addr:    cfg.Host + ":" + cfg.ServerPort,

			// Good idea, good live!!!
			//WriteTimeout: 10 * time.Second,
			//ReadTimeout:  10 * time.Second,
		}
	}

	// Config to upload our server
	confServerHttpTest = &http.Server{

		Handler: rTest,
		Addr:    cfg.Host + ":4001",

		// Good idea, good live!!!
		//WriteTimeout: 10 * time.Second,
		//ReadTimeout:  10 * time.Second,
	}

	// Defining whether it is https or http,
	// if it is http leave the calls
	// without access keys
	if cfg.Schema == "https" {

		go func() {

			log.Fatal(confServerHttps.ListenAndServeTLS(cfg.Pem, cfg.Key))

		}()

		log.Fatal(confServerHttpTest.ListenAndServe())

	} else {

		go func() { log.Fatal(confServerHttp.ListenAndServe()) }()

		log.Fatal(confServerHttpTest.ListenAndServe())
	}
}
