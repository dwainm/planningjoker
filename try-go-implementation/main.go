//go:generate go run github.com/Khan/genqlient genqlient.yaml

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"golang.org/x/oauth2"
	"golang.org/x/sys/windows"
)

var client graphql.Client

type project struct{
	Id string
	Title string
}

type issue struct {
	Id string
	Title string
	Estimate int 
}

var username = "dwainm"
func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGithubToken()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client = graphql.NewClient("https://api.github.com/graphql", httpClient)

	// client.Log = func(s string) { log.Println(s) }
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

	app.Get("/project/:id", func (c *fiber.Ctx) error {
		project, err := getProject(c.Params("id"))
		if err != nil {
			log.Println(err)
			return err
		}
		// // Render index - start with views directory
        return c.Render("project", fiber.Map{
			"Title":"Project:" + project.Title ,
			"Issues": getIssuesFor(project),
        })
	})
	log.Fatal(app.Listen(":3000"))
}

func getGithubToken() string {
	return os.Getenv("PLANNING_JOKER_GH_TOKEN")
}

func getProjects() map[string]*project{
	projects, err := apiGetProjects(context.Background(), client, username);
	if( err != nil){
		log.Print(err)	
	}
	tranformedProjects := make(map[string]*project)
	for i, p := range projects.User.ProjectsV2.Nodes {
		tranformedProjects[projects.User.ProjectsV2.Nodes[i].Id] = &project{Id:p.Id, Title: p.Title}
	}
	return tranformedProjects
}

func getProject( Id string ) (*project, error) {
	projects := getProjects()
	proj, ok := projects[Id];
	if ok {
		return proj, nil;
	} else {
		return nil, errors.New("Project not found")
	}
}

func getIssuesFor(proj *project) map[string]issue {
	issues, err := apiGetProjectIssues(context.Background(), client, proj.Id)
	fields, err := apiGetAllFields(context.Background(), client)
	log.Println("== Isues data:")
	log.Println(*issues)
	log.Println(issues)
	log.Println(err)
	log.Println("== Field data:")
	log.Println(fields.GetNode())

	// // make a request
	// req := graphql.NewRequest(`
	// `)

	// // set any variables
	// req.Var("project_id", proj.Id)

	// // set header fields
	// req.Header.Set("Cache-Control", "no-cache")

	// // run it and capture the response
	// var respData struct {
	// 	data struct {
	// 		node struct {
	// 			items struct{
	// 			Nodes [] struct{
	// 				Id string
	// 				FieldValues struct{
	// 				Nodes [] struct{
	// 				Title string
	// 			}
	// 		}
	// 	} 
	// }
	// err := client.Run(context.Background(), req, &respData)
	// if  err != nil {
	// 	return make(map[string]*project)
	// }

	// projects := make(map[string]*project)

	// for i, p := range respData.User.ProjectsV2.Nodes {
	// 	projects[respData.User.ProjectsV2.Nodes[i].Id] = &project{Id:p.Id, Title: p.Title}
	// }

	// return projects
	// // variables := map[string]interface{}{
	// // 	"owner": githubv4.String(owner),
	// // 	"name":  githubv4.String(name),
	// // }
	return make(map[string]issue)

}
