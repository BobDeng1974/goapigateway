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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/jeffotoni/goapigateway/conf"
)

// Our global variables
var (
	err           error
	returns       string
	confServer    *http.Server
	AUTHORIZATION = `bc8c154ebabc6f3da724e9x5fef79238`

	// A directory is created by
	// the key we are simulating
	acessekey = "123456"
)

// To return the messages
type JsonPostTest1 struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func redirect(w http.ResponseWriter, req *http.Request) {

	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	log.Printf("redirect to: %s", target)

	http.Redirect(w, req, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
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

	if cfg.Schema == "https" {

		// This method will only serve to redirect everything you get
		// on port 80 to port 443, we will ensure that all access
		// will come from https
		go http.ListenAndServe(":"+cfg.PortRedirect, http.HandlerFunc(redirect))
	}

	///create route
	router := mux.NewRouter().StrictSlash(true)

	//router.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")

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

			if r.Method == "POST" || r.Method == "PUT" || r.Method == "GET" {

				// When the receipt is in json format
				if r.Header.Get("Content-Type") == "application/json" {

					// Treating the sending body
					// to json and transforming
					// into objects
					jsonRequest(w, r)

				} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" { // application/x-www-form-urlencoded default POST

					// Treating the sending
					// body to get native
					// postform
					postFormRequest(w, r)

				} else {

					// Uploading files
					// in 2 receive formats
					UploadFileEasy(w, r)

				}

			} else {

				msgjson := conf.JsonMsg(500, "Not authorized / Allowed method POST")
				fmt.Fprintln(w, msgjson)
			}
		})

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
		confServer = &http.Server{

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
		confServer = &http.Server{

			Handler: router,
			Addr:    cfg.Host + ":" + cfg.ServerPort,

			// Good idea, good live!!!
			//WriteTimeout: 10 * time.Second,
			//ReadTimeout:  10 * time.Second,
		}
	}

	// Defining whether it is https or http,
	// if it is http leave the calls
	// without access keys
	if cfg.Schema == "https" {

		log.Fatal(confServer.ListenAndServeTLS(cfg.Pem, cfg.Key))

	} else {

		log.Fatal(confServer.ListenAndServe())
	}

}

// Method responsible for capturing what arrives
// in the request transform into objects, what
// it receives is a string of json type
// coming from the client
func jsonRequest(w http.ResponseWriter, r *http.Request) {

	// Seeking our structure json test
	objJson := JsonPostTest1{}

	// Now let's decode what's coming in on the
	// request in json for objects
	// we've created with objJson
	errj := json.NewDecoder(r.Body).Decode(&objJson)

	if errj == nil {

		// Seeking fields direct from our object
		color.Yellow("When the receipt is in json format..")
		fmt.Println("email: ", objJson.Email)
		fmt.Println("password: ", objJson.Password)

		msgjson := conf.JsonMsg(200, "ok")
		fmt.Fprintln(w, msgjson)

	} else {

		fmt.Println("Error ..", errj)

		msgjson := conf.JsonMsg(500, "error: "+fmt.Sprintf("%s", errj))
		fmt.Fprintln(w, msgjson)
	}
}

// Method responsible for capturing what
// arrives in the request natively using
// PostFormValue, fields coming from the client.
func postFormRequest(w http.ResponseWriter, r *http.Request) {

	color.Green("When the receipt is a default")

	// Decoding when the post did not come
	// as json, content-type came as application / x-www-form-urlencoded
	// When the receipt is a default
	fmt.Println("email: ", r.PostFormValue("email"))
	fmt.Println("password: ", r.PostFormValue("password"))

	msgjson := conf.JsonMsg(200, "ok")
	fmt.Fprintln(w, msgjson)
}

// Method UploadFileEasy responsible for simulating our types of uploads,
// types are: multipart / form-data using form or option -F | --form
// of curl, application / octet-stream using --data-binary
func UploadFileEasy(w http.ResponseWriter, r *http.Request) {

	nameFileUp := r.Header.Get("Name-File")

	// This header was defined by our restful
	// server so we can understand and know
	// that the submitted type is a binary upload
	if nameFileUp != "" {

		uploadBinary(w, r)

	} else {

		errup := r.ParseMultipartForm(100000)
		if errup != nil {

			log.Printf("Error: Content-type or submitted format is incorrect to upload MultiForm  %s\n", errup)
			msgjson := conf.JsonMsg(500, errup.Error())
			fmt.Fprintln(w, msgjson)

			return
		}

		//get a ref to the
		//parsed multipart form
		multi := r.MultipartForm
		if len(multi.File["fileupload[]"]) > 0 {

			fmt.Println("size array files: ", len(multi.File["fileupload[]"]))

			fmt.Printf("map %+v\n", multi.File)

			// Call our method to save the
			// sending of multiple files to disk
			uploadFormFileMulti(w, r, multi.File["fileupload[]"])

		} else {

			// Upload only 1 file,
			// saving a file to disk
			uploadFormFile(w, r)

		}
	}
}

