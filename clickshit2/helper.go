package clickshit2

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func clickhouseTime(t time.Time) string {
	return t.Format(`2006-01-02T15:04:05.999999`)
}

func SQLRowsToMap(rows *sql.Rows) (ds []map[string]any, err error) {
	// 获取列名
	columns, err := rows.Columns()
	if err != nil {
		err = fmt.Errorf("get columns: %w", err)
		return
	}

	// 迭代每一行数据
	for rows.Next() {
		// 创建一个用于存储列值的切片
		values := make([]any, len(columns))
		for i := range values {
			values[i] = new(any)
		}
		// 将列值扫描到values切片中
		err := rows.Scan(values...)
		if err != nil {
			return nil, fmt.Errorf("scan values: %w", err)
		}

		// 创建一个用于存储结果的空map
		row := make(map[string]any)

		// 将列名和对应的值存储到结果map中
		for i, column := range columns {
			if values[i] == nil {
				continue
			}
			row[column] = *(values[i].(*any))
		}

		ds = append(ds, row)
	}

	// 检查是否有错误发生
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows last check: %w", err)
	}

	return
}

func SqlxRowsToMap(rows *sqlx.Rows) (ds []map[string]any, err error) {
	for rows.Next() {
		rowMap := make(map[string]any)
		if err := rows.MapScan(rowMap); err != nil {
			return nil, fmt.Errorf("rows map scan: %w", err)
		}

		for k, v := range rowMap {
			switch val := v.(type) {
			case time.Time:
				rowMap[k] = clickhouseTime(val)
			case []uint8:
				if len(val) == 0 {
					delete(rowMap, k)
				} else {
					newSlice := make([]any, 0, len(val))
					for _, it := range val {
						newSlice = append(newSlice, it)
					}
					rowMap[k] = newSlice
				}
			}
		}

		ds = append(ds, rowMap)
	}

	return
}
