CREATE TABLE IF NOT EXISTS student (
    roll bigint NOT NULL,
    name text,
    dep_id bigint,
    CONSTRAINT fk_student
       FOREIGN KEY(dep_id)
          REFERENCES departments(id)
);