package rss

import (
	"strings"
	"net/http"
	"io/ioutil"
	"golang.org/x/net/html/charset"
	"log"
)

func check(err error) {
	if err != nil {
		log.Println("Failed to parse rss: ", err)
	}
}

func ParseRss(url string, args ...string) (out string) {
	var begin, end string
	if len(args) == 2 {
		begin = args[0]
		end = args[1]
	} else {
		begin = "<![CDATA["
		end = "]]>"
	}

	resp, receivingError := http.Get(url)
	check(receivingError)
	defer resp.Body.Close()

	utf8, decodingError := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	check(decodingError)

	body, parsingError := ioutil.ReadAll(utf8)
	check(parsingError)

	code := string(body)
	var from, to int
	for true {
		from = strings.Index(code, begin) + len(begin)
		to = strings.Index(code, end)
		if from == -1 || to == -1 { break }
		if from - to < 4000 {
			out += code[from:to] + "\\end\\\n"
			code = code[to+len(end):]
		}
	}

	filter := []struct{
		old string
		new string
	}{
		{"<a>", ""}, 		{"</a>", ""},
		{"<p>", ""}, 		{"</p>", ""},
		{"<em>", ""},		{"</em>", ""},
		{"<pre>", ""},		{"</pre>", ""},
		{"<code>", ""},	{"</code>", ""},
		{"<li>", ""},		{"</li>", ""},
		{"<ol>", ""},		{"</ol>", ""},
		{"<ul>", ""},		{"</ul>", ""},
		{"<br>", "\n"},	{"<a href=", ""},
		{"&lt;", "\""},		{"&gt;", "\""},
		{"&quot;", "\""},	{">", " "},
	}

	for _, replacement := range filter {
		out = strings.Replace(out, replacement.old, replacement.new, -1)
	}
	return
}