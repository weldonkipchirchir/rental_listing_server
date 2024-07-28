ALTER TABLE notifications 
ADD COLUMN admin_id INT,
ADD CONSTRAINT fk_admin_id FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE;

ALTER TABLE notifications
ALTER COLUMN user_id DROP NOT NULL;
