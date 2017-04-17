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

// This method uploadFormFile only receives files coming in
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/jeffotoni/goapigateway/conf"
)

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
