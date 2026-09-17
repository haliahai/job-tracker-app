CREATE TABLE application_status_history (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    application_id BIGINT UNSIGNED NOT NULL,
    status_id      TINYINT UNSIGNED NOT NULL,
    changed_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes          VARCHAR(500),
    PRIMARY KEY (id),
    KEY idx_status_history_application (application_id),
    KEY idx_status_history_status (status_id),
    CONSTRAINT fk_status_history_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_status_history_status
        FOREIGN KEY (status_id) REFERENCES application_statuses (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
