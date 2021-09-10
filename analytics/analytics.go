// Package analytics provides functionalities for processing
// Github event data.
package analytics

import "sort"

type Actor struct {
	ID       uint64
	Username string
}

type Event struct {
	ID      uint64
	Type    EventType
	ActorID uint64
	RepoID  uint64
}

type Commit struct {
	Sha     string
	Message string
	EventID uint64
}

type Repo struct {
	ID   uint64
	Name string
}

type Store interface {
	GetUsers(f func(Actor) bool) ([]Actor, error)
	GetEvents(f func(Event) bool) ([]Event, error)
	GetRepos(f func(Repo) bool) ([]Repo, error)
}

//
type ListOptions struct {
	limit         int
	sortCriterion []SortCriteria
}

// Analytics processes Github event data.
type Analytics struct {
	store       Store
	listOptions ListOptions
}

func New(store Store) *Analytics {
	return &Analytics{store: store}
}

type EventType string
type SortCriteria string

const (
	CommitsPushed SortCriteria = "commitsPushed"
	PrCreated     SortCriteria = "prCreated"
)

const (
	PullRequestEvent EventType = "PullRequestEvent"
	PushEvent        EventType = "PushEvent"
	WatchEvent       EventType = "WatchEvent"
)

func Limit(size int) func(*Analytics) error {
	return func(a *Analytics) error {
		return a.setListOptionsLimit(size)
	}
}

func Sort(sortCriterion []SortCriteria) func(*Analytics) error {
	return func(a *Analytics) error {
		return a.setListOptionsSortCriterion(sortCriterion)
	}
}

func (a *Analytics) setListOptionsLimit(size int) error {
	a.listOptions.limit = size
	return nil
}

func (a *Analytics) setListOptionsSortCriterion(sortCriterion []SortCriteria) error {
	a.listOptions.sortCriterion = sortCriterion
	return nil
}

func (a *Analytics) buildList(sortCriterion []SortCriteria) []EventType {
	sortToEventType := map[SortCriteria]EventType{
		CommitsPushed:            PushEvent,
		PrCreated:                PullRequestEvent,
		SortCriteria(WatchEvent): WatchEvent,
	}

	var filterEventTypes []EventType
	for _, c := range sortCriterion {
		filterEventTypes = append(filterEventTypes, sortToEventType[c])
	}

	return filterEventTypes
}

func (a *Analytics) parseListOptions(options []func(*Analytics) error) error {
	for _, option := range options {
		err := option(a)
		if err != nil {
			return err
		}
	}

	return nil
}

type EventsCountByUserID struct {
	UserID      uint64
	EventsCount int
}

func (a *Analytics) ListUsers(options ...func(*Analytics) error) ([]Actor, error) {
	err := a.parseListOptions(options)
	if err != nil {
		return nil, err
	}

	filterEventTypes := a.buildList(a.listOptions.sortCriterion)

	events, err := a.store.GetEvents(func(e Event) bool {
		for _, evt := range filterEventTypes {
			if evt == e.Type {
				return true
			}
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	eventsCollection := make([]Element, len(events))
	for i, e := range events {
		eventsCollection[i] = Element{Value: e}
	}

	eventsByUserID := a.GroupBy(eventsCollection, func(el Element) interface{} {
		evt, _ := el.Value.(Event)
		return evt.ActorID
	})

	// To store the keys in slice in sorted order
	keys := make([]uint64, len(eventsByUserID))
	i := 0
	for k := range eventsByUserID {
		userID, _ := k.(uint64)
		keys[i] = userID
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	listOfEventsCountByUserID := make([]EventsCountByUserID, len(eventsByUserID))
	idx := 0
	for _, k := range keys {
		v := eventsByUserID[k]
		listOfEventsCountByUserID[idx] = EventsCountByUserID{
			UserID:      k,
			EventsCount: len(v),
		}
		idx++
	}

	sort.SliceStable(listOfEventsCountByUserID, func(i, j int) bool {
		return listOfEventsCountByUserID[i].EventsCount > listOfEventsCountByUserID[j].EventsCount
	})

	users, err := a.store.GetUsers(func(a Actor) bool {
		_, ok := eventsByUserID[a.ID]
		return ok
	})
	if err != nil {
		return nil, err
	}

	var sortedUsers []Actor
	for _, v := range listOfEventsCountByUserID {
		for _, u := range users {
			if v.UserID == u.ID {
				sortedUsers = append(sortedUsers, u)
			}
		}
	}

	topNUsers := sortedUsers[0:a.listOptions.limit]
	return topNUsers, nil
}

type EventsCountByRepoID struct {
	RepoID      uint64
	EventsCount int
}

func (a *Analytics) ListRepos(options ...func(*Analytics) error) ([]Repo, error) {
	err := a.parseListOptions(options)
	if err != nil {
		return nil, err
	}

	filterEventTypes := a.buildList(a.listOptions.sortCriterion)
	events, err := a.store.GetEvents(func(e Event) bool {
		for _, evt := range filterEventTypes {
			if evt == e.Type {
				return true
			}
		}
		return false
	})
	if err != nil {
		return nil, err
	}

	eventsCollection := make([]Element, len(events))
	for i, e := range events {
		eventsCollection[i] = Element{Value: e}
	}

	eventsByRepoID := a.GroupBy(eventsCollection, func(el Element) interface{} {
		evt, _ := el.Value.(Event)
		return evt.RepoID
	})

	// To store the keys in slice in sorted order
	keys := make([]uint64, len(eventsByRepoID))
	i := 0
	for k := range eventsByRepoID {
		userID, _ := k.(uint64)
		keys[i] = userID
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	listOfEventsCountByRepoID := make([]EventsCountByRepoID, len(eventsByRepoID))
	idx := 0
	for _, k := range keys {
		v := eventsByRepoID[k]
		listOfEventsCountByRepoID[idx] = EventsCountByRepoID{
			RepoID:      k,
			EventsCount: len(v),
		}
		idx++
	}

	sort.SliceStable(listOfEventsCountByRepoID, func(i, j int) bool {
		return listOfEventsCountByRepoID[i].EventsCount > listOfEventsCountByRepoID[j].EventsCount
	})

	repos, err := a.store.GetRepos(func(r Repo) bool {
		_, ok := eventsByRepoID[r.ID]
		return ok
	})
	if err != nil {
		return nil, err
	}

	var sortedRepos []Repo
	for _, v := range listOfEventsCountByRepoID {
		for _, r := range repos {
			if v.RepoID == r.ID {
				sortedRepos = append(sortedRepos, r)
			}
		}
	}

	topNRepos := sortedRepos[0:a.listOptions.limit]
	return topNRepos, nil
}

type Element struct {
	Value interface{}
}

func (a Analytics) GroupBy(ls []Element, keyGetter func(item Element) interface{}) map[interface{}][]Element {
	grouped := make(map[interface{}][]Element)

	for _, el := range ls {
		key := keyGetter(el)

		_, ok := grouped[key]
		if !ok {
			grouped[key] = []Element{}
			grouped[key] = append(grouped[key], el)
		} else {
			grouped[key] = append(grouped[key], el)
		}
	}

	return grouped
}
