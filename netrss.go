package netrss

import (
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	Feed    []byte
	fetched bool
	url     url.URL
}

func (nr *NetRss) fetchSourceFeed() bool {
	resp, err := http.Get(nr.Address)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}
	nr.Feed = body
	return true
}

func (nr *NetRss) ParseFeedContent() (Rss2, bool) {
	v := Rss2{}

	if nr.Address == "" {
		log.Println("Missing address...")
		return v, false
	}

	if nr.fetched == false {
		nr.fetchSourceFeed()
		nr.fetched = true
	}

	// TODO: need a decoder in case of wrong charset
	err := xml.Unmarshal(nr.Feed, &v)

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
/*	import (
/*		"github.com/erpe/netrss"
/*		"encoding/xml"
/*	)
/*
/*	func main() {
/*		np := netrss.NetRss{ Address: 'https://netzpolitik.org/rss' }
/*		rss2, b := np.ParseFeedContent()
/*		if b == false {
/*			log.Println("parseFeedContent returned false...")
/*		}
/*		for _, e := range rss2.ItemList {
/*			fmt.Println(e.Title)
/*		}
/*	}
*/
