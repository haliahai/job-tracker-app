CREATE TABLE contacts (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    application_id BIGINT UNSIGNED NOT NULL,
    name           VARCHAR(255) NOT NULL,
    role           ENUM('recruiter','referral','hiring_manager','interviewer','other') NOT NULL,
    company        VARCHAR(255),
    email          VARCHAR(255),
    phone          VARCHAR(50),
    linkedin_url   VARCHAR(500),
    notes          TEXT,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_contacts_application (application_id),
    CONSTRAINT fk_contacts_application
        FOREIGN KEY (application_id) REFERENCES applications (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
