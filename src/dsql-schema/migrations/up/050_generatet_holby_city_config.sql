
INSERT INTO ORGANISATIONS (ods, cellar, organisation) VALUES ('HCH01', '01KSM3WRPK2D9K04RS92MBYSHT', 'Holby City NHS Foundation Trust');

INSERT INTO ENVIRONMENTS (cellar, vault, environment, is_enabled) VALUES ('01KSM3WRPK2D9K04RS92MBYSHT', '1B4E3DF588D92ED9A63DCE0016', 'PRODUCTION', 'Y');
INSERT INTO ENVIRONMENTS (cellar, vault, environment, is_enabled) VALUES ('01KSM3WRPK2D9K04RS92MBYSHT', '4DED63CFF6D4773B0D66358754', 'MOCK', 'Y');
INSERT INTO ENVIRONMENTS (cellar, vault, environment, is_enabled) VALUES ('01KSM3WRPK2D9K04RS92MBYSHT', '28DDBC80B91F0AD0429414BC88', 'CERT', 'Y');
INSERT INTO ENVIRONMENTS (cellar, vault, environment, is_enabled) VALUES ('01KSM3WRPK2D9K04RS92MBYSHT', 'E5DD1CF54E8DF31FF37D9CC01E', 'TRAINING', 'Y');
INSERT INTO ENVIRONMENTS (cellar, vault, environment, is_enabled) VALUES ('01KSM3WRPK2D9K04RS92MBYSHT', 'EC0E11BB07C0D50E0AB3F2F10B', 'BUILD', 'Y');

INSERT INTO FACILITIES (cellar, institution) VALUES ('01KSM3WRPK2D9K04RS92MBYSHT', 'Holby City Hospital');
INSERT INTO FACILITIES (cellar, institution) VALUES ('01KSM3WRPK2D9K04RS92MBYSHT', 'Wyvern District Hospital');

INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (1, 1, 'AE', 'Accident & Emergency');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (1, 2, 'MAU', 'Medical Assessment Unit');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (1, 3, 'WARD1', 'Ward 1');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (1, 4, 'WARD2', 'Ward 2');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (1, 5, 'WARD3', 'Ward 3');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (1, 6, 'WARD4', 'Ward 4');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (1, 7, 'WARD5', 'Ward 5');

INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (2, 1, 'MU', 'Maternity Unit');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (2, 2, 'SCBU', 'Special Care Baby Unit');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (2, 3, 'WARD1', 'Ward 1');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (2, 4, 'WARD3', 'Ward 3');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (2, 5, 'WARD3', 'Ward 4');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (2, 6, 'WARD4', 'Ward 5');
INSERT INTO PRACTICE_SETTING (facility_id, sort_order, setting_code, setting_name) VALUES (2, 7, 'WARD6', 'Ward 6');