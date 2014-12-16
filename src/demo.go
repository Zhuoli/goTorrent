package src

import (
	"fmt"
//	"io"
	"net/http"
	"os"
	"strings"
)

func downloadFromUrl(url string) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	fmt.Println("Body length:")
	fmt.Println(response.ContentLength)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
//
//	n, err := io.Copy(output, response.Body)
//	if err != nil {
//		fmt.Println("Error while downloading", url, "-", err)
//		return
//	}
//
//	fmt.Println(n, "bytes downloaded.")
}

func main() {
	url := "http://download.nextag.com/apache/maven/maven-3/3.2.3/binaries/apache-maven-3.2.3-bin.tar.gz"
	downloadFromUrl(url)
}