package job

import (
	"fmt"
	"github.com/sacloud/sackerel/job/core"
	"net/http"
)

// HealthCheckServer ヘルスチェック用Webサーバー
type HealthCheckServer struct {
	Port int
}

// NewHealthCheckServer HealthCheckServerの新規作成
func NewHealthCheckServer(port int) *HealthCheckServer {
	return &HealthCheckServer{
		Port: port,
	}
}

// ListenAndServe ヘルスチェック用に指定ポートでリッスンし、httpリクエストに対し応答する
func (s *HealthCheckServer) ListenAndServe() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	if s.Port == 0 {
		s.Port = core.DefaultHealthCheckWebServerPort
	}

	go http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
}
