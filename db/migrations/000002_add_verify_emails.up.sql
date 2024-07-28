-- Create verify_emails for users
CREATE TABLE "user_verify_emails" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "secret_code" varchar NOT NULL,
    "is_used" bool NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "expired_at" timestamptz NOT NULL DEFAULT now() + interval '15 minutes',
    CONSTRAINT fk_user_verify_emails_users_email
        FOREIGN KEY ("email") REFERENCES "users" ("email") ON DELETE CASCADE
);

-- Add is_email_verified column to users table
ALTER TABLE "users" ADD COLUMN "is_email_verified" bool NOT NULL DEFAULT false;

-- Create verify_emails for admins
CREATE TABLE "admin_verify_emails" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "secret_code" varchar NOT NULL,
    "is_used" bool NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "expired_at" timestamptz NOT NULL DEFAULT now() + interval '15 minutes',
    CONSTRAINT fk_admin_verify_emails_admins_email
        FOREIGN KEY ("email") REFERENCES "admins" ("email") ON DELETE CASCADE
);

-- Add is_email_verified column to admins table
ALTER TABLE "admins" ADD COLUMN "is_email_verified" bool NOT NULL DEFAULT false;
