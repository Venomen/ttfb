# ttfb
tool in golang for testing website speed and body & headers search

### Description
After using over and over "while curl" loops, I made simple equivalent in go. 
Just for fun and making life easier ;-)
 
### Instalation

- get it: `go get github.com/Venomen/ttfb`
- install it: `go install github.com/Venomen/ttfb`
- link it: `ln -s ~/go/bin/ttfb /usr/local/bin/ttfb`

### Usage

- run from your `~/go/bin/ttfb` or after symlink just `ttfb`
- remember to edit your `~/.ttfbEnv` file (1st install will copy default conf) 
- config options: 

"<b>url</b>" - provide url for testing, with http:// or https:// <br>
"<b>search</b>" - you can grep some data in Header & Body from your request (wildcards available)

### Functions
- http client
- dot .env configuration
- cli for config options
- measuring connection speed; time to first byte