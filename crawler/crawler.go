package crawler

import (
	"fmt"
	"github.com/andream16/go-amazon-itemlookup/scraper"
	"github.com/andream16/go-amazon-itemlookup/configuration"
	"github.com/andream16/go-amazon-itemlookup/model"
	"github.com/andream16/go-amazon-itemlookup/request"
	"github.com/andream16/go-amazon-itemlookup/querybuilder"
)

func Crawl(config configuration.Configuration) {
	fmt.Println("Getting items from api . . .")
	items, itemsError := request.GetItems(config); if itemsError != nil {
		panic(itemsError)
	}
	fmt.Println(fmt.Sprintf("Successfully got %d items. Now starting to range into them . . .", len(items.Items)))
	for _, item := range items.Items {
		var amazonEntry model.Amazon
		q, e := querybuilder.BuildRequest(config, item.Item)
		if e == nil {
			q.MainCategories = append(q.MainCategories, q.ChildrenCategories...)
			if len(q.MainCategories) > 0 {
				fmt.Println("Adding amazonEntry.Categories . . .")
				amazonEntry.Categories = model.Categories{Item: item.Item, Categories: q.MainCategories}
			}
			if len(q.ReviewsUrl) > 0 && q.HasReviews {
				fmt.Println("Starting top scrape for reviews . . .")
				reviews, reviewsError := scraper.ReviewsScraper(q.ReviewsUrl, 0)
				if reviewsError != nil {
					return
				} else {
					fmt.Println(fmt.Sprintf("Successfully got %d reviews for item %s!", len(reviews), item.Item))
					amazonEntry.Reviews.Item = item.Item
					for _, review := range reviews {
						review.Item = item.Item
						amazonEntry.Reviews.Reviews = append(amazonEntry.Reviews.Reviews,  review)
					}
				}
			}
			if len(item.Item) > 0 && len(q.Manufacturer) > 0 {
				fmt.Println("Adding amazonEntry.Item . . .")
				amazonEntry.Item = model.Item{Item : item.Item, Manufacturer: q.Manufacturer, Title: q.Title, Image: q.Image, URL: q.URL}
				fmt.Println("Adding amazonEntry.Manufacturer . . .")
				amazonEntry.Manufacturer = model.Manufacturer{Manufacturer: q.Manufacturer}
				request.AddAmazonEntry(config, amazonEntry)
			}
		}
	}
	fmt.Println("Crawling has finished.")
}