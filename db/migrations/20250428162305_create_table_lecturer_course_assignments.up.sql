CREATE TABLE lecturer_course_assignments (
  lecturer_teaching_code varchar(30),
  course_code VARCHAR(3),
  PRIMARY KEY (lecturer_teaching_code, course_code),
  FOREIGN KEY (lecturer_teaching_code) REFERENCES lecturer_teachings(code) ON DELETE CASCADE,
  FOREIGN KEY (course_code) REFERENCES courses(course_code) ON DELETE CASCADE
) ENGINE = InnoDB;