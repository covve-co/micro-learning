CREATE TABLE courses (
  id            SERIAL  PRIMARY KEY,
  title         text,
  description   text,
  template_name text                  NOT NULL,
  num_sections  integer               NOT NULL
);
