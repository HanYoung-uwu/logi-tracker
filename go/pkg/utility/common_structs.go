package utility

type Name struct {
	Name string `json:"name"`
}

// 0 is super admin, 1 is clan admin, 2 is ordinary memenber, 3 is temporary account for invitation links,
// 4 is clan admin invitation links

type Account struct {
	Name       string
	Permission int
	Clan       string
}
