package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"mlibrary-htmx/pkg/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	template *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}

func main() {
	err := database.InitDb()
	if err != nil {
		log.Fatalf("Shit: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/css", "css")

	t := &Template{
		template: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t
	e.GET("/", RedirectToBase)

	e.GET("/books", GetAllBooks)

	e.GET("/books/new", HandleNewBook)
	e.POST("/books/new", CreateNewBook)
	e.PUT("/books/new/:id", UpdateExistingBook)

	e.GET("/books/:id", HandleExistingBook)
	e.DELETE("/books/:id", HandleDeleteBook)
	e.GET("/books/show/:id", HandleShowBook)

	e.GET("/upload", GetUploadPage)

	e.GET("/download", Download)
	e.POST("/upload", Upload)

	// e.HTTPErrorHandler = customHTTPErrorHandler
	e.Logger.Fatal(e.Start(":4444"))
}
