package errorcode

const (
	// UNKNOWNERROR  .
	UNKNOWNERROR uint64 = 99
	// NOTINIT  .
	NOTINIT uint64 = 100
	// DBERROR .
	DBERROR uint64 = 101
	// ParameterError .
	ParameterError uint64 = 102
	// NOTEXISTRECORD .
	NOTEXISTRECORD uint64 = 103
	// ParameterIllegal .
	ParameterIllegal uint64 = 104
	// MarshalError .
	MarshalError uint64 = 105

	// ACCOUNTNOTMATCH .
	ACCOUNTNOTMATCH uint64 = 201
	// NOTRUNPPMQCONNECT .
	NOTRUNPPMQCONNECT uint64 = 202
	// SUBREJECT .
	SUBREJECT uint64 = 203
	// PUBREJECT .
	PUBREJECT uint64 = 204
)
