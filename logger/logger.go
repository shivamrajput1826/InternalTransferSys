package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type CustomLogger struct {
	logInstance zerolog.Logger
}

type LogOptions struct {
	MethodName string      `json:"methodName,omitempty"`
	Details    interface{} `json:"details,omitempty"`
}

func buildMessage(method, suffix string) string {
	b := strings.Builder{}
	if len(method) > 0 {
		b.WriteString(fmt.Sprintf("%s: ", method))
	}
	b.WriteString(suffix)
	return b.String()
}

func CreateLogger(logContext string) (cl *CustomLogger) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	cl = &CustomLogger{
		logInstance: zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().
			Timestamp().
			Str("logContext", logContext).
			Logger(),
	}
	return
}

func (cl *CustomLogger) WithFiberContext(ctx *fiber.Ctx) (nl *CustomLogger) {
	nl = &CustomLogger{
		logInstance: cl.logInstance,
	}
	contextKeys := []string{"custId", "requestId"}

	for _, key := range contextKeys {
		if ctx.Locals(key) != nil {
			nl.logInstance = nl.logInstance.With().
				Str(key, ctx.Locals(key).(string)).
				Logger()
		}
	}
	return
}

func (cl *CustomLogger) Debug(message string, details interface{}) {
	d := details
	if details == nil {
		d = ""
	}
	jsonBytes, _ := json.Marshal(d)
	cl.logInstance.Debug().Str("details", string(jsonBytes)).Msg(message)
}

func (cl *CustomLogger) Error(option LogOptions) {
	temp := cl.logInstance.Error()
	if option.Details != nil {
		if err, ok := option.Details.(error); ok {
			temp = temp.Str("details", err.Error())
		}
	}
	temp.
		Msg(buildMessage(option.MethodName, "error"))
}

func (cl *CustomLogger) LogResponse(options LogOptions) {
	d := options.Details
	if d == nil {
		d = ""
	}
	jsonBytes, _ := json.Marshal(d)
	cl.logInstance.Info().
		Str("details", string(jsonBytes)).
		Msg(buildMessage(options.MethodName, "response"))
}

func (cl *CustomLogger) LogRequest(options LogOptions) {
	d := options.Details
	if d == nil {
		d = ""
	}
	jsonBytes, _ := json.Marshal(d)
	cl.logInstance.Info().
		Str("details", string(jsonBytes)).
		Msg(buildMessage(options.MethodName, "request"))
}
