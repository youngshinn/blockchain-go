package global

var from string

func FROM() string {
	if from == "" {
		return ""
	}
	return from
}

func SetFrom(f string) {
	from = f
}
