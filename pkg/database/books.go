package database

import (
  "time"
  "fmt"
  "database/sql"
)

type Book struct {
  Id int
  CreatedDate time.Time
  Lccn string
  Isbn string
  Title string
  Author string
  CopyrightDate time.Time
  CopyrightDateString string
  Publisher string
  Location string
  Genre string
  Pages string
}

type BookCsv struct {
  Lccn string  
  Isbn string
  Title string
  Author string
  CopyrightDate time.Time
  Publisher string
  Location string
  Genre string
  Pages string
}

type ErrorMap = map[string]string

const GET_BOOK_LIST_QUERY = "SELECT * FROM master_books"
const PAGINATE_BOOK_LIST_QUERY = "SELECT * FROM master_books WHERE id > $1 ORDER BY id LIMIT 25;"
const PAGINATE_BOOK_LIST_SORT_BY_QUERY = "SELECT * FROM master_books WHERE id > $1 ORDER BY %s LIMIT 25;"
const GET_BOOK_BY_ID_QUERY = "SELECT * FROM master_books WHERE id = $1"
const DELETE_BOOK_BY_ID_QUERY = "DELETE FROM master_books WHERE id = $1"
const FILTER_BOOKS_QUERY = `SELECT * FROM master_books
WHERE (
 lccn like $1 or
 isbn like $1 or 
 title like $1 or 
 author like $1 or  
 copyright_date like $1 or
 publisher like $1 or
 location like $1 or
 genre like $1
)`

// 9 values
const INSERT_BOOK_QUERY = `INSERT INTO master_books (lccn, isbn, title, author, copyright_date, publisher, location, genre, pages) values (?,?,?,?,?,?,?,?,?)`
// 6 Values. Ending with id
const UPDATE_BOOK_QUERY = `UPDATE master_books SET lccn = ?, isbn = ?, title = ?, author = ?, copyright_date = ?, publisher = ?, location = ?, genre = ?, pages = ? WHERE id = ?` 

func GetBooksList() ([]Book, error) {
  res, err := Db.Query(GET_BOOK_LIST_QUERY)
  if err != nil {
    return nil, err
  }
  defer res.Close()

  var books []Book

  for res.Next() {
    var lccn sql.NullString
    var isbn sql.NullString
    var title sql.NullString
    var author sql.NullString
    var copyright_date_string sql.NullString
    var publisher sql.NullString
    var location sql.NullString
    var genre sql.NullString
    var pages sql.NullString
    var _created_at string
    var _id int

    err = res.Scan(
      &_id,
      &_created_at,
      &lccn,
      &isbn,
      &title,
      &author,
      &copyright_date_string,
      &publisher,
      &location,
      &genre,
      &pages,
    )
    if err != nil {
      return nil, fmt.Errorf("unable to scan db row: %v", err)
    }

    layout := "2006-01-02T15:04:05Z"
    copyright_date, err := time.Parse(layout, getValidNullStr(copyright_date_string))
    if err != nil {
      return nil, fmt.Errorf("error parsing date string: %v", err)
    }

    book := Book{
      Lccn: getValidNullStr(lccn),
      Isbn: getValidNullStr(isbn),
      Title: getValidNullStr(title),
      Author: getValidNullStr(author),
      CopyrightDate: copyright_date,
      CopyrightDateString: copyright_date.Format("2006-01-02"),
      Publisher: getValidNullStr(publisher),
      Location: getValidNullStr(location),
      Genre: getValidNullStr(genre),
      Pages: getValidNullStr(pages),
      Id: _id,
    }

    books = append(books, book)
  }

  return books, nil
}

