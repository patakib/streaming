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

CREATE TABLE city_stats
(
    city VARCHAR(100) NOT NULL,
    popular_sport VARCHAR(100) NOT NULL,
    popular_sport_count INT,
    most_active_user INT,
    most_active_user_activity_count INT
);