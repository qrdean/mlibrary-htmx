ALTER TABLE master_books
RENAME COLUMN publish_date TO copyright_date;
ALTER TABLE master_books
ADD COLUMN publisher TEXT DEFAULT NULL;
ALTER TABLE master_books
ADD COLUMN location TEXT DEFAULT NULL;
ALTER TABLE master_books
ADD COLUMN genre TEXT DEFAULT NULL;
ALTER TABLE master_books
ADD COLUMN pages TEXT DEFAULT NULL;

-- DROP TABLE IF EXISTS master_books;
-- CREATE TABLE IF NOT EXISTS master_books (
--   id INTEGER PRIMARY KEY,
--   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--   lccn TEXT DEFAULT NULL,
--   isbn TEXT DEFAULT NULL,
--   title TEXT DEFAULT NULL,
--   author TEXT DEFAULT NULL,
--   copyright_date DATE DEFAULT NULL
--   publisher TEXT DEFAULT NULL,
--   location TEXT DEFAULT NULL,
--   genre TEXT DEFAULT NULL,
--   pages TEXT DEFAULT NULL,
-- )
