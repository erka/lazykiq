package sidekiq

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SortedEntry represents a job stored in a Sidekiq sorted set (dead, retry, schedule).
// It embeds a JobRecord for the job data and adds the sorted set score (timestamp).
type SortedEntry struct {
	*JobRecord        // the actual job data (embedded for method promotion)
	Score      float64 // sorted set score (timestamp)
}

// NewSortedEntry creates a SortedEntry from raw JSON data and score.
func NewSortedEntry(value string, score float64) *SortedEntry {
	return &SortedEntry{
		JobRecord: NewJobRecord(value, ""),
		Score:     score,
	}
}

// At returns the timestamp as Unix seconds (same as score for dead/retry/schedule).
func (se *SortedEntry) At() int64 {
	return int64(se.Score)
}

// getSortedSetJobs fetches jobs from a sorted set with pagination.
// If reverse is true, returns highest scores first (ZREVRANGE), otherwise lowest first (ZRANGE).
func (c *Client) getSortedSetJobs(ctx context.Context, key string, start, count int, reverse bool) ([]*SortedEntry, int64, error) {
	size, err := c.redis.ZCard(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return nil, 0, err
	}

	if size == 0 {
		return nil, 0, nil
	}

	end := int64(start + count - 1)
	var results []redis.Z
	if reverse {
		results, err = c.redis.ZRevRangeWithScores(ctx, key, int64(start), end).Result()
	} else {
		results, err = c.redis.ZRangeWithScores(ctx, key, int64(start), end).Result()
	}
	if err != nil {
		return nil, size, err
	}

	entries := make([]*SortedEntry, len(results))
	for i, z := range results {
		value, _ := z.Member.(string)
		entries[i] = NewSortedEntry(value, z.Score)
	}

	return entries, size, nil
}

// GetDeadJobs fetches dead jobs with pagination (newest first).
func (c *Client) GetDeadJobs(ctx context.Context, start, count int) ([]*SortedEntry, int64, error) {
	return c.getSortedSetJobs(ctx, "dead", start, count, true)
}

// GetRetryJobs fetches retry jobs with pagination (earliest retry first).
func (c *Client) GetRetryJobs(ctx context.Context, start, count int) ([]*SortedEntry, int64, error) {
	return c.getSortedSetJobs(ctx, "retry", start, count, false)
}

// GetScheduledJobs fetches scheduled jobs with pagination (earliest execution time first).
func (c *Client) GetScheduledJobs(ctx context.Context, start, count int) ([]*SortedEntry, int64, error) {
	return c.getSortedSetJobs(ctx, "schedule", start, count, false)
}
