CREATE TABLE completions (
    user_id     integer   REFERENCES users(id)    NOT NULL,
    course_id   integer   REFERENCES courses(id)  NOT NULL,
    timestamp   timestamp                         DEFAULT now(),
    PRIMARY KEY (user_id, course_id, timestamp)
);
