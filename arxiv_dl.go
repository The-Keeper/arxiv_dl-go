package main

import (
    "flag"				// for using command-line arguments
	"errors"			// for having custom error handling
	"fmt"				// for outputing and formatting strings
	"io"				// for writing data to file
	"net/http"			// for getting data from the website
	"os"				// for file and directory operations
    "strings"			// for splitting strings
	"path/filepath"		// for joining paths
	"sync"				// for goroutine WaitGroups
	"github.com/codeclysm/extract/v3" // for having source files extracted
	"bytes"				// required indirectly for extract library
	"context"			
	"io/ioutil"
)

// downloads a URL to fileName
func downloadFile(URL, fileName string) error {
	// get the response from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Response code: %d", response.StatusCode))
	}
	// make an empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// write data to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// extracts arguments from comma-delimited string and ensures that formatting is correct
func getArgsFromString(nstr string) []string {
    arrSpl := strings.Split(nstr, ",")

    // filter out empty strings
    var res []string
	for _, s := range arrSpl {
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}

// constants declaration
const DL_DIR = "dl"
const EXTR_DIR = "extracted"

// WaitGroup for the downloading queue
var wg sync.WaitGroup

// Renamer object for archive extraction
var Renamer extract.Renamer

// downloads the archive and extracts
func arcDownload(n string) {
	defer wg.Done()
	dlURL := fmt.Sprintf("https://arxiv.org/e-print/%s", n)
	arcPath := filepath.Join(DL_DIR, n)
	extrPath := filepath.Join(EXTR_DIR, n)
	downloadFile(dlURL, arcPath)
	fmt.Printf("\nDownload finished: %s", dlURL)

	data, _ := ioutil.ReadFile(arcPath)
	buffer := bytes.NewBuffer(data)
	
	extract.Archive(context.Background(), buffer, extrPath, Renamer)
	fmt.Printf("\nSource extracted: %s", n)
}

func main() {
    // variables declaration  
    var fn_str string    
 
    // flags declaration using flag package
    flag.StringVar(&fn_str, "n", "", "comma-separated list of arxiv ids")
    flag.Parse()  // after declaring flags we need to call it
	
    arr_n := getArgsFromString(fn_str)
    fmt.Printf("%v", arr_n)

	// create download directory
	os.Mkdir(DL_DIR, os.ModePerm)

	// add files to WaitGroup
	for _, element := range arr_n {
		wg.Add(1)
		go arcDownload(element)
	}

	// wait for downloading to finish
	wg.Wait()
}