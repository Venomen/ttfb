# ttfb
tool in golang for testing website speed and body & headers search

### Description
After using over and over "while curl" loops, I made simple equivalent in go. You don't have to use Web Browser with Debug, curl or some other networking stuff. 
Just for fun and making life easier ;-)
 
### Instalation

- go min. version: 1.14.5
- get it: `go get github.com/Venomen/ttfb`
- install it: `go install github.com/Venomen/ttfb`
- link it: `ln -s ~/go/bin/ttfb /usr/local/bin/ttfb`

### Usage

- run from your `~/go/bin/ttfb` or after symlink just `ttfb`
- remember to edit your `~/.ttfbEnv` file (1st install will copy default conf) 
- config options: 

`--url` - provide url for testing, with http:// or https:// <br>
`--search` - you can grep some data in Header & Body from your request (wildcards available) <br>
`--no-cache` - disable HTTP keep-alive and DNS cache (each request uses a new connection, --cache is by default)

### Functions
- http client
- cache on/off (dns reusing transport & keep-alive mode)
- dot .env configuration
- cli for config options
- searching source code elements
- measuring connection speed; time to first byte
