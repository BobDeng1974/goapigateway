# goapigateway

This program aims to test the calls that a restful server could receive from a specific client. Our client that will trigger the requests is the Aws  ervice, Api Gateway. 

The objective is to test all incoming messages from the Aws Api Gateway and implement them as optimally as possible, showing how to  andle each type of request, whether it is PUT, POST, GET, DELETE, HEAD, OPTIONS. It also aims to document and show how to mount a server restful  integrated with Api Gateway from Aws.

Every test can be automated, but in the beginning for didactic reasons  we will do all the manual process, so that we can thoroughly understand  its operation and in the second moment propose an automation of process step.

## Used libraries:
- https://github.com/gorilla/mux - Implements a request router and dispatcher for matching incoming requests
- https://github.com/fatih/color - Help define the colors in the terminal


Body of main function StartTestServer

```go

// This method StartTestServer that will start our server,
// and mount our handler so we can work everything
// that arrives and everything that can come out.
func StartTestServer() {

	cfg := Config()

	color.Cyan("Testing services")
	color.Yellow("successfully...")

	postest := cfg.Schema + "://" + cfg.ServerHost + ":" + cfg.ServerPort + "/postest"

	color.Red("POST " + postest)
	color.Red("GET  " + postest)

	color.Yellow("Starting service...")
	color.Green("Host: " + cfg.ServerHost)
	color.Green("Schema: " + cfg.Schema)
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
			fmt.Println("Protocolo: ", r.Proto)
			fmt.Println("ProtoMajor: ", r.ProtoMajor)
			fmt.Println("ProtoMinor: ", r.ProtoMinor)
			fmt.Println("GetBody: ", r.GetBody)
			fmt.Println("Body: ", r.Body)

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

					msgjson := JsonMsg(500, "Set Content-Type correctly: Allowed: application / x-www-form-urlencoded, application / json")
					fmt.Fprintln(w, msgjson)

				}

			} else {

				msgjson := JsonMsg(500, "Not authorized / Allowed method POST")
				fmt.Fprintln(w, msgjson)
			}
		})

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

## Example of curl client simulating the Api Gateway Aws

```
curl -H 'Authorization:tyladfadiwkxceieixweiex747' --form "{}" http://localhost:9001/postest

```
