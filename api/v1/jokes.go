package v1

// UserName provides User name
type UserName struct {
	First string `json:"first_name,omitempty"`
	Last  string `json:"last_name,omitempty"`
}

// JokeResponse provides a joke response
type JokeResponse struct {
	ResponseType string   `json:"type"`
	Value        JokeInfo `json:"value"`
}

// JokeInfo provides a joke info
type JokeInfo struct {
	ID         string   `json:"id,omitempty"`
	Joke       string   `json:"joke,omitempty"`
	Categories []string `json:"categories,omitempty"`
}
