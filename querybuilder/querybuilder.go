package querybuilder

import (
	"github.com/andream16/go-amazon-itemlookup/configuration"
	"github.com/andream16/go-amazon-itemlookup/model"
	"github.com/andream16/go-amazon-itemlookup/util"
	"time"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/xml"
)

type Query struct {
	xml.Name
	Title			string	  `xml:"Items>Item>ItemAttributes>Title"`
	URL             string    `xml:"Items>Item>DetailPageURL"`
	Image			string 	  `xml:"Items>Item>LargeImage>URL"`
	Manufacturer    string    `xml:"Items>Item>ItemAttributes>Manufacturer"`
	MainCategories []string `xml:"Items>Item>BrowseNodes>BrowseNode>Name"`
	ChildrenCategories []string `xml:"Items>Item>BrowseNodes>BrowseNode>Children>BrowseNode>Name"`
	ReviewsUrl string `xml:"Items>Item>CustomerReviews>IFrameURL"`
	HasReviews bool `xml:"Items>Item>CustomerReviews>HasReviews"`
}

func NewQueryModel(configuration configuration.Configuration, itemId string, timestamp time.Time) string {
	lightQueryModel := model.LightQueryModel{}
	lightQueryModel.ItemId = itemId
	lightQueryModel.AssociateTag = configuration.Credentials.AssociateTag
	lightQueryModel.AWSAccessKeyId = configuration.Credentials.AWSAccessKeyId
	lightQueryModel.ResponseGroup = configuration.Remote.ResponseGroup
	lightQueryModel.Operation = configuration.Remote.Operation
	lightQueryModel.Timestamp = timestamp
	return util.BuildQuery(lightQueryModel, configuration, configuration.Credentials.AWSSecretKey)
}

func BuildRequest(configuration configuration.Configuration, itemPid string) (Query, error) {
	fmt.Println(fmt.Sprintf("Starting to build request for item %s . . .", itemPid))
	queryModel := NewQueryModel(configuration, itemPid, time.Now())
	fmt.Println(fmt.Sprintf("Got queryModel %s. Now retrieving the item . . .", queryModel))
	item, err := retrieveItem(queryModel)
	fmt.Println(fmt.Sprintf("Successfully retrieved item=%s response=%s. Now unmarshalling response into query model . . .", itemPid, item))
	var q Query
	xml.Unmarshal([]byte(item), &q)
	if err != nil {
		fmt.Println(fmt.Sprintf("Got error while unmarshalling response into query model for item %s. Error: %s", itemPid, err.Error()))
		return q, err
	}
	fmt.Println(fmt.Sprintf("Successfully unmarshalled response into query model for item %s", itemPid))
	return q, nil
}

func retrieveItem(url string) (string, error) {
	response, responseError := http.Get(url)
	if responseError != nil {
		fmt.Println(fmt.Sprintf("Unable to get data for request %s, got error %s", url, responseError))
		return "", fmt.Errorf("GET error: %v", responseError)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("Unable to get data for request %s, got status response %d", url, response.StatusCode))
		return "", fmt.Errorf("Status error: %v", response.StatusCode)
	}
	data, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		fmt.Println(fmt.Sprintf("Unable to get data for request %s, Error while reading response body %s", url, readErr))
		return "", fmt.Errorf("Read body: %v", readErr)
	}
	return string(data), nil
}