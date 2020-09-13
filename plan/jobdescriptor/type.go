// Package jobdescriptor provides a type to store information of a job in
package jobdescriptor

// JobDescriptor contains all relevant information on a job.
type JobDescriptor struct {
	// Action contains the code of the action to perform.
	// -1 = delete
	// 0 = disabled
	// 1 = move
	// 2 = mkdir
	Action int
	// SourcePath contains the path of the original file
	SourcePath string
	// DstPath contains the destination path if applicable
	DstPath string
}
