package codeforeces

const (
	StatusUrl        = "https://codeforces.com/api/user.status"
	ContestStatusUrl = "https://codeforces.com/api/contest.status"
)

type ResultList struct {
	Status string
	Result []*StatuResult
}

type StatuResult struct {
	Id                  int64
	ContestId           int
	CreationTimeSeconds int64
	RelativeTimeSeconds int64
	Problem             *Problem
	Author              *Author
	ProgrammingLanguage string
	Verdict             string
	Testset             string
	PassedTestCount     int
	TimeConsumedMillis  int
	MemoryConsumedBytes int64
}

type Problem struct {
	ContestId int
	Index     string
	Name      string
	Type      string
	rating    int
	Points    float64
	Tags      []string
}

type Author struct {
	ContestId        int
	Members          []*Members
	ParticipantType  string
	Ghost            bool
	StartTimeSeconds int64
}

type Members struct {
	Handle string
}
