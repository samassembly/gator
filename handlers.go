package main

import (
	"fmt"
	"time"
	"context"
	"github.com/samassembly/gator/internal/database"
	"github.com/samassembly/gator/internal/rss"
	"github.com/google/uuid"

)

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

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Errorf("Failed to fetch feed: %v", err)
	}
	fmt.Printf("%v", feed)
	return nil
}

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

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get feeds from database: %v", err)
	}
	fmt.Printf("%v", feeds)
	return nil
}

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

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to remove users: %v\n", err)
	}
	return nil
}