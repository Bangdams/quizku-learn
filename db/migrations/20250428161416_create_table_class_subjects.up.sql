CREATE TABLE class_subjects (
  class_id INT NOT NULL,
  course_code VARCHAR(3) NOT NULL,
  PRIMARY KEY (class_id, course_code),
  FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
  FOREIGN KEY (course_code) REFERENCES courses(course_code) ON DELETE CASCADE
) ENGINE = InnoDB;