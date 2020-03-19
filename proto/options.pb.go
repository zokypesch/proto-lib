// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/options.proto

package options

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

var E_HttpMode = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50056,
	Name:          "httpMode",
	Tag:           "bytes,50056,opt,name=httpMode",
	Filename:      "proto/options.proto",
}

var E_Agregator = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50062,
	Name:          "agregator",
	Tag:           "bytes,50062,opt,name=agregator",
	Filename:      "proto/options.proto",
}

var E_Withelist = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50067,
	Name:          "withelist",
	Tag:           "varint,50067,opt,name=withelist",
	Filename:      "proto/options.proto",
}

var E_Foo = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50001,
	Name:          "foo",
	Tag:           "bytes,50001,opt,name=foo",
	Filename:      "proto/options.proto",
}

var E_IsRepo = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50057,
	Name:          "isRepo",
	Tag:           "varint,50057,opt,name=isRepo",
	Filename:      "proto/options.proto",
}

var E_IsEs = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50064,
	Name:          "isEs",
	Tag:           "varint,50064,opt,name=isEs",
	Filename:      "proto/options.proto",
}

var E_IgnoreFieldDb = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50058,
	Name:          "ignoreFieldDb",
	Tag:           "varint,50058,opt,name=ignoreFieldDb",
	Filename:      "proto/options.proto",
}

var E_IsPrimaryKey = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50059,
	Name:          "isPrimaryKey",
	Tag:           "varint,50059,opt,name=isPrimaryKey",
	Filename:      "proto/options.proto",
}

var E_Required = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50060,
	Name:          "required",
	Tag:           "varint,50060,opt,name=required",
	Filename:      "proto/options.proto",
}

var E_RequiredType = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50061,
	Name:          "required_type",
	Tag:           "bytes,50061,opt,name=required_type,json=requiredType",
	Filename:      "proto/options.proto",
}

var E_Fulltext = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50063,
	Name:          "fulltext",
	Tag:           "varint,50063,opt,name=fulltext",
	Filename:      "proto/options.proto",
}

var E_ErrDesc = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50065,
	Name:          "errDesc",
	Tag:           "bytes,50065,opt,name=errDesc",
	Filename:      "proto/options.proto",
}

var E_ForeignKey = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50066,
	Name:          "foreignKey",
	Tag:           "bytes,50066,opt,name=foreignKey",
	Filename:      "proto/options.proto",
}

var E_AssociateKey = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50078,
	Name:          "associateKey",
	Tag:           "bytes,50078,opt,name=associateKey",
	Filename:      "proto/options.proto",
}

var E_Integration = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         50068,
	Name:          "integration",
	Tag:           "varint,50068,opt,name=integration",
	Filename:      "proto/options.proto",
}

var E_ProtoFileLoc = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50069,
	Name:          "protoFileLoc",
	Tag:           "bytes,50069,opt,name=protoFileLoc",
	Filename:      "proto/options.proto",
}

var E_GrpcMethod = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50070,
	Name:          "grpcMethod",
	Tag:           "bytes,50070,opt,name=grpcMethod",
	Filename:      "proto/options.proto",
}

var E_GrpcAddress = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50071,
	Name:          "grpcAddress",
	Tag:           "bytes,50071,opt,name=grpcAddress",
	Filename:      "proto/options.proto",
}

var E_ProtoDomain = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50072,
	Name:          "protoDomain",
	Tag:           "bytes,50072,opt,name=protoDomain",
	Filename:      "proto/options.proto",
}

var E_GrpcPort = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50073,
	Name:          "grpcPort",
	Tag:           "bytes,50073,opt,name=grpcPort",
	Filename:      "proto/options.proto",
}

var E_GrpcRequestName = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50074,
	Name:          "grpcRequestName",
	Tag:           "bytes,50074,opt,name=grpcRequestName",
	Filename:      "proto/options.proto",
}

var E_GrpcResponseName = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50075,
	Name:          "grpcResponseName",
	Tag:           "bytes,50075,opt,name=grpcResponseName",
	Filename:      "proto/options.proto",
}

var E_GrpcRequestMessage = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50076,
	Name:          "grpcRequestMessage",
	Tag:           "bytes,50076,opt,name=grpcRequestMessage",
	Filename:      "proto/options.proto",
}

