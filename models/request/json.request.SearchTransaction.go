package request

type SearchTransaction struct {
	Line        string `json:"line"`
	Amount      string `json:"amount"`
	IdTRN       string `json:"IdTRN"`
	IdTRNClient string `json:"IdTRNClient"`
}
