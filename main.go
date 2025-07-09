package main

import (
    "bufio"
    "errors"
    "flag"
    "fmt"
    "github.com/joho/godotenv"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"
)

// main metadata
var version = "0.0.2"
var copyRights = "deregowski.net (c) 2020"

var urlStr = ""
var search = ""

var home, _ = os.UserHomeDir()
var homeInside = filepath.Join(home, "/.ttfbEnv")

// Shared transport for connection reuse and DNS cache (default mode)
var sharedTransport = &http.Transport{}

// Checks if a file exists and is not a directory
func configExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

// Loads environment variable from .ttfbEnv file
func goDotEnvVariable(key string) string {
    loadEnv := godotenv.Load(homeInside)
    if loadEnv != nil {
        log.Fatalf("No ~/.ttfbEnv file! Please create it (see README.md) and run again!")
    }
    return os.Getenv(key)
}

// Validates if the provided string is a valid URL
func validateURL(u string) error {
    parsed, err := url.ParseRequestURI(u)
    if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
        return errors.New("invalid URL: must start with http:// or https://")
    }
    return nil
}

// Validates if the search string is not empty
func validateSearch(s string) error {
    if strings.TrimSpace(s) == "" {
        return errors.New("search string cannot be empty")
    }
    return nil
}

// Reads user data from CLI flags, env file, or prompts
func enterData() {
    reader := bufio.NewReader(os.Stdin)

    // CLI flags
    flagUrl := flag.String("url", "", "URL to test (overrides .ttfbEnv and positional argument)")
    flagSearch := flag.String("search", "", "Regex to search in HTML (overrides .ttfbEnv)")
    noCache := flag.Bool("no-cache", false, "Disable HTTP keep-alive and DNS cache (each request uses a new connection)")
    showVersion := flag.Bool("version", false, "Show version and exit")
    showHelp := flag.Bool("help", false, "Show help and exit")
    flag.Usage = func() {
        fmt.Println("Usage: ttfb [--url URL] [--search REGEX] [--no-cache] [domain]")
        fmt.Println("Options:")
        fmt.Println("  --url        URL to test (overrides .ttfbEnv and positional argument)")
        fmt.Println("  --search     Regex to search in HTML (overrides .ttfbEnv)")
        fmt.Println("  --no-cache   Disable HTTP keep-alive and DNS cache (each request uses a new connection)")
        fmt.Println("  --version    Show version and exit")
        fmt.Println("  --help       Show this help message and exit")
        fmt.Println("\nYou can also provide the domain as the first positional argument.")
        fmt.Println("\nIf .ttfbEnv file exists in your home directory, values from it will be used as defaults.")
        fmt.Println("Order of precedence: CLI flag > positional argument > .ttfbEnv > interactive prompt")
    }

    // Custom flag error handling: show help on unknown flags
    flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
    flag.CommandLine.SetOutput(io.Discard) // Suppress default error output

    if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
        fmt.Println("Unknown or invalid flag detected.")
        flag.Usage()
        // Use os.Exit(0) to avoid showing "exit status 1"
        os.Exit(0)
    }

    // Restore output for further prints
    flag.CommandLine.SetOutput(os.Stderr)

    if *showHelp {
        flag.Usage()
        os.Exit(0)
    }
    if *showVersion {
        fmt.Printf("version [%s] by %s\n", version, copyRights)
        os.Exit(0)
    }

    // Load from env file
    dotenvUrl := goDotEnvVariable("url")
    dotenvSearch := goDotEnvVariable("search")

    // If user provided a domain as the first argument (without flag)
    if *flagUrl == "" && len(flag.Args()) > 0 {
        *flagUrl = flag.Args()[0]
    }

    // Set url and search with priority: CLI > positional > env > prompt
    if *flagUrl != "" {
        urlStr = *flagUrl
    } else {
        urlStr = dotenvUrl
    }
    if *flagSearch != "" {
        search = *flagSearch
    } else {
        search = dotenvSearch
    }

    fmt.Println("ttfb, testing your slow website since 2020 ;-)")
    fmt.Println("---------------------")
    fmt.Println("You can test custom url by .ttfbEnv, CLI flag or as a positional argument")
    fmt.Println("There is also 'search' if you need some HTML body & headers info")
    fmt.Println("Use --no-cache to disable keep-alive and DNS cache for each request")
    fmt.Println("---------------------")

    // Prompt for url if still empty or invalid
    for {
        if err := validateURL(urlStr); err != nil {
            fmt.Printf("-> 'url' is invalid or empty - please provide with http:// or https:// !\n")
            text, _ := reader.ReadString('\n')
            urlStr = strings.TrimSpace(text)
        } else {
            break
        }
    }

    // Prompt for search if still empty or invalid
    for {
        if err := validateSearch(search); err != nil {
            fmt.Printf("-> 'search' is empty - please provide (you can use wildcards - for ex. '.*com')?\n")
            data, _ := reader.ReadString('\n')
            search = strings.TrimSpace(data)
        } else {
            break
        }
    }

    // Store noCache flag for use in run()
    if *noCache {
        os.Setenv("TTFB_NO_CACHE", "1")
    } else {
        os.Unsetenv("TTFB_NO_CACHE")
    }
}

