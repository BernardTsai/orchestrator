package file

//------------------------------------------------------------------------------

// ROOTDIR points to the root directory of the file system
const ROOTDIR = "/tmp/data"

// DATADIR is the name of the subdirectory holding the component and instance information
const DATADIR = ".data"

// COMPFILE is the name of the file which holds the component information
const COMPFILE = ".component"

// Controller manages the lifecycle of a file
type Controller struct {
}

//------------------------------------------------------------------------------
