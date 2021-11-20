package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
)

func NewService() *Service {
	ctx := context.Background()
	conn, _ := pgx.Connect(ctx, "postgres://root:root@localhost/test")
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	return &Service{conn, rdb, ctx, log.Default()}
}

type Service struct {
	conn *pgx.Conn
	rdb  *redis.Client
	ctx  context.Context
	log  *log.Logger
}

func (s *Service) GetTask(id int32) interface{} {
	data := s.getTaskFromDB(id)
	val, err := json.Marshal(data)
	if err != nil {
		s.log.Println(err.Error())
	}

	return string(val)
}

func (s *Service) GetTaskWithRedis(id int32) interface{} {
	key := strconv.Itoa(int(id))
	val, err := s.rdb.Get(s.ctx, key).Result()
	if err != nil {
		s.log.Println(err.Error())
	}
	if val != "" {
		return val
	}

	fromDB := s.getTaskFromDB(id)
	valByte, err := json.Marshal(fromDB)
	if err != nil {
		s.log.Println(err.Error())
	}
	valString := string(valByte)

	err = s.rdb.Set(s.ctx, key, valString, 0).Err()
	if err != nil {
		s.log.Println(err.Error())
	}

	return valString
}

type Task struct {
	ID          int32  `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

func (s *Service) getTaskFromDB(id int32) []Task {
	rows, _ := s.conn.Query(s.ctx, "select * from tasks where id=$1", id)
	tasks := make([]Task, 0)

	for rows.Next() {
		rawValues, err := rows.Values()
		if err != nil {
			s.log.Println(err.Error())
		}
		id := rawValues[0].(int32)
		description := rawValues[1].(string)
		task := Task{id, description}
		tasks = append(tasks, task)
	}

	return tasks
}

func main() {
	s := NewService()
	fmt.Println(s.GetTaskWithRedis(1))
}
