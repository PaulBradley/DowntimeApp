CREATE TABLE IF NOT EXISTS users (
    cellar          CHAR(26)        NOT NULL,
    subject         CHAR(50)        NOT NULL,
    email           VARCHAR(100)    NOT NULL,
    family_name     VARCHAR(50)     NOT NULL,
    given_name      VARCHAR(50)     NOT NULL,
    is_enabled      CHAR(1)         NOT NULL DEFAULT 'Y',
    last_login      TIMESTAMP       NULL,
    last_updated    TIMESTAMP       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    non_prod_access CHAR(1)         NOT NULL DEFAULT 'N',
    PRIMARY KEY (cellar, subject)
);

CREATE INDEX ASYNC IF NOT EXISTS idx_users_cellar_subject ON users (cellar, family_name);
CREATE INDEX ASYNC IF NOT EXISTS idx_users_cellar_is_enabled ON users (cellar, is_enabled);

COMMENT ON  TABLE users IS '
// As Amazon Aurora DSQL does not support foreign key constraints, adding
// DBML (https://dbml.dbdiagram.io/docs) hints as comments to indicate
// relationships between tables for documentation purposes.

Ref: users.cellar > organisations.cellar
';