package main

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// important data

var version = "0.0.1"
var copyRights = "deregowski.net (c) 2020"

var url = ""
var search = ""

var home, _ = os.UserHomeDir()
var homeInside = filepath.Join(home, "/.ttfbEnv")

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.

func configExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func goDotEnvVariable(key string) string {

	// load .ttfbEnv file
	loadEnv := godotenv.Load(homeInside)

	if loadEnv != nil {
		log.Fatalf("No ~/.ttfbEnv file! Please create it (description in README.md) and run again!")
	}

	return os.Getenv(key)
}

// provide your data

func enterData() {

	dotenvUrl := goDotEnvVariable("url")
	url = dotenvUrl

	dotenvSearch := goDotEnvVariable("search")
	search = dotenvSearch

	fmt.Println("ttfb, testing your slow website since 2020 ;-)")
	fmt.Println("---------------------")
	fmt.Println("You can test custom url by .ttfbEnv or inline")
	fmt.Println("There is also 'search' if you need some HTML body & headers info")
	fmt.Println("---------------------")

	reader := bufio.NewReader(os.Stdin)

	if url == "" {
		// define url
		for {
			fmt.Print("-> 'url' in .ttfbEnv is empty - please provide with http:// or https:// !\n")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			url = strings.Replace(text, "\n", "", -1)
			break
		}

	}

	if search == "" {
		// define search
		for {
			fmt.Print("-> 'search' in .ttfbEnv is empty - please provide (you can use wildcards - for ex. '.*com')?\n")
			data, _ := reader.ReadString('\n')
			// convert CRLF to LF
			search = strings.Replace(data, "\n", "", -1)
			break
		}
	}
}

// run all

func run() {

	enterData()

	// setting the start test time..
	currentTime := time.Now()
	fmt.Printf("Starting ttfb to %s \n", url)

	// ..removing dots
	var cutStartTime = currentTime.Format("15.04.05.000")
	startTime := strings.Replace(cutStartTime, ".", "", -1)

	// setting http client
	var client http.Client
	resp, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		// making HTTP request
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		// searching provided html data
		re := regexp.MustCompile(search)
		matches := re.FindStringSubmatch(bodyString)
		fmt.Println("HTML Search Results:", matches)

		// setting the end test time..
		endTime := time.Now()

		// ..removing dots
		var cutEndTime = endTime.Format("15.04.05.000")
		stopTime := strings.Replace(cutEndTime, ".", "", -1)

		// convert time to int
		stoptimeInt, _ := strconv.Atoi(stopTime)
		starttimeInt, _ := strconv.Atoi(startTime)

		// calculate final time
		final := stoptimeInt - starttimeInt
		finalDuration := time.Duration(final) * time.Millisecond

		fmt.Printf("\nYour test took %s\n", finalDuration)

	} else {
		fmt.Printf("Ups... no connection to %s. Please check your internet\n", url)
	}

}

func main() {

	if len(os.Args) > 1 {
		firstArg := os.Args[1]

		if firstArg == "version" || firstArg == "help" {
			fmt.Printf("version [%s] by %s\n", version, copyRights)
			os.Exit(0)
		}
	}

	if configExists(homeInside) {
		_ = "all ok, go forward"
	} else {
		fmt.Println("Config file does not exist, linking default to ~/.ttfbEnv")
		fmt.Println("Please edit it after all.")
		os.Link(".ttfbEnv", homeInside)
		// TODO: handling copy-file error
	}

	run()
}