var E_GrpcResponseMessage = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         50077,
	Name:          "grpcResponseMessage",
	Tag:           "bytes,50077,opt,name=grpcResponseMessage",
	Filename:      "proto/options.proto",
}

func init() {
	proto.RegisterExtension(E_HttpMode)
	proto.RegisterExtension(E_Agregator)
	proto.RegisterExtension(E_Withelist)
	proto.RegisterExtension(E_Foo)
	proto.RegisterExtension(E_IsRepo)
	proto.RegisterExtension(E_IsEs)
	proto.RegisterExtension(E_IgnoreFieldDb)
	proto.RegisterExtension(E_IsPrimaryKey)
	proto.RegisterExtension(E_Required)
	proto.RegisterExtension(E_RequiredType)
	proto.RegisterExtension(E_Fulltext)
	proto.RegisterExtension(E_ErrDesc)
	proto.RegisterExtension(E_ForeignKey)
	proto.RegisterExtension(E_AssociateKey)
	proto.RegisterExtension(E_Integration)
	proto.RegisterExtension(E_ProtoFileLoc)
	proto.RegisterExtension(E_GrpcMethod)
	proto.RegisterExtension(E_GrpcAddress)
	proto.RegisterExtension(E_ProtoDomain)
	proto.RegisterExtension(E_GrpcPort)
	proto.RegisterExtension(E_GrpcRequestName)
	proto.RegisterExtension(E_GrpcResponseName)
	proto.RegisterExtension(E_GrpcRequestMessage)
	proto.RegisterExtension(E_GrpcResponseMessage)
}

func init() { proto.RegisterFile("proto/options.proto", fileDescriptor_options_94601d2cd59a601b) }

