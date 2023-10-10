package utils

import (
	"log"
	"net/url"
	constants "webscrapper/constants/v2"
)

func CreateQueryMapCopy(key string, school string) map[string]string {
	ret := map[string]string{}

	for k, v := range constants.ConstantLinks[school][key] {
		ret[k] = v
	}

	return ret
}

func FormatUrl(key string, school string) (string, error) {

	endpoint := constants.ConstantLinks[school]["endpoint"]["url"]
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return "", err
	}

	// Encode the POST data as a URL-encoded string
	postData := url.Values{}
	for key, val := range constants.ConstantLinks[school][key] {
		postData.Set(key, val)
	}
	postDataStr := postData.Encode()
	formattedURL := urlObj.String() + postDataStr
	return formattedURL, nil
}

func FormatDynamicUrl(data map[string]string, school string) (string, error) {

	endpoint := constants.ConstantLinks[school]["endpoint"]["url"]
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		log.Println("Error parsing URL:", err)
		return "", err
	}

	// Encode the POST data as a URL-encoded string
	postData := url.Values{}
	for key, val := range data {
		postData.Set(key, val)
	}
	postDataStr := postData.Encode()
	formattedURL := urlObj.String() + postDataStr
	return formattedURL, nil
}
