# goapigateway

This program is a restful server, its purpose is to receive multiple POST requests, GET so that we can test the various ways to send files to a restful server. Our goal is to discover the different ways upload and receive upload so that we can implement our file server.

The Amazon is allowed to send binaries to its Api Gateway, and it is possible to use lambda functions so that the entire upload process is done by the Api Gateway without necessarily needing to send direct to a restful server, but our goal is to send direct to our server Restful

We will use curl as our client to test our submissions and we will also use the Amazon Api Gateway and see if it is possible to send a binary to our restful server directly without using lambda functions.

Each test can be automated, but at the beginning for didactic reasons we will do the entire process manual so that we can fully understand its operation.

## Used libraries:
- https://github.com/gorilla/mux - Implements a request router and dispatcher for matching incoming requests
- https://github.com/fatih/color - Help define the colors in the terminal


## http.ListenAndServe

In golang everything becomes simpler, check the call to go up the door of our rest server, using http protocol

```
http.ListenAndServe(":8080", nil)

Using Mux

confServerHttpTest.ListenAndServe()

```

To upload a port using https for example would

```
http.ListenAndServeTLS(":8081", "cert.crt", "yourkey.key", nil)

Using Mux

confServerHttps.ListenAndServeTLS(server.crt, yourke.Key)

```

# Using OpenSSL

You can generate the files with the openssl native linux command.

We will generate to authenticate in some server in the internet like rapidssl, gogetssl and thousands of others spread in the web, to simulate the production environment.

We will not use startssl for example that is free for 1 year, it is no longer recognized and our https on our server golang will go from error.

The -nodes option is not the English word "nodes", but it is "in DES".

It means that OpenSSL will not encrypt your private key.

This will prevent from initializing the server it prompts for password.


```
openssl req -nodes -newkey rsa:2048 -keyout yourkey.key -out serve.csr

```

Now, is to send your csr to the server that will authenticate your key and issue your certificate for use of https, it will resend to you the files and what we will use on our golang server is .crt (Alternate synonymous most common Among * nix systems .pem (pubkey)).


Look at the definition of the files:

-  .crt — Alternate synonymous most common among *nix systems .pem (pubkey).
-  .csr — Certficate Signing Requests (synonymous most common among *nix systems).
-  .cer — Microsoft alternate form of .crt, you can use MS to convert .crt to .cer (DER encoded .cer, or base64[PEM] encoded .cer). 
-  .pem = The PEM extension is used for different types of X.509v3 files which contain ASCII (Base64) armored data prefixed with a «—– BEGIN …» line. These files may also bear the cer or the crt extension.
- .der — The DER extension is used for binary DER encoded certificates.


# Generating the Certficate Signing Request

There is another possibility that is to generate your own key signal.
You can also auto-sign your keys for local testing.

I do not recommend you want to test and simulate a production-like environment, but it's good to learn and test.

Generation of self-signed(x509) / public key (.pem|.crt) based on the private (.key)

```
openssl req -new -sha256 -key server.key -out server.csr
openssl x509 -req -sha256 -in server.csr -signkey server.key -out server.crt -days 3650

```

# Methods of our server

Body of main type struct

```go

// Structure of our server configurations
type Configs struct {
	Domain     string `json:"domain"`
	Process    string `json:"process"`
	Ping       string `json:"ping"`
	ServerPort string `json:"serverport"`
	Host       string `json:"host"`
	Schema     string `json:"shcema"`
	ServerHost string `json:"serverhost"`
	UploadSize int64  `json:"uploadsize"`
	PathLocal  string `json:"pathlocal"`
}

```

Body of main variables

```go
// Our global variables
var (
	err           error
	returns       string
	confServer    *http.Server
	AUTHORIZATION = `bc8c154ebabc6f3da724e9x5fef79238`
	objason       Configs
)

```

Body of main function ConfigJson

```go

// This method ConfigJson sets up our
// server variables from our struct
func ConfigJson() string {

	// Defining the values of our config
	data := &Configs{Domain: "localhost", Process: "2", Ping: "ok", ServerPort: "9001", Host: "", Schema: "http", ServerHost: "localhost", UploadSize: 100, PathLocal: "uploads"}

	// Converting our struct into json format
	cjson, err := json.Marshal(data)
	if err != nil {
		// handle err
	}

	return string(cjson)
}

```

Body of main function UploadFileEasy

```go


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

		uploadFormFile(w, r)

	}
}

```

