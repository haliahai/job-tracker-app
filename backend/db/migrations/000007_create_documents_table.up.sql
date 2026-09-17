CREATE TABLE documents (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    application_id  BIGINT UNSIGNED NOT NULL,
    document_type   ENUM('resume','cover_letter','portfolio','other') NOT NULL,
    version_label   VARCHAR(150) NOT NULL,
    file_path       VARCHAR(1000),
    notes           TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_documents_application (application_id),
    CONSTRAINT fk_documents_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
