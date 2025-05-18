package ovh

type aliasesPost struct {
	From      string `json:"from"`
	LocalCopy bool   `json:"localCopy"`
	To        string `json:"to"`
}
