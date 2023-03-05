package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func getEnv() string {
    if len(os.Args) < 2 {
        fmt.Printf("Please provide a folder to serve. Exiting now\n")
        os.Exit(1)
    }

    FOLDER := os.Args[1]

    return FOLDER
}

func fileExists(p string) bool {
    _, notExistErr := os.Stat(p)
    return !os.IsNotExist(notExistErr) 
}

func handleResponse(w http.ResponseWriter, r *http.Request) {
    file := r.URL.Path
    folder := getEnv()
    if file[len(file) - 1] == "/"[0] {
        file = strings.Join([]string{file, "index.html"}, "")
    } 

    filePath := strings.Join([]string {folder, file}, "")
    splitPath := strings.Split(filePath, "/")
    splitFilename := strings.Split(splitPath[len(splitPath) - 1], ".")
    fileExtension := splitFilename[len(splitFilename) - 1] 


    fmt.Printf("got request to %s\n", filePath)
    fmt.Printf("file ext is %s\n", fileExtension)
    
    if fileExists(filePath) {
        content, err := ioutil.ReadFile(filePath)
        if err != nil {
            content, err := ioutil.ReadFile(strings.Join([]string{filePath, "/index.html"}, ""))
            if err != nil {
                io.WriteString(w, "500 server error")
                fmt.Printf("File %s does not exist\n", strings.Join([]string{filePath, "/index.html"}, ""))
            } else {
                io.WriteString(w, string(content))
            }
        } else {
            w.Header().Add("test-header", "this is a test header")
            w.Header().Add("Content-Type", strings.Join([]string{"text/", fileExtension, "; charset=utf-8"}, ""))
            io.WriteString(w, string(content))
        }
    } else {
        fmt.Printf("File does not exist\n")
        io.WriteString(w, "404 not found")
    }
}


func main() {
    FOLDER := getEnv()

    if !fileExists(FOLDER) {
        fmt.Printf("Please provide a path to an existing folder. Exiting now\n")
        os.Exit(1)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/", handleResponse)

    err := http.ListenAndServe(":3333", mux)

    if errors.Is(err, http.ErrServerClosed) {
        fmt.Printf("server closed\n")
    } else if err != nil {
        fmt.Printf("error starting server: %s\n", err)
        os.Exit(1)
    }
}
