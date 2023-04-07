package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var client *githubv4.Client

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGithubToken()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client = githubv4.NewClient(httpClient)

	// Views engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func (c *fiber.Ctx) error {
		// Render index - start with views directory
        return c.Render("index", fiber.Map{
            "Title": fmt.Sprintf(" Met structs! Dag se, %s!", "USER"),
			"Projects": getProjects(),
        })
	})

	log.Fatal(app.Listen(":3000"))
}

func getGithubToken() string {
	return os.Getenv("PLANNING_JOKER_GH_TOKEN")
}

func getProjects() []struct { Id string; Title string} {
		var projectQuery struct {
		User struct {
			ProjectsV2 struct{
				Nodes [] struct{
					Id string
					Title string
				}
			}  `graphql:"projectsV2(first: 20)"`
		}  `graphql:"user(login: \"dwainm\")"`
	}

	// variables := map[string]interface{}{
	// 	"owner": githubv4.String(owner),
	// 	"name":  githubv4.String(name),
	// }

	err := client.Query(context.Background(), &projectQuery, nil)
	if err != nil {
		log.Println(err)
	}

	return projectQuery.User.ProjectsV2.Nodes
}
