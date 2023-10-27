package watchdog

type watchdogCheckModel struct {
	Url      string
	FilePath string
}

var SupportedServices = map[string]watchdogCheckModel{
	"server": watchdogCheckModel{
		Url:      "https://localhost:8080/watchdog/start",
		FilePath: "server_live",
	},
}
