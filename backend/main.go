package backend

import (
	"fmt"
	"net/http"
	"time"
	conf "wba-bc-project-05/backend/config"
	"wba-bc-project-05/backend/controller"
	"wba-bc-project-05/backend/router"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	// config 초기화
	cf := conf.NewConfig("./config/config.toml")

	if ctl, err := controller.NewCTL(cf); err != nil {
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
	// 종료 대기
	if err := g.Wait(); err != nil {
		panic(err)
	}
}
