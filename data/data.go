package data

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/dikaeinstein/ghanalytics/analytics"
)

type Store struct {
	commits []analytics.Commit
	events  []analytics.Event
	repos   []analytics.Repo
	users   []analytics.Actor
}

func NewStore(actorsCSVFile, commitsCSVFile, eventsCSVFile, reposCSVFile io.Reader) (*Store, error) {
	users, err := loadUsers(actorsCSVFile)
	if err != nil {
		return nil, err
	}

	commits, err := loadCommits(commitsCSVFile)
	if err != nil {
		return nil, err
	}

	events, err := loadEvents(eventsCSVFile)
	if err != nil {
		return nil, err
	}

	repos, err := loadRepos(reposCSVFile)
	if err != nil {
		return nil, err
	}

	return &Store{
		commits: commits,
		events:  events,
		repos:   repos,
		users:   users,
	}, nil
}

func loadUsers(csvFile io.Reader) ([]analytics.Actor, error) {
	var users []analytics.Actor
	reader := csv.NewReader(csvFile)

	// skip header
	_, err := reader.Read()
	if err == io.EOF {
		return users, nil
	}
	if err != nil {
		return users, err
	}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return users, err
		}

		if len(line) < 2 {
			return users, fmt.Errorf("Invalid file structure")
		}

		userID, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			return users, err
		}
		users = append(users, analytics.Actor{
			ID:       userID,
			Username: line[1],
		})
	}

	ids := make(map[uint64]bool)
	var dedupedUsers []analytics.Actor
	for _, u := range users {
		if _, ok := ids[u.ID]; !ok {
			ids[u.ID] = true
			dedupedUsers = append(dedupedUsers, u)
		}
	}

	return dedupedUsers, nil
}

func loadCommits(csvFile io.Reader) ([]analytics.Commit, error) {
	var commits []analytics.Commit
	reader := csv.NewReader(csvFile)

	// skip header
	_, err := reader.Read()
	if err == io.EOF {
		return commits, nil
	}
	if err != nil {
		return commits, err
	}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return commits, err
		}

		if len(line) < 3 {
			return commits, fmt.Errorf("Invalid file structure")
		}

		eventID, err := strconv.ParseUint(line[2], 10, 64)
		if err != nil {
			return commits, err
		}
		commits = append(commits, analytics.Commit{
			Sha:     line[0],
			Message: line[1],
			EventID: eventID,
		})
	}

	shas := make(map[string]bool)
	var dedupedCommits []analytics.Commit
	for _, c := range commits {
		if _, ok := shas[c.Sha]; !ok {
			shas[c.Sha] = true
			dedupedCommits = append(dedupedCommits, c)
		}
	}

	return dedupedCommits, nil
}

func loadEvents(csvFile io.Reader) ([]analytics.Event, error) {
	var events []analytics.Event
	reader := csv.NewReader(csvFile)

	// skip header
	_, err := reader.Read()
	if err == io.EOF {
		return events, nil
	}
	if err != nil {
		return events, err
	}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return events, err
		}

		if len(line) < 4 {
			return events, fmt.Errorf("Invalid file structure")
		}

		eventID, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			return events, err
		}
		actorID, err := strconv.ParseUint(line[2], 10, 64)
		if err != nil {
			return events, err
		}
		repoID, err := strconv.ParseUint(line[3], 10, 64)
		if err != nil {
			return events, err
		}
		events = append(events, analytics.Event{
			ID:      eventID,
			Type:    analytics.EventType(line[1]),
			ActorID: actorID,
			RepoID:  repoID,
		})
	}

	ids := make(map[uint64]bool)
	var dedupedEvents []analytics.Event
	for _, e := range events {
		if _, ok := ids[e.ID]; !ok {
			ids[e.ID] = true
			dedupedEvents = append(dedupedEvents, e)
		}
	}

	return dedupedEvents, nil
}

func loadRepos(csvFile io.Reader) ([]analytics.Repo, error) {
	var repos []analytics.Repo
	reader := csv.NewReader(csvFile)

	// skip header
	_, err := reader.Read()
	if err == io.EOF {
		return repos, nil
	}
	if err != nil {
		return repos, err
	}

	for {
		line, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return repos, err
		}

		if len(line) < 2 {
			return repos, fmt.Errorf("Invalid file structure")
		}

		repoID, err := strconv.ParseUint(line[0], 10, 64)
		if err != nil {
			return repos, err
		}
		repos = append(repos, analytics.Repo{
			ID:   repoID,
			Name: line[1],
		})
	}

	ids := make(map[uint64]bool)
	var dedupedRepos []analytics.Repo
	for _, r := range repos {
		if _, ok := ids[r.ID]; !ok {
			ids[r.ID] = true
			dedupedRepos = append(dedupedRepos, r)
		}
	}

	return dedupedRepos, nil
}

func (s *Store) GetUsers(f func(analytics.Actor) bool) ([]analytics.Actor, error) {
	var matchingUsers []analytics.Actor

	for _, u := range s.users {
		if matching := f(u); matching {
			matchingUsers = append(matchingUsers, u)
		}
	}

	return matchingUsers, nil
}

func (s *Store) GetEvents(f func(analytics.Event) bool) ([]analytics.Event, error) {
	var matchingEvents []analytics.Event

	for _, e := range s.events {
		if matching := f(e); matching {
			matchingEvents = append(matchingEvents, e)
		}
	}

	return matchingEvents, nil
}

func (s *Store) GetRepos(f func(analytics.Repo) bool) ([]analytics.Repo, error) {
	var matchingRepos []analytics.Repo

	for _, r := range s.repos {
		if matching := f(r); matching {
			matchingRepos = append(matchingRepos, r)
		}
	}

	return matchingRepos, nil
}