func PaginateBooks(lastId int) ([]Book, error) {
  res, err := Db.Query(PAGINATE_BOOK_LIST_QUERY, lastId)
  if err != nil {
    return nil, err
  }
  defer res.Close()

  var books []Book

  for res.Next() {
    var lccn sql.NullString
    var isbn sql.NullString
    var title sql.NullString
    var author sql.NullString
    var copyright_date_string sql.NullString
    var publisher sql.NullString
    var location sql.NullString
    var genre sql.NullString
    var pages sql.NullString
    var _created_at string
    var _id int

    err = res.Scan(
      &_id,
      &_created_at,
      &lccn,
      &isbn,
      &title,
      &author,
      &copyright_date_string,
      &publisher,
      &location,
      &genre,
      &pages,
    )
    if err != nil {
      return nil, fmt.Errorf("unable to scan db row: %v", err)
    }

    layout := "2006-01-02T15:04:05Z"
    copyright_date, err := time.Parse(layout, getValidNullStr(copyright_date_string))
    if err != nil {
      return nil, fmt.Errorf("error parsing date string: %v", err)
    }

    book := Book{
      Lccn: getValidNullStr(lccn),
      Isbn: getValidNullStr(isbn),
      Title: getValidNullStr(title),
      Author: getValidNullStr(author),
      CopyrightDate: copyright_date,
      CopyrightDateString: copyright_date.Format("2006-01-02"),
      Publisher: getValidNullStr(publisher),
      Location: getValidNullStr(location),
      Genre: getValidNullStr(genre),
      Pages: getValidNullStr(pages),
      Id: _id,
    }

    books = append(books, book)
  }

  return books, nil
}

func SortAndPaginateBooks(lastId int, sortBy string) ([]Book, error) {
  queryString := fmt.Sprintf(PAGINATE_BOOK_LIST_SORT_BY_QUERY, sortBy)
  res, err := Db.Query(queryString, lastId)
  if err != nil {
    return nil, err
  }
  defer res.Close()

  var books []Book

  for res.Next() {
    var lccn sql.NullString
    var isbn sql.NullString
    var title sql.NullString
    var author sql.NullString
    var copyright_date_string sql.NullString
    var publisher sql.NullString
    var location sql.NullString
    var genre sql.NullString
    var pages sql.NullString
    var _created_at string
    var _id int

    err = res.Scan(
      &_id,
      &_created_at,
      &lccn,
      &isbn,
      &title,
      &author,
      &copyright_date_string,
      &publisher,
      &location,
      &genre,
      &pages,
    )
    if err != nil {
      return nil, fmt.Errorf("unable to scan db row: %v", err)
    }

    layout := "2006-01-02T15:04:05Z"
    copyright_date, err := time.Parse(layout, getValidNullStr(copyright_date_string))
    if err != nil {
      return nil, fmt.Errorf("error parsing date string: %v", err)
    }

    book := Book{
      Lccn: getValidNullStr(lccn),
      Isbn: getValidNullStr(isbn),
      Title: getValidNullStr(title),
      Author: getValidNullStr(author),
      CopyrightDate: copyright_date,
      CopyrightDateString: copyright_date.Format("2006-01-02"),
      Publisher: getValidNullStr(publisher),
      Location: getValidNullStr(location),
      Genre: getValidNullStr(genre),
      Pages: getValidNullStr(pages),
      Id: _id,
    }

    books = append(books, book)
  }

  return books, nil
}

func GetBookById(id int) (*Book, error) {
  res, err := Db.Query(GET_BOOK_BY_ID_QUERY, id)
  if err != nil {
    return nil, err
  }
  defer res.Close()

  var book Book

  if res.Next() {
    var lccn sql.NullString
    var isbn sql.NullString
    var title sql.NullString
    var author sql.NullString
    var copyright_date_string sql.NullString
    var publisher sql.NullString
    var location sql.NullString
    var genre sql.NullString
    var pages sql.NullString
    var _created_at string
    var _id int

    err = res.Scan(
      &_id,
      &_created_at,
      &lccn,
      &isbn,
      &title,
      &author,
      &copyright_date_string,
      &publisher,
      &location,
      &genre,
      &pages,
    )
    if err != nil {
      return nil, fmt.Errorf("unable to scan db row: %v", err)
    }

    layout := "2006-01-02T15:04:05Z"
    copyright_date, err := time.Parse(layout, getValidNullStr(copyright_date_string))
    if err != nil {
      return nil, fmt.Errorf("error parsing date string: %v", err)
    }

    book = Book{
      Lccn: getValidNullStr(lccn),
      Isbn: getValidNullStr(isbn),
      Title: getValidNullStr(title),
      Author: getValidNullStr(author),
      CopyrightDate: copyright_date,
      CopyrightDateString: copyright_date.Format("2006-01-02"),
      Publisher: getValidNullStr(publisher),
      Location: getValidNullStr(location),
      Genre: getValidNullStr(genre),
      Pages: getValidNullStr(pages),
      Id: _id,
    }
  }

  return &book, nil
}

