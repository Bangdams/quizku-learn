CREATE TABLE lecturer_course_assignments (
  lecturer_teaching_code varchar(30),
  course_id INT,
  PRIMARY KEY (lecturer_teaching_code, course_id),
  FOREIGN KEY (lecturer_teaching_code) REFERENCES lecturer_teachings(code) ON DELETE CASCADE,
  FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE = InnoDB;