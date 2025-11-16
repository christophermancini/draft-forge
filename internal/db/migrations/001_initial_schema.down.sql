-- Drop triggers
DROP TRIGGER IF EXISTS update_project_stats_updated_at ON project_stats;
DROP TRIGGER IF EXISTS update_ai_credits_updated_at ON ai_credits;
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order
DROP TABLE IF EXISTS project_stats;
DROP TABLE IF EXISTS agent_runs;
DROP TABLE IF EXISTS ai_usage_log;
DROP TABLE IF EXISTS ai_credits;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