// This method uploadBinary only receives files coming in binary format,
// it will copy to disk what is coming via http
func uploadBinary(w http.ResponseWriter, r *http.Request) {

	cfg := conf.Config()

	nameFileUp := r.Header.Get("Name-File")

	// Upload octet-stream
	if nameFileUp != "" {

		// Creating the structure if it does not exist
		pathUpKeyUser := cfg.PathLocal + "/" + acessekey
		existPath, _ := os.Stat(pathUpKeyUser)
		if existPath == nil {
			// create path
			os.MkdirAll(pathUpKeyUser, 0777)
		}

		// Setting the path
		pathUpKeyUserFull := pathUpKeyUser + "/" + nameFileUp

		// In amazon does not receive multipart / form-data only application
		// / octet-stream ie --data-binary or instead of --form nameupload = @,
		// then we implement the 2 forms for our upload test
		ff, _ := os.OpenFile(pathUpKeyUserFull, os.O_WRONLY|os.O_CREATE, 0777)
		defer ff.Close()

		// Copying the contents of http
		// body to our file
		sizef, _ := io.Copy(ff, r.Body)

		// Writing in our response
		w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", sizef)))

		color.Red("File name: %s\n", nameFileUp)
		color.Yellow("copied: %v bytes\n", sizef)
		color.Yellow("copied: %v Kb\n", sizef/1024)
		color.Yellow("copied: %v Mb\n", sizef/1048576)

		msgjson := conf.JsonMsg(200, "ok upload size: "+fmt.Sprintf("%d bytes are recieved.\n", sizef)+" name file: "+nameFileUp)
		fmt.Fprintln(w, msgjson)

	}
}

// This method uploadFormFile only receives files coming in
// the multipart / form-data format, ie comes from a form
// sent by our client
func uploadFormFile(w http.ResponseWriter, r *http.Request) {

	cfg := conf.Config()

	// Validating the upload FormFile
	errup := r.ParseMultipartForm(32 << 20)
	if errup != nil {

		log.Printf("Error: Content-type or submitted format is incorrect to upload  %s\n", errup)
		msgjson := conf.JsonMsg(500, errup.Error())
		fmt.Fprintln(w, msgjson)

		return
	}

	// Upload multipart/form-data
	sizeMaxUpload := r.ContentLength / 1048576 ///Mb

	if sizeMaxUpload > cfg.UploadSize {

		fmt.Println("The maximum upload size: ", cfg.UploadSize, "Mb is large: ", sizeMaxUpload, "Mb", " in bytes: ", r.ContentLength)

		msgjson := conf.JsonMsg(500, "Unsupported file size max:"+fmt.Sprintf("%v", cfg.UploadSize)+"Mb")
		fmt.Fprintln(w, msgjson)

	} else {

		// Looking for the file in the FormFile method
		file, handler, errf := r.FormFile("fileupload")

		if errf != nil {

			log.Println(errf.Error())
			//http.Error(w, errf.Error(), http.StatusBadRequest)

			msgjson := conf.JsonMsg(500, errf.Error())
			fmt.Fprintln(w, msgjson)
			return
		}
		defer file.Close()

		if errf != nil {
			color.Red("Error big file, try again!")
			http.Error(w, "Error parsing uploaded file: "+errf.Error(), http.StatusBadRequest)
			return
		}

		defer file.Close()

		///create dir to key
		pathUpKeyUser := cfg.PathLocal + "/" + acessekey

		existPath, _ := os.Stat(pathUpKeyUser)

		if existPath == nil {

			// create path
			os.MkdirAll(pathUpKeyUser, 0777)
		}

		pathUserAcess := cfg.PathLocal + "/" + acessekey + "/" + handler.Filename

		// create files
		f, _ := os.OpenFile(pathUserAcess, os.O_WRONLY|os.O_CREATE, 0777)
		defer f.Close()

		// Copying the FormFile file
		// to our local disk file
		sizef, _ := io.Copy(f, file)

		// Capturing sending of the request from the
		// sending example curl -F or --form,
		// curl -F "email=jeff@yes.com" -F "password = 3939x9393"
		color.Yellow("Get Form Data\n")
		fmt.Println("email: ", r.PostFormValue("email"))
		fmt.Println("password: ", r.PostFormValue("password"))

		//To display results on server
		name := strings.Split(handler.Filename, ".")
		color.Red("File name: %s\n", name[0])
		color.Yellow("extension: %s\n", name[1])

		color.Yellow("size file: %v\n", sizeMaxUpload)
		color.Yellow("allowed: %v\n", cfg.UploadSize)

		color.Yellow("copied: %v bytes\n", sizef)
		color.Yellow("copied: %v Kb\n", sizef/1024)
		color.Yellow("copied: %v Mb\n", sizef/1048576)

		msgjson := conf.JsonMsg(200, "ok upload size: "+fmt.Sprintf("%d bytes are recieved.\n", sizef)+" name file: "+handler.Filename)
		fmt.Fprintln(w, msgjson)
	}
}

