package request

import (
	"github.com/andream16/go-amazon-itemlookup/configuration"
	"github.com/andream16/go-amazon-itemlookup/model"
	"github.com/andream16/go-amazon-itemlookup/util"
	"net/http"
	"strings"
	"encoding/json"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
)

type Q struct {
	Page uint `json:"page"`
	Size uint `json:"size"`
}

// Adds an Amazon Entry given an Item
func AddAmazonEntry(config configuration.Configuration, amazonRequest model.Amazon) error {
	fmt.Println(fmt.Sprintf("Posting new amazon entry for item %s . . .", amazonRequest.Item.Item))
	body, bodyError := json.Marshal(amazonRequest); if bodyError != nil {
		fmt.Println(fmt.Sprintf("Unable to marshal new amazon entry for item %s, got error: %s", amazonRequest.Item.Item, bodyError.Error()))
		return bodyError
	}
	response, requestErr := http.Post("http://" + createEndPointURL(config, config.Api.Endpoints.Amazon), "application/json", bytes.NewBuffer(body)); if requestErr != nil {
		fmt.Println(fmt.Sprintf("Unable to post new amazon entry for item %s, got error: %s", amazonRequest.Item.Item, requestErr.Error()))
		return requestErr
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Unable to post new amazon entry for item %s, got status code: %d", amazonRequest.Item.Item, response.StatusCode))
		return errors.New(fmt.Sprintf("Unable to post amazon entry for item %s", amazonRequest.Item.Item))
	}
	fmt.Println(fmt.Sprintf("Successfully posted new amazon entry for item %s. Returning.", amazonRequest.Item.Item))
	return nil
}

// Gets items given page and size
func GetItems(config configuration.Configuration) (model.Items, error) {
	baseUrl := appendQueryParams(createEndPointURL(config, config.Api.Endpoints.Item), Q{ 1, 157 })
	fmt.Println(fmt.Sprintf("Using %s url for request", baseUrl))
	response, responseError := http.Get("http://" + baseUrl); if responseError != nil {
		return model.Items{}, responseError
	}
	defer response.Body.Close()
	return unmarshalItems(response), nil
}

func createApiURL(config configuration.Configuration) string {
	return strings.Join([]string{config.Api.Host, config.Api.Port}, ":")
}

func createEndPointURL(config configuration.Configuration, endPoint string) string {
	return strings.Join([]string{createApiURL(config), endPoint}, "/")
}

func appendQueryParams(base string, queryParameters interface{}) string {
	return strings.Join([]string{base, strings.ToLower(util.QueryModelToQueryString(queryParameters))}, "?")
}

func unmarshalItems(r *http.Response) model.Items {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var items model.Items
	err = json.Unmarshal(body, &items)
	if err != nil {
		panic(err)
	}
	return items
}