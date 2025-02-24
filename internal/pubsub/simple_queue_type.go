package pubsub

type SimpleQueueType int

const (
	Durable   SimpleQueueType = 0
	Transient SimpleQueueType = 1
)
