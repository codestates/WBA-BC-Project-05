package main

import (
	"fmt"
	"net/http"
	"time"
	"wba-bc-project-05/backend/controller"
	"wba-bc-project-05/backend/model"
	"wba-bc-project-05/backend/router"
	conf "wba-bc-project-05/config"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	// config 초기화
	cf := conf.NewConfig("../config/config.toml")

	if md, err := model.NewModel(cf.DB.Host); err != nil {
		panic(fmt.Errorf("controller.NewCTL error: %v", err))
	} else if ctl, err := controller.NewCTL(cf, md); err != nil {
		// 컨트롤러 초기화
		panic(fmt.Errorf("controller.NewCTL error: %v", err))
	} else if rt, err := router.NewRouter(ctl); err != nil {
		// 라우터 초기화
		panic(fmt.Errorf("router.NewRouter error: %v", err))
	} else {
		// 웹서버 설정
		mapi := &http.Server{
			Addr:           cf.Web.Port,
			Handler:        rt.Idx(),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		// 웹서버 실행
		g.Go(func() error {
			return mapi.ListenAndServe()
		})
	}
	fmt.Println("wait!", time.Now().Unix())

	// 종료 대기
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
