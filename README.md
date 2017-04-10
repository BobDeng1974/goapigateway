# goapigateway

This program is a restful server, its purpose is to receive multiple POST requests, GET so that we can test the various ways to send files to a restful server. Our goal is to discover the different ways upload and receive upload so that we can implement our file server.

The Amazon is allowed to send binaries to its Api Gateway, and it is possible to use lambda functions so that the entire upload process is done by the Api Gateway without necessarily needing to send direct to a restful server, but our goal is to send direct to our server Restful

We will use curl as our client to test our submissions and we will also use the Amazon Api Gateway and see if it is possible to send a binary to our restful server directly without using lambda functions.

Each test can be automated, but at the beginning for didactic reasons we will do the entire process manual so that we can fully understand its operation.

## Used libraries:
- https://github.com/gorilla/mux - Implements a request router and dispatcher for matching incoming requests
- https://github.com/fatih/color - Help define the colors in the terminal


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

Body of main function StartTestServer

```go

// This method StartTestServer that will start our server,
// and mount our handler so we can work everything
// that arrives and everything that can come out.
func StartTestServer() {

	// Instancing config
	cfg := Config()

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
	color.Green("Port: " + cfg.ServerPort)

	///create route
	router := mux.NewRouter().StrictSlash(true)

	router.Headers("Content-Type", "application/json",
		"X-Requested-With", "XMLHttpRequest")

	// Every time trying to access our api without a
	// method it fires to the root and sends a welcome message
	router.Handle("/", http.FileServer(http.Dir("msg")))

	// This handler is that we will test all the possibilities
	// that it can receive when the method is post coming from the api gateway of aws
	router.
		HandleFunc("/postest", func(w http.ResponseWriter, r *http.Request) {

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
			fmt.Println("AWS-Trace: ", r.Header.Get("X-Amzn-Trace-Id"))
			fmt.Println("AWS-Api-Id: ", r.Header.Get("X-Amzn-Apigateway-Api-Id"))

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

			if r.Method == "POST" || r.Method == "PUT" || r.Method == "GET" {

				// When the receipt is in json format
				if r.Header.Get("Content-Type") == "application/json" {

					objJson := JsonPostTest1{}
					errj := json.NewDecoder(r.Body).Decode(&objJson)

					if errj == nil {

						color.Yellow("When the receipt is in json format..")
						email := objJson.Email
						password := objJson.Password

						fmt.Println("email: ", email)
						fmt.Println("password: ", password)

						msgjson := JsonMsg(200, "ok")
						fmt.Fprintln(w, msgjson)

					} else {

						fmt.Println("Error ..", errj)
					}
				} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" { // application/x-www-form-urlencoded default POST

					// When the receipt is a default
					color.Green("When the receipt is a default")
					fmt.Println("email: ", r.PostFormValue("email"))
					fmt.Println("password: ", r.PostFormValue("password"))

					msgjson := JsonMsg(200, "ok")
					fmt.Fprintln(w, msgjson)

				} else {
					//else if r.Header.Get("Content-Type") == "image/jpg" || r.Header.Get("Content-Type") == "image/png" || r.Header.Get("Content-Type") == "application/octet-stream" {

					UploadFileEasy(w, r)

				}

			} else {

				msgjson := JsonMsg(500, "Not authorized / Allowed method POST")
				fmt.Fprintln(w, msgjson)
			}
		})

	// Config to upload our server
	confServer = &http.Server{

		Handler: router,
		Addr:    cfg.Host + ":" + cfg.ServerPort,

		// Good idea, good live!!!
		//WriteTimeout: 10 * time.Second,
		//ReadTimeout:  10 * time.Second,
	}

	log.Fatal(confServer.ListenAndServe())
}

```

Body of main function UploadFileEasy

```go

// Method UploadFileEasy responsible for simulating our types of uploads,
// types are: multipart / form-data using form or option -F | --form
// of curl, application / octet-stream using --data-binary
func UploadFileEasy(w http.ResponseWriter, r *http.Request) {

	cfg := conf.Config()

	// A directory is created by
	// the key we are simulating
	acessekey := "123456"

	// This header was defined by our restful
	// server so we can understand and know
	// that the submitted type is a binary upload
	nameFileUp := r.Header.Get("Name-File")

	// Upload octet-stream
	if nameFileUp != "" {

		pathUpKeyUser := cfg.PathLocal + "/" + acessekey
		existPath, _ := os.Stat(pathUpKeyUser)
		if existPath == nil {
			// create path
			os.MkdirAll(pathUpKeyUser, 0777)
		}

		pathUpKeyUserFull := pathUpKeyUser + "/" + nameFileUp

		// In amazon does not receive multipart / form-data only application
		// / octet-stream ie --data-binary or instead of --form nameupload = @,
		// then we implement the 2 forms for our upload test
		ff, _ := os.OpenFile(pathUpKeyUserFull, os.O_WRONLY|os.O_CREATE, 0777)
		defer ff.Close()
		sizef, _ := io.Copy(ff, r.Body)
		w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", sizef)))

		color.Red("File name: %s\n", nameFileUp)
		color.Yellow("copied: %v bytes\n", sizef)
		color.Yellow("copied: %v Kb\n", sizef/1024)
		color.Yellow("copied: %v Mb\n", sizef/1048576)

		msgjson := conf.JsonMsg(200, "ok upload size: "+fmt.Sprintf("%d bytes are recieved.\n", sizef)+" name file: "+nameFileUp)
		fmt.Fprintln(w, msgjson)

	} else {

		// Upload multipart/form-data
		sizeMaxUpload := r.ContentLength / 1048576 ///Mb

		if sizeMaxUpload > cfg.UploadSize {

			fmt.Println("The maximum upload size: ", cfg.UploadSize, "Mb is large: ", sizeMaxUpload, "Mb", " in bytes: ", r.ContentLength)
			fmt.Fprintln(w, "", 500, "Unsupported file size max: ", cfg.UploadSize, "Mb")

		} else {

			errup := r.ParseMultipartForm(32 << 20)
			if errup != nil {
				log.Printf("ERROR UPLOAD PARSE: %s\n", errup)
				http.Error(w, errup.Error(), 500)
				return
			}

			file, handler, errf := r.FormFile("fileupload")
			if errf != nil {
				log.Println(errf.Error())
				http.Error(w, errf.Error(), http.StatusBadRequest)
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
			sizef, _ := io.Copy(f, file)

			//up_size := fmt.Sprintf("%v", r.ContentLength)

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

## Example using Api Gateway Aws

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

## Example POST text using Api Gateway Aws

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
-H "Content-Type: application/json" \
-H "Authorization: jeff b7d03a6947b217efb6f3ec3bd3504582" \
-d -d "email=jefferson&password=3838373773" 

```

##Example using form data passing fields through url of Api Gateway Aws

The content-type is not defined in this submission, Amazon's Api Gateway will not allow

```
curl -i -X POST https://xxxxxx.execute-api.us-xxx-x.amazonaws.com/goapigateway \
-u API_KEY:383883jef903xxxx838xxxx \
-H "Authorization: jeff b7d03a6947b217efb6f3ec3bd3504582" \
-d 'email=jeff@&password=38xxx8w8e'

```