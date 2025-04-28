CREATE TABLE lecturer_course_assignments (
  id INT NOT NULL AUTO_INCREMENT,
  lecturer_teaching_code varchar(30) NOT NULL UNIQUE,
  course_id INT NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (lecturer_teaching_code) REFERENCES lecturer_teachings(code) ON DELETE CASCADE,
  FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE = InnoDB;