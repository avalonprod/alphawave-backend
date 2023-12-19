package types

type TasksDTO struct {
	ActiveTasks   []TaskDTO `json:"activeTasks"`
	DeletedTasks  []TaskDTO `json:"deletedTasks"`
	FinishedTasks []TaskDTO `json:"finishedTasks"`
}

type TaskDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Status   string `json:"status"`
	Priority string `json:"priority"`
	Index    int    `json:"index"`
}

type UpdatePositionDTO struct {
	Id string `json:"id"`
}

type TasksCreateDTO struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Index    int    `json:"index"`
}

type UpdateTaskDTO struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
	Index    int    `json:"index"`
}
