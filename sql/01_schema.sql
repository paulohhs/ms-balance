CREATE TABLE IF NOT EXISTS accounts (
  id         VARCHAR(255)  NOT NULL,
  balance    DECIMAL(15,2) NOT NULL DEFAULT 0,
  updated_at DATETIME      NOT NULL,
  PRIMARY KEY (id)
);
