package http

import (
	"fmt"

	"github.com/openshift/geard/containers"
	"github.com/openshift/geard/git"
	"github.com/openshift/geard/git/http/remote"
	gitjobs "github.com/openshift/geard/git/jobs"
	"github.com/openshift/geard/http"
	"github.com/openshift/geard/http/client"
	"github.com/openshift/geard/jobs"
	"github.com/openshift/go-json-rest"
)

type HttpExtension struct{}

func (h *HttpExtension) Routes() http.ExtensionMap {
	return http.ExtensionMap{
		&remote.HttpCreateRepositoryRequest{}: HandleCreateRepositoryRequest,
		&remote.HttpGitArchiveContentRequest{
			GitArchiveContentRequest: gitjobs.GitArchiveContentRequest{Ref: "*"},
		}: HandleGitArchiveContentRequest,
	}
}

func (h *HttpExtension) HttpJobFor(job interface{}) (exc client.RemoteExecutable, err error) {
	return remote.HttpJobFor(job)
}

func HandleCreateRepositoryRequest(conf *http.HttpConfiguration, context *http.HttpContext, r *rest.Request) (interface{}, error) {
	repositoryId, errg := containers.NewIdentifier(r.PathParam("id"))
	if errg != nil {
		return nil, errg
	}
	// TODO: convert token into a safe clone spec and commit hash
	return &gitjobs.CreateRepositoryRequest{
		git.RepoIdentifier(repositoryId),
		r.URL.Query().Get("source"),
		context.Id,
	}, nil
}

func HandleGitArchiveContentRequest(conf *http.HttpConfiguration, context *http.HttpContext, r *rest.Request) (interface{}, error) {
	repoId, errr := containers.NewIdentifier(r.PathParam("id"))
	if errr != nil {
		return nil, jobs.SimpleError{jobs.ResponseInvalidRequest, fmt.Sprintf("Invalid repository identifier: %s", errr.Error())}
	}
	ref, errc := gitjobs.NewGitCommitRef(r.PathParam("ref"))
	if errc != nil {
		return nil, jobs.SimpleError{jobs.ResponseInvalidRequest, fmt.Sprintf("Invalid commit ref: %s", errc.Error())}
	}

	return &gitjobs.GitArchiveContentRequest{
		git.RepoIdentifier(repoId),
		ref,
	}, nil
}
