ALTER TABLE notifications 
DROP CONSTRAINT fk_sender_admin_id,
DROP COLUMN sender_admin_id;

ALTER TABLE notifications 
DROP CONSTRAINT fk_sender_user_id,
DROP COLUMN sender_user_id;


