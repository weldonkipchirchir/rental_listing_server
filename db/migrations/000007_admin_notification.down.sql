ALTER TABLE notifications 
DROP CONSTRAINT fk_admin_id,
DROP COLUMN admin_id;

ALTER TABLE notifications
ALTER COLUMN user_id SET NOT NULL;