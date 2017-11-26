package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/cdipaolo/sentiment"
	"github.com/andream16/go-amazon-itemlookup/model"
	"errors"
	"strconv"
	"strings"
	"fmt"
)

var monthAsStringToMonthAsNumber = map[string]string{
						"January":"01", "February":"02", "March":"03", "April":"04", "May":"05",
						"June":"06", "July":"07", "August":"08", "September":"09", "October":"10",
						"November":"11", "December":"12",
}

func ReviewsScraper(url string, it int) ([]model.Review, error) {
	var reviews []model.Review
	sentimentModel, err := sentiment.Restore()
	if err != nil {
		return reviews, errors.New("Unable to initialize Model!")
	}
	reviewsUrl, reviewsUrlError := findReviewsUrl(url, 0)
	if reviewsUrlError != nil {
		return reviews, reviewsUrlError
	}
	fmt.Println(fmt.Sprintf("Successfully got review url=%s. Now counting reviews pages . . .", reviewsUrl))
	reviewsPagesNumber, reviewsPagesNumberError := countReviewsPages(reviewsUrl, 0)
	if reviewsPagesNumberError != nil || reviewsPagesNumber == 0 {
		return reviews, reviewsPagesNumberError
	}
	fmt.Println(fmt.Sprintf("Found %d pages! Now starting to scrape each page . . .", reviewsPagesNumber)); reviewsUrl += "&pageNumber=1"; c := make(chan []model.Review); index := 0
	for page := 1; page <= reviewsPagesNumber; page++ {
		if page > 1 {
			reviewsUrl = strings.Replace(reviewsUrl, "pageNumber=" + strconv.Itoa(page-1), "pageNumber=" + strconv.Itoa(page), 1)
		}
		fmt.Println(fmt.Sprintf("Scraping page=%d . . .", page))
		go scrapeReviewsByPage(reviewsUrl, &sentimentModel, c, page, 0)
	}
	for {
		currReviews := <- c
		fmt.Println("Got response number " + strconv.Itoa(index))
		reviews = append(reviews, currReviews...)
		if index == reviewsPagesNumber - 1 {
			fmt.Println("Returning . . .")
			return reviews, nil
		}
		index++
	}
	if len(reviews) > 0 {
		close(c)
		return reviews, nil
	}
	close(c)
	return reviews, errors.New("No Reviews found for that item!")
}

// Find Reviews Url given the url of the main reviews page. Eventually returns an error if url is not found.
func findReviewsUrl(url string, it int) (string, error) {
	fmt.Println("Looking for main reviews Url . . .")
	var reviewsUrl string
	var ok = false
	baseDoc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to get reviewUrl for url %s", url))
		return "", err
	}
	baseDoc.Find(".crIFrame .crIframeReviewList .small b a").Each(func(i int, s *goquery.Selection) {
		reviewsUrl, ok = s.Attr("href")
	})
	if ok && len(reviewsUrl) > 0 {
		fmt.Println(fmt.Sprintf("Main review url found at iteration=%d", it))
		return reviewsUrl, nil
	} else if !ok && it<=10 {
		it++
		fmt.Println("No main reviews url found! Trying again . .")
		findReviewsUrl(url, it)
	} else {
		fmt.Println(fmt.Sprintf("No main review url found for url %s. Returning.", url))
	}
	return "", errors.New("Unable to find main reviews url!")
}

// Finds the total number of review pages for an item. Eventually returns an error if no pages are found.
func countReviewsPages(reviewsUrl string, it int) (int, error) {
	reviewsDoc, reviewsDocErr := goquery.NewDocument(reviewsUrl)
	if reviewsDocErr != nil {
		fmt.Println(fmt.Sprintf("Unable to open document to count pages. Error: %s", reviewsDocErr.Error()))
		return 0, reviewsDocErr
	}
	reviewsPagesNumber, e := strconv.Atoi(reviewsDoc.Find("ul.a-pagination").Children().Last().Prev().Text())
	if e != nil || reviewsPagesNumber == 0 {
		if it <= 10 {
			fmt.Println(fmt.Sprintf("No number of review pages found at iteration=%d! Trying again. . .", it))
			it++
			countReviewsPages(reviewsUrl, it)
		} else {
			fmt.Println(fmt.Sprintf("No number of review pages found for url=%s. Returning.", reviewsUrl))
			return 0, errors.New("No Number of reviews pages found!")
		}
	}
	return reviewsPagesNumber, nil
}

// Finds all reviews in a page, parses them and adds them into a slice of reviews.
// Eventually returns an error if no reviews are found.
func scrapeReviewsByPage (pageUrl string, sentimentModel *sentiment.Models, c chan []model.Review, page int, it int) {
	var reviews []model.Review
	reviewsDoc, reviewsDocError := goquery.NewDocument(pageUrl)
	if reviewsDocError != nil {
		fmt.Println(fmt.Sprintf("Unable to open document for pageUrl=%s at iteration %d for page=%d. Trying again . . .", pageUrl, it, page))
		it++
		go scrapeReviewsByPage(pageUrl, sentimentModel, c, page, it)
	} else {
		reviewsDoc.Find("div .a-section .review").Each(func(i int, reviewsSelection *goquery.Selection) {
			var r = model.Review{}
			stars, hasStars := reviewsSelection.Find("div .a-section .celwidget .a-row .a-link-normal").Attr("title")
			if hasStars {
				r.Stars, _ = strconv.ParseFloat(strings.TrimSpace(strings.Split(stars, " ")[0]), 64)
			}
			r.Content = reviewsSelection.Find("div .a-section .celwidget .review-data .review-text").Text()
			date := strings.Split(strings.Replace(reviewsSelection.Find("div .a-section .celwidget .review-date").Text(), ",", "", -1), " ")
			r.Date = date[3] + "-" + monthAsStringToMonthAsNumber[date[1]] + "-" + date[2]
			r.Sentiment = sentimentModel.SentimentAnalysis(r.Content, sentiment.English).Score
			fmt.Println("Appending new review to reviews . . .")
			reviews = append(reviews, r)
		})
		if len(reviews) == 0 && it <= 10 {
			fmt.Println(fmt.Sprintf("No reviews found for pageUrl=%s page=%d iteration=%d", pageUrl, page, it))
			it++
			go scrapeReviewsByPage(pageUrl, sentimentModel, c, page, it)
		} else {
			fmt.Println(fmt.Sprintf("Returning %d reviews for pageUrl=%s", len(reviews), pageUrl))
			c <- reviews
		}
	}
}