// This method is responsible for receiving multiple
// files from uploading and writing the various files to disk.
func uploadFormFileMulti(w http.ResponseWriter, r *http.Request, files []*multipart.FileHeader) {

	cfg := conf.Config()

	fmt.Println("size array: ", len(files))

	///create dir to key
	pathUpKeyUser := cfg.PathLocal + "/" + acessekey

	// Checking if the folder exists,
	// otherwise create the folder where
	// the files will be uploaded
	existPath, _ := os.Stat(pathUpKeyUser)
	if existPath == nil {
		// create path
		os.MkdirAll(pathUpKeyUser, 0777)
	}

	var uploadSize int64
	nameFileString := ""

	// Reading multipart.File Header to loop
	// around the vector, and grab file by file and write to disk
	for i, _ := range files {

		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("multi name: ", files[i].Filename)
		pathUserAcess := cfg.PathLocal + "/" + acessekey + "/" + files[i].Filename

		// copy file and write
		f, _ := os.Create(pathUserAcess)
		defer f.Close()

		//copy the uploaded file to the destination file
		if sizef, err := io.Copy(f, file); err != nil {

			msgerr := fmt.Sprintf("Copy MultiForm: %s", err.Error()) + " " + fmt.Sprintf("%s", http.StatusInternalServerError)
			msgjson := conf.JsonMsg(500, msgerr)
			fmt.Fprintln(w, msgjson)
			return

		} else {

			uploadSize += sizef
			nameFileString = nameFileString + "; " + files[i].Filename
		}
	}

	// Capturing sending of the request from the
	// sending example curl -F or --form,
	// curl -F "email=jeff@yes.com" -F "password = 3939x9393"
	color.Yellow("Get Form Data\n")
	fmt.Println("email: ", r.PostFormValue("email"))
	fmt.Println("password: ", r.PostFormValue("password"))

	//To display results upload
	color.Red("File name: %s\n", nameFileString)
	color.Yellow("size file: %v\n", uploadSize)
	color.Yellow("allowed: %v\n", cfg.UploadSize)

	msgjson := conf.JsonMsg(200, "ok upload size: "+fmt.Sprintf("%d bytes are recieved.\n", uploadSize)+" name file: "+nameFileString)
	fmt.Fprintln(w, msgjson)
}

// Only presents our
// environment variables
func showMsg() {

	// Instancing config
	cfg := conf.Config()

	// Presentation on the
	// screen of the start of our server
	color.Cyan("Testing services")
	color.Yellow("successfully...")
	postest := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/postest"
	color.Red("POST " + postest)
	color.Red("GET  " + postest)
	color.Yellow("Starting service...")
	color.Green("Host: " + cfg.ServerHost)
	color.Green("Schema: " + cfg.Schema)
	color.Green("Port: " + cfg.ServerPort)
}

// Presenting on the screen some of
// the variables that we will need
// to handle when receiving a requisition
func showMsgHandler(r *http.Request) {

	// Showing some important variables
	// that come from our requests
	fmt.Println("Fired method ..")
	fmt.Println("Header: ", r.Header)
	fmt.Println("Host: ", r.Host)
	fmt.Println("Method: ", r.Method)
	fmt.Println("RemoteAddr: ", r.RemoteAddr)
	fmt.Println("RequestURI: ", r.RequestURI)
	fmt.Println("Response: ", r.Response)
	fmt.Println("URL: ", r.URL)
	fmt.Println("TLS: ", r.TLS)
	fmt.Println("Agent: ", r.UserAgent())
	fmt.Println("ContentLength: ", r.ContentLength)
	fmt.Println("Content-type: ", r.Header.Get("Content-Type"))
	fmt.Println("Autorization: ", r.Header.Get("Authorization"))

	// AWS
	fmt.Println("AWS-Api-Id: ", r.Header.Get("X-Amzn-Apigateway-Api-Id"))
	fmt.Println("x-amzn-RequestId: ", r.Header.Get("x-amzn-RequestId"))
	fmt.Println("X-Amzn-Trace-Id: ", r.Header.Get("X-Amzn-Trace-Id"))
	fmt.Println("X-Cache: ", r.Header.Get("X-Cache"))
	fmt.Println("Via: ", r.Header.Get("Via"))
	fmt.Println("X-Amz-Cf-Id: ", r.Header.Get("X-Amz-Cf-Id"))

	// Some important variables
	fmt.Println("Protocolo: ", r.Proto)
	fmt.Println("ProtoMajor: ", r.ProtoMajor)
	fmt.Println("ProtoMinor: ", r.ProtoMinor)
	fmt.Println("GetBody: ", r.GetBody)
	fmt.Println("Body: ", r.Body)

	// upload octet-stream
	// Name-File
	fmt.Println("Name-File: ", r.Header.Get("Name-File"))

	// Basic authentication
	KEY, KEY_PASS, _ := r.BasicAuth()
	fmt.Println("KEY:", KEY, "PASS: ", KEY_PASS)

}
