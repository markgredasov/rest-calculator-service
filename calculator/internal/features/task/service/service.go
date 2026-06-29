package tasks_service

type TasksService struct {
	tasksRepository TasksRepository
}

type TasksRepository interface{}

func NewTasksService(repo TasksRepository) TasksService {
	return TasksService{
		tasksRepository: repo,
	}
}
