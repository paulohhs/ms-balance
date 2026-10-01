CREATE TABLE IF NOT EXISTS clients (
  id         VARCHAR(255) NOT NULL,
  name       VARCHAR(255) NOT NULL,
  email      VARCHAR(255) NOT NULL,
  created_at DATETIME     NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS accounts (
  id         VARCHAR(255)  NOT NULL,
  client_id  VARCHAR(255)  NOT NULL,
  balance    DECIMAL(15,2) NOT NULL DEFAULT 0,
  created_at DATETIME      NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_accounts_client FOREIGN KEY (client_id) REFERENCES clients (id)
);

CREATE TABLE IF NOT EXISTS transactions (
  id              VARCHAR(255)  NOT NULL,
  account_id_from VARCHAR(255)  NOT NULL,
  account_id_to   VARCHAR(255)  NOT NULL,
  amount          DECIMAL(15,2) NOT NULL,
  created_at      DATETIME      NOT NULL,
  PRIMARY KEY (id)
);
