package application

type Input struct {
	IgnoreVendors bool
	OutputType    OutputType
	RootDirectory string
}

type OutputType string

var (
	OutputTypeMarkdown OutputType = "markdown"
)

type ProjectParameters struct {
	GoVersion   string
	ProjectName string
	ProjectUrl  string
}
