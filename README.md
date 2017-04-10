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