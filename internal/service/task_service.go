package service

import (
	"context"
	"encoding/json"
	"errors"

	"example.com/pz9-redis-cache/internal/cache"
	"example.com/pz9-redis-cache/internal/config"
	"example.com/pz9-redis-cache/internal/task"
	"github.com/redis/go-redis/v9"
)

type TaskService struct {
	repo  *task.Repo
	redis *redis.Client
	cfg   config.Config
}

func NewTaskService(repo *task.Repo, redisClient *redis.Client, cfg config.Config) *TaskService {
	return &TaskService{repo: repo, redis: redisClient, cfg: cfg}
}

func (s *TaskService) ListTasks(ctx context.Context) ([]task.Task, error) {
	key := cache.TasksListKey()
	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var tasks []task.Task
		if json.Unmarshal([]byte(cached), &tasks) == nil {
			return tasks, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		// Redis is optional for this educational service.
	}

	tasks := s.repo.List()
	if data, err := json.Marshal(tasks); err == nil {
		ttl := cache.TTLWithJitter(s.cfg.TaskCacheTTL, s.cfg.TaskCacheJitter)
		_ = s.redis.Set(ctx, key, data, ttl).Err()
	}
	return tasks, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id int64) (task.Task, error) {
	key := cache.TaskByIDKey(id)
	cached, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		var t task.Task
		if json.Unmarshal([]byte(cached), &t) == nil {
			return t, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		// Redis is an optimization in this practice; repository remains source of truth.
	}

	t, err := s.repo.GetByID(id)
	if err != nil {
		return task.Task{}, err
	}

	if data, err := json.Marshal(t); err == nil {
		ttl := cache.TTLWithJitter(s.cfg.TaskCacheTTL, s.cfg.TaskCacheJitter)
		_ = s.redis.Set(ctx, key, data, ttl).Err()
	}
	return t, nil
}

func (s *TaskService) CreateTask(ctx context.Context, title string, description string) (task.Task, error) {
	t := s.repo.Create(title, description)
	_ = s.redis.Del(ctx, cache.TasksListKey()).Err()
	return t, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, t task.Task) error {
	if err := s.repo.Update(t); err != nil {
		return err
	}
	_ = s.redis.Del(ctx, cache.TaskByIDKey(t.ID), cache.TasksListKey()).Err()
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	_ = s.redis.Del(ctx, cache.TaskByIDKey(id), cache.TasksListKey()).Err()
	return nil
}
