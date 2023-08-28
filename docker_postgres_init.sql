CREATE TABLE user_activity
(
    user_id INT NOT NULL,
    activity_id INT NOT NULL,
    birth_year INT NOT NULL,
    location VARCHAR(100) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    activity_type VARCHAR(100) NOT NULL,
    intensity INT,
    duration INT
);