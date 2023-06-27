package helpers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

func GetImageNameFromUrl(url string) string {
	segments := strings.Split(url, "/")
	filename := segments[len(segments)-1]
	filename = strings.Split(filename, "?")[0]
	return filename
}

func DownloadImage(url string, filename string) {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
}

func ReplaceRemoteImagesWithLocal(base string, filename string) {
	filepath := path.Join(base, filename)
	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

	r := regexp.MustCompile(`!\[.*?]\((.+?)\)`)
	matches := r.FindAllStringSubmatch(string(file), -1)

	// mkdir if not exists
	p := path.Join(base, "images")

	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := os.MkdirAll(p, 0755); err != nil {
			log.Fatal(err)
		}
	}

	for _, v := range matches {
		imageName := GetImageNameFromUrl(v[1])
		DownloadImage(v[1], path.Join(p, imageName))
		//fmt.Println(imageName)
		file = []byte(strings.ReplaceAll(string(file), v[1], fmt.Sprintf("./images/%s", imageName)))
	}
	if err := os.WriteFile(filepath, file, 0o644); err != nil {
		log.Fatal(err)
	}
}
