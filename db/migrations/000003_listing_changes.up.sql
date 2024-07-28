-- Remove the available_from and available_to columns from the listings table
ALTER TABLE listings
DROP COLUMN available_from,
DROP COLUMN available_to;

-- Add the available and imageLink columns to the listings table
ALTER TABLE listings
ADD COLUMN available BOOLEAN DEFAULT TRUE,
ADD COLUMN imageLink VARCHAR(255);
