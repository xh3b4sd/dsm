package project

var (
	description = "Data structure manipulation. Used in e.g. deployment pipelines."
	gitSHA      = "n/a"
	name        = "dsm"
	source      = "https://github.com/xh3b4sd/dsm"
	version     = "n/a"
)

func Description() string {
	return description
}

func GitSHA() string {
	return gitSHA
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
