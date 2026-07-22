package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/AbdullahBasir/aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("command must include duration time")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("could not parse duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for range ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			log.Printf("could not scrape feeds: %v", err)
		}
	}
	return nil
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could not fetch next feed: %w", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}

	fmt.Printf("\n--- %s ---\n", feed.Name)
	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		} else {
			log.Printf("could not parse date: %v", err)
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: true},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Print("duplicate key value violates unique constraint")
			log.Printf("couldn't create post: %v", err)
			continue
		}
	}
	return nil
}
