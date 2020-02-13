package grpc_excel

import (
	"context"
	"errors"
	excel "github.com/deltrinos/tpl21/excel/proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"os"
	"time"
)

const (
	MaxGrpcSendMsgSize    = 1024 * 1024 * 100
	MaxGrpcReceiveMsgSize = 1024 * 1024 * 10
)

type GrpcExcelClient struct {
	conn    *grpc.ClientConn
	c       excel.ExcelServiceClient
	handler string
	list    []string
	macros  []string
}

type ExcelHandler string

func New(server string) (*GrpcExcelClient, error) {
	client := &GrpcExcelClient{}
	var err error

	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	client.conn, err = grpc.DialContext(ctx, server,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(MaxGrpcReceiveMsgSize),
			grpc.MaxCallSendMsgSize(MaxGrpcSendMsgSize),
		),
	)
	if err != nil {
		log.Error().Err(err).Msgf("grpc.Dial error: %v", err)
		return nil, err
	}
	client.c = excel.NewExcelServiceClient(client.conn)
	return client, nil
}

func Default() *GrpcExcelClient {
	excelConnStr := os.Getenv("EXCEL_CONN_STR")
	if excelConnStr == "" {
		excelConnStr = "localhost:9999"
	}
	return GetInstance(excelConnStr)
}

func GetInstance(addr string) *GrpcExcelClient {
	client, err := New(addr)
	if err != nil {
		log.Error().Err(err).Msgf("Can't connect to grpc_excel_server %v", err)
		return nil
	} else {
		log.Info().Msgf("connected to grpc_excel_server %v", addr)
	}
	return client
}

func (client *GrpcExcelClient) Disconnect() {
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil {
			log.Error().Err(err).Msgf("conn.Close error: %v", err)
		}
	} else {
		log.Error().Msgf("client conn is nil !")
	}
}

func (client *GrpcExcelClient) Ping(handler string) error {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	_, err := client.c.Ping(ctx, &excel.HandlerRq{
		Handler: handler,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to Ping error %v", err)
		return err
	}
	return nil
}

func (client *GrpcExcelClient) Open(fileName, handler string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	res, err := client.c.Open(ctx, &excel.OpenRq{
		Filename:   fileName,
		Handler:    handler,
		WithSheets: true,
	})
	if err != nil {
		log.Error().Msgf("list error %v", err)
		return "", err
	}
	return res.Handler, nil
}

func (client *GrpcExcelClient) Close(handler string, saveIt bool) error {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	_, err := client.c.Close(ctx, &excel.HandlerRq{
		Handler: handler,
		SaveIt:  saveIt,
	})
	if err != nil {
		log.Error().Msgf("fail to client.c.Close %v", err)
		return err
	}
	return nil
}

func (client *GrpcExcelClient) List() ([]string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	r, err := client.c.List(ctx, &excel.Empty{})
	if err != nil {
		log.Error().Msgf("list error %v", err)
		return nil, err
	}
	client.list = r.Filename
	return client.list, nil
}

func (client *GrpcExcelClient) ListMacros(handler string) ([]string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	r, err := client.c.GetMacroNames(ctx, &excel.ListMacrosRq{
		Handler:    handler,
		MacroNames: []string{},
	})
	if err != nil {
		log.Error().Msgf("ListMacros error %v", err)
		return nil, err
	}
	client.macros = r.MacroNames
	return client.macros, nil
}

func (client *GrpcExcelClient) RunMacros(handler string, macroToRun string) (int, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	r, err := client.c.RunMacro(ctx, &excel.MacroNameRq{
		Handler: handler,
		Name:    macroToRun,
	})
	if err != nil {
		log.Error().Msgf("RunMacro(%s) error %v", macroToRun, err)
		return 0, err
	}
	return int(r.Status), nil
}

func (client *GrpcExcelClient) GetNamedCells(handler string, names []string) ([]*excel.NamedRange, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	res, err := client.c.GetNamedCells(ctx, &excel.GetNamedCellsRq{
		Handler: handler,
		Names:   names,
	})
	if err != nil {
		log.Error().Msgf("GetNamedCells error %v", err)
		return nil, err
	}
	return res.Names, nil
}

func (client *GrpcExcelClient) GetSheetCells(handler, sheetName string) ([]*excel.NamedRange, error) {
	return client.GetCells(handler, sheetName, "begin", "end")
}

func (client *GrpcExcelClient) GetCells(handler, sheetName, begin, end string) ([]*excel.NamedRange, error) {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	res, err := client.c.GetCells(ctx, &excel.GetCellsRq{
		Handler: handler,
		Sheet:   sheetName,
		Begin:   begin,
		End:     end,
	})
	if err != nil {
		log.Error().Msgf("GetCells error %v", err)
		return nil, err
	}
	return res.Names, nil
}

func (client *GrpcExcelClient) SetValues(handler string, ranges []*excel.NamedRange) error {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	_, err := client.c.SetValues(ctx, &excel.SetSetValuesRq{
		Handler: handler,
		Ranges:  ranges,
	})
	if err != nil {
		log.Error().Err(err).Msgf("SetValues error %v", err)
		return err
	}
	return nil
}

func (client *GrpcExcelClient) Upload(name string, fileBytes []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	res, err := client.c.Upload(ctx, &excel.UploadRq{
		Name:  name,
		Bytes: fileBytes,
	})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to Upload error %v", err)
		return err
	}
	if !res.IsOk {
		return errors.New(res.Handler)
	}
	return nil
}
