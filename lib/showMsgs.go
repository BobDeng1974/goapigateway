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
package lib

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/jeffotoni/goapigateway/conf"
)

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
