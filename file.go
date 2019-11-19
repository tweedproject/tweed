package tweed

type File struct {
	Filename    string `json:"filename"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Contents    string `json:"contents"`
}
