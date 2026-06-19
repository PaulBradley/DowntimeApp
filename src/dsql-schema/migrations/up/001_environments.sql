CREATE TABLE IF NOT EXISTS environments (
    cellar          CHAR(26)    NOT NULL,
    vault           CHAR(26)    NOT NULL,
    environment     VARCHAR(11) NOT NULL,
    is_enabled      CHAR(1)     NOT NULL DEFAULT 'Y',
    is_in_downtime  CHAR(1)     NOT NULL DEFAULT 'N',
    banner_colour   CHAR(6)     NOT NULL DEFAULT '005EB8',
    PRIMARY KEY (vault)
);

CREATE INDEX ASYNC IF NOT EXISTS idx_environments_cellar ON environments (cellar);
CREATE INDEX ASYNC IF NOT EXISTS idx_environments_vault ON environments (vault, is_enabled);

COMMENT ON  TABLE environments IS '
// As Amazon Aurora DSQL does not support foreign key constraints, adding
// DBML (https://dbml.dbdiagram.io/docs) hints as comments to indicate
// relationships between tables for documentation purposes.

Ref: environments.cellar > organisations.cellar
';