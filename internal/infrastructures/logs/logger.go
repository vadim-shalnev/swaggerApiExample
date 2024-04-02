package logs

import (
	"github.com/vadim-shalnev/swaggerApiExample/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(conf config.AppConf, sync zapcore.WriteSyncer) *zap.Logger {

	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	atom := zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	zapConf.Level = atom
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConf.EncoderConfig),
		sync,
		atom,
	)
	logger := zap.New(core)

	return logger.Named("geoservis")
}
