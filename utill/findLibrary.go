package utill

import (
	"io/ioutil"
	"library/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FindLibrary(c *gin.Context) []byte {

	url := config.AppConfig.BaseURL + ":" + config.AppConfig.Port + "/library"

	req, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	return body
}