var fileDescriptor_options_94601d2cd59a601b = []byte{
	// 513 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0x4b, 0x6f, 0x13, 0x31,
	0x10, 0xc7, 0x85, 0x8a, 0x4a, 0x3b, 0xb4, 0x02, 0xa5, 0x17, 0x84, 0x04, 0xf4, 0xc8, 0x29, 0x39,
	0x54, 0x1c, 0x6a, 0x10, 0xa8, 0x22, 0xad, 0x84, 0x4a, 0x69, 0x89, 0xb8, 0x23, 0x27, 0x3b, 0x71,
	0x2c, 0x6d, 0x76, 0x5c, 0xdb, 0x11, 0xe4, 0x0b, 0xf0, 0x7e, 0xbf, 0xdf, 0xdf, 0x89, 0x6f, 0x84,
	0xbc, 0x3b, 0x6e, 0xb6, 0x34, 0x92, 0xf7, 0xb6, 0x4e, 0xfe, 0xbf, 0xdf, 0xcc, 0x8e, 0xbd, 0x86,
	0x35, 0x63, 0xc9, 0x53, 0x87, 0x8c, 0xd7, 0x54, 0xb8, 0x76, 0xb9, 0xba, 0xb8, 0xae, 0x88, 0x54,
	0x8e, 0x9d, 0x72, 0xd5, 0x9f, 0x0c, 0x3b, 0x19, 0xba, 0x81, 0xd5, 0xc6, 0x93, 0xad, 0x12, 0xe2,
	0x06, 0x2c, 0x8d, 0xbc, 0x37, 0x7b, 0x94, 0x61, 0xeb, 0x72, 0xbb, 0x8a, 0xb7, 0x63, 0xbc, 0xbd,
	0x87, 0x7e, 0x44, 0xd9, 0x7e, 0xe5, 0xbc, 0xf0, 0xec, 0xe9, 0xc2, 0xfa, 0xa9, 0xab, 0xcb, 0xbd,
	0x23, 0x42, 0xdc, 0x84, 0x65, 0xa9, 0x2c, 0x2a, 0xe9, 0xc9, 0x26, 0xf1, 0x37, 0x8c, 0xcf, 0x90,
	0xc0, 0x3f, 0xd2, 0x7e, 0x84, 0xb9, 0x76, 0x3e, 0xc9, 0x7f, 0x2c, 0xf9, 0xa5, 0xde, 0x0c, 0x11,
	0x1b, 0xb0, 0x30, 0x24, 0x6a, 0x5d, 0x99, 0x43, 0x3a, 0x27, 0x15, 0x46, 0xf4, 0xef, 0x93, 0xaa,
	0x74, 0x48, 0x8b, 0x4d, 0x58, 0xd4, 0xae, 0x87, 0xa6, 0x01, 0xf7, 0x9c, 0x4b, 0x32, 0x20, 0xae,
	0xc1, 0x69, 0xed, 0xb6, 0x5d, 0x1a, 0x7c, 0xc7, 0x60, 0x19, 0x17, 0xdb, 0xb0, 0xaa, 0x55, 0x41,
	0x16, 0x77, 0x34, 0xe6, 0x59, 0xb7, 0xdf, 0xba, 0x74, 0x82, 0x2f, 0xff, 0x89, 0xf4, 0x0b, 0xa6,
	0x8f, 0x53, 0xe2, 0x36, 0xac, 0x68, 0x77, 0x60, 0xf5, 0x58, 0xda, 0xe9, 0x2e, 0x4e, 0x53, 0x96,
	0x97, 0x6c, 0x39, 0x06, 0x89, 0xeb, 0xb0, 0x64, 0xf1, 0x70, 0xa2, 0x2d, 0x66, 0x29, 0xc1, 0x2b,
	0x16, 0x1c, 0x01, 0xa2, 0x0b, 0xab, 0xf1, 0xf9, 0xa1, 0x9f, 0x1a, 0x4c, 0x19, 0x5e, 0xf3, 0x96,
	0xaf, 0x44, 0xea, 0xc1, 0xd4, 0x60, 0x68, 0x61, 0x38, 0xc9, 0x73, 0x8f, 0x8f, 0x7d, 0x4a, 0xf0,
	0x36, 0xb6, 0x10, 0x01, 0xb1, 0x09, 0x67, 0xd0, 0xda, 0x2e, 0xba, 0x41, 0x8a, 0x7d, 0xcf, 0xc5,
	0x63, 0x5e, 0xdc, 0x02, 0x18, 0x92, 0x45, 0xad, 0x8a, 0x06, 0xd3, 0xfb, 0xc0, 0x74, 0x0d, 0x09,
	0x1b, 0x20, 0x9d, 0xa3, 0x81, 0x96, 0x1e, 0x1b, 0x28, 0xfe, 0xc4, 0xb7, 0xaf, 0x43, 0x62, 0x0b,
	0xce, 0xea, 0xc2, 0xa3, 0xb2, 0x32, 0xa4, 0x52, 0x8e, 0x4f, 0x3c, 0x80, 0x3a, 0x13, 0xfa, 0x28,
	0xd3, 0x3b, 0x3a, 0xc7, 0xbb, 0x94, 0x1c, 0xc4, 0xe7, 0xd8, 0x47, 0x1d, 0x0a, 0xd3, 0x50, 0xd6,
	0x0c, 0xaa, 0x6f, 0x2c, 0xa5, 0xf8, 0x12, 0xa7, 0x31, 0x43, 0xc2, 0x8b, 0x84, 0xd5, 0x56, 0x96,
	0x59, 0x74, 0x2e, 0x65, 0xf8, 0xca, 0x86, 0x3a, 0x13, 0x14, 0x65, 0xba, 0x4b, 0x63, 0xa9, 0x93,
	0xb3, 0xf8, 0x16, 0x15, 0x35, 0x26, 0x1c, 0xa6, 0x60, 0x3c, 0x20, 0x9b, 0x3c, 0x4c, 0xdf, 0xe3,
	0xfd, 0x15, 0x01, 0x71, 0x07, 0xce, 0x85, 0xe7, 0x1e, 0x1e, 0x4e, 0xd0, 0xf9, 0x7b, 0x72, 0x9c,
	0x3c, 0xd1, 0x3f, 0xd8, 0xf1, 0x3f, 0x27, 0x76, 0xe1, 0x7c, 0xf5, 0x93, 0x33, 0x54, 0x38, 0x6c,
	0xe2, 0xfa, 0xc9, 0xae, 0x13, 0xa0, 0xd8, 0x87, 0x56, 0xcd, 0xcf, 0x77, 0x4b, 0x4a, 0xf7, 0x8b,
	0x75, 0x73, 0x50, 0x71, 0x1f, 0xd6, 0xea, 0x45, 0x1a, 0x1a, 0x7f, 0xb3, 0x71, 0x1e, 0xdb, 0x5f,
	0x2c, 0x99, 0x8d, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x7f, 0xc9, 0xd1, 0x61, 0x79, 0x06, 0x00,
	0x00,
}
