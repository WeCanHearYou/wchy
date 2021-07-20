CREATE TABLE IF NOT EXISTS webhooks (
  id SERIAL PRIMARY KEY,
  name VARCHAR(60) NOT NULL,
  type SMALLINT NOT NULL,
  status SMALLINT NOT NULL,
  url VARCHAR(256) NOT NULL,
  content TEXT NULL,
  http_method VARCHAR(16) NOT NULL,
  additional_http_headers JSONB NULL,
  tenant_id INT NOT NULL,
  FOREIGN KEY (tenant_id) REFERENCES tenants (id)
);