ALTER TABLE notifications 
ADD COLUMN sender_admin_id INT,
ADD CONSTRAINT fk_sender_admin_id FOREIGN KEY (sender_admin_id) REFERENCES admins(id) ON DELETE CASCADE;
ALTER TABLE notifications 
ADD COLUMN sender_user_id INT,
ADD CONSTRAINT fk_sender_user_id FOREIGN KEY (sender_user_id) REFERENCES users(id) ON DELETE CASCADE;
