package main

import (
	"fmt"
	"time"
	"context"
	"log"
	"github.com/samassembly/gator/internal/database"
	"github.com/samassembly/gator/internal/rss"
	"github.com/google/uuid"
)

//set a registered db user as the current use in the configuration
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("User not registered: %v", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

//register a new user in the db
func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()
	name := cmd.Args[0]

	create_args := database.CreateUserParams{
		ID: id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name: name,
	}

	user, err := s.db.CreateUser(context.Background(), create_args)
	if err != nil {
		return fmt.Errorf("Failed to create user: %v\n", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %v\n", err)
	}

	fmt.Printf("User created Successfully: %v\n", user)
	return nil
}

//retrieve users in the db
func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to retrieve users: %v\n", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

//fetch the feed for a given url
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

//add a specified feed to the db 
func handlerAddFeed(s *state, cmd command, user database.User) error {
	userid := user.ID

	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()
	name := cmd.Args[0]
	url := cmd.Args[1]

	create_args := database.CreateFeedParams{
		ID: id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		Name: name,
		Url: url,
		UserID: userid, 
	}

	feed, err := s.db.CreateFeed(context.Background(), create_args)
	if err != nil {
		return fmt.Errorf("Failed to add feed to database: %v\n", err)
	}

	feedid := id
	id = uuid.New()
	created_at = time.Now()
	updated_at = time.Now()
	follow_args := database.CreateFeedFollowParams{
		ID: id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		UserID: userid,
		FeedID: feedid,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), follow_args)
	if err != nil {
		return fmt.Errorf("Could not create feed_follow: %v", err)
	}

	fmt.Printf("Feed added to Database: %v\n", feed)
	return nil
}

//return feeds in database
func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds from database: %v", err)
	}
	fmt.Printf("%v", feeds)
	return nil
}

//create an entry in the feed_follow table for the current user given a url
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Could not retrieve feed: %v", err)
	}

	id := uuid.New()
	created_at := time.Now()
	updated_at := time.Now()
	userid := user.ID
	feedid := feed.ID

	follow_args := database.CreateFeedFollowParams{
		ID: id,
		CreatedAt: created_at,
		UpdatedAt: updated_at,
		UserID: userid,
		FeedID: feedid,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), follow_args)
	if err != nil {
		return fmt.Errorf("Could not create feed_follow: %v", err)
	}
	fmt.Printf("%s is now following %s", s.cfg.CurrentUserName, feed.Name)
	return nil
}

//return the names of followed feeds for the current user
func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	
	userid := user.ID
	followed_feeds, err := s.db.GetFeedFollowsForUser(context.Background(), userid)
	if err != nil {
		return fmt.Errorf("Failed to get followed feeds: %v\n", err)
	}

	for _, followed_feed := range followed_feeds {
		fmt.Printf("- %s\n", followed_feed.FeedName)
	}
	return nil
}

//unfollow a specified feed for the current user
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}
	url := cmd.Args[0]
	userid := user.ID
	feed, err := s.db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error retrieving feed: %v\n", err)
	}
	feedid := feed.ID

	unfollow_args := database.UnfollowParams{
		UserID: userid,
		FeedID: feedid, 
	}
	err = s.db.Unfollow(context.Background(), unfollow_args)
	if err != nil {
		return fmt.Errorf("Could not unfollow: %v", err)
	}
	return nil
}

//reset the database
func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to remove users: %v\n", err)
	}
	return nil
}