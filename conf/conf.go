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
package conf

import "encoding/json"

// Structure of our server configurations
type Configs struct {
	Domain       string `json:"domain"`
	Process      string `json:"process"`
	Ping         string `json:"ping"`
	ServerPort   string `json:"serverport"`
	Host         string `json:"host"`
	Schema       string `json:"shcema"`
	ServerHost   string `json:"serverhost"`
	UploadSize   int64  `json:"uploadsize"`
	PathLocal    string `json:"pathlocal"`
	PortRedirect string `json:"portredirect"`
	Pem          string `json:"pem"`
	Key          string `json:key`
}

// Our global variables
var (
	objason Configs
)

// This method ConfigJson sets up our
// server variables from our struct
func ConfigJson() string {

	// Defining the values of our config
	data := &Configs{

		Domain:       "localhost",
		Process:      "2",
		Ping:         "ok",
		ServerPort:   "443",
		Host:         "",
		Schema:       "https",
		ServerHost:   "localhost",
		UploadSize:   100,
		PathLocal:    "uploads",
		PortRedirect: "80",
		Pem:          "../certs/pem.crt",
		Key:          "../certs/key.key",
	}

	// Converting our struct into json format
	cjson, err := json.Marshal(data)
	if err != nil {
		// handle err
	}

	return string(cjson)
}

// This method Config returns the objects
// of our config so that it can be accessed
func Config() *Configs {

	// We are implementing singleton for this object,
	// every time it is instantiated it does
	// not redo it it only resends the object
	if objason.Domain != "" {

		return &objason

	} else {

		jsonT := []byte(ConfigJson())
		json.Unmarshal(jsonT, &objason)

		return &objason
	}
}

// This method Message is to return our messages
// in json, ie the client will
// receive messages in json format
type Message struct {
	Code int    `json:code`
	Msg  string `json:msg`
}

// This method is a simplified abstraction
// so that we can send them to our client
// when making a request
func JsonMsg(codeInt int, msgText string) string {

	data := &Message{Code: codeInt, Msg: msgText}

	djson, err := json.Marshal(data)
	if err != nil {
		// handle err
	}

	return string(djson)
}
