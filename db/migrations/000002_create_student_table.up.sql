CREATE TABLE IF NOT EXISTS students (
    rollnumber bigint NOT NULL,
    name text,
    departmentId bigint REFERENCES departments(id),
    PRIMARY KEY(rollnumber)
);