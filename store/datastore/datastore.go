package datastore

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/drone/drone/model"
	"github.com/russross/meddler"
)

// datastore is an implementation of a model.Store built on top
// of the sql/database driver with a relational database backend.
type datastore struct {
	*sql.DB
}

func (db *datastore) GetUser(id int64) (*model.User, error) {
	var usr = new(model.User)
	var err = meddler.Load(db, userTable, usr, id)
	return usr, err
}

func (db *datastore) GetUserLogin(login string) (*model.User, error) {
	var usr = new(model.User)
	var err = meddler.QueryRow(db, usr, rebind(userLoginQuery), login)
	return usr, err
}

func (db *datastore) GetUserList() ([]*model.User, error) {
	var users = []*model.User{}
	var err = meddler.QueryAll(db, &users, rebind(userListQuery))
	return users, err
}

func (db *datastore) GetUserFeed(listof []*model.RepoLite) ([]*model.Feed, error) {
	var (
		feed []*model.Feed
		args []interface{}
		stmt string
	)
	switch meddler.Default {
	case meddler.PostgreSQL:
		stmt, args = toListPosgres(listof)
	default:
		stmt, args = toList(listof)
	}
	err := meddler.QueryAll(db, &feed, fmt.Sprintf(userFeedQuery, stmt), args...)
	return feed, err
}

func (db *datastore) GetUserCount() (int, error) {
	var count int
	var err = db.QueryRow(rebind(userCountQuery)).Scan(&count)
	return count, err
}

func (db *datastore) CreateUser(user *model.User) error {
	return meddler.Insert(db, userTable, user)
}

func (db *datastore) UpdateUser(user *model.User) error {
	return meddler.Update(db, userTable, user)
}

func (db *datastore) DeleteUser(user *model.User) error {
	var _, err = db.Exec(rebind(userDeleteStmt), user.ID)
	return err
}

func (db *datastore) GetRepo(id int64) (*model.Repo, error) {
	var repo = new(model.Repo)
	var err = meddler.Load(db, repoTable, repo, id)
	return repo, err
}

func (db *datastore) GetRepoName(name string) (*model.Repo, error) {
	var repo = new(model.Repo)
	var err = meddler.QueryRow(db, repo, rebind(repoNameQuery), name)
	return repo, err
}

func (db *datastore) GetRepoListOf(listof []*model.RepoLite) ([]*model.Repo, error) {
	var (
		repos []*model.Repo
		args  []interface{}
		stmt  string
	)
	switch meddler.Default {
	case meddler.PostgreSQL:
		stmt, args = toListPosgres(listof)
	default:
		stmt, args = toList(listof)
	}
	err := meddler.QueryAll(db, &repos, fmt.Sprintf(repoListOfQuery, stmt), args...)
	return repos, err
}

func (db *datastore) GetRepoCount() (int, error) {
	var count int
	var err = db.QueryRow(rebind(repoCountQuery)).Scan(&count)
	return count, err
}

func (db *datastore) CreateRepo(repo *model.Repo) error {
	return meddler.Insert(db, repoTable, repo)
}

func (db *datastore) UpdateRepo(repo *model.Repo) error {
	return meddler.Update(db, repoTable, repo)
}

func (db *datastore) DeleteRepo(repo *model.Repo) error {
	var _, err = db.Exec(rebind(repoDeleteStmt), repo.ID)
	return err
}

func (db *datastore) GetKey(repo *model.Repo) (*model.Key, error) {
	var key = new(model.Key)
	var err = meddler.QueryRow(db, key, rebind(keyQuery), repo.ID)
	return key, err
}

func (db *datastore) CreateKey(key *model.Key) error {
	return meddler.Save(db, keyTable, key)
}

func (db *datastore) UpdateKey(key *model.Key) error {
	return meddler.Save(db, keyTable, key)
}

func (db *datastore) DeleteKey(key *model.Key) error {
	var _, err = db.Exec(rebind(keyDeleteStmt), key.ID)
	return err
}

