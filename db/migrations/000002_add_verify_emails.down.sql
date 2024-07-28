-- Drop foreign key constraints on user_verify_emails and admin_verify_emails tables
ALTER TABLE "user_verify_emails" DROP CONSTRAINT IF EXISTS fk_user_verify_emails_users_email;
ALTER TABLE "admin_verify_emails" DROP CONSTRAINT IF EXISTS fk_admin_verify_emails_admins_email;

-- Drop user_verify_emails and admin_verify_emails tables
DROP TABLE IF EXISTS "user_verify_emails";
DROP TABLE IF EXISTS "admin_verify_emails";
