package internal

type SiteInfo struct {
	ID            int    `json:"id"`
	URL           string `json:"url"`
	CheckInterval int    `json:"check_interval"`
}

type UserSites struct {
	UserID int        `json:"user_id"`
	Sites  []SiteInfo `json:"sites"`
}
