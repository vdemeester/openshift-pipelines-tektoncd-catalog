package config

type Config struct {
	Tasks []Task
}

type Task struct {
	Repository string
}
