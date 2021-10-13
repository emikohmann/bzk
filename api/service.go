package api

import (
	"bytes"
	"fmt"
	"github.com/mercadolibre/go-meli-toolkit/goutils/apierrors"
	"github.com/mercadolibre/go-meli-toolkit/goutils/logger"
	"github.com/mercadolibre/go-meli-toolkit/restful/rest"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Service interface {
	Run(test LoadTest) (Result, apierrors.ApiError)
}

type ServiceImpl struct {
}

func NewServiceImpl() ServiceImpl {
	return ServiceImpl{}
}

func (serviceImpl ServiceImpl) Run(test LoadTest) (Result, apierrors.ApiError) {
	headers := make(http.Header)
	for key, value := range test.Target.Headers {
		headers.Add(key, value)
	}

	builder := rest.RequestBuilder{
		Headers:        headers,
		BaseURL:        fmt.Sprintf("%s://%s", test.Target.Protocol, test.Target.BaseURL),
		EnableCache:    false,
		DisableTimeout: true,
		FollowRedirect: true,
	}

	task := func(result chan int, group *sync.WaitGroup) {
		defer group.Done()
		var response *rest.Response
		switch test.Target.Method {
		case http.MethodGet:
			response = builder.Get(func() string {
				args := make([]interface{}, 0)
				for _, values := range test.Target.Endpoint.Context {
					args = append(args, values[rand.Intn(len(values))])
				}
				return fmt.Sprintf(test.Target.Endpoint.Format, args...)
			}())
		}
		if response == nil || response.Err != nil {
			result <- http.StatusInternalServerError
			return
		}
		result <- response.StatusCode
	}

	stepResults := make([]StepResult, 0)
	start := time.Now()

	for i, step := range test.Steps {
		var (
			stepStart                   = time.Now()
			stepTimer                   = time.NewTimer(time.Duration(step.DurationSec) * time.Second)
			currentRPM                  = float64(step.RPMFrom)
			rpmUpdaterInterval          = 1 * time.Second
			rpmUpdater                  = time.NewTicker(rpmUpdaterInterval)
			ascendingStep               = step.RPMTo > step.RPMFrom
			rpmUpdaterIntervalVariation = func() float64 {
				var stepDiff float64
				if ascendingStep {
					stepDiff = float64(step.RPMTo) - float64(step.RPMFrom)
				} else {
					stepDiff = float64(step.RPMFrom) - float64(step.RPMTo)
				}
				return stepDiff / float64(step.DurationSec)
			}()
			currentAwait = float64(60) / currentRPM
			stepTicker    = time.NewTicker(time.Duration(currentAwait*1000) * time.Millisecond)
			stepResult    = StepResult{
				CallCount:       0,
				Status:          make(map[int]int),
				DurationSeconds: 0,
			}
			stepGroup   = &sync.WaitGroup{}
			stepChannel = make(chan int)
		)
		logger.Infof("Running step [%d]", i)
		go func() {
			for msg := range stepChannel {
				stepResult.Status[msg]++
			}
		}()
		go func() {
			for {
				select {
				case <-rpmUpdater.C:
					if (ascendingStep && currentRPM < float64(step.RPMTo)) ||
						(!ascendingStep && currentRPM > float64(step.RPMTo)) {
						if ascendingStep {
							currentRPM += rpmUpdaterIntervalVariation
						} else {
							currentRPM -= rpmUpdaterIntervalVariation
						}
						currentAwait = float64(60) / currentRPM
						if currentAwait < rpmUpdaterInterval.Seconds() && currentAwait > 0 {
							stepTicker.Reset(time.Duration(currentAwait*1000) * time.Millisecond)
						}
					}
				}
			}
		}()
	StepCalls:
		for {
			select {
			case <-stepTimer.C:
				stepResult.DurationSeconds = time.Since(stepStart).Seconds()
				stepResults = append(stepResults, stepResult)
				break StepCalls
			default:
				stepGroup.Add(1)
				var buffer bytes.Buffer
				for code, count := range stepResult.Status {
					buffer.WriteString(fmt.Sprintf("| %d (%d)", code, count))
				}
				logger.Infof("[%d] From %d RPM to %d RPM (current %.2f RPM) [%d]: %s", i, step.RPMFrom, step.RPMTo, currentRPM, stepResult.CallCount, buffer.String())
				go task(stepChannel, stepGroup)
				<-stepTicker.C
				if currentAwait > 0 {
					stepTicker.Reset(time.Duration(currentAwait*1000) * time.Millisecond)
				}
				stepResult.CallCount++
			}
		}
		stepGroup.Wait()
		close(stepChannel)
	}
	return Result{
		StepResults:     stepResults,
		DurationSeconds: time.Since(start).Seconds(),
	}, nil
}
