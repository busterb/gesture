// Google Image Search functionality
package gis

import (
	"encoding/json"
	"errors"
	"gesture/plugin"
	"math/rand"
	"net/url"
	"strings"
	"gesture/util"
)

// lol types
type GisPlugin bool

// lol types
func NewPlugin() GisPlugin {
	return GisPlugin(false)
}

func (gis GisPlugin) Call(mc plugin.MessageContext) (bool, error) {
	if mc.Command() == "gis" {
		if len(mc.CommandArgs()) > 0 {
			link, err := search(strings.Join(mc.CommandArgs(), " "))
			if err != nil {
				return false, err
			} else {
				mc.Ftfy(link)
			}
		}
	}
	return false, nil
}

// these structs really tie the room together, man
type gisResult struct {
	Url string
}
type gisResponse struct {
	ResponseData struct {
		Results []gisResult
	}
	Results []gisResult
}

// Search queries google for some images, and then randomly selects one
func search(search string) (result string, err error) {
	url := "http://ajax.googleapis.com/ajax/services/search/images?v=1.0&q=" + url.QueryEscape(search)
	body, err := util.GetUrl(url)
	if err != nil {
		return "", err
	}
	var gisResponse gisResponse
	json.Unmarshal(body, &gisResponse)
	if len(gisResponse.ResponseData.Results) > 0 {
		indexes := rand.Perm(len(gisResponse.ResponseData.Results))
		for i := 0; i < len(indexes); i++ {
			imageUrl := gisResponse.ResponseData.Results[indexes[i]].Url
			if isImage(imageUrl) {
				return imageUrl, nil
			}
		}
	}
	return "", errors.New("No image could be found for " + search)
}

// returns true if the url ends with some well known suffixes
func isImage(url string) bool {
	suffixes := []string{".jpg", ".jpeg", ".gif", ".png", ".bmp"}
	lowered := strings.ToLower(url)
	for _, suffix := range suffixes {
		if strings.HasSuffix(lowered, suffix) {
			return true
		}
	}
	return false
}
