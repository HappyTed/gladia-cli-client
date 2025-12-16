package entities

type TaskStatus string

const (
	DONE        TaskStatus = "done"
	IN_PROGRESS TaskStatus = "in progress"
	AWAITS      TaskStatus = "awaits"
	ERROR       TaskStatus = "error"
)
