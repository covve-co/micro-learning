CREATE TABLE options (
  question_id integer   REFERENCES questions(id)  NOT NULL,
  position    integer,
  content     text                                NOT NULL, 
  correct     boolean                             NOT NULL,
  PRIMARY KEY (question_id, position)
);
