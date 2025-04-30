CREATE TABLE class_subjects (
  class_id INT NOT NULL,
  course_id INT NOT NULL,
  PRIMARY KEY (class_id, course_id),
  FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE,
  FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
) ENGINE = InnoDB;