CREATE TABLE IF NOT EXISTS patients (
    vault               CHAR(26)        NOT NULL,
    practice_setting_id BIGINT          NOT NULL,
    mrn                 VARCHAR(10)     NOT NULL,
    given_name          VARCHAR(255)    NOT NULL,
    family_name         VARCHAR(255)    NOT NULL,
    admission_date_time TIMESTAMP       NOT NULL,
    discharge_date_time TIMESTAMP       NULL,
    PRIMARY KEY (vault, mrn)
);

CREATE INDEX ASYNC IF NOT EXISTS idx_patients_vault_mrn ON patients (vault, mrn);
CREATE INDEX ASYNC IF NOT EXISTS idx_patients_vault_practice_setting_family_name ON patients (vault, practice_setting_id, family_name);