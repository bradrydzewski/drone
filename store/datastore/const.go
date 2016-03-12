package datastore

const (
	userTable  = "users"
	repoTable  = "repos"
	buildTable = "builds"
	jobTable   = "jobs"
	keyTable   = "keys"
	logTable   = "logs"
)

const userLoginQuery = `
SELECT *
FROM users
WHERE user_login=?
LIMIT 1
`

const userListQuery = `
SELECT *
FROM users
ORDER BY user_login ASC
`

const userCountQuery = `
SELECT count(1)
FROM users
`

const userDeleteStmt = `
DELETE FROM users
WHERE user_id=?
`

const userFeedQuery = `
SELECT
 repo_owner
,repo_name
,repo_full_name
,build_number
,build_event
,build_status
,build_created
,build_started
,build_finished
,build_commit
,build_branch
,build_ref
,build_refspec
,build_remote
,build_title
,build_message
,build_author
,build_email
,build_avatar
FROM
 builds b
,repos r
WHERE b.build_repo_id = r.repo_id
  AND r.repo_full_name IN (%s)
ORDER BY b.build_id DESC
LIMIT 25
`

// repos

const repoNameQuery = `
SELECT *
FROM repos
WHERE repo_full_name = ?
LIMIT 1;
`

const repoListQuery = `
SELECT *
FROM repos
WHERE repo_id IN (
	SELECT DISTINCT build_repo_id
	FROM builds
	WHERE build_author = ?
)
ORDER BY repo_full_name
`

const repoListOfQuery = `
SELECT *
FROM repos
WHERE repo_full_name IN (%s)
ORDER BY repo_name
`

const repoCountQuery = `
SELECT COUNT(*) FROM repos
`

const repoDeleteStmt = `
DELETE FROM repos
WHERE repo_id = ?
`

// keys

const keyQuery = "SELECT * FROM `keys` WHERE key_repo_id=? LIMIT 1"

const keyDeleteStmt = "DELETE FROM `keys` WHERE key_id=?"

// builds

const buildListQuery = `
SELECT *
FROM builds
WHERE build_repo_id = ?
ORDER BY build_number DESC
LIMIT 50
`

const buildNumberQuery = `
SELECT *
FROM builds
WHERE build_repo_id = ?
  AND build_number = ?
LIMIT 1;
`

const buildLastQuery = `
SELECT *
FROM builds
WHERE build_repo_id = ?
  AND build_branch  = ?
  AND build_event   = 'push'
ORDER BY build_number DESC
LIMIT 1
`

const buildLastBeforeQuery = `
SELECT *
FROM builds
WHERE build_repo_id = ?
  AND build_branch  = ?
  AND build_id < ?
ORDER BY build_number DESC
LIMIT 1
`

const buildCommitQuery = `
SELECT *
FROM builds
WHERE build_repo_id = ?
  AND build_commit  = ?
  AND build_branch  = ?
LIMIT 1
`

const buildRefQuery = `
SELECT *
FROM builds
WHERE build_repo_id = ?
  AND build_ref     = ?
LIMIT 1
`

const buildNumberLast = `
SELECT MAX(build_number)
FROM builds
WHERE build_repo_id = ?
`

// jobs

const jobListQuery = `
SELECT *
FROM jobs
WHERE job_build_id = ?
ORDER BY job_number ASC
`

const jobNumberQuery = `
SELECT *
FROM jobs
WHERE job_build_id = ?
AND   job_number = ?
LIMIT 1
`

// logs

const logQuery = `
SELECT *
FROM logs
WHERE log_job_id=?
LIMIT 1
`
