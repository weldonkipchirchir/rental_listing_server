ALTER TABLE listings
DROP COLUMN IF EXISTS imageLink;
ALTER TABLE listings 
  ADD COLUMN imageLinks TEXT[];