package ecspresso

import "log"

var (
	SortTaskDefinitionForDiff    = sortTaskDefinitionForDiff
	SortServiceDefinitionForDiff = sortServiceDefinitionForDiff
	EqualString                  = equalString
	ToNumberCPU                  = toNumberCPU
	ToNumberMemory               = toNumberMemory
	CalcDesiredCount             = calcDesiredCount
	ParseTags                    = parseTags
	ParseRoleArn                 = parseRoleArn
	IsLongArnFormat              = isLongArnFormat
	ECRImageURLRegex             = ecrImageURLRegex
	NewLogger                    = newLogger
	NewLogFilter                 = newLogFilter
)

func (d *App) SetLogger(logger *log.Logger) {
	d.logger = logger
}

func SetLogger(logger *log.Logger) {
	commonLogger = logger
}
