package orm

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestGorm(t *testing.T) {

	db := MustNewMysql(&Config{
		DataSource:   "root:200212..@tcp(127.0.0.1:3306)/dev_bach_pet_lhp?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
		MaxLifetime:  60,
	})

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * time.Duration(i) * 10)
			// 执行 SQL 查询，获取所有表名
			db.Statement.WithContext(context.Background())
			// db.WithContext(context.Background())
			db.InstanceSet("key", i)
			rows, err := db.Raw("SHOW TABLES").Rows()
			if err != nil {
				t.Fatalf("failed to execute query: %v", err)
			}
			defer rows.Close()

			var tableName string
			for rows.Next() {
				if err = rows.Scan(&tableName); err != nil {
					t.Fatalf("failed to scan row: %v", err)
				}
				// fmt.Println(tableName)
			}
			if err = rows.Err(); err != nil {
				t.Fatalf("error during iteration: %v", err)
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
