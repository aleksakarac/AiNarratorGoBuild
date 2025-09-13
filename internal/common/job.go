package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusRunning    JobStatus = "running"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCanceled   JobStatus = "canceled"
	JobStatusProcessing JobStatus = "processing"
)

type JobType string

const (
	JobTypeNarration JobType = "narration"
	JobTypeMixing    JobType = "mixing"
)

type Job struct {
	ID        string    `json:"id"`
	Type      JobType   `json:"type"`
	Status    JobStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Input and output paths for audio files
	InputFilePath  string `json:"input_file_path"`
	OutputFilePath string `json:"output_file_path"`
	// Specific fields for narration jobs
	TextContent string `json:"text_content,omitempty"`
	VoiceID     string `json:"voice_id,omitempty"`
	// Specific fields for mixing jobs
	BackgroundAudioFilePath string `json:"background_audio_file_path,omitempty"`
	Volume float64 `json:"volume,omitempty"`
	// Error message if job failed
	Error string `json:"error,omitempty"`
}

func NewJob(jobType JobType, inputPath, outputPath string) *Job {
	return &Job{
		ID:        uuid.New().String(),
		Type:      jobType,
		Status:    JobStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		InputFilePath:  inputPath,
		OutputFilePath: outputPath,
	}
}

func (j *Job) MarshalBinary() ([]byte, error) {
	return json.Marshal(j)
}

func (j *Job) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, j)
}

type JobUpdate struct {
	ID     string    `json:"id"`
	Status JobStatus `json:"status,omitempty"`
	Error  string    `json:"error,omitempty"`
}

func (ju *JobUpdate) String() string {
	return fmt.Sprintf("JobUpdate{ID: %s, Status: %s, Error: %s}", ju.ID, ju.Status, ju.Error)
}

// Event types for job updates
const (
	JobEventCreated  = "job.created"
	JobEventUpdated  = "job.updated"
	JobEventProgress = "job.progress"
)

type JobEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewJobEvent(eventType string, payload interface{}) (JobEvent, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return JobEvent{}, fmt.Errorf("failed to marshal event payload: %w", err)
	}
	return JobEvent{
		Type:    eventType,
		Payload: data,
	},
}

func (je *JobEvent) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(je.Payload, v)
}

// WorkerPoolStatus represents the status of the worker pool
type WorkerPoolStatus struct {
	TotalWorkers    int `json:"total_workers"`
	IdleWorkers     int `json:"idle_workers"`
	ProcessingJobs  int `json:"processing_jobs"`
	QueuedJobs      int `json:"queued_jobs"`
	FailedJobsToday int `json:"failed_jobs_today"`
}

// WorkerStatus represents the status of an individual worker
type WorkerStatus struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"` // e.g., "idle", "busy", "error"
	CurrentJob string    `json:"current_job,omitempty"`
	LastActive time.Time `json:"last_active"`
}

// SystemStatus combines overall system health with worker pool and individual worker statuses
type SystemStatus struct {
	OverallHealth string             `json:"overall_health"` // e.g., "healthy", "degraded", "unhealthy"
	WorkerPool    WorkerPoolStatus   `json:"worker_pool"`
	Workers       []WorkerStatus     `json:"workers"`
	LastUpdated   time.Time          `json:"last_updated"`
	Messages      []string           `json:"messages,omitempty"` // General messages or alerts
}

// NewSystemStatus creates a new SystemStatus with default values
func NewSystemStatus() SystemStatus {
	return SystemStatus{
		OverallHealth: "healthy",
		WorkerPool: WorkerPoolStatus{
			TotalWorkers:    0,
			IdleWorkers:     0,
			ProcessingJobs:  0,
			QueuedJobs:      0,
			FailedJobsToday: 0,
		},
		Workers:     []WorkerStatus{},
		LastUpdated: time.Now(),
		Messages:    []string{},
	}
}

// AddMessage adds a message to the SystemStatus
func (ss *SystemStatus) AddMessage(msg string) {
	ss.Messages = append(ss.Messages, msg)
}

// UpdateWorkerPoolStatus updates the worker pool metrics
func (ss *SystemStatus) UpdateWorkerPoolStatus(total, idle, processing, queued, failed int) {
	ss.WorkerPool.TotalWorkers = total
	ss.WorkerPool.IdleWorkers = idle
	ss.WorkerPool.ProcessingJobs = processing
	ss.WorkerPool.QueuedJobs = queued
	ss.WorkerPool.FailedJobsToday = failed
}

// UpdateWorkerStatus updates or adds an individual worker\'s status
func (ss *SystemStatus) UpdateWorkerStatus(workerID, status, currentJob string) {
	found := false
	for i, worker := range ss.Workers {
		if worker.ID == workerID {
			ss.Workers[i].Status = status
			ss.Workers[i].CurrentJob = currentJob
			ss.Workers[i].LastActive = time.Now()
			found = true
			break
		}
	}
	if !found {
		ss.Workers = append(ss.Workers, WorkerStatus{
			ID:         workerID,
			Status:     status,
			CurrentJob: currentJob,
			LastActive: time.Now(),
		})
	}
}

// RemoveWorkerStatus removes a worker from the status list
func (ss *SystemStatus) RemoveWorkerStatus(workerID string) {
	for i, worker := range ss.Workers {
		if worker.ID == workerID {
			ss.Workers = append(ss.Workers[:i], ss.Workers[i+1:]...)
			break
		}
	}
}

// UpdateOverallHealth sets the overall health status
func (ss *SystemStatus) UpdateOverallHealth(health string) {
	ss.OverallHealth = health
}

// RefreshLastUpdated sets the LastUpdated timestamp to the current time
func (ss *SystemStatus) RefreshLastUpdated() {
	ss.LastUpdated = time.Now()
}
