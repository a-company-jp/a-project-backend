CREATE TABLE IF NOT EXISTS Users
(
    user_id         VARCHAR(255) PRIMARY KEY,
    firebase_uid    VARCHAR(255) NOT NULL UNIQUE,
    username        VARCHAR(255) NOT NULL,
    firstname       VARCHAR(255),
    lastname        VARCHAR(255),
    firstname_kana  VARCHAR(255),
    lastname_kana   VARCHAR(255),
    status_message  TEXT,
    icon_image_hash VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS Tags
(
    tag_id   VARCHAR(255) PRIMARY KEY,
    tag_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS UserTags
(
    user_id VARCHAR(255),
    tag_id  VARCHAR(255),
    PRIMARY KEY (user_id, tag_id),
    FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES Tags (tag_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Milestones
(
    milestone_id VARCHAR(255) PRIMARY KEY,
    user_id      VARCHAR(255),
    title        VARCHAR(255) NOT NULL,
    content      TEXT,
    image_hash   VARCHAR(255),
    begin_date   INT,
    finish_date  INT,
    FOREIGN KEY (user_id) REFERENCES Users (user_id) ON DELETE CASCADE
);
