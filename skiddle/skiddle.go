package skiddle

import ("fmt"
)
const (
    BaseURL = "https://www.skiddle.com/api/v1"
    APIKey  = "3da8b2c4d34b51e8a45e86cc08280ad9"
)

func EventDetailsURL(eventID uint) string {
    return fmt.Sprintf("%s/events/%d/?api_key=%s", BaseURL, eventID, APIKey)
}

func EventSearchURL(params map[string]string) string {
    url := fmt.Sprintf("%s/events/search/?api_key=%s", BaseURL, APIKey)
    for key, value := range params {
        url += fmt.Sprintf("&%s=%s", key, value)
    }
    return url
}