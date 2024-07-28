ALTER TABLE notifications 
DROP CONSTRAINT fk_booking,
DROP COLUMN email,
DROP COLUMN booking_id,
DROP COLUMN subject;