package event

type QueueManger interface {
	Add(taskId string, queue *Queue)
	Get(taskId string) *Queue
	Tap(taskId string) *Queue
	Close(taskId string)
	CreateOrTap(taskId string) *Queue
}