Body of main function uploadBinary

```go

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

		// Copying the contents of http body to our file
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

```

Body of main function uploadFormFile

```go

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

		// copy file and write
		f, _ := os.OpenFile(pathUserAcess, os.O_WRONLY|os.O_CREATE, 0777)
		defer f.Close()

		// Copying the FormFile file to our local disk file
		sizef, _ := io.Copy(f, file)

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

```

Body of main function uploadFormFileMulti

```go

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

```

We created two routes one for homologation and one for production

```go
///create route
router := mux.NewRouter().StrictSlash(true)

// Opening an escape port for Homologation
rTest := mux.NewRouter().StrictSlash(true)

```

Below we define the individual config servers for when https config is different

```go

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

```

## Example of curl client Sending file in binary format

```
curl -i --request POST http://localhost:9001/postest \
-u "API_KEY:383883jef903xxxx838xxxx" \
-H "Accept: binary/octet-stream" \
-H "Content-Type: binary/octet-stream" \
-H "Name-File: file1.jpg" \
--data-binary "@file1.jpg"

```
## Example of curl client Sending file on multipart/form-data

```
curl -i -X POST http://localhost:9001/postest \
-u "API_KEY:383883jef903xxxx838xxxx" \
--form "fileupload=@files/file1.jpg" 

```

## Example using Api Gateway Aws  / Limit 10 MB per request

```
curl -i --request POST https://xxxxxx.execute-api.us-xxx-x.amazonaws.com/goapigateway \
-u "API_KEY:xxxxxx38383xxxx" \
-H "Accept: binary/octet-stream" \
-H "Content-Type: binary/octet-stream" \
-H "Name-File: file1.jpg" \
--data-binary "@files/file1.jpg"

```

## Example Using json to tesar our post and get

```
curl -i -X POST localhost:9001/postest \
-u API_KEY:383883jef903xxxx838xxxx \
-H "Content-Type: application/json" \
-H "Authorization: jeff b7d03a6947b217efb6f3ec3bd3504582" \
-d '{"email":"jeffotoni@yes.com","password":"3838373773"}'

```

## Example Post text using Api Gateway Aws  / Limit 10 MB per request

```
curl -i -X POST https://xxxxxx.execute-api.us-xxx-x.amazonaws.com/goapigateway \
-u API_KEY:383883jef903xxxx838xxxx \
-H "Content-Type: application/json" \
-H "Authorization: jeff b7d03a6947b217efb6f3ec3bd3504582" \
-d '{"email":"jeffotoni@yes.com","password":"3838373773"}'

```

## Example using form data passing fields through url


```
curl -i -X POST localhost:9001/postest \
-u API_KEY:383883jef903xxxx838xxxx \
-H "Authorization: jeff b7d03a6947b217efb6f3ec3bd3504582" \
-d -d "email=jefferson&password=3838373773" 

```

## Example using form data passing fields through url of Api Gateway Aws / Limit 10 MB per request

The content-type is not defined in this submission, Amazon's Api Gateway will not allow

```
curl -i -X POST https://xxxxxx.execute-api.us-xxx-x.amazonaws.com/goapigateway \
-u API_KEY:383883jef903xxxx838xxxx \
-H "Authorization: jeff b7d03a6947b217efb6f3ec3bd3504582" \
-d 'email=jeff@&password=38xxx8w8e'

```

## In a single request we send form fields and an upload file

The Amazon Gateway Api unfortunately does not accept --form for only binary file uploads then it will 
not work using the upload to Amazon.

```
curl -i --request POST localhost:9001/postest \
-u "API_KEY:383883jef903xxxx838xxxx" \
--form 'email=jefferson&password=3838373773' \
--form "fileupload=@files/file1.jpg"
```

## In a single request we send form fields and multiple files to upload

Unfortunately Amazon does not support --form / or better multipart / form-data for sending files only --data-binary

```
curl -i --request POST localhost:9001/postest \
-u "API_KEY:383883jef903xxxx838xxxx" \
--form 'email=jefferson&password=3838373773' \
--form "fileupload[]=@files/file1.jpg" \
--form "fileupload[]=@files/file2.pdf"

```

## We are using --cacert to read our bundle key so you can allow access to our https

```
curl --cacert budle.pem -i \
-X POST https://myapp/postest \
-u "API_KEY:383883jef903xxxx838xxxx" \
--form 'email=jefferson&password=3838373773' \
--form "fileupload[]=@bigfile1.jpg"

```