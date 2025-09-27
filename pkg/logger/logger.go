package logger

import (
	"log"

	"go.uber.org/zap"
)

// Log глобальный экземпляро
var Log *zap.SugaredLogger

// Init инициализация
func Init() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't init logger: %v", err)
	}
	Log = zapLogger.Sugar()
}

// Sync очищает буферы
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
