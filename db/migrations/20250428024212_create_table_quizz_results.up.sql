CREATE TABLE quizz_results (
  id INT NOT NULL AUTO_INCREMENT,
  user_id INT NOT NULL,
  quizz_id INT NOT NULL,
  score INT NOT NULL,
  status ENUM('gagal', 'lulus') NOT NULL,
  correct_answer_count INT NOT NULL,
  incorrect_answer_count INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (quizz_id) REFERENCES quizzes(id) ON DELETE CASCADE
) ENGINE = InnoDB;