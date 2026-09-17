CREATE TABLE interview_stages (
    id                     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    application_id         BIGINT UNSIGNED NOT NULL,
    stage_name             VARCHAR(150) NOT NULL,
    stage_type             ENUM('phone_screen','technical','behavioral','system_design','take_home','onsite','panel','final','other') NOT NULL,
    scheduled_at           DATETIME DEFAULT NULL,
    completed_at           DATETIME DEFAULT NULL,
    format                 ENUM('phone','video','onsite','async') DEFAULT NULL,
    interviewer_contact_id BIGINT UNSIGNED DEFAULT NULL,
    outcome                ENUM('pending','passed','failed','cancelled','no_show') NOT NULL DEFAULT 'pending',
    feedback_notes         TEXT,
    created_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_interview_stages_application (application_id),
    KEY idx_interview_stages_scheduled_at (scheduled_at),
    KEY idx_interview_stages_interviewer (interviewer_contact_id),
    CONSTRAINT fk_interview_stages_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_interview_stages_contact
        FOREIGN KEY (interviewer_contact_id) REFERENCES contacts (id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
