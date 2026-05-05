package logger

func (l *Logger_struct) Logger_info_method(message string) {
	l.instance.Info(message)
}
