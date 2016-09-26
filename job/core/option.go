package core

import (
	"fmt"
	"time"
)

var (
	// DefaultJobQueueBufSize ジョブキューのデフォルトサイズ
	DefaultJobQueueBufSize = 50

	// DefaultSakuraAPIReqQueueBufSize さくらのクラウドAPIキューのデフォルトサイズ
	DefaultSakuraAPIReqQueueBufSize = 10

	// DefaultMackerelAPIReqQueueBufSize MackerelAPIキューのデフォルトサイズ
	DefaultMackerelAPIReqQueueBufSize = 10

	// DefaultThrottledAPIReqQueueBufSize スロットリング対象APIキューのデフォルトサイズ
	DefaultThrottledAPIReqQueueBufSize = 0

	// DefaultAPICallInterval スロットリングAPI待機時間デフォルト値
	DefaultAPICallInterval = 500 * time.Millisecond

	// DefaultTimerJobInterval タイマージョブのデフォルト起動間隔
	DefaultTimerJobInterval = 2 * time.Minute

	// DefaultMetricsHistoryPeriod メトリクス収集時の過去分取得範囲デフォルト値
	DefaultMetricsHistoryPeriod = 15 * time.Minute

	// DefaultHealthCheckWebServerPort ヘルスチェック用Webサーバーのデフォルトポート
	DefaultHealthCheckWebServerPort = 39700

	// DefaultIgnoreTag sackerelでの連携対象外を示すマーカータグ名称
	DefaultIgnoreTag = "@mackerel-ignore"
)

// Option sackerel動作オプション
type Option struct {
	SakuraCloudOption           *SakuraCloudOption
	MackerelOption              *MackerelOption
	JobQueueBufSize             int
	ThrottledAPIReqQueueBufSize int
	TimerJobInterval            time.Duration
	MetricsHistoryPeriod        time.Duration
	APICallInterval             time.Duration
	HealthCheckWebServerPort    int
	DisableHealthCheck          bool
	SkipInit                    bool
	TraceLog                    bool
	InfoLog                     bool
	WarnLog                     bool
	ErrorLog                    bool
	IgnoreTag                   string
}

// NewOption Optionの新規作成
func NewOption() *Option {
	return &Option{
		SakuraCloudOption: NewSakuraCloudOption(),
		MackerelOption:    NewMackerelOption(),
		IgnoreTag:         DefaultIgnoreTag,
	}
}

type optionValidator interface {
	validate() []error
}

// Validate オプション値の妥当性検証
func (o *Option) Validate() []error {
	var errors []error
	validators := []optionValidator{o.SakuraCloudOption, o.MackerelOption}
	for _, validator := range validators {
		errs := validator.validate()
		errors = append(errors, errs...)
	}
	return errors
}

// QueueBufSizes ジョブキューのサイズ取得
func (o *Option) QueueBufSizes() (int, int, int, int) {
	return o.JobQueueBufSize, o.ThrottledAPIReqQueueBufSize, o.SakuraCloudOption.ReuqestQueueBufSize, o.MackerelOption.RequestQueueBufSize
}

// ----------------------------------------------

// SakuraCloudDefaultZones さくらのクラウド対象ゾーンのデフォルト値
var SakuraCloudDefaultZones = []string{"is1b", "tk1a"}

// SakuraCloudOption さくらのクラウド用の動作オプション
type SakuraCloudOption struct {
	AccessToken         string
	AccessTokenSecret   string
	Zone                []string
	IgnoreTags          []string
	TraceMode           bool
	ReuqestQueueBufSize int
}

// NewSakuraCloudOption SakuraCloudOptionの新規作成
func NewSakuraCloudOption() *SakuraCloudOption {
	return &SakuraCloudOption{
		Zone: SakuraCloudDefaultZones,
	}
}

func (o *SakuraCloudOption) validate() []error {
	var errors []error
	if o.AccessToken == "" {
		errors = append(errors, fmt.Errorf("[%s] is required", "sakuracloud-access-token"))
	}
	if o.AccessTokenSecret == "" {
		errors = append(errors, fmt.Errorf("[%s] is required", "sakuracloud-access-token-secret"))
	}

	return errors
}

// ----------------------------------------------

// MackerelOption Mackerel用の動作オプション
type MackerelOption struct {
	APIKey              string
	RequestQueueBufSize int
	TraceMode           bool
}

// NewMackerelOption MackerelOptionの新規作成
func NewMackerelOption() *MackerelOption {
	return &MackerelOption{}
}

func (o *MackerelOption) validate() []error {
	var errors []error
	if o.APIKey == "" {
		errors = append(errors, fmt.Errorf("[%s] is required", "mackerel-api-key"))
	}

	return errors
}
