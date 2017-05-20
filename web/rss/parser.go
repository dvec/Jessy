package rss

import (
	"strings"
	"net/http"
	"io/ioutil"
	"golang.org/x/net/html/charset"
	"main/conf"
)

const (
	begin = "<![CDATA[" //The expression to begin useful data
	end   = "]]>" //The expression to end useful data
)

//Values that will be replaced
var filter = []struct{
	old string //Old value
	new string //New value
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
	{"&lt;", "\""},	{"&gt;", "\""},
	{"&quot;", "\""},	{">", " "},
}

//A function that loads and parses RSS data
func GetRSSData(url string) ([]string, error) {
	var out []string

	//Receiving data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//Reading data
	decodedData, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type")) //Setting the encoding
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(decodedData)
	if err != nil {
		return nil, err
	}
	code := string(body)

	//Parsing useful data
	var from, to int
	for {
		from = strings.Index(code, begin) + len(begin)
		to = strings.Index(code, end)
		if from == -1 || to == -1 { break }
		if to - from < conf.MaxMessageLen {
			out = append(out, code[from:to])
		}
		code = code[to+len(end):]
	}

	//Filtering a data
	for _, replacement := range filter {
		for index, story := range out {
			out[index] = strings.Replace(story, replacement.old, replacement.new, -1)
		}
	}

	return out, nil
}