func (db *datastore) GetBuild(id int64) (*model.Build, error) {
	var build = new(model.Build)
	var err = meddler.Load(db, buildTable, build, id)
	return build, err
}

func (db *datastore) GetBuildNumber(repo *model.Repo, num int) (*model.Build, error) {
	var build = new(model.Build)
	var err = meddler.QueryRow(db, build, rebind(buildNumberQuery), repo.ID, num)
	return build, err
}

func (db *datastore) GetBuildRef(repo *model.Repo, ref string) (*model.Build, error) {
	var build = new(model.Build)
	var err = meddler.QueryRow(db, build, rebind(buildRefQuery), repo.ID, ref)
	return build, err
}

func (db *datastore) GetBuildCommit(repo *model.Repo, sha, branch string) (*model.Build, error) {
	var build = new(model.Build)
	var err = meddler.QueryRow(db, build, rebind(buildCommitQuery), repo.ID, sha, branch)
	return build, err
}

func (db *datastore) GetBuildLast(repo *model.Repo, branch string) (*model.Build, error) {
	var build = new(model.Build)
	var err = meddler.QueryRow(db, build, rebind(buildLastQuery), repo.ID, branch)
	return build, err
}

func (db *datastore) GetBuildLastBefore(repo *model.Repo, branch string, num int64) (*model.Build, error) {
	var build = new(model.Build)
	var err = meddler.QueryRow(db, build, rebind(buildLastBeforeQuery), repo.ID, branch, num)
	return build, err
}

func (db *datastore) GetBuildList(repo *model.Repo) ([]*model.Build, error) {
	var builds = []*model.Build{}
	var err = meddler.QueryAll(db, &builds, rebind(buildListQuery), repo.ID)
	return builds, err
}

func (db *datastore) CreateBuild(build *model.Build, jobs ...*model.Job) error {
	var number int
	db.QueryRow(rebind(buildNumberLast), build.RepoID).Scan(&number)
	build.Number = number + 1
	build.Created = time.Now().UTC().Unix()
	build.Enqueued = build.Created
	err := meddler.Insert(db, buildTable, build)
	if err != nil {
		return err
	}
	for i, job := range jobs {
		job.BuildID = build.ID
		job.Number = i + 1
		job.Enqueued = build.Created
		err = meddler.Insert(db, jobTable, job)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *datastore) UpdateBuild(build *model.Build) error {
	return meddler.Update(db, buildTable, build)
}

func (db *datastore) GetJob(id int64) (*model.Job, error) {
	var job = new(model.Job)
	var err = meddler.Load(db, jobTable, job, id)
	return job, err
}

func (db *datastore) GetJobNumber(build *model.Build, num int) (*model.Job, error) {
	var job = new(model.Job)
	var err = meddler.QueryRow(db, job, rebind(jobNumberQuery), build.ID, num)
	return job, err
}

func (db *datastore) GetJobList(build *model.Build) ([]*model.Job, error) {
	var jobs = []*model.Job{}
	var err = meddler.QueryAll(db, &jobs, rebind(jobListQuery), build.ID)
	return jobs, err
}

func (db *datastore) CreateJob(job *model.Job) error {
	return meddler.Insert(db, jobTable, job)
}

func (db *datastore) UpdateJob(job *model.Job) error {
	return meddler.Update(db, jobTable, job)
}

func (db *datastore) ReadLog(job *model.Job) (io.ReadCloser, error) {
	var log = new(model.Log)
	var err = meddler.QueryRow(db, log, rebind(logQuery), job.ID)
	var buf = bytes.NewBuffer(log.Data)
	return ioutil.NopCloser(buf), err
}

func (db *datastore) WriteLog(job *model.Job, r io.Reader) error {
	var log = new(model.Log)
	var err = meddler.QueryRow(db, log, rebind(logQuery), job.ID)
	if err != nil {
		log = &model.Log{JobID: job.ID}
	}
	log.Data, _ = ioutil.ReadAll(r)
	return meddler.Save(db, logTable, log)
}
