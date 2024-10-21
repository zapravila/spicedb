// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: google/type/expr.proto

package _type

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Represents a textual expression in the Common Expression Language (CEL)
// syntax. CEL is a C-like expression language. The syntax and semantics of CEL
// are documented at https://github.com/google/cel-spec.
//
// Example (Comparison):
//
//	title: "Summary size limit"
//	description: "Determines if a summary is less than 100 chars"
//	expression: "document.summary.size() < 100"
//
// Example (Equality):
//
//	title: "Requestor is owner"
//	description: "Determines if requestor is the document owner"
//	expression: "document.owner == request.auth.claims.email"
//
// Example (Logic):
//
//	title: "Public documents"
//	description: "Determine whether the document should be publicly visible"
//	expression: "document.type != 'private' && document.type != 'internal'"
//
// Example (Data Manipulation):
//
//	title: "Notification string"
//	description: "Create a notification string with a timestamp."
//	expression: "'New message received at ' + string(document.create_time)"
//
// The exact variables and functions that may be referenced within an expression
// are determined by the service that evaluates it. See the service
// documentation for additional information.
type Expr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Textual representation of an expression in Common Expression Language
	// syntax.
	Expression string `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
	// Optional. Title for the expression, i.e. a short string describing
	// its purpose. This can be used e.g. in UIs which allow to enter the
	// expression.
	Title string `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	// Optional. Description of the expression. This is a longer text which
	// describes the expression, e.g. when hovered over it in a UI.
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	// Optional. String indicating the location of the expression for error
	// reporting, e.g. a file name and a position in the file.
	Location string `protobuf:"bytes,4,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *Expr) Reset() {
	*x = Expr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_google_type_expr_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Expr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Expr) ProtoMessage() {}

func (x *Expr) ProtoReflect() protoreflect.Message {
	mi := &file_google_type_expr_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Expr.ProtoReflect.Descriptor instead.
func (*Expr) Descriptor() ([]byte, []int) {
	return file_google_type_expr_proto_rawDescGZIP(), []int{0}
}

func (x *Expr) GetExpression() string {
	if x != nil {
		return x.Expression
	}
	return ""
}

func (x *Expr) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Expr) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Expr) GetLocation() string {
	if x != nil {
		return x.Location
	}
	return ""
}

var File_google_type_expr_proto protoreflect.FileDescriptor

var file_google_type_expr_proto_rawDesc = []byte{
	0x0a, 0x16, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x2f, 0x65, 0x78,
	0x70, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x22, 0x7a, 0x0a, 0x04, 0x45, 0x78, 0x70, 0x72, 0x12, 0x1e, 0x0a,
	0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x72, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x9b, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x74, 0x79, 0x70, 0x65, 0x42, 0x09, 0x45, 0x78, 0x70, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x50, 0x01, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x7a, 0x65, 0x64, 0x2f, 0x73, 0x70, 0x69, 0x63, 0x65, 0x64, 0x62, 0x2f, 0x70,
	0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x74, 0x79, 0x70, 0x65, 0xa2, 0x02, 0x03, 0x47, 0x54, 0x58, 0xaa, 0x02, 0x0b, 0x47, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0xca, 0x02, 0x0b, 0x47, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x5c, 0x54, 0x79, 0x70, 0x65, 0xe2, 0x02, 0x17, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x5c,
	0x54, 0x79, 0x70, 0x65, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x0c, 0x47, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x3a, 0x3a, 0x54, 0x79, 0x70, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_google_type_expr_proto_rawDescOnce sync.Once
	file_google_type_expr_proto_rawDescData = file_google_type_expr_proto_rawDesc
)

func file_google_type_expr_proto_rawDescGZIP() []byte {
	file_google_type_expr_proto_rawDescOnce.Do(func() {
		file_google_type_expr_proto_rawDescData = protoimpl.X.CompressGZIP(file_google_type_expr_proto_rawDescData)
	})
	return file_google_type_expr_proto_rawDescData
}

var file_google_type_expr_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_google_type_expr_proto_goTypes = []any{
	(*Expr)(nil), // 0: google.type.Expr
}
var file_google_type_expr_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_google_type_expr_proto_init() }
func file_google_type_expr_proto_init() {
	if File_google_type_expr_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_google_type_expr_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Expr); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_google_type_expr_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_google_type_expr_proto_goTypes,
		DependencyIndexes: file_google_type_expr_proto_depIdxs,
		MessageInfos:      file_google_type_expr_proto_msgTypes,
	}.Build()
	File_google_type_expr_proto = out.File
	file_google_type_expr_proto_rawDesc = nil
	file_google_type_expr_proto_goTypes = nil
	file_google_type_expr_proto_depIdxs = nil
}
