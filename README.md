> **THIS REPO IS UNDER HEAVY DEVELOPMENT...**

# PathFinder

PathFinder – the ultimate web crawler script designed for lightning-fast, concurrent, and recursive URL scraping. With its cutting-edge multithreading architecture, PathFinder ensures rapid URL extraction while guaranteeing that each page is visited only once. This powerful tool can effortlessly extract paths from various HTML tags, including `a`, `form`, `iframe`, `img`, `embed`, and more.

## Usage

TL;DR;

```bash
# Run path extraction
pathfinder -u http://example.com --threads 25

# Show help
pathfinder -h
```

### Initial URL

`-u <url>`, `--url <url` [required]

Use this parameter to set the start page for the script. By default, the script extracts all paths that refer to the domain provided in this parameter and its subdomains. External websites are considered as out of the scope, they are not visited by the script.

### Threads

`-t <num>`, `--threads <num>` [default: 10]

Use this parameter to set the number of threads that will extract data concurrently.

### Output file

SOON...

### Disable recursive scraping

`--no-recursion` [default: false]

Use this parameter to disable recursive scraping. No other page will be visited except the one you provided using `-u <url>` parameter. Only one page will be visited. It actually disables what's coolest about PathFinder.

### Disable subdomains scraping

`--no-subdomains` [default: false]

Use this parameter to disable scraping of subdomains to the URL provided using `-u <url>` parameter.

Example (`pathfinder -u http://test.example.com`):

- `http://test.example.com/index.php` - is scraped.
- `http://api.test.example.com/index.php` - is scraped.
- `http://example.com/index.php` - is not scraped.

Example (`pathfinder -u http://test.example.com --no-subdomains`):

- `http://test.example.com/index.php` - is scraped.
- `http://api.test.example.com/index.php` - is not scraped.
- `http://example.com/index.php` - is not scraped.

### Enable scraping of static assets

`--with-assets` [default: false]

Use this parameter to enable scraping paths of static assets like CSS, JavaScript, images, fonts and so on. This is disabled by default because it usually generates too much noise.
