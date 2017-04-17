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
	"log"
	"net/http"

	"github.com/jeffotoni/goapigateway/conf"
)

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
