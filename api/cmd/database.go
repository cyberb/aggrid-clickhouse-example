package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"
	"net"
	"time"
)

const (
	CreateTableSql = `
CREATE TABLE IF NOT EXISTS transaction
(
id UUID DEFAULT generateUUIDv4(),
year UInt16 DEFAULT (rand() % 100) + 1900,
amount Float64 DEFAULT randNormal(0, 100000), 
field1 String DEFAULT concat('test', toString(floor(randUniform(1, 5) % 10))),  
field2 String DEFAULT concat('test', toString(floor(randUniform(1, 5) % 10))),  
field3 String DEFAULT concat('test', toString(floor(randUniform(1, 5) % 10))),  
field4 String DEFAULT concat('test', toString(floor(randUniform(1, 5) % 10))),  
field5 String DEFAULT concat('test', toString(floor(randUniform(1, 5) % 10))),  
field6 String DEFAULT concat('test', toString(floor(randUniform(1, 5) % 10)))
)
ENGINE = Memory
`
	InsertData = `
insert into transaction (id) select generateUUIDv4() from numbers(100000)
`
)

type Database struct {
	conn   driver.Conn
	logger *zap.Logger
}

func NewDatabase(logger *zap.Logger) *Database {
	return &Database{
		logger: logger,
	}
}

func (d *Database) Connect() error {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 38,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      16,
		MaxCompressionBuffer: 19240,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "example", Version: "0.1"},
			},
		},
	})
	if err != nil {
		return err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return err
	}
	d.conn = conn

	return nil
}

func (d *Database) StartUpdater() {
	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := d.Update()
				if err != nil {
					d.logger.Error("update", zap.Error(err))
				}
			}
		}
	}()
}

func (d *Database) Update() error {
	fmt.Println("")

	d.logger.Info("update trigger")
	err := d.conn.Exec(context.Background(), `
alter table transaction update
amount = randNormal(0, 100000),
field1 = concat('test', toString(floor(randUniform(10, 20)%3))),
field2 = concat('test', toString(floor(randUniform(20, 30)%3))),
field3 = concat('test', toString(floor(randUniform(40, 50)%3))),
field4 = concat('test', toString(floor(randUniform(60, 70)%3))),
field5 = concat('test', toString(floor(randUniform(80, 90)%3))),
field6 = concat('test', toString(floor(randUniform(90, 100)%3)))
where toString(id) like concat('%',toString(rand()%10), '%')
`)
	return err
}

func (d *Database) GenerateData() error {
	err := d.conn.Exec(context.Background(), CreateTableSql)
	if err != nil {
		return err
	}
	row := d.conn.QueryRow(context.Background(), "select count(*) as size from transaction")
	if row.Err() != nil {
		return row.Err()
	}
	var count uint64
	err = row.Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		err := d.conn.Exec(context.Background(), InsertData)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Database) GetData(request Request) ([]Row, uint64, error) {
	d.logger.Info("get data", zap.Any("request", request))

	query := "SELECT "
	if request.IsGrouping() {
		query += request.RowGroupCols.ToSelectSql()
		query += request.ValueCols.ToSql()
	} else {
		query += " * "
	}
	query += " FROM transaction WHERE 1=1 "
	query += request.FilterModel.ToSql()
	query += request.ToGroupWhereSql()
	if request.IsGrouping() {

		query += request.RowGroupCols.ToGroupSql()
	}
	query += request.ToOrderSql()
	query += " LIMIT ?, ?"

	var total uint64
	ctxWProfile := clickhouse.Context(
		context.Background(),
		clickhouse.WithProfileInfo(func(info *clickhouse.ProfileInfo) {
			if info != nil {
				total = info.RowsBeforeLimit
			}
		}))
	rows, err := d.conn.Query(ctxWProfile, query, request.StartRow, request.EndRow-request.StartRow)
	if err != nil {
		d.logger.Error("error", zap.Error(err))
		return nil, 0, err
	}
	results := []Row{}
	for rows.Next() {
		var row Row
		if err := rows.ScanStruct(&row); err != nil {
			d.logger.Error("error", zap.Error(err))
			return nil, 0, err
		}
		//Logger.Info("row", zap.Int32("year", row.Year), zap.Float64("amount", row.Amount))
		results = append(results, row)
	}
	err = rows.Close()
	if err != nil {
		d.logger.Error("error", zap.Error(err))
		return nil, 0, err

	}
	err = rows.Err()
	if err != nil {
		d.logger.Error("error", zap.Error(err))
		return nil, 0, err
	}

	d.logger.Info("")
	d.logger.Info("get data", zap.Int("size", len(results)))

	return results, total, nil
}
