package main

import (
	"fmt"
	"log"
	"reflect"

	. "github.com/emicklei/proto"
)

type encoder struct {
	ptype string
	stack []map[string]interface{}
}

func (c *encoder) VisitMessage(m *Message)         { c.ptype = "Message" }
func (c *encoder) VisitService(v *Service)         { c.ptype = "Service" }
func (c *encoder) VisitSyntax(s *Syntax)           { c.ptype = "Syntax" }
func (c *encoder) VisitPackage(p *Package)         { c.ptype = "Package" }
func (c *encoder) VisitOption(o *Option)           { c.ptype = "Option" }
func (c *encoder) VisitImport(i *Import)           { c.ptype = "Import" }
func (c *encoder) VisitNormalField(i *NormalField) { c.ptype = "NormalField" }
func (c *encoder) VisitEnumField(i *EnumField)     { c.ptype = "EnumField" }
func (c *encoder) VisitEnum(e *Enum)               { c.ptype = "Enum" }
func (c *encoder) VisitComment(e *Comment)         { c.ptype = "Comment" }
func (c *encoder) VisitOneof(o *Oneof)             { c.ptype = "Oneof" }
func (c *encoder) VisitOneofField(o *OneOfField)   { c.ptype = "OneOfField" }
func (c *encoder) VisitReserved(rs *Reserved)      { c.ptype = "Reserved" }
func (c *encoder) VisitRPC(rpc *RPC)               { c.ptype = "RPC" }
func (c *encoder) VisitMapField(f *MapField)       { c.ptype = "MapField" }
func (c *encoder) VisitGroup(g *Group)             { c.ptype = "Group" }
func (c *encoder) VisitExtensions(e *Extensions)   { c.ptype = "Extensions" }

func toJSON(e Visitee) {
	c := new(encoder)
	c.stack = []map[string]interface{}{}
	e.Accept(c)
	fmt.Println(c.ptype)
}

func (c *encoder) Struct(v reflect.Value) error {
	// top := map[string]interface{}{}
	// top[v.S]
	return nil
}
func (c *encoder) StructField(s reflect.StructField, v reflect.Value) error {
	log.Println(s, v)
	return nil
}
