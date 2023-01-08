package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gomarkdown/markdown"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"strings"
)

var BuildIndexPresent bool
var BuildConfigs []BuildConfig

type BuildConfig struct {
	PostDate    string `yaml:"postDate"`
	PublishDate string `yaml:"publishDate"`
	Title       string `yaml:"title"`
	Category    string `yaml:"category"`
	ContentPath string `yaml:"contentPath"`
	ConfigName  string `yaml:"configName"`
}

func main() {
	// Nuke build/
	err := os.RemoveAll(path.Join("./build"))
	if err != nil {
		log.Fatalln(err)
	}

	err = os.Mkdir("build", 0755)
	if err != nil {
		log.Fatalln(err)
	}

	// Iterate over content/ files
	items, _ := os.ReadDir("./content")
	log.Println("Iterating configuration files...")
	for _, item := range items {
		yfile, err := os.ReadFile(path.Join("./content", item.Name()))
		if err != nil {
			log.Fatalln(err)
		}

		config := BuildConfig{}
		err = yaml.Unmarshal(yfile, &config)
		if err != nil {
			log.Fatalln(err)
		}

		// If the config file doesn't specify a path to the content, attempt to find a markdown file with the same name
		// in the includes/ directory, otherwise we'll have to drop this configuration
		if config.ContentPath == "" {
			log.Printf("No content path defined in config file '%s', attempting auto discovery \n", item.Name())

			parts := strings.Split(item.Name(), ".")
			discoverName := parts[:len(parts)-1][0] + ".md"

			log.Printf("Searching includes for content file '%s'", discoverName)
			_, err := os.Stat(path.Join("./includes", discoverName))

			if err != nil && errors.Is(err, os.ErrNotExist) {
				log.Printf("Unable to auto discover content for config file '%s', skipping", item.Name())
				continue
			} else if err != nil {
				log.Printf("Unknown error occurred during content discovery for config file '%s', skipping", item.Name())
				continue
			} else {
				log.Printf("Content auto-discover successful for config file '%s', located include '%s'", item.Name(), discoverName)
				config.ContentPath = discoverName
			}
		}

		config.ConfigName = item.Name()
		BuildConfigs = append(BuildConfigs, config)
	}

	// Commence build process
	footer, _ := os.ReadFile("./template/footer.html")
	normalize, _ := os.ReadFile("./template/normalize.css")
	templateCss, _ := os.ReadFile("./template/style.css")
	titlePlaceholder := []byte("||PAGE_TITLE||")

	for _, item := range BuildConfigs {
		log.Printf("Preparing to build config %s\n", item.ConfigName)
		// At this point if we fail to load the include after all the sanity checking we've done
		// then the file is just flat out missing, nothing we can do about that. Skip and move on
		content, err := os.ReadFile(path.Join("./includes", item.ContentPath))
		if err != nil {
			log.Printf("Failed to load include '%s' for config: '%s'", item.ContentPath, item.ConfigName)
			continue
		}

		// Header has dynamic replacements, will load each time to avoid mutable state
		header, _ := os.ReadFile("./template/header.html")
		header = bytes.Replace(header, titlePlaceholder, []byte(item.Title), -1)

		html := markdown.ToHTML(content, nil, nil)

		parts := strings.Split(item.ConfigName, ".")
		outputName := parts[:len(parts)-1][0]

		// create directory for each config so css file can be copied in
		err = os.Mkdir("./build/"+outputName, 0755)
		if err != nil {
			return
		}

		outputHtml := fmt.Sprintf("./build/%s/index.html", outputName)
		outputCss := fmt.Sprintf("./build/%s/style.css", outputName)
		outputContent := [][]byte{header, html, footer}
		outputCssContent := [][]byte{normalize, templateCss}

		err = os.WriteFile(outputHtml, bytes.Join(outputContent, nil), 0755)
		err = os.WriteFile(outputCss, bytes.Join(outputCssContent, nil), 0755)

		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Successfully built config %s\n", item.ConfigName)
	}
}
