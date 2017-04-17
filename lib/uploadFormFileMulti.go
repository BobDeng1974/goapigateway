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
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/jeffotoni/goapigateway/conf"
)

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
