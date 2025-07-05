ALTER TABLE projects
ADD CONSTRAINT unique_user_project_name UNIQUE (name);