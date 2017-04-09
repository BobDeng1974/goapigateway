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

	// Every time trying to access our api without a
	// method it fires to the root and sends a welcome message
	router.Handle("/", http.FileServer(http.Dir("msg")))

	// This handler is that we will test all the possibilities
	// that it can receive when the method is post coming from the api gateway of aws
	router.
		HandleFunc("/postest", func(w http.ResponseWriter, r *http.Request) {

			fmt.Println("Fired method ..")

			if r.Method == "POST" {

				msgjson := JsonMsg(200, "ok")
				fmt.Fprintln(w, msgjson)

			} else if r.Method == "GET" {

				msgjson := JsonMsg(500, "Not authorized / Allowed method POST!")
				fmt.Fprintln(w, msgjson)

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
