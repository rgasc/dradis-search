# Dradis search
This is a small program written in Go to search for issues in Dradis by their title and content. It returns a list of links to the issues that are found, along with their title.

## Installation
Download a prebuilt binary from the [releases](https://github.com/rgasc/dradis-search/releases) page and run it as shown in the [usage](#Usage) section.

Or, clone this repository:
```shell
git clone https://github.com/rgasc/dradis-search.git
```

And compile it:
```shell
go build -o dradis-search *.go
```

## Usage
Run the binary as follows:
```shell
./dradis-search -baseurl="https://YOUR_BASE_URL/pro" -apikey="YOUR_API_KEY" -term="SEARCH TERM"
```
