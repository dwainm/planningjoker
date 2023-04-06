package main

import (
	"context"
	"log"
	"os"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGithubToken()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	// query{
    // organization(login: "ORGANIZATION"){
    //   projectV2(number: NUMBER) {
    //     id
    //   }
    // }
  // }'
	var query struct {
		Viewer struct{
			Login string
			WebsiteURL string
		}
	}

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Println(err)
	}

	// Views engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func (c *fiber.Ctx) error {
		// Render index - start with views directory
        return c.Render("index", fiber.Map{
            "Title": fmt.Sprintf("Hello, %s!", query.Viewer.Login),
        })
	})

	log.Fatal(app.Listen(":3000"))
}

func getGithubToken() string {
	return os.Getenv("PLANNING_JOKER_GH_TOKEN")
}
