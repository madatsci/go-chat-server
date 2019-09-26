package providers

import "go.uber.org/zap"

func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewDevelopment()

	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}
