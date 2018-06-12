CREATE TABLE users (
  id          SERIAL  PRIMARY KEY,
  staff_no    text    UNIQUE        NOT NULL, 
  name        text                  NOT NULL, 
  nric        char(9)               NOT NULL, 
  password    text,
  registered  boolean               DEFAULT FALSE,
  token       text    UNIQUE        NOT NULL
);
