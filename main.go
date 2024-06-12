package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/inits"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/comm/log"
	"github.com/WeixinCloud/wxcloudrun-wxcomponent/routers"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/sync/errgroup"
)

type Ticket struct {
	Value string `json:"ticket"`
}

func postTicket(url string, ticket string) error {
	data := Ticket{
		Value: ticket,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	return nil
}

func main() {
	log.Infof("system begin")
	if err := inits.Init(); err != nil {
		log.Errorf("inits failed, err:%v", err)
		return
	}
	log.Infof("inits.Init Succ")

	var g errgroup.Group

	// 内部服务
	g.Go(func() error {
		r := routers.InnerServiceInit()
		if err := r.Run("127.0.0.1:8081"); err != nil {
			log.Error("startup inner service failed, err:%v", err)
			return err
		}
		return nil
	})

	// 外部服务
	g.Go(func() error {
		r := routers.Init()
		if err := r.Run(":80"); err != nil {
			log.Error("startup service failed, err:%v", err)
			return err
		}
		return nil
	})

	// url := "http://192.168.2.18:5095/api/wxcallback/ticket"

	// go func() {
	// 	for {
	// 		ticket := wxbase.GetTicket()

	// 		log.Infof("Fetched ticket: %s", ticket)

	// 		err := postTicket(url, ticket)
	// 		if err != nil {
	// 			log.Errorf("Failed to post ticket: %v", err)
	// 		} else {
	// 			log.Info("Ticket posted successfully")
	// 		}

	// 		time.Sleep(5 * time.Second)
	// 	}
	// }()

	if err := g.Wait(); err != nil {
		log.Error(err)
	}
}
