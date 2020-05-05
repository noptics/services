package golog

var loggers []Logger
var level int

func Add(l Logger) {
	loggers = append(loggers, l)
}

func Init() error {
	for _, l := range loggers {
		err := l.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func Finish() error {
	for _, l := range loggers {
		err := l.Finish()
		if err != nil {
			return err
		}
	}

	return nil
}

func Logw(level int, message string, data ...interface{}) {
	for _, l := range loggers {
		l.Infow(message, data...)
	}
}

func Info(data ...interface{}) {
	for _, l := range loggers {
		l.Info(data...)
	}
}

func Infow(message string, data ...interface{}) {
	for _, l := range loggers {
		l.Infow(message, data...)
	}
}

func Error(data ...interface{}) {
	for _, l := range loggers {
		l.Info(data...)
	}
}

func Errorw(message string, data ...interface{}) {
	for _, l := range loggers {
		l.Infow(message, data...)
	}
}

func Debug(data ...interface{}) {
	for _, l := range loggers {
		l.Info(data...)
	}
}

func Debugw(message string, data ...interface{}) {
	for _, l := range loggers {
		l.Infow(message, data...)
	}
}
