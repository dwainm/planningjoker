package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/machinebox/graphql"
	"golang.org/x/oauth2"
)

var client *graphql.Client

type project struct{
	Id string
	Title string
}

type issue struct {
	Id string
	Title string
	Estimate int 
}

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGithubToken()},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client = graphql.NewClient("https://api.github.com/graphql", graphql.WithHTTPClient(httpClient))

	client.Log = func(s string) { log.Println(s) }
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
		// Render index - start with views directory
        return c.Render("project", fiber.Map{
			"Title":"Project:" + project.Title ,
			"Projects": getProjects(),
        })
	})
	log.Fatal(app.Listen(":3000"))
}

func getGithubToken() string {
	return os.Getenv("PLANNING_JOKER_GH_TOKEN")
}

func getProjects() map[string]*project {

	// make a request
	req := graphql.NewRequest(`
	query($username: String!){
		user(login: $username) {
			projectsV2(first: 20) {
				nodes {
					id
					title
				}
			}
		}
	}
	`)

	// set any variables
	req.Var("username", "dwainm")

	// set header fields
	req.Header.Set("Cache-Control", "no-cache")

	// run it and capture the response
	var respData struct {
		User struct {
			ProjectsV2 struct{
				Nodes [] struct{
					Id string
					Title string
				}
			}
		} 
	}
	err := client.Run(context.Background(), req, &respData)
	if  err != nil {
		return make(map[string]*project)
	}

	projects := make(map[string]*project)

	for i, p := range respData.User.ProjectsV2.Nodes {
		projects[respData.User.ProjectsV2.Nodes[i].Id] = &project{Id:p.Id, Title: p.Title}
	}

	return projects
}

func getProject( Id string ) (*project, error){
	projects := getProjects()
	proj, ok := projects[Id];
	if ok {
		return proj, nil;
	} else {
		return nil, errors.New("Project not found")
	}
}

// func getIssuesFor(project) *issue {
		// var issueQuery struct {
		// Node struct {
		// 	ProjectsV2 struct{
		// 		Nodes [] struct{
		// 			Id string
		// 			Title string
		// 		}
		// 	}  `graphql:"... on ProjectV2"`
		// }  `graphql:"node(id: \"PVT_kwHOABolQs2J4g\")"`
	// }

	// variables := map[string]interface{}{
	// 	"owner": githubv4.String(owner),
	// 	"name":  githubv4.String(name),
	// }

  // query{
  //   node(id: "PROJECT_ID") {
  //       ... on ProjectV2 {
  //         items(first: 200) {
  //           nodes{
  //             id
  //             fieldValues(first: 8) {
  //               nodes{                
  //                 ... on ProjectV2ItemFieldTextValue {
  //                   text
  //                   field {
  //                     ... on ProjectV2FieldCommon {
  //                       name
  //                     }
  //                   }
  //                 }
  //                 ... on ProjectV2ItemFieldDateValue {
  //                   date
  //                   field {
  //                     ... on ProjectV2FieldCommon {
  //                       name
  //                     }
  //                   }
  //                 }
  //                 ... on ProjectV2ItemFieldSingleSelectValue {
  //                   name
  //                   field {
  //                     ... on ProjectV2FieldCommon {
  //                       name
  //                     }
  //                   }
  //                 }
  //               }              
  //             }
  //             content{              
  //               ... on DraftIssue {
  //                 title
  //                 body
  //               }
  //               ...on Issue {
  //                 title
  //                 assignees(first: 10) {
  //                   nodes{
  //                     login
  //                   }
  //                 }
  //               }
  //             }
  //           }
  //         }
  //       }
  //     }
  //   }
// }
