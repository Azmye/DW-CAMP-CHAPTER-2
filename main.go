package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	template *template.Template
}

type Project struct {
	ProjectName string
	StartDate   string
	EndDate     string
	Description string
	TechIcon    map[string]string
	Image       string
}

var projectsData = []Project{
	{
		ProjectName: "Dumbways Way App",
		StartDate:   "12 Jan 2023",
		EndDate:     "15 Jan 2023",
		Description: "App Project that can make you're coding life easier, this app built with React, and NodeJs.",
		TechIcon: map[string]string{
			"Laravel": "",
			"Go":      "",
			"NodeJs":  "on",
			"ReactJs": "on",
		},

		Image: "https://source.unsplash.com/random/900*700?tech,programming",
	},
	{
		ProjectName: "Scheduler.IO",
		StartDate:   "5 Mar 2022",
		EndDate:     "15 Mar 2022",
		Description: "App Project that can make you're  life easier, this app built with React, and Golang.",
		TechIcon: map[string]string{
			"Laravel": "on",
			"Go":      "on",
			"NodeJs":  "",
			"ReactJs": "",
		},
		Image: "https://source.unsplash.com/random/900*700?games,football",
	},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Static("/public", "public")

	t := &Template{
		template: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	e.GET("/", home)
	e.GET("/contact-form", contactForm)
	e.GET("/project-form", projectForm)
	e.POST("/project-add", projectAdd)
	e.POST("/project-edit", projectEdit)
	e.GET("/project-edit-form/:id", projectEditForm)
	e.GET("/project-detail/:id", projectDetail)
	e.GET("/project-delete/:id", projectDelete)

	fmt.Println("Server Berlajalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	projects := map[string]interface{}{
		"Projects": projectsData,
	}

	return c.Render(http.StatusOK, "index.html", projects)
}

func contactForm(c echo.Context) error {
	return c.Render(http.StatusOK, "contact-form.html", nil)
}

func projectForm(c echo.Context) error {

	var Projects = Project{
		ProjectName: "",
		StartDate:   "",
		EndDate:     "",
		Description: "",
	}

	sendDatas := map[string]interface{}{
		"Project": Projects,
		"Button":  `<button type="submit" class="btn btn-dark rounded-5 px-4 py-1">Submit</button>`,
		"Action":  "/project-add",
	}

	return c.Render(http.StatusOK, "project-form.html", sendDatas)
}

func projectAdd(c echo.Context) error {
	projectName := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	description := c.FormValue("description")
	techIcon := map[string]string{
		"Laravel": c.FormValue("laravel"),
		"Go":      c.FormValue("go"),
		"NodeJs":  c.FormValue("nodeJs"),
		"ReactJs": c.FormValue("reactJs"),
	}
	image := " https://source.unsplash.com/random/900*700?programming, tech, game"

	println("Project Name : " + projectName)
	println("Start Date : " + startDate)
	println("End Date : " + endDate)
	println("Description :" + description)
	for k, v := range techIcon {
		fmt.Printf("%s: %s\n", k, v)
	}

	var newProject = Project{
		ProjectName: projectName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: description,
		TechIcon:    techIcon,
		Image:       image,
	}

	projectsData = append(projectsData, newProject)

	return c.Redirect(http.StatusMovedPermanently, "/")

}

func projectDetail(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDets = Project{}

	for idx, data := range projectsData {
		if id == idx {
			ProjectDets = Project{
				ProjectName: data.ProjectName,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
				TechIcon:    data.TechIcon,
				Image:       data.Image,
			}
		}
	}

	sendProjectDets := map[string]interface{}{
		"Project": ProjectDets,
	}

	return c.Render(http.StatusOK, "project-detail.html", sendProjectDets)
}

func projectDelete(c echo.Context) error {
	id, _ := strconv.Atoi("id")

	projectsData = append(projectsData[:id], projectsData[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func projectEditForm(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var Projects = Project{}

	for idx, data := range projectsData {
		if id == idx {
			Projects = Project{
				ProjectName: data.ProjectName,
				StartDate:   data.StartDate,
				EndDate:     data.EndDate,
				Description: data.Description,
				TechIcon:    data.TechIcon,
				Image:       data.Image,
			}
		}
	}

	sendDatas := map[string]interface{}{
		"Project": Projects,
		"Button":  `<button type="submit" class="btn btn-dark rounded-5 px-4 py-1">Edit</button>`,
		"Action":  "/project-edit",
	}

	return c.Render(http.StatusOK, "project-form.html", sendDatas)
}

func projectEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	projectsData[id].ProjectName = c.FormValue("projectName")
	projectsData[id].StartDate = c.FormValue("startDate")
	projectsData[id].EndDate = c.FormValue("endDate")
	projectsData[id].Description = c.FormValue("description")
	projectsData[id].TechIcon = map[string]string{
		"Laravel": c.FormValue("laravel"),
		"Go":      c.FormValue("go"),
		"NodeJs":  c.FormValue("nodeJs"),
		"ReactJs": c.FormValue("reactJs"),
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}
