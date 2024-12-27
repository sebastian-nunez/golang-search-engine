package handler

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type settingsForm struct {
	// Number of URLs per hour
	URLsPerHour int `form:"urlsPerHour"`
	// Enable/disable searching with the crawlers
	SearchOn bool `form:"searchOn"`
	// Add new URLs
	AddNewURLs bool `form:"addNewUrls"`
}
