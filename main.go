package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang-htmx/pkg/database"

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

type BookContent struct {
	Header Header
	Books  []database.Book
	Params map[string]interface{}
}

type Header struct {
	Title string
}

type NewBookPage struct {
	Header   Header
	Book     *database.Book
	Message  string
	Existing bool
	Errors   map[string]string
}

func RedirectToBase(c echo.Context) error {
	basePath := "/books"
	return c.Redirect(http.StatusFound, basePath)
}

func sort_and_paginate_books(page int, sortBy string) ([]database.Book, error) {
	maxsize := 25
	iteration := maxsize * page
	initial := iteration - maxsize
	books, err := database.SortAndPaginateBooks(initial, sortBy)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func paginate_books(page int) ([]database.Book, error) {
	maxsize := 25
	iteration := maxsize * page
	initial := iteration - maxsize
	books, err := database.PaginateBooks(initial)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func GetAllBooks(c echo.Context) error {
	searchParam := c.QueryParam("q")
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
	}

	sortParam := c.QueryParam("sort-by")

	// count := len()
	params := make(map[string]interface{})
	params["page"] = page
	params["next"] = page + 1
	params["prev"] = page - 1
	// params["count"] = count
	params["search"] = searchParam

	if c.Request().Header.Get("HX-Trigger") == "sort-by" {
		if sortParam != "" {
			fmt.Println("sort param", sortParam)
			books, err := sort_and_paginate_books(page, sortParam)
			if err != nil {
        c.Logger().Error(err)
				return err
			}

			return c.Render(http.StatusOK, "book", BookContent{
				Header: Header{
					Title: "Books",
				},
				Books:  books,
				Params: params,
			})
		}

		books, err := paginate_books(page)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.Render(http.StatusOK, "book", BookContent{
			Header: Header{
				Title: "Books",
			},
			Books:  books,
			Params: params,
		})
	}

	if searchParam != "" {
		if c.Request().Header.Get("HX-Trigger") == "search" {
			books, err := database.FilterBook(searchParam)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.Render(http.StatusOK, "book", BookContent{
				Header: Header{
					Title: "Books",
				},
				Books:  books,
				Params: params,
			})
		}
	}

	if c.Request().Header.Get("HX-Trigger") == "search" {
		books, err := paginate_books(page)
		if err != nil {
      c.Logger().Error(err)
			return err
		}

		return c.Render(http.StatusOK, "book", BookContent{
			Header: Header{
				Title: "Books",
			},
			Books:  books,
			Params: params,
		})
	}

	books, err := paginate_books(page)
	if err != nil {
    c.Logger().Error(err)
		return err
	}

	return c.Render(http.StatusOK, "books", BookContent{
		Header: Header{
			Title: "Books",
		},
		Books:  books,
		Params: params,
	})
}

func HandleNewBook(c echo.Context) error {
	return c.Render(http.StatusOK, "new-book", NewBookPage{
		Header: Header{
			Title: "Create Book",
		},
		Book:     nil,
		Existing: false,
		Errors:   map[string]string{},
	})
}

func HandleExistingBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	book, err := database.GetBookById(id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.Render(http.StatusOK, "new-book", NewBookPage{
		Header: Header{
			Title: "Update Book",
		},
		Book:     book,
		Existing: true,
		Errors:   map[string]string{},
	})
}

func CreateNewBook(c echo.Context) error {
	newBook := database.Book{
		Isbn:      c.FormValue("isbn"),
		Lccn:      c.FormValue("lccn"),
		Title:     c.FormValue("title"),
		Author:    c.FormValue("author"),
		Publisher: c.FormValue("publisher"),
		Location:  c.FormValue("location"),
		Genre:     c.FormValue("genre"),
		Pages:     c.FormValue("pages"),
		Id:        -1,
	}

	if c.FormValue("copyright-date") == "" {
		errors := make(database.ErrorMap)
		errors["publish_date"] = "Copyright Date Required"
		return c.Render(http.StatusOK, "new-book-template", NewBookPage{
			Book:     &newBook,
			Existing: false,
			Errors:   errors,
		})
	}
	layout := "2006-01-02"
	publish_date, err := time.Parse(layout, c.FormValue("copyright-date"))
	if err != nil {
		return fmt.Errorf("error parsing date string: %v", err)
	}

	newBook.CopyrightDate = publish_date
	newBook.CopyrightDateString = publish_date.Format("2006-01-02")
	newBook.Id = -1

	errorMap, err := newBook.Save()
	if err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusOK, "new-book-template", NewBookPage{
			Message:  "An Internal Error Occurred",
			Book:     &newBook,
			Existing: false,
			Errors:   errorMap,
		})
	}

	if len(errorMap) > 0 {
		return c.Render(http.StatusOK, "new-book-template", NewBookPage{
			Book:     &newBook,
			Existing: false,
			Errors:   errorMap,
		})
	}

	return c.Render(http.StatusOK, "new-book-template", NewBookPage{
		Message:  "Book Created",
		Book:     &newBook,
		Existing: true,
		Errors:   map[string]string{},
	})
}

func UpdateExistingBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	newBook := database.Book{
		Isbn:      c.FormValue("isbn"),
		Lccn:      c.FormValue("lccn"),
		Title:     c.FormValue("title"),
		Author:    c.FormValue("author"),
		Publisher: c.FormValue("publisher"),
		Location:  c.FormValue("location"),
		Genre:     c.FormValue("genre"),
		Pages:     c.FormValue("pages"),
		Id:        id,
	}
	if c.FormValue("copyright-date") == "" {
		errorMap := make(map[string]string)
		errorMap["publish_date"] = "Copyright Date Required"
		return c.Render(http.StatusOK, "new-book-template", NewBookPage{
			Message:  "",
			Book:     &newBook,
			Existing: false,
			Errors:   errorMap,
		})
	}

	layout := "2006-01-02"
	publish_date, err := time.Parse(layout, c.FormValue("copyright-date"))
	if err != nil {
		return fmt.Errorf("error parsing date string: %v", err)
	}

	newBook.CopyrightDate = publish_date

	errorMap, err := newBook.Save()
	if err != nil {
		c.Logger().Error(err)
		return c.Render(http.StatusOK, "new-book-template", NewBookPage{
			Message:  "Internal Error Occurred",
			Book:     &newBook,
			Existing: true,
			Errors:   map[string]string{},
		})
	}

	if len(errorMap) > 0 {
		return c.Render(http.StatusOK, "new-book-template", NewBookPage{
			Message:  "",
			Book:     &newBook,
			Existing: true,
			Errors:   errorMap,
		})
	}

	return c.Render(http.StatusOK, "new-book-template", NewBookPage{
		Message:  "Book Updated",
		Book:     &newBook,
		Existing: true,
		Errors:   errorMap,
	})
}

func HandleShowBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	book, err := database.GetBookById(id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.Render(http.StatusOK, "show-book", NewBookPage{
		Header: Header{
			Title: "Show Book",
		},
		Book:     book,
		Existing: true,
		Errors:   map[string]string{},
	})
}

func HandleDeleteBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	err = database.DeleteBook(id)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/books")
}

func GetUploadPage(c echo.Context) error {
	return c.Render(http.StatusOK, "upload", BookContent{
		Header: Header{
			Title: "Upload Book",
		},
	})
}

func Download(c echo.Context) error {
	return c.File("./csv_temp/book_template.csv")
}

func Upload(c echo.Context) error {
	// Handle the book upload logistics here
	file, err := c.FormFile("file")
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	contentType := file.Header.Get("Content-Type")
	if contentType != "text/csv" {
		return c.HTML(http.StatusOK, fmt.Sprintf("<p>File is not correct type %s. Please use a .csv</p>", contentType))
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		c.Logger().Error(err)
		return err
	}

	err = handleBookUpload(dst.Name())
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.HTML(http.StatusOK, fmt.Sprintf("<p>File %s uploaded</p>", file.Filename))
}

func handleBookUpload(filename string) error {
	dst, err := os.Open(filename)
	if err != nil {
		return err
	}

	csvReader := csv.NewReader(dst)
	data, err := csvReader.ReadAll()
	if err != nil {
		return err
	}

	var bookRows []database.BookCsv
	for i, line := range data {
		if i > 0 {
			var bookRow database.BookCsv
			for j, field := range line {
				if j == 0 {
					bookRow.Isbn = field
				}
				if j == 1 {
					bookRow.Lccn = field
				}
				if j == 2 {
					bookRow.Title = field
				}
				if j == 3 {
					bookRow.Author = field
				}
				if j == 4 {
					if field != "" {
						layout := "01-02-2006"
						publish_date, err := time.Parse(layout, field)
						if err != nil {
							fmt.Println(err)
              bookRow.CopyrightDate = time.Time{}
            } else {
              bookRow.CopyrightDate = publish_date
            }
					}
				}
				if j == 5 {
					bookRow.Publisher = field
				}
				if j == 6 {
					bookRow.Location = field
				}
				if j == 7 {
					bookRow.Genre = field
				}
				if j == 8 {
					bookRow.Pages = field
				}
			}
			bookRows = append(bookRows, bookRow)
		}
	}
	err = database.BulkInsert(bookRows)
	if err != nil {
		return err
	}
	return nil
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
