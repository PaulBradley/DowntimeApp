CREATE TABLE IF NOT EXISTS organisations (
    ods             VARCHAR(10)     NOT NULL,
    cellar          CHAR(26)        NOT NULL,
    organisation    VARCHAR(255)    NOT NULL,
    PRIMARY KEY (ods)
);

CREATE INDEX ASYNC IF NOT EXISTS idx_organisations_cellar ON organisations (cellar);

COMMENT ON  TABLE organisations IS '
// As Amazon Aurora DSQL does not support foreign key constraints, adding
// DBML (https://dbml.dbdiagram.io/docs) hints as comments to indicate
// relationships between tables for documentation purposes.

Ref: organisations.cellar < facilities.cellar
';