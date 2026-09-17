CREATE TABLE application_statuses (
    id           TINYINT UNSIGNED NOT NULL AUTO_INCREMENT,
    code         VARCHAR(30)  NOT NULL,
    label        VARCHAR(50)  NOT NULL,
    sort_order   TINYINT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_application_statuses_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO application_statuses (code, label, sort_order) VALUES
    ('wishlist',            'Wishlist',              10),
    ('applied',             'Applied',               20),
    ('phone_screen',        'Phone Screen',          30),
    ('technical_interview', 'Technical Interview',   40),
    ('onsite_interview',    'Onsite / Final Round',  50),
    ('offer',               'Offer Received',        60),
    ('accepted',            'Offer Accepted',        70),
    ('rejected',            'Rejected',              80),
    ('withdrawn',           'Withdrawn',             90),
    ('ghosted',             'Ghosted',              100);
