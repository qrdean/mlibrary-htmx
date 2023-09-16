ALTER TABLE master_books RENAME TO temp_master_books;
CREATE TABLE IF NOT EXISTS master_books (
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  lccn TEXT DEFAULT NULL,
  isbn TEXT DEFAULT NULL,
  title TEXT DEFAULT NULL,
  author_first TEXT DEFAULT NULL,
  author_last TEXT DEFAULT NULL,
  copyright_date DATE DEFAULT NULL,
  publisher TEXT DEFAULT NULL,
  location TEXT DEFAULT NULL,
  genre TEXT DEFAULT NULL,
  pages TEXT DEFAULT NULL
);
INSERT INTO master_books(
  id,
  created_at,
  lccn,
  isbn,
  title,
  author_last,
  copyright_date,
  publisher,
  location,
  genre,
  pages
)
SELECT id, created_at, lccn, isbn, author, copyright_date, publisher, location, genre, pages
FROM temp_master_books;

DROP TABLE temp_master_books;
