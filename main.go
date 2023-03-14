package main

import (
	"Personal-website/connection"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	template *template.Template
}

type ProjectDb struct {
	ID           int
	ProjectName  string
	StartDate    time.Time
	EndDate      time.Time
	StartDateStr string
	EndDateStr   string
	Description  string
	TechIcon     []string
	Image        string
	Difference   string
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

func main() {
	connection.DatabaseConnect()
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
	data, _ := connection.Conn.Query(context.Background(), "SELECT * FROM tb_projects")

	var results []ProjectDb

	for data.Next() {
		var each = ProjectDb{}

		err := data.Scan(&each.ID, &each.ProjectName, &each.StartDate, &each.EndDate, &each.Description, &each.TechIcon, &each.Image)

		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
		}

		each.Difference = dateDifference(each.StartDate, each.EndDate)

		words := strings.Fields(each.Description)
		if len(words) > 18 {
			each.Description = strings.Join(words[:18], " ") + "..."
		}

		results = append(results, each)
	}
	projects := map[string]interface{}{
		"Projects": results,
	}

	return c.Render(http.StatusOK, "index.html", projects)
}

func contactForm(c echo.Context) error {
	return c.Render(http.StatusOK, "contact-form.html", nil)
}

func projectForm(c echo.Context) error {

	var Projects = ProjectDb{
		ProjectName:  "",
		StartDateStr: "",
		EndDateStr:   "",
		Description:  "",
	}

	sendDatas := map[string]interface{}{
		"Project": Projects,
		"Button":  `<button type="submit" class="btn btn-dark rounded-5 px-4 py-1">Submit</button>`,
		"Action":  "/project-add",
	}

	return c.Render(http.StatusOK, "project-form.html", sendDatas)
}

func projectAdd(c echo.Context) error {
	name := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	description := c.FormValue("description")
	techonologies := []string{}

	if c.FormValue("nodeJs") == "on" {
		techonologies = append(techonologies, "NodeJs")
	}
	if c.FormValue("reactJs") == "on" {
		techonologies = append(techonologies, "ReactJs")
	}
	if c.FormValue("javascript") == "on" {
		techonologies = append(techonologies, "Javascript")
	}
	if c.FormValue("go") == "on" {
		techonologies = append(techonologies, "Go")
	}

	startDateVal, _ := time.Parse("2006-01-02", startDate)
	endDateVal, _ := time.Parse("2006-01-02", endDate)

	image := "https://source.unsplash.com/random/900*700?games,football"

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_projects (name, start_date, end_date, description, technologies, image) VALUES ($1, $2, $3, $4, $5, $6)", name, startDateVal, endDateVal, description, techonologies, image)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")

}

func projectDetail(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var projectDets = ProjectDb{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(&projectDets.ID, &projectDets.ProjectName, &projectDets.StartDate, &projectDets.EndDate, &projectDets.Description, &projectDets.TechIcon, &projectDets.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	projectDets.StartDateStr = projectDets.StartDate.Format("02 Jan 2006")
	projectDets.EndDateStr = projectDets.EndDate.Format("02 Jan 2006")

	projectDets.Difference = dateDifference(projectDets.StartDate, projectDets.EndDate)

	sendProjectDets := map[string]interface{}{
		"Project": projectDets,
	}

	return c.Render(http.StatusOK, "project-detail.html", sendProjectDets)
}

func projectDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM public.tb_projects WHERE id=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func projectEditForm(c echo.Context) error {

	id, _ := strconv.Atoi(c.Param("id"))

	var projectDets = ProjectDb{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, name, start_date, end_date, description, technologies, image FROM tb_projects WHERE id=$1", id).Scan(&projectDets.ID, &projectDets.ProjectName, &projectDets.StartDate, &projectDets.EndDate, &projectDets.Description, &projectDets.TechIcon, &projectDets.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	projectDets.StartDateStr = projectDets.StartDate.Format("2006-01-02")
	projectDets.EndDateStr = projectDets.EndDate.Format("2006-01-02")

	sendDatas := map[string]interface{}{
		"Project": projectDets,
		"Button":  `<button type="submit" class="btn btn-dark rounded-5 px-4 py-1">Edit</button>`,
		"Action":  "/project-edit",
	}

	return c.Render(http.StatusOK, "project-form.html", sendDatas)
}

func projectEdit(c echo.Context) error {
	id := c.FormValue("id")

	name := c.FormValue("projectName")
	startDate := c.FormValue("startDate")
	endDate := c.FormValue("endDate")
	description := c.FormValue("description")
	techonologies := []string{}

	if c.FormValue("nodeJs") == "on" {
		techonologies = append(techonologies, "NodeJs")
	}
	if c.FormValue("reactJs") == "on" {
		techonologies = append(techonologies, "ReactJs")
	}
	if c.FormValue("javascript") == "on" {
		techonologies = append(techonologies, "Javascript")
	}
	if c.FormValue("go") == "on" {
		techonologies = append(techonologies, "Go")
	}

	startDateVal, _ := time.Parse("2006-01-02", startDate)
	endDateVal, _ := time.Parse("2006-01-02", endDate)

	image := "https://source.unsplash.com/random/900*700?games,football"

	_, err := connection.Conn.Exec(context.Background(), "UPDATE public.tb_projects SET name=$2, start_date=$3, end_date=$4, description=$5, technologies=$6, image=$7 WHERE id=$1", id, name, startDateVal, endDateVal, description, techonologies, image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func dateDifference(startDate time.Time, endDate time.Time) string {
	var difference string

	diff := endDate.Sub(startDate)
	getDiffMonths := diff.Hours() / 24 / 30
	getDiffWeeks := diff.Hours() / 24 / 7
	getDiffDays := diff.Hours() / 24
	getDiffHours := diff.Hours()

	if getDiffMonths >= 1 {
		difference = fmt.Sprint(int(getDiffMonths), " Months")
	} else if getDiffWeeks >= 1 {
		difference = fmt.Sprint(int(getDiffWeeks), " Weeks")
	} else if getDiffDays >= 1 {
		difference = fmt.Sprint(int(getDiffDays), " Days")
	} else {
		difference = fmt.Sprint(int(getDiffHours), " Hours")
	}

	return difference
}
