// package to ease fetching/parsing of rss-feeds
package netrss

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"html/template"
	"log"
	"net/http"
)

type Rss2 struct {
	XMLName     xml.Name `xml:"rss"`
	Version     string   `xml:"version,attr"`
	Title       string   `xml:"channel>title"`
	Link        string   `xml:"channel>link"`
	Description string   `xml:"channel>description`
	PubDate     string   `xml:"channel>pubDate"`
	ItemList    []Item   `xml:"channel>item"`
}

type Item struct {
	Title       string        `xml:"title"`
	Link        string        `xml:"link"`
	Description template.HTML `xml:"description"`
	Content     template.HTML `xml:"encoded"`
	PubDate     string        `xml:"pubDate"`
	Comments    string        `xml:"comments"`
}

type NetRss struct {
	Address string
}

func (nr *NetRss) ParseFeedContent() (Rss2, bool) {
	v := Rss2{}

	if nr.Address == "" {
		log.Println("Missing address...")
		panic("Missing Address")
	}

	resp, err := http.Get(nr.Address)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// because we dunno respons' charset
	// we convert in advance
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&v)

	if err != nil {
		log.Println(err)
		panic(err)
	}

	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
	}

	if v.Version == "2.0" {
		for i, _ := range v.ItemList {
			if v.ItemList[i].Content != "" {
				v.ItemList[i].Description = v.ItemList[i].Content
			}
		}
		return v, true
	}
	log.Println("not RSS 2.0")
	return v, false
}

/*
//	import (
//		"github.com/erpe/netrss"
//		"encoding/xml"
//	)
//
//	func main() {
//		np := netrss.NetRss{ Address: 'https://netzpolitik.org/rss' }
//		rss2, b := np.ParseFeedContent()
//		if b == false {
//			log.Println("parseFeedContent returned false...")
//		}
//		for _, e := range rss2.ItemList {
//			fmt.Println(e.Title)
//		}
//	}
*/
