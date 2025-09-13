package common

import (
	"time"
)

type WorkerConfig struct {
	ID string `json:"id"`
	// Add other worker-specific configuration fields here
}

// Worker represents a worker instance
type Worker interface {
	ID() string
	Run() error
	Stop()
	Status() WorkerStatus
	// Add other worker-specific methods here
}

// WorkerPool manages a pool of workers
type WorkerPool interface {
	AddWorker(worker Worker)
	RemoveWorker(workerID string)
	Start()
	Stop()
	SubmitJob(job *Job) error
	Status() WorkerPoolStatus
	// Add methods for job management, e.g., GetJobStatus, CancelJob
}

// JobQueue defines the interface for a job queuing system
type JobQueue interface {
	Enqueue(job *Job) error
	Dequeue() (*Job, error)
	Peek() (*Job, error)
	Size() (int64, error)
	// Add methods for job status updates, error handling, etc.
}

// EventPublisher defines the interface for publishing system events
type EventPublisher interface {
	Publish(event JobEvent) error
}

// EventSubscriber defines the interface for subscribing to system events
type EventSubscriber interface {
	Subscribe(eventType string, handler func(event JobEvent))
	Unsubscribe(eventType string, handler func(event JobEvent))
}

// Database defines the interface for database operations
type Database interface {
	SaveJob(job *Job) error
	GetJob(id string) (*Job, error)
	UpdateJob(job *Job) error
	DeleteJob(id string) error
	ListJobs(filter JobStatus, limit, offset int) ([]*Job, error)
	// Add methods for worker status, system status, etc.
}

// Config defines the interface for application configuration
type Config interface {
	Get(key string) (string, error)
	Set(key, value string) error
	// Add methods for loading/saving config, etc.
}

// NarratorService defines the interface for text-to-speech narration
type NarratorService interface {
	Narrate(text, voiceID, outputPath string) (*Job, error)
	// Add methods for listing available voices, etc.
}

// AudioMixerService defines the interface for audio mixing operations
type AudioMixerService interface {
	Mix(inputPath, backgroundPath, outputPath string, volume float64) (*Job, error)
	// Add methods for audio effects, etc.
}

// FileStorageService defines the interface for file storage operations
type FileStorageService interface {
	Upload(localFilePath, remotePath string) (string, error)
	Download(remotePath, localFilePath string) error
	Delete(remotePath string) error
	GetURL(remotePath string) (string, error)
}

// NewWorkerStatus creates a new WorkerStatus instance
func NewWorkerStatus(id string, status string) WorkerStatus {
	return WorkerStatus{
		ID:         id,
		Status:     status,
		LastActive: time.Now(),
	}
}
