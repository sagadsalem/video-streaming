package main

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
)

func main() {
	// define gorilla router
	router := mux.NewRouter().StrictSlash(true)
	// the index router
	router.HandleFunc("/", Index).Methods("GET")

	// the file system structure
	/* /storage/{cdn-disk}/posts/{post-hash}/{video-hash}/{ts-steam-index} */

	// the base streaming router
	router.HandleFunc("/storage/{cdn:disk[0-9]+}/posts/{post}/{video}/", Streaming).Methods("GET")
	// the other streaming router
	router.HandleFunc("/storage/{cdn:disk[0-9]+}/posts/{post}/{video}/{segment:index[0-9]+.ts}", Streaming).Methods("GET")

	log.Info("server is now running on http://localhost:8000")
	// if anything goes wrong with the server PANIC!
	panic(http.ListenAndServe(":8000", router))

}

// Index Handler to view the index.html page
func Index(w http.ResponseWriter, r *http.Request) {
	// parsing templates before execute them
	index, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatalf("could not parsing the template %v \n", err)
	}
	// execute template to render it
	data := struct {
		Post  string
		Video string
		Cdn   string
	}{
		Post:  "fc7e987f23de5bd6562b",
		Video: "7c0063cad659",
		Cdn:   "disk1",
	}
	if err := index.Execute(w, data); err != nil {
		log.Fatalf("could not parsing the template %v\n", err)
		return
	}
	// log success
	log.Info("template parsed successfully")
	return
}

// Streaming Handler to handle the video streaming
func Streaming(w http.ResponseWriter, r *http.Request) {
	// parse all querystring & parameters to vars
	vars := mux.Vars(r)

	post := vars["post"]
	video := vars["video"]
	cdn := vars["cdn"]

	path := fmt.Sprintf("storage/%s/posts/%s/%s", cdn, post, video)

	var file, contentType string
	segment, ok := vars["segment"]
	if !ok {
		contentType = "application/x-mpegURL"
		file = fmt.Sprintf("%s/%s", path, "index.m3u8")
	} else {
		contentType = "video/MP2T"
		file = fmt.Sprintf("%s/%s", path, segment)
	}
	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, file)
}
