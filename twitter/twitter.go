package twitter

import (
    "context"
    "fmt"
    "log"
    "time"

    twitterscraper "github.com/ThallesP/twitter-scraper-openaccount"
)

func Twitterscrapering(eventName string) ([]map[string]interface{}, error) {
    // Create a new scraper instance
    scraper := twitterscraper.New()
    scraper.WithDelay(5) // Add 5 second delay between requests

    // Login to Twitter account
    err := scraper.Login("eventconnect22", "Universe12%")
    if err != nil {
        log.Printf("Error logging into Twitter: %v", err)
        return nil, fmt.Errorf("error logging into Twitter: %w", err)
    }

    // Search for tweets containing specific keywords related to the event
    keywords := eventName
    var tweets []map[string]interface{}
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    for tweet := range scraper.SearchTweets(ctx, keywords, 50) {
        if tweet.Error != nil {
            log.Printf("Error fetching tweet: %v", tweet.Error)
            continue
        }
        tweetData := map[string]interface{}{
            "text": tweet.Text,
            "user": map[string]interface{}{
                "name":        tweet.Username,
                "screen_name": tweet.UserID,
            },
        }
        tweets = append(tweets, tweetData)
    }

    if len(tweets) == 0 {
        log.Printf("No tweets found for event: %s", eventName)
        return nil, fmt.Errorf("no tweets found for event: %s", eventName)
    }

    return tweets, nil
}