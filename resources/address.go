package resources

type Address struct {
	Name    string `json:"name"`
	Line1   string `json:"line_1"`
	Line2   string `json:"line_2"`
	Line3   string `json:"line_3"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
	Phone   string `json:"phone_number"`
}
