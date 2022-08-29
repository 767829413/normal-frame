package apm

import (
	"sync"

	"github.com/767829413/normal-frame/fork/SkyAPM/go2sky"
	"github.com/767829413/normal-frame/pkg/apm/reporter"

	//"go-api-frame/pkg/SkyAPM/go2sky/reporter"
	"github.com/767829413/normal-frame/internal/pkg/logger"
	"github.com/767829413/normal-frame/internal/pkg/options"
	"github.com/767829413/normal-frame/pkg/util"
)

var (
	re     go2sky.Reporter
	once   sync.Once
	tracer *tracerInc
)

type tracerInc struct {
	mutex  sync.Mutex
	Tracer *go2sky.Tracer
}

func GetApmTracer(opts *options.ApmOptions) *tracerInc {
	if opts != nil && !opts.Enabled {
		return nil
	}
	if opts == nil && tracer == nil {
		return nil
	}
	var err error
	once.Do(func() {
		// 	"github.com/767829413/normal-frame/fork/SkyAPM/go2sky/reporter"
		//re, err := reporter.NewGRPCReporter(confer.GetGlobalConfig().APM.Addr)
		re, err = reporter.NewSidecarReporter(opts.Address)
		if err != nil {
			return
		}
		tmpTra, err := go2sky.NewTracer(util.GetUniqueID(), go2sky.WithReporter(re))
		if err != nil {
			logger.LogErrorf(nil, logger.LogNameAmq, "once.Do GetApmTracer: %v", err)
			return
		}
		tracer = &tracerInc{Tracer: tmpTra}
		//defer re.Close()
	})
	return tracer
}

func (t *tracerInc) Close() error {
	defer t.mutex.Unlock()
	t.mutex.Lock()
	if re != nil {
		re.Close()
	}
	return nil
}
