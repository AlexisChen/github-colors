package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	colorful "github.com/lucasb-eyer/go-colorful"
	yaml "gopkg.in/yaml.v2"
)

func main() {

	// get and decode yaml
	baseURL := "https://raw.githubusercontent.com"
	req, err := http.Get(baseURL + "/github/linguist/master/lib/linguist/languages.yml")
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]map[string]interface{}
	err = yaml.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	// generate svgs and readme
	readme := "# Github Language Colors\n\n"
	link := "[![](./svgs/%s.svg)](%s)\n"

	// sort languages by name
	keys := make([]string, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// find languages with colors
	for _, lang := range keys {
		meta := data[lang]
		if meta["color"] != nil {
			color, ok := meta["color"].(string)
			if ok {

				url := fmt.Sprintf("https://github.com/trending?l=%s", lang)

				// check if color is light or dark to determine text color
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

				// create svg images
				cleaned := strings.Replace(lang, " ", "-", -1)
				cleaned = strings.Replace(cleaned, "'", "-", -1)

				t, err := template.ParseFiles("./svg.tmpl")
				if err != nil {
					log.Fatal(err)
				}

				var svg bytes.Buffer

				templateData := struct {
					LangName  string
					LangColor string
					TextColor string
				}{
					lang,
					color,
					textColor,
				}
				err = t.Execute(&svg, templateData)
				if err != nil {
					log.Fatal(err)
				}

				err = ioutil.WriteFile("./svgs/"+cleaned+".svg", svg.Bytes(), 0644)
				if err != nil {
					log.Fatal(err)
				}

				// encode any spaces
				lang = strings.Replace(lang, " ", "%20", -1)
				url = strings.Replace(url, " ", "%20", -1)

				// encode any single quotes
				lang = strings.Replace(lang, "'", "&apos;", -1)
				url = strings.Replace(url, "'", "&apos;", -1)

				// add language to readme
				readme += fmt.Sprintf(link, cleaned, url)
			}
		}
	}

	// create README.md
	err = ioutil.WriteFile("./README.md", []byte(readme), 0644)
	if err != nil {
		log.Fatal(err)
	}

}
