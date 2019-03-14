package main

import (
	"context"
	"fmt"
	"os"

	"text/template"

	"github.com/google/go-github/github"
	"github.com/idec-net/github2idec/g2i"
)

func main() {
	ctx := context.Background()
	config := loadConfig(filePath)
	gClient := g2i.NewGithubClient(config, ctx)
	eventsURL, err := gClient.EventsURL()
	if err != nil {
		panic(err)
	}

	fmt.Println(eventsURL)

	// events, err := gClient.GetEvents(eventsURL)
	// if err != nil {
	// 	panic(err)
	// }

	issues, err := gClient.GetIssues()
	if err != nil {
		panic(err)
	}

	helloMessage := struct {
		Sources, Name, Owner, Repo string
		Issues                     []*github.Issue
	}{
		Sources: "https://github.com/idec-net/github2idec",
		Name:    "Gdec",
		Owner:   config.Github.RepoOwner,
		Repo:    config.Github.Repo,
		Issues:  issues,
	}
	t, err := template.New("hello_message.tpl").
		ParseFiles("../templates/hello_message.tpl")
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, helloMessage)
	if err != nil {
		panic(err)
	}
	// for _, event := range events {
	// 	if config.IsEventProcessable(*event.Type) {
	// 		b, err := event.RawPayload.MarshalJSON()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		var comment g2i.GithubIssueComment
	// 		err = json.Unmarshal(b, &comment)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		fmt.Printf("===============\nAuthor: %s\n", comment.Comment.User.Login)
	// 		fmt.Printf("Issue: %s\nUrl: %s\n", comment.Issue.Title, comment.Comment.URL)
	// 		fmt.Printf("CreatedAt: %s\n", comment.Comment.CreatedAt.Format("2006 Jan 2 15:04:05 UTC"))
	// 		fmt.Printf("%s\n===============\n\n", comment.Comment.Body)
	// 	}
	// }
}
