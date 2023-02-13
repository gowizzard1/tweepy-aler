package main

import (
	"context"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"github.com/twitter-mention-notifier/emailNotification"
	"golang.org/x/time/rate"
	"log"
	"net/url"
	"os"
	"sync"
	"time"
)

const keyword = "golang"
const numWorkers = 5

type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	SMTPServer     string
	SMTPPort       string
	SMTPEmail      string
	SMTPPassword   string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPassword := os.Getenv("SMTP_PORT")
	config := &emailNotification.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		AccessToken:    accessToken,
		AccessSecret:   accessSecret,
		SMTPServer:     smtpServer,
		SMTPPort:       smtpPort,
		SMTPEmail:      smtpEmail,
		SMTPPassword:   smtpPassword,
	}

	// rate limiter to avoid exceeding Twitter API rate limits
	limiter := rate.NewLimiter(rate.Every(time.Minute), 1)

	// Twitter API client
	//httpClient := oauth1.NewConfig(config.ConsumerKey, config.ConsumerSecret).Client(oauth1.NoContext, oauth1.NewToken(config.AccessToken, config.AccessSecret))
	client := anaconda.NewTwitterApi(config.AccessToken, config.AccessSecret)

	// Stream tweets containing keyword
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := client.PublicStreamFilter(url.Values{"track": []string{keyword}})

	// Keep track of processed tweets to avoid duplicates
	processedTweets := make(map[int64]struct{})
	cache := make(map[int64]time.Time)

	// Create a channel for incoming tweets and a wait group for worker goroutines
	tweetChan := make(chan *anaconda.Tweet, 100)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, tweetChan, config, limiter)
	}

	// Stream tweets and add them to the channel
	go func() {
		defer close(tweetChan)
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-stream.C:
				status := v.(anaconda.Tweet)
				if _, processed := processedTweets[status.Id]; !processed {
					tweetChan <- &status
					processedTweets[status.Id] = struct{}{}
					cache[status.Id] = time.Now()
				}
			}
		}
	}()

	// Wait for all worker goroutines to finish
	wg.Wait()
	stream.Stop()
}

func worker(id int, wg *sync.WaitGroup, tweetChan chan *anaconda.Tweet, config *emailNotification.Config, limiter *rate.Limiter) {
	defer wg.Done()
	for tweet := range tweetChan {
		if limiter.Allow() {
			email := &emailNotification.Email{
				SenderID: "Twitter API",
				ToIDs:    []string{"recipient@example.com"},
				Subject:  "New tweet containing keyword",
				Message:  fmt.Sprintf("A new tweet with keyword '%s' was just posted: %s", keyword, tweet.Text),
			}
			if err := emailNotification.SendEmail(config, email); err != nil {
				log.Printf("Worker %d: Failed to send email notification: %s", id, err)
			} else {
				log.Printf("Worker %d: Sent email notification for tweet: %s", id, tweet.Text)
			}
		}
	}
}
