package auth

type Request struct {
	AuthID           string `json:"auth_id"`
	AuthExpires      int    `json:"auth_expires"`
	RefreshID        string `json:"refresh_id"`
	MemberID         string `json:"member_id"`
	Status           string `json:"status"`
	Placement        string `json:"placement"`
	PlacementOptions string `json:"placement_options"`
}
