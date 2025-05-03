CREATE TABLE refresh_tokens (
  user_id INT NOT NULL,
  token TEXT NOT NULL,
  status_logout TINYINT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);