func DeleteBook(id int) error {
  _, err := Db.Exec(DELETE_BOOK_BY_ID_QUERY, id)
  if err != nil {
    return fmt.Errorf("unable to delete contact from db: %v", err)
  }

  return nil
}

func FilterBook(q string) ([]Book, error) {
  res, err := Db.Query(FILTER_BOOKS_QUERY, "%" + q + "%")

  if err != nil {
    return nil, fmt.Errorf("unable to query db: %v", err)
  }
  defer res.Close()

  var books []Book
  for res.Next() {
    var lccn sql.NullString
    var isbn sql.NullString
    var title sql.NullString
    var author sql.NullString
    var copyright_date_string sql.NullString
    var publisher sql.NullString
    var location sql.NullString
    var genre sql.NullString
    var pages sql.NullString
    var _created_at string
    var _id int

    err = res.Scan(
      &_id,
      &_created_at,
      &lccn,
      &isbn,
      &title,
      &author,
      &copyright_date_string,
      &publisher,
      &location,
      &genre,
      &pages,
    )
    if err != nil {
      return nil, fmt.Errorf("unable to scan db row: %v", err)
    }

    layout := "2006-01-02T15:04:05Z"
    copyright_date, err := time.Parse(layout, getValidNullStr(copyright_date_string))
    if err != nil {
      return nil, fmt.Errorf("error parsing date string: %v", err)
    }

    book := Book{
      Lccn: getValidNullStr(lccn),
      Isbn: getValidNullStr(isbn),
      Title: getValidNullStr(title),
      Author: getValidNullStr(author),
      CopyrightDate: copyright_date,
      CopyrightDateString: copyright_date.Format("2006-01-02"),
      Publisher: getValidNullStr(publisher),
      Location: getValidNullStr(location),
      Genre: getValidNullStr(genre),
      Pages: getValidNullStr(pages),
      Id: _id,
    }

    books = append(books, book)
  }

  return books, nil
}

func (b *Book) validate() ErrorMap {
  errors := make(ErrorMap)
  if b.Isbn == "" {
    errors["isbn"] = "Isbn Required"
  }
  if b.Lccn == "" {
    errors["lccn"] = "Lccn Required"
  }
  if b.Author == "" {
    errors["author"] = "Author Name Required"
  }
  if b.Title == "" {
    errors["title"] = "Title Required"
  }

  return errors
}

func (b *Book) Save() (ErrorMap, error) {
  errors := b.validate()
  if len(errors) > 0 {
    return errors, nil
  }

  var err error
  if b.Id == -1 {
    _, err = Db.Exec(INSERT_BOOK_QUERY, b.Lccn, b.Isbn, b.Title, b.Author, b.CopyrightDate.Format("2006-01-02"), b.Publisher, b.Location, b.Genre, b.Pages)
  } else {
    _, err = Db.Exec(UPDATE_BOOK_QUERY, b.Lccn, b.Isbn, b.Title, b.Author, b.CopyrightDate.Format("2006-01-02"), b.Publisher, b.Location, b.Genre, b.Pages, b.Id)
  }

  return errors, err
}

func BulkInsert(bookCsv []BookCsv) error {
  tx, err := Db.Begin()
  if err != nil {
    return err
  }
  stmt, err := tx.Prepare(INSERT_BOOK_QUERY)
  if err != nil {
    return err
  }
  defer stmt.Close()
  for _, line := range bookCsv {
    _, err = stmt.Exec(line.Lccn, line.Isbn, line.Title, line.Author, line.CopyrightDate, line.Publisher, line.Location, line.Genre, line.Pages)  
    if err != nil {
      tx.Rollback()
      return err 
    }
  }
  err = tx.Commit()
  if err != nil {
    return err
  }
  return nil
} 

func getValidNullStr(nullString sql.NullString) string {
  if nullString.Valid {
    return nullString.String
  } else {
    return ""
  }
}

