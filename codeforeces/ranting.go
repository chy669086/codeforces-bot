package codeforeces

const RatingUrl = "https://codeforces.com/api/user.rating"

type RantingResultList struct {
	Status string
	Result []*RantingResult
}

type RantingResult struct {
	ContestId               int
	ContestName             string
	Handle                  string
	Rank                    int
	RatingUpdateTimeSeconds int64
	OldRating               int
	NewRating               int
}
