package interop

import (
	"github.com/drone/drone-go/drone"
	"github.com/drone/drone/core"
)

func fromStage(from *drone.Stage) *core.Stage {
	return &core.Stage{
		ID: from.ID,
		// RepoID:    from.RepoID,
		BuildID:   from.BuildID,
		Number:    from.Number,
		Name:      from.Name,
		Kind:      from.Kind,
		Type:      from.Type,
		Status:    from.Status,
		Error:     from.Error,
		ErrIgnore: from.ErrIgnore,
		ExitCode:  from.ExitCode,
		Machine:   from.Machine,
		OS:        from.OS,
		Arch:      from.Arch,
		Variant:   from.Variant,
		Kernel:    from.Kernel,
		Limit:     from.Limit,
		Started:   from.Started,
		Stopped:   from.Stopped,
		Created:   from.Created,
		Updated:   from.Updated,
		Version:   from.Version,
		OnSuccess: from.OnSuccess,
		OnFailure: from.OnFailure,
		DependsOn: from.DependsOn,
		Labels:    from.Labels,
		Steps:     fromSteps(from.Steps),
	}
}

func toStage(from *core.Stage) *drone.Stage {
	return &drone.Stage{
		ID: from.ID,
		// RepoID:    from.RepoID,
		BuildID:   from.BuildID,
		Number:    from.Number,
		Name:      from.Name,
		Kind:      from.Kind,
		Type:      from.Type,
		Status:    from.Status,
		Error:     from.Error,
		ErrIgnore: from.ErrIgnore,
		ExitCode:  from.ExitCode,
		Machine:   from.Machine,
		OS:        from.OS,
		Arch:      from.Arch,
		Variant:   from.Variant,
		Kernel:    from.Kernel,
		Limit:     from.Limit,
		Started:   from.Started,
		Stopped:   from.Stopped,
		Created:   from.Created,
		Updated:   from.Updated,
		Version:   from.Version,
		OnSuccess: from.OnSuccess,
		OnFailure: from.OnFailure,
		DependsOn: from.DependsOn,
		Labels:    from.Labels,
		Steps:     toSteps(from.Steps),
	}
}

func toSteps(from []*core.Step) []*drone.Step {
	var dest []*drone.Step
	for _, s := range from {
		dest = append(dest, toStep(s))
	}
	return dest
}

func toStep(from *core.Step) *drone.Step {
	return &drone.Step{
		ID:        from.ID,
		StageID:   from.StageID,
		Number:    from.Number,
		Name:      from.Name,
		Status:    from.Status,
		Error:     from.Error,
		ErrIgnore: from.ErrIgnore,
		ExitCode:  from.ExitCode,
		Started:   from.Started,
		Stopped:   from.Stopped,
		Version:   from.Version,
	}
}

func fromSteps(from []*drone.Step) []*core.Step {
	var dest []*core.Step
	for _, s := range from {
		dest = append(dest, fromStep(s))
	}
	return dest
}

func fromStep(from *drone.Step) *core.Step {
	return &core.Step{
		ID:        from.ID,
		StageID:   from.StageID,
		Number:    from.Number,
		Name:      from.Name,
		Status:    from.Status,
		Error:     from.Error,
		ErrIgnore: from.ErrIgnore,
		ExitCode:  from.ExitCode,
		Started:   from.Started,
		Stopped:   from.Stopped,
		Version:   from.Version,
	}
}

func toRepo(from *core.Repository) *drone.Repo {
	return &drone.Repo{
		ID:         from.ID,
		UID:        from.UID,
		UserID:     from.UserID,
		Namespace:  from.Namespace,
		Name:       from.Name,
		Slug:       from.Slug,
		SCM:        from.SCM,
		HTTPURL:    from.HTTPURL,
		SSHURL:     from.SSHURL,
		Link:       from.Link,
		Branch:     from.Branch,
		Private:    from.Private,
		Visibility: from.Visibility,
		Active:     from.Active,
		Config:     from.Config,
		Trusted:    from.Trusted,
		Protected:  from.Protected,
		Timeout:    from.Timeout,
	}
}

func toBuild(from *core.Build) *drone.Build {
	return &drone.Build{
		ID:           from.ID,
		RepoID:       from.RepoID,
		Trigger:      from.Trigger,
		Number:       from.Number,
		Parent:       from.Parent,
		Status:       from.Status,
		Error:        from.Error,
		Event:        from.Event,
		Action:       from.Action,
		Link:         from.Link,
		Timestamp:    from.Timestamp,
		Title:        from.Title,
		Message:      from.Message,
		Before:       from.Before,
		After:        from.After,
		Ref:          from.Ref,
		Fork:         from.Fork,
		Source:       from.Source,
		Target:       from.Target,
		Author:       from.Author,
		AuthorName:   from.AuthorName,
		AuthorEmail:  from.AuthorEmail,
		AuthorAvatar: from.AuthorAvatar,
		Sender:       from.Sender,
		Params:       from.Params,
		Deploy:       from.Deploy,
		Started:      from.Started,
		Finished:     from.Finished,
		Created:      from.Created,
		Updated:      from.Updated,
		Version:      from.Version,
	}
}

func toNetrc(from *core.Netrc) *drone.Netrc {
	if from == nil {
		return nil
	}
	return &drone.Netrc{
		Machine:  from.Machine,
		Login:    from.Login,
		Password: from.Password,
	}
}

func fromLines(from []*drone.Line) []*core.Line {
	var dest []*core.Line
	for _, s := range from {
		dest = append(dest, fromLine(s))
	}
	return dest
}

func fromLine(from *drone.Line) *core.Line {
	return &core.Line{
		Number:    from.Number,
		Message:   from.Message,
		Timestamp: from.Timestamp,
	}
}