// Measures the time to first byte for a given URL
func measureTTFB(url string, noCache bool) (time.Duration, *http.Response, error) {
    start := time.Now()
    var client *http.Client
    if noCache {
        // New transport, no keep-alive, no DNS cache
        client = &http.Client{
            Transport: &http.Transport{
                DisableKeepAlives: true,
            },
        }
    } else {
        // Shared transport (keep-alive and DNS cache for process lifetime)
        client = &http.Client{
            Transport: sharedTransport,
        }
    }
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return 0, nil, err
    }
    resp, err := client.Do(req)
    if err != nil {
        return 0, nil, err
    }
    // Read a single byte to ensure we measure TTFB
    buf := make([]byte, 1)
    _, err = resp.Body.Read(buf)
    ttfb := time.Since(start)
    if err != nil && err != io.EOF {
        resp.Body.Close()
        return ttfb, resp, err
    }
    return ttfb, resp, nil
}

// Main logic
func run() {
    enterData()
    noCache := os.Getenv("TTFB_NO_CACHE") == "1"
    fmt.Printf("Starting ttfb to %s \n", urlStr)
    if noCache {
        fmt.Println("HTTP keep-alive and DNS cache are DISABLED (--no-cache)")
    } else {
        fmt.Println("HTTP keep-alive and DNS cache are ENABLED (default)")
    }

    ttfb, resp, err := measureTTFB(urlStr, noCache)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    fmt.Printf("TTFB: %s\n", ttfb)

    if resp.StatusCode == http.StatusOK {
        bodyBytes, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }
        bodyString := string(bodyBytes)

        re := regexp.MustCompile(search)
        matches := re.FindStringSubmatch(bodyString)
        fmt.Println("HTML Search Results:", matches)
    } else {
        fmt.Printf("Ups... no connection to %s. Please check your internet\n", urlStr)
    }
}

func main() {
    // Check config file and copy if needed
    if configExists(homeInside) {
        // All ok
    } else {
        fmt.Println("Config file does not exist, linking default to ~/.ttfbEnv")
        fmt.Print("Please edit it after all.\n\n")
        err := os.Link(".ttfbEnv", homeInside)
        if err != nil {
            // Try to copy file if hard link fails
            src, errOpen := os.Open(".ttfbEnv")
            if errOpen != nil {
                log.Fatalf("Cannot open .ttfbEnv: %v", errOpen)
            }
            defer src.Close()
            dst, errCreate := os.Create(homeInside)
            if errCreate != nil {
                log.Fatalf("Cannot create %s: %v", homeInside, errCreate)
            }
            defer dst.Close()
            _, errCopy := io.Copy(dst, src)
            if errCopy != nil {
                log.Fatalf("Cannot copy config file: %v", errCopy)
            }
            fmt.Println("Copied .ttfbEnv to ~/.ttfbEnv")
        }
    }

	run()

}