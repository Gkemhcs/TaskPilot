ALTER TABLE tasks
ADD CONSTRAINT unique_title_project
UNIQUE (project_id, title);