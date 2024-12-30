package handler

type loginPayload struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type settingsPayload struct {
	// Number of URLs per hour
	URLsPerHour int `form:"urlsPerHour"`
	// Enable/disable searching with the crawlers
	SearchOn bool `form:"searchOn"`
	// Add new URLs
	AddNewURLs bool `form:"addNewUrls"`
}

type searchPayload struct {
	Query string `form:"query"`
}
