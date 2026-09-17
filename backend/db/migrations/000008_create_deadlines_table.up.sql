CREATE TABLE deadlines (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    application_id  BIGINT UNSIGNED NOT NULL,
    deadline_type   VARCHAR(100) NOT NULL,
    due_date        DATETIME NOT NULL,
    is_completed    BOOLEAN NOT NULL DEFAULT FALSE,
    notes           VARCHAR(500),
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_deadlines_application (application_id),
    KEY idx_deadlines_due_date (due_date),
    CONSTRAINT fk_deadlines_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
