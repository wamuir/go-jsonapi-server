package backend

var schema = map[string][]byte{
	"CreateIndexEdges": {
		0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x20, 0x55, 0x4e, 0x49, 0x51, 0x55,
		0x45, 0x20, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x20, 0x49, 0x46, 0x20, 0x4e,
		0x4f, 0x54, 0x20, 0x45, 0x58, 0x49, 0x53, 0x54, 0x53, 0x20, 0x65, 0x64,
		0x67, 0x65, 0x73, 0x5f, 0x69, 0x64, 0x78, 0x31, 0x0a, 0x20, 0x20, 0x20,
		0x20, 0x4f, 0x4e, 0x20, 0x65, 0x64, 0x67, 0x65, 0x73, 0x20, 0x28, 0x66,
		0x72, 0x6f, 0x6d, 0x5f, 0x72, 0x6f, 0x77, 0x69, 0x64, 0x2c, 0x20, 0x6b,
		0x65, 0x79, 0x2c, 0x20, 0x74, 0x6f, 0x5f, 0x72, 0x6f, 0x77, 0x69, 0x64,
		0x29, 0x0a,
	},
	"CreateIndexVertices": {
		0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x20, 0x55, 0x4e, 0x49, 0x51, 0x55,
		0x45, 0x20, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x20, 0x49, 0x46, 0x20, 0x4e,
		0x4f, 0x54, 0x20, 0x45, 0x58, 0x49, 0x53, 0x54, 0x53, 0x20, 0x76, 0x65,
		0x72, 0x74, 0x69, 0x63, 0x65, 0x73, 0x5f, 0x69, 0x64, 0x78, 0x31, 0x0a,
		0x20, 0x20, 0x20, 0x20, 0x4f, 0x4e, 0x20, 0x76, 0x65, 0x72, 0x74, 0x69,
		0x63, 0x65, 0x73, 0x20, 0x28, 0x74, 0x79, 0x70, 0x65, 0x2c, 0x20, 0x69,
		0x64, 0x29, 0x0a,
	},
	"CreateTableEdges": {
		0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x20, 0x54, 0x41, 0x42, 0x4c, 0x45,
		0x20, 0x49, 0x46, 0x20, 0x4e, 0x4f, 0x54, 0x20, 0x45, 0x58, 0x49, 0x53,
		0x54, 0x53, 0x20, 0x65, 0x64, 0x67, 0x65, 0x73, 0x20, 0x28, 0x0a, 0x20,
		0x20, 0x20, 0x20, 0x72, 0x6f, 0x77, 0x69, 0x64, 0x20, 0x49, 0x4e, 0x54,
		0x45, 0x47, 0x45, 0x52, 0x20, 0x50, 0x52, 0x49, 0x4d, 0x41, 0x52, 0x59,
		0x20, 0x4b, 0x45, 0x59, 0x2c, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x66, 0x72,
		0x6f, 0x6d, 0x5f, 0x72, 0x6f, 0x77, 0x69, 0x64, 0x20, 0x49, 0x4e, 0x54,
		0x45, 0x47, 0x45, 0x52, 0x20, 0x4e, 0x4f, 0x54, 0x20, 0x4e, 0x55, 0x4c,
		0x4c, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x52, 0x45,
		0x46, 0x45, 0x52, 0x45, 0x4e, 0x43, 0x45, 0x53, 0x20, 0x76, 0x65, 0x72,
		0x74, 0x69, 0x63, 0x65, 0x73, 0x28, 0x72, 0x6f, 0x77, 0x69, 0x64, 0x29,
		0x0a, 0x09, 0x4f, 0x4e, 0x20, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x20,
		0x43, 0x41, 0x53, 0x43, 0x41, 0x44, 0x45, 0x2c, 0x0a, 0x20, 0x20, 0x20,
		0x20, 0x74, 0x6f, 0x5f, 0x72, 0x6f, 0x77, 0x69, 0x64, 0x20, 0x49, 0x4e,
		0x54, 0x45, 0x47, 0x45, 0x52, 0x20, 0x4e, 0x4f, 0x54, 0x20, 0x4e, 0x55,
		0x4c, 0x4c, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x52,
		0x45, 0x46, 0x45, 0x52, 0x45, 0x4e, 0x43, 0x45, 0x53, 0x20, 0x76, 0x65,
		0x72, 0x74, 0x69, 0x63, 0x65, 0x73, 0x28, 0x72, 0x6f, 0x77, 0x69, 0x64,
		0x29, 0x0a, 0x09, 0x4f, 0x4e, 0x20, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
		0x20, 0x43, 0x41, 0x53, 0x43, 0x41, 0x44, 0x45, 0x2c, 0x0a, 0x20, 0x20,
		0x20, 0x20, 0x6b, 0x65, 0x79, 0x20, 0x54, 0x45, 0x58, 0x54, 0x20, 0x4e,
		0x4f, 0x54, 0x20, 0x4e, 0x55, 0x4c, 0x4c, 0x2c, 0x0a, 0x20, 0x20, 0x20,
		0x20, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x49, 0x4e,
		0x54, 0x45, 0x47, 0x45, 0x52, 0x2c, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x6d,
		0x65, 0x74, 0x61, 0x20, 0x54, 0x45, 0x58, 0x54, 0x0a, 0x29, 0x0a,
	},
	"CreateTableVertices": {
		0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x20, 0x54, 0x41, 0x42, 0x4c, 0x45,
		0x20, 0x49, 0x46, 0x20, 0x4e, 0x4f, 0x54, 0x20, 0x45, 0x58, 0x49, 0x53,
		0x54, 0x53, 0x20, 0x76, 0x65, 0x72, 0x74, 0x69, 0x63, 0x65, 0x73, 0x20,
		0x28, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x72, 0x6f, 0x77, 0x69, 0x64, 0x20,
		0x49, 0x4e, 0x54, 0x45, 0x47, 0x45, 0x52, 0x20, 0x50, 0x52, 0x49, 0x4d,
		0x41, 0x52, 0x59, 0x20, 0x4b, 0x45, 0x59, 0x2c, 0x0a, 0x20, 0x20, 0x20,
		0x20, 0x74, 0x79, 0x70, 0x65, 0x20, 0x54, 0x45, 0x58, 0x54, 0x20, 0x4e,
		0x4f, 0x54, 0x20, 0x4e, 0x55, 0x4c, 0x4c, 0x2c, 0x0a, 0x20, 0x20, 0x20,
		0x20, 0x69, 0x64, 0x20, 0x54, 0x45, 0x58, 0x54, 0x20, 0x4e, 0x4f, 0x54,
		0x20, 0x4e, 0x55, 0x4c, 0x4c, 0x2c, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x61,
		0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x20, 0x54, 0x45,
		0x58, 0x54, 0x2c, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x6d, 0x65, 0x74, 0x61,
		0x20, 0x54, 0x45, 0x58, 0x54, 0x0a, 0x29, 0x0a,
	},
}
