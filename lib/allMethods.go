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
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/jeffotoni/goapigateway/conf"
)

func allBodyExec(w http.ResponseWriter, r *http.Request) {

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
}

// Will not work with https
// Redirecting everything that arrives on port 80 to 443
func redirect(w http.ResponseWriter, req *http.Request) {

	// remove/add not default ports from req.Host
	target := "http://app." + req.Host + ":4001" + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	log.Printf("redirect to: %s", target)

	msgjson := conf.JsonMsg(500, "This host is not allowed")
	fmt.Fprintln(w, msgjson)

	// Redirection
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
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
