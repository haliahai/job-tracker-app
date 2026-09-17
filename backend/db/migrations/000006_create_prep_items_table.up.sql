CREATE TABLE prep_items (
    id                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    application_id      BIGINT UNSIGNED NOT NULL,
    interview_stage_id  BIGINT UNSIGNED DEFAULT NULL,
    category            ENUM('behavioral','technical','system_design','coding','domain_knowledge','company_specific','other') NOT NULL DEFAULT 'other',
    kind                ENUM('to_prepare','asked_to_me','i_asked_them') NOT NULL DEFAULT 'to_prepare',
    content              TEXT NOT NULL,
    my_answer_notes      TEXT,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_prep_items_application (application_id),
    KEY idx_prep_items_stage (interview_stage_id),
    CONSTRAINT fk_prep_items_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_prep_items_stage
        FOREIGN KEY (interview_stage_id) REFERENCES interview_stages (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
