# PathFinder 🕵🏻‍♂️🌐

PathFinder – the ultimate web crawler script designed for lightning-fast, concurrent, and recursive URL scraping. Cutting-edge multithreading architecture ensures rapid URL extraction while guaranteeing that each page is visited only once. This tool extracts URLs from various HTML tags, including `a`, `form`, `iframe`, `img`, `embed`, and more. External URLs, relative paths and subdomains are supported as well.

![Example PathFinder usage](_img/img-01.jpg)

**Usage**: It is a great tool for discovering new web paths and subdomains, creating a site map, gathering OSINT information. It might be very useful for bug hunters and pentesters. 🔥👾🔥

## Installation

There are 2 options:

1. Download latest binary from [GitHub releases](https://github.com/Print3M/pathfinder/releases).
2. Build manually:

#### Easy install
```bash
go install -v github.com/Print3M/pathfinder@latest
```

OR

```bash
# Download and build the source code
git clone https://github.com/Print3M/pathfinder
cd pathfinder/
go build

# Run
./pathfinder --help
```

## Flags
```bash
PathFinder is a crawler script for concurrent and recursive scraping of URLs from any website.

Usage:
  PathFinder [flags]

Flags:
  -H, --headers stringArray   Add HTTP header (one -H must contain only one header)
  -h, --help                  help for PathFinder
      --no-externals          Disable external URLs scraping
      --no-recursion          Disable recursive scraping
      --no-subdomains         Disable subdomains scraping
  -o, --output string         Output file
  -q, --quiet                 Disable printing scraped URLs on the screen
  -r, --rate uint             Number of requests per second
  -t, --threads uint          Number of concurrent threads (default 10)
  -u, --url string            URL to start
      --with-assets           Enable asset URLs scraping (images, CSS, JS etc.)
```

## How to use it?

TL;DR;

```bash
# Run URLs extraction
pathfinder -u http://example.com --threads 25

# Show help
pathfinder -h
```

### Initial URL

`-u <url>`, `--url <url` [required]

Use this parameter to set the start page for the script. By default, the script extracts all URLs that refer to the domain provided in this parameter and its subdomains. External URLs are scraped but not visited by the script.

### Threads

`-t <num>`, `--threads <num>` [default: 10]

Use this parameter to set the number of threads that will extract data concurrently.

## Rate limiting

`-r <reqs_per_sec>`, `--rate <reqs_per_sec>` [default: none]

Use this parameter to specify max number of requests per second. It's a total number of requests per second - number of threads doesn't matter. By default, requests are sent as fast as possible! 🚄💨

### Output file

`-o <file>`, `--output <file>` [default: none]

Use this parameter to specify output file where URLs will be saved. However, they will be still printed out on the screen if you do not use quiet mode (`-q` or `--quite`).

Output is saved to a file on the fly, so even if you stop executing the script the downloaded URLs will be saved.  

### Add HTTP header

`-H <header>`, `--header <header>` [default: none]

Use this parameter to specify custom HTTP headers. They are used with every scraping request. One `-H` parameter must contain only one HTTP header but you can use it multiple times. Headers

Example:

`./pathfinder ... -H "Authorization: test" -H "Cookies: cookie1=choco; cookie2=yummy"`

### Quiet mode

`-q`, `--quiet` [default: false]

Use this parameter to disable printing out scraped URLs on the screen.

### Disable recursive scraping

`--no-recursion` [default: false]

Use this parameter to disable recursive scraping. No other page will be visited except the one you provided using `-u <url>` parameter. Only one page will be visited. It actually disables what's coolest about PathFinder.

### Disable subdomains scraping

`--no-subdomains` [default: false]

Use this parameter to disable scraping of subdomains of the URL provided using `-u <url>` parameter.

Example (`-u http://test.example.com`):

- `http://test.example.com/index.php` - ✅ scraped
- `http://api.test.example.com/index.php` - ✅ scraped
- `http://example.com/index.php` - ❌ not scraped.

Example (`-u http://test.example.com --no-subdomains`):

- `http://test.example.com/index.php` - ✅ scraped.
- `http://api.test.example.com/index.php` - ❌ not scraped.
- `http://example.com/index.php` - ❌ not scraped.

### Disable externals scraping

`--no-externals` [default: false]

External URLs are not visited anyway, but using this parameter you filter out all external URLs from the output.

### Enable scraping of static assets

`--with-assets` [default: false]

Use this parameter to enable scraping URLs of static assets like CSS, JavaScript, images, fonts and so on. This is disabled by default because it usually generates too much noise.

### User agent

User agent is randomly changed on each request from a set of predefined strings.
