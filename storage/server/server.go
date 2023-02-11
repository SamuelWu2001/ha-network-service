package server

import (
	"context"
	// "encoding/json"
	// "fmt"

	"github.com/eddielin0926/ha-network-service/grpcpb/storage"
	// "github.com/eddielin0926/ha-network-service/storage/models"
	// "google.golang.org/protobuf/encoding/protojson"
	// "gorm.io/gorm"

	//new
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type storageServer struct {
	storage.UnimplementedStorageServer
	// DB *gorm.DB
	//new
	DB *sql.DB
}

// func NewStorageServer(db *gorm.DB) *storageServer {
// 	return &storageServer{DB: db}
// }
//new
func NewStorageServer(db *sql.DB) *storageServer {
	return &storageServer{DB: db}
}

// func (s *storageServer) SaveRecord(ctx context.Context, in *storage.Record) (*storage.Response, error) {
// 	s.DB.Create(&models.Record{
// 		Location:  in.Location,
// 		Timestamp: in.Timestamp,
// 		Signature: in.Signature,
// 		Material:  in.Material,
// 		A:         in.A,
// 		B:         in.B,
// 		C:         in.C,
// 		D:         in.D,
// 	})
// 	return &storage.Response{Status: storage.Status_SUCCESS}, nil
// }

//new
func (s *storageServer) SaveRecord(ctx context.Context, in *storage.Record) (*storage.Response, error) {
	_, err := s.DB.Exec("insert INTO data(Location,Timestamp,Signature,Material,A,B,C,D) values(?,?,?,?,?,?,?,?)", in.Location, in.Timestamp, in.Signature, in.Material, in.A, in.B, in.C, in.D)
	if err != nil {
		log.Fatalf("Insert data failed,err:%v", err)
	}
	return &storage.Response{Status: storage.Status_SUCCESS}, nil
}

// func (s *storageServer) GetRecords(ctx context.Context, in *storage.Query) (*storage.RecordsResponse, error) {
// 	var data []models.Record
// 	s.DB.Where("location = ? AND timestamp LIKE ?", in.GetLocation(), in.GetDate()+"%").Find(&data)
// 	fmt.Println(data)
// 	u := protojson.UnmarshalOptions{
// 		AllowPartial:   true,
// 		DiscardUnknown: true,
// 	}

// 	var records []*storage.Record
// 	for _, e := range data {
// 		r := &storage.Record{}
// 		bin, _ := json.Marshal(e)
// 		u.Unmarshal(bin, r)
// 		records = append(records, r)
// 	}

// 	return &storage.RecordsResponse{Records: records}, nil
// }

//new
func (s *storageServer) GetRecords(ctx context.Context, in *storage.Query) (*storage.RecordsResponse, error) {
	// var data []models.Record
	rows, _ := s.DB.Query("select * from customer where location=? AND timestamp LIKE ?", in.GetLocation(), in.GetDate()+"%")
	// fmt.Println(data)
	defer rows.Close()
	// if err != nil {
	// 	fmt.Printf("Query failed,err:%v\n", err)
	// 	return
	// }
	var records []*storage.Record
	for rows.Next() {
		r := &storage.Record{}
		_ = rows.Scan(&r.Location, &r.Timestamp, &r.Signature, &r.Material, &r.A, &r.B, &r.C, &r.D )
		// if err != nil {
		// 	fmt.Printf("Scan failed,err:%v\n", err)
		// 	return
		// }
		records = append(records, r)
	}

	return &storage.RecordsResponse{Records: records}, nil
	// rows, err := s.DB.Query("select * from customer where location=? AND timestamp LIKE ?", in.GetLocation(), in.GetDate()+"%")
	// defer rows.Close()
}

// func (s *storageServer) GetReport(ctx context.Context, in *storage.Query) (*storage.Report, error) {
// 	report := &storage.Report{}
// 	var data []models.Record
// 	s.DB.Where("location = ? AND timestamp LIKE ?", in.GetLocation(), in.GetDate()+"%").Find(&data)
// 	for _, e := range data {
// 		report.A += e.A
// 		report.B += e.B
// 		report.C += e.C
// 		report.D += e.D
// 		report.Material += e.Material
// 	}
// 	report.Location = in.GetLocation()
// 	report.Date = in.GetDate()
// 	report.Count = int32(len(data))
// 	return report, nil
// }

func (s *storageServer) GetReport(ctx context.Context, in *storage.Query) (*storage.Report, error) {
	report := &storage.Report{}
	// var data []models.Record
	// s.DB.Where("location = ? AND timestamp LIKE ?", in.GetLocation(), in.GetDate()+"%").Find(&data)
	rows, _ := s.DB.Query("select * from customer where location=? AND timestamp LIKE ?", in.GetLocation(), in.GetDate()+"%")
	defer rows.Close()
	// if err != nil {
	// 	fmt.Printf("Query failed,err:%v\n", err)
	// 	return
	// }
	count := 0
	for rows.Next() {
		r := &storage.Report{}
		_ = rows.Scan(&r.A, &r.B, &r.C, &r.D, &r.Material)
		report.A += r.A
		report.B += r.B
		report.C += r.C
		report.D += r.D
		report.Material += r.Material
		count += 1
	}
	report.Location = in.GetLocation()
	report.Date = in.GetDate()
	report.Count = int32(count)
	return report, nil
}
