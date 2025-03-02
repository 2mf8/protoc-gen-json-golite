// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal_gengo

import (
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/compiler/protogen"
)

/*
	 func genReflectFileDescriptor(gen *protogen.Plugin, g *protogen.GeneratedFile, f *fileInfo) {
		g.P("var ", f.GoDescriptorIdent, " ", protoreflectPackage.Ident("FileDescriptor"))
		g.P()
		if len(f.allEnums) > 0 {
			g.P("var ", enumTypesVarName(f), " = make([]", protoimplPackage.Ident("EnumInfo"), ",", len(f.allEnums), ")")
		}
		if len(f.allMessages) > 0 {
			g.P("var ", messageTypesVarName(f), " = make([]", protoimplPackage.Ident("MessageInfo"), ",", len(f.allMessages), ")")
		}

		// Generate a unique list of Go types for all declarations and dependencies,
		// and the associated index into the type list for all dependencies.
		var goTypes []string
		var depIdxs []string
		seen := map[protoreflect.FullName]int{}
		genDep := func(name protoreflect.FullName, depSource string) {
			if depSource != "" {
				line := fmt.Sprintf("%d, // %d: %s -> %s", seen[name], len(depIdxs), depSource, name)
				depIdxs = append(depIdxs, line)
			}
		}
		genEnum := func(e *protogen.Enum, depSource string) {
			if e != nil {
				name := e.Desc.FullName()
				if _, ok := seen[name]; !ok {
					line := fmt.Sprintf("(%s)(0), // %d: %s", g.QualifiedGoIdent(e.GoIdent), len(goTypes), name)
					goTypes = append(goTypes, line)
					seen[name] = len(seen)
				}
				if depSource != "" {
					genDep(name, depSource)
				}
			}
		}
		genMessage := func(m *protogen.Message, depSource string) {
			if m != nil {
				name := m.Desc.FullName()
				if _, ok := seen[name]; !ok {
					line := fmt.Sprintf("(*%s)(nil), // %d: %s", g.QualifiedGoIdent(m.GoIdent), len(goTypes), name)
					if m.Desc.IsMapEntry() {
						// Map entry messages have no associated Go type.
						line = fmt.Sprintf("nil, // %d: %s", len(goTypes), name)
					}
					goTypes = append(goTypes, line)
					seen[name] = len(seen)
				}
				if depSource != "" {
					genDep(name, depSource)
				}
			}
		}

		// This ordering is significant.
		// See filetype.TypeBuilder.DependencyIndexes.
		type offsetEntry struct {
			start int
			name  string
		}
		var depOffsets []offsetEntry
		for _, enum := range f.allEnums {
			genEnum(enum.Enum, "")
		}
		for _, message := range f.allMessages {
			genMessage(message.Message, "")
		}
		depOffsets = append(depOffsets, offsetEntry{len(depIdxs), "field type_name"})
		for _, message := range f.allMessages {
			for _, field := range message.Fields {
				source := string(field.Desc.FullName())
				genEnum(field.Enum, source+":type_name")
				genMessage(field.Message, source+":type_name")
			}
		}
		depOffsets = append(depOffsets, offsetEntry{len(depIdxs), "extension extendee"})
		for _, extension := range f.allExtensions {
			source := string(extension.Desc.FullName())
			genMessage(extension.Extendee, source+":extendee")
		}
		depOffsets = append(depOffsets, offsetEntry{len(depIdxs), "extension type_name"})
		for _, extension := range f.allExtensions {
			source := string(extension.Desc.FullName())
			genEnum(extension.Enum, source+":type_name")
			genMessage(extension.Message, source+":type_name")
		}
		depOffsets = append(depOffsets, offsetEntry{len(depIdxs), "method input_type"})
		for _, service := range f.Services {
			for _, method := range service.Methods {
				source := string(method.Desc.FullName())
				genMessage(method.Input, source+":input_type")
			}
		}
		depOffsets = append(depOffsets, offsetEntry{len(depIdxs), "method output_type"})
		for _, service := range f.Services {
			for _, method := range service.Methods {
				source := string(method.Desc.FullName())
				genMessage(method.Output, source+":output_type")
			}
		}
		depOffsets = append(depOffsets, offsetEntry{len(depIdxs), ""})
		for i := len(depOffsets) - 2; i >= 0; i-- {
			curr, next := depOffsets[i], depOffsets[i+1]
			depIdxs = append(depIdxs, fmt.Sprintf("%d, // [%d:%d] is the sub-list for %s",
				curr.start, curr.start, next.start, curr.name))
		}
		if len(depIdxs) > math.MaxInt32 {
			panic("too many dependencies") // sanity check
		}

		g.P("var ", depIdxsVarName(f), " = []int32{")
		for _, s := range depIdxs {
			g.P(s)
		}
		g.P("}")

		g.P("func init() { ", initFuncName(f.File), "() }")

		g.P("func ", initFuncName(f.File), "() {")
		g.P("if ", f.GoDescriptorIdent, " != nil {")
		g.P("return")
		g.P("}")

		// Ensure that initialization functions for different files in the same Go
		// package run in the correct order: Call the init funcs for every .proto file
		// imported by this one that is in the same Go package.
		for i, imps := 0, f.Desc.Imports(); i < imps.Len(); i++ {
			impFile := gen.FilesByPath[imps.Get(i).Path()]
			if impFile.GoImportPath != f.GoImportPath {
				continue
			}
			g.P(initFuncName(impFile), "()")
		}

		if len(f.allMessages) > 0 {
			// Populate MessageInfo.OneofWrappers.
			for _, message := range f.allMessages {
				if len(message.Oneofs) > 0 {
					idx := f.allMessagesByPtr[message]
					typesVar := messageTypesVarName(f)

					// Associate the wrapper types by directly passing them to the MessageInfo.
					g.P(typesVar, "[", idx, "].OneofWrappers = []any {")
					for _, oneof := range message.Oneofs {
						if !oneof.Desc.IsSynthetic() {
							for _, field := range oneof.Fields {
								g.P("(*", unexportIdent(field.GoIdent, message.isOpaque()), ")(nil),")
							}
						}
					}
					g.P("}")
				}
			}
		}

		g.P("type x struct{}")
		g.P("out := ", protoimplPackage.Ident("TypeBuilder"), "{")
		g.P("File: ", protoimplPackage.Ident("DescBuilder"), "{")
		g.P("GoPackagePath: ", reflectPackage.Ident("TypeOf"), "(x{}).PkgPath(),")
		// Avoid a copy of the descriptor. This means modification of the
		// RawDescriptor byte slice will crash the program. But generated
		// RawDescriptors are never supposed to be modified anyway.
		g.P("NumEnums: ", len(f.allEnums), ",")
		g.P("NumMessages: ", len(f.allMessages), ",")
		g.P("NumExtensions: ", len(f.allExtensions), ",")
		g.P("NumServices: ", len(f.Services), ",")
		g.P("},")
		g.P("GoTypes: ", goTypesVarName(f), ",")
		g.P("DependencyIndexes: ", depIdxsVarName(f), ",")
		if len(f.allEnums) > 0 {
			g.P("EnumInfos: ", enumTypesVarName(f), ",")
		}
		if len(f.allMessages) > 0 {
			g.P("MessageInfos: ", messageTypesVarName(f), ",")
		}
		if len(f.allExtensions) > 0 {
			g.P("ExtensionInfos: ", extensionTypesVarName(f), ",")
		}
		g.P("}.Build()")
		g.P(f.GoDescriptorIdent, " = out.File")

		// Set inputs to nil to allow GC to reclaim resources.
		g.P(goTypesVarName(f), " = nil")
		g.P(depIdxsVarName(f), " = nil")
		g.P("}")
	}
*/
func fileVarName(f *protogen.File, suffix string) string {
	prefix := f.GoDescriptorIdent.GoName
	_, n := utf8.DecodeRuneInString(prefix)
	prefix = strings.ToLower(prefix[:n]) + prefix[n:]
	return prefix + "_" + suffix
}
func goTypesVarName(f *fileInfo) string {
	return fileVarName(f.File, "goTypes")
}
func depIdxsVarName(f *fileInfo) string {
	return fileVarName(f.File, "depIdxs")
}
func enumTypesVarName(f *fileInfo) string {
	return fileVarName(f.File, "enumTypes")
}
func messageTypesVarName(f *fileInfo) string {
	return fileVarName(f.File, "msgTypes")
}
func extensionTypesVarName(f *fileInfo) string {
	return fileVarName(f.File, "extTypes")
}
func initFuncName(f *protogen.File) string {
	return fileVarName(f, "init")
}
