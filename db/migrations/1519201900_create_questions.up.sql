CREATE TABLE questions (
  id        SERIAL  PRIMARY KEY,
  course_id integer REFERENCES courses(id)  NOT NULL,
  title     text                            NOT NULL, 
  content   text                            NOT NULL, 
  image     text
);
