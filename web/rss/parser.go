package rss

import (
	"strings"
	"net/http"
	"io/ioutil"
	"golang.org/x/net/html/charset"
	"log"
	"main/conf"
)

func check(err error) {
	if err != nil {
		log.Println("Failed to parse rss: ", err)
	}
}

func ParseRss(url string, args ...string) (out []string) {
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
		if to - from < conf.MAX_MESSAGE_LEN {
			out = append(out, code[from:to])
		}
		code = code[to+len(end):]
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
		for index, story := range out {
			out[index] = strings.Replace(story, replacement.old, replacement.new, -1)
		}
	}

	return
}