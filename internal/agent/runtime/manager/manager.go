package manager

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"ospnet/internal/agent/db"
	"ospnet/internal/agent/runtime/docker"
)

type StartRequest struct {
	Image string
	Name  string
	Port  int
}

type Manager struct {
	db     *gorm.DB
	docker *docker.Client
	logger *log.Logger
}

func New(db *gorm.DB, dockerClient *docker.Client, logger *log.Logger) *Manager {
	return &Manager{db: db, docker: dockerClient, logger: logger}
}

func (m *Manager) StartContainer(ctx context.Context, req StartRequest) (db.Containers, error) {
	if strings.TrimSpace(req.Image) == "" || strings.TrimSpace(req.Name) == "" || req.Port <= 0 {
		return db.Containers{}, fmt.Errorf("image, name and port are required")
	}

	record := db.Containers{
		ID:     uuid.NewString(),
		Image:  req.Image,
		Name:   req.Name,
		Port:   req.Port,
		Status: "pending",
	}

	err := gorm.G[db.Containers](m.db).Create(ctx, &record)

	if err != nil {
		return db.Containers{}, fmt.Errorf("failed to create container record: %w", err)
	}

	containerID, err := m.docker.RunContainer(ctx, docker.RunOptions{
		Image: req.Image,
		Name:  req.Name,
		Port:  req.Port,
	})

	if err != nil {
		_, err := gorm.G[db.Containers](m.db).Where("id = ?", containerID).Update(ctx, "Status", "error")

		return db.Containers{}, fmt.Errorf("failed to start container: %w", err)
	}

	_, err = gorm.G[db.Containers](m.db).Where("id = ?", containerID).Update(ctx, "Status", "running")

	if err != nil {
		return db.Containers{}, fmt.Errorf("container started but status update failed: %w", err)
	}

	record.DockerID = containerID
	record.Status = "running"
	return record, nil
}

func (m *Manager) StopContainer(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("container name is required")
	}

	record, err := gorm.G[db.Containers](m.db).Where("name = ?", name).First(ctx)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to lookup container %s: %w", name, err)
	}

	targetID := record.DockerID
	if targetID == "" {
		items, listErr := m.docker.FindContainersByName(ctx, name)
		if listErr != nil {
			return fmt.Errorf("failed to find docker container by name: %w", listErr)
		}
		if len(items) == 0 {
			return fmt.Errorf("container %s not found", name)
		}
		targetID = items[0].ID
	}

	if err := m.docker.StopContainer(ctx, targetID); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", name, err)
	}

	if record.ID != "" {
		_, err = gorm.G[db.Containers](m.db).Where("id = ?", record.ID).Update(ctx, "Status", "stopped")

		if err != nil {
			m.logger.Printf("failed to update db status for %s: %v", name, err)
		}
	}

	return nil
}

func (m *Manager) GetContainers(ctx context.Context) ([]db.Containers, error) {
	containers, err := gorm.G[db.Containers](m.db).Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}
	return containers, nil
}

func (m *Manager) Reconcile(ctx context.Context) error {
	desired, err := m.GetContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to read desired containers: %w", err)
	}

	actual, err := m.docker.ListContainers(ctx)
	if err != nil {
		return fmt.Errorf("failed to read docker containers: %w", err)
	}

	actualByName := map[string][]docker.Container{}
	for _, container := range actual {
		if container.Name == "" {
			continue
		}
		actualByName[container.Name] = append(actualByName[container.Name], container)
	}

	for _, expected := range desired {
		if expected.Status != "running" {
			continue
		}

		matches := actualByName[expected.Name]
		runningMatches := make([]docker.Container, 0, len(matches))
		for _, item := range matches {
			if item.Running {
				runningMatches = append(runningMatches, item)
			}
		}

		if len(runningMatches) > 1 {
			for i := 1; i < len(runningMatches); i++ {
				if stopErr := m.docker.StopContainer(ctx, runningMatches[i].ID); stopErr != nil {
					m.logger.Printf("failed to stop duplicate container %s (%s): %v", expected.Name, runningMatches[i].ID, stopErr)
				}
			}
		}

		if len(runningMatches) >= 1 {
			running := runningMatches[0]
			if expected.DockerID != running.ID {
				if updateErr := m.UpdateStatus(ctx, expected.ID, "running"); updateErr != nil {
					m.logger.Printf("failed to update running container id for %s: %v", expected.Name, updateErr)
				}
			}
			continue
		}

		for _, existing := range matches {
			if stopErr := m.docker.StopContainer(ctx, existing.ID); stopErr != nil {
				m.logger.Printf("failed to remove stale container %s (%s): %v", expected.Name, existing.ID, stopErr)
			}
		}

		_, startErr := m.docker.RunContainer(ctx, docker.RunOptions{
			Image: expected.Image,
			Name:  expected.Name,
			Port:  expected.Port,
		})
		if startErr != nil {
			if updateErr := m.UpdateStatus(ctx, expected.ID, "error"); updateErr != nil {
				m.logger.Printf("failed to set error status for %s: %v", expected.Name, updateErr)
			}
			m.logger.Printf("failed to restart missing container %s: %v", expected.Name, startErr)
			continue
		}

		if updateErr := m.UpdateStatus(ctx, expected.ID, "running"); updateErr != nil {
			m.logger.Printf("failed to persist restarted container %s: %v", expected.Name, updateErr)
		}
	}

	return nil
}

func (m *Manager) UpdateStatus(ctx context.Context, id string, status string) error {
	_, err := gorm.G[db.Containers](m.db).Where("id = ?", id).Update(ctx, "Status", status)
	return err
}
