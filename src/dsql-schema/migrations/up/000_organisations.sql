CREATE TABLE IF NOT EXISTS organisations (
    ods             VARCHAR(10)     NOT NULL,
    vault           CHAR(26)        NOT NULL,
    organisation    VARCHAR(255)    NOT NULL,
    PRIMARY KEY (ods)
);

CREATE INDEX ASYNC IF NOT EXISTS idx_organisations_vault ON organisations (vault);