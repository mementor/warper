package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// characters used for short-urls
	SYMBOLS = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"

	// someone set us up the bomb !!
	BASE = int64(len(SYMBOLS))
)

var config struct {
	Port     string
	Temp     string
	Basedir  string
	Storedir string
}

func init() {
	osTMP := os.TempDir()

	flag.StringVar(&config.Port, "port", "8080", "port number, default: 8080")
	flag.StringVar(&config.Temp, "temp", osTMP, "")
	flag.StringVar(&config.Basedir, "basedir", "warper", "dir inside tmp dir")
	flag.Parse()

	config.Storedir = filepath.Join(config.Temp, config.Basedir)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/", root).Methods("GET")
	r.HandleFunc("/{filename}", uploadHandler).Methods("PUT")
	r.HandleFunc("/{filename}", getHandler).Methods("GET")

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: r,
	}
	s.ListenAndServe()
	fmt.Println("Exit...")
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	filename := vars["filename"]

	reader, contentType, contentLength, err := readFile(filename)
	if err != nil {
		if err.Error() == "The specified key does not exist." {
			http.Error(w, "File not found", 404)
			return
		} else {
			fmt.Printf("%s", err.Error())
			http.Error(w, "Could not retrieve file.", 500)
			return
		}
	}
	defer reader.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatUint(contentLength, 10))
	w.Header().Set("Connection", "close")

	if _, err = io.Copy(w, reader); err != nil {
		fmt.Printf("%s", err.Error())
		http.Error(w, "Error occured copying to output stream", 500)
		return
	}
}

func readFile(filename string) (reader io.ReadCloser, contentType string, contentLength uint64, err error) {
	path := filepath.Join(config.Storedir, filename)
	// content type , content length
	if reader, err = os.Open(path); err != nil {
		return
	}
	var fi os.FileInfo
	if fi, err = os.Lstat(path); err != nil {
		return
	}
	contentLength = uint64(fi.Size())
	contentType = mime.TypeByExtension(filepath.Ext(filename))

	return
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileext := filepath.Ext(vars["filename"])
	contentLength := r.ContentLength
	contentType := r.Header.Get("Content-Type")
	token := Encode(10000000 + int64(rand.Intn(1000000000)))
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(vars["filename"]))
	}

	var reader io.Reader
	reader = r.Body

	w.Header().Set("Content-Type", "text/plain")

	writeFile(token, fileext, reader, contentType, uint64(contentLength))

	fmt.Fprintf(w, "http://%s/%s%s\n", r.Host, token, fileext)
}

func writeFile(token string, fileext string, reader io.Reader, contentType string, contentLength uint64) error {
	var f io.WriteCloser
	var err error
	filename := token + fileext
	if f, err = os.OpenFile(filepath.Join(config.Storedir, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		fmt.Printf("%s", err)
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, reader); err != nil {
		return err
	}

	return nil
}

func Encode(number int64) string {
	rest := number % BASE
	// strings are a bit weird in go...
	result := string(SYMBOLS[rest])
	if number-rest != 0 {
		newnumber := (number - rest) / BASE
		result = Encode(newnumber) + result
	}
	return result
}
