CREATE TABLE IF NOT EXISTS dept_staffs_mapping (
    id bigint NOT NULL,
    departmentId bigint REFERENCES departments(id),
    staffId bigint REFERENCES staffs(id),
    PRIMARY KEY(id)
);