-- Add the available_from and available_to columns back to the listings table
ALTER TABLE listings
ADD COLUMN available_from DATE,
ADD COLUMN available_to DATE;

-- Remove the available and imageLink columns from the listings table
ALTER TABLE listings
DROP COLUMN available,
DROP COLUMN imageLink;
