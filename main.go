package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	colorful "github.com/lucasb-eyer/go-colorful"
	yaml "gopkg.in/yaml.v2"
)

func main() {

	// get and parse yaml
	req, err := http.Get("https://raw.githubusercontent.com/github/linguist/master/lib/linguist/languages.yml")
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]map[string]interface{}
	err = yaml.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	readme := "# Github Language Colors\n\n"
	link := "[![](https://placeholder.pics/svg/600x50/%s/%s-FFF/%s)](%s)"

	keys := make([]string, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// find langs with colors
	for _, lang := range keys {
		meta := data[lang]
		if meta["color"] != nil {
			color, ok := meta["color"].(string)
			if ok {
				url := fmt.Sprintf("https://github.com/trending?l=%s", lang)
				c, err := colorful.Hex(color)
				if err != nil {
					log.Fatal(err)
				}
				_, _, l := c.Hcl()
				fmt.Printf("%s %s %.04f\n", lang, color[1:], l)
				textColor := "#FFF"
				if l > 0.7 {
					textColor = "#000"
				}

				// encode any spaces
				lang = strings.Replace(lang, " ", "%20", -1)
				url = strings.Replace(url, " ", "%20", -1)

				// encode any single quotes
				lang = strings.Replace(lang, "'", "&apos;", -1)
				url = strings.Replace(url, "'", "&apos;", -1)
				readme += fmt.Sprintf(link, color[1:], textColor[1:], lang, url)
			}
		}
	}

	// generate README.md
	err = ioutil.WriteFile("./README.md", []byte(readme), 0644)
	if err != nil {
		log.Fatal(err)
	}

}
