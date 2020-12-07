-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS event (
    id INT(11) AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    user_id INT(11) NOT NULL,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    notification_date DATETIME NOT NULL,
    UNIQUE (user_id, start_date)
) ENGINE=INNODB;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE event;
