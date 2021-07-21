ALTER TABLE tenants ADD COLUMN is_email_auth_allowed BOOLEAN;
UPDATE tenants SET is_email_auth_allowed = TRUE;
ALTER TABLE tenants ALTER COLUMN is_email_auth_allowed SET NOT NULL;