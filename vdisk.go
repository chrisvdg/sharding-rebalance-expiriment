package main

import (
	"fmt"
)

// defines the get shard algorithm
var (
	getShard = getShardIndexModulo
)

// Errors
var (
	ErrShardIndexNotFound = fmt.Errorf("could not find shard index")
	ErrShardNotHealthy    = fmt.Errorf("shard is not healthy")
)

// NewVdisk constructs a new vdisk
func NewVdisk(shardCount int) *Vdisk {
	var vdisk Vdisk
	for i := 0; i < shardCount; i++ {
		vdisk.Shards = append(vdisk.Shards, NewShard())
	}

	return &vdisk
}

// Vdisk vdisk represents a vdisk
type Vdisk struct {
	Shards []*Shard
}

// SetBlock sets a block in a vdisk
func (vdisk *Vdisk) SetBlock(blockIndex int, data byte) error {
	shardIndex, err := getShard(vdisk, blockIndex)
	if err != nil {
		return err
	}
	s := vdisk.Shards[shardIndex]
	if !s.OK() {
		return ErrShardNotHealthy
	}

	s.SetBlock(blockIndex, data)
	return nil
}

// GetBlock gets the data from a block in a vdisk
func (vdisk *Vdisk) GetBlock(blockIndex int) (byte, error) {
	shardIndex, err := getShard(vdisk, blockIndex)
	if err != nil {
		return 0, err
	}

	s := vdisk.Shards[shardIndex]
	if !s.OK() {
		return 0, ErrShardNotHealthy
	}

	return s.GetBlock(blockIndex)
}

// FailShard set a shard to unhealthy and redistributes the data of the failed shard
func (vdisk *Vdisk) FailShard(shardIndex int) error {
	if shardIndex >= len(vdisk.Shards) {
		return ErrShardIndexNotFound
	}

	s := vdisk.Shards[shardIndex]

	s.SetHealth(false)
	for blockIndex, data := range s.data {
		err := vdisk.SetBlock(blockIndex, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// HealthyShards returns the count of healthy shard in a vdisk
func (vdisk *Vdisk) HealthyShards() int {
	healthyCount := 0

	for _, shard := range vdisk.Shards {
		if shard.OK() {
			healthyCount++
		}
	}

	return healthyCount
}

// PrintShardingState prints out the current block count for each shard
func (vdisk *Vdisk) PrintShardingState() {
	fmt.Println("\n\t--- Current vdisk state ---")
	for i := range vdisk.Shards {
		s := vdisk.Shards[i]
		blocks := s.BlockCount()
		var health string
		if s.OK() {
			health = "healthy"
		} else {
			health = "unhealthy"
		}

		fmt.Printf("Shard %d is %s and has %d blocks\n", i, health, blocks)
	}
}

// GetShardIndex returns a shardindex for a given blockindex
func getShardIndexModulo(vdisk *Vdisk, blockIndex int) (int, error) {
	return blockIndex % len(vdisk.Shards), nil
}

func getShardGeertsAlgo(vdisk *Vdisk, blockIndex int) (int, error) {
	shardCount := len(vdisk.Shards)
	shardIndex := blockIndex % shardCount
	if vdisk.Shards[shardIndex].OK() {
		return shardIndex, nil
	}

	shardIndex = hash(blockIndex)
	shardCounter := 0

	for i := 0; i < shardCount; i++ {
		if !vdisk.Shards[i].OK() {
			continue
		}
		if shardCounter == shardIndex {
			return i, nil
		}
		shardCounter++
	}
	return 0, ErrShardIndexNotFound
}

func hash(blockIndex int) int {

	return 0
}
