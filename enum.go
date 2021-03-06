// Copyright (c) 2017 Ernest Micklei
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proto

import "strconv"

// Enum definition consists of a name and an enum body.
type Enum struct {
	Comment  *Comment
	Name     string
	Elements []Visitee
}

// Accept dispatches the call to the visitor.
func (e *Enum) Accept(v Visitor) {
	v.VisitEnum(e)
}

// Doc is part of Documented
func (e *Enum) Doc() *Comment {
	return e.Comment
}

// addElement is part of elementContainer
func (e *Enum) addElement(v Visitee) {
	e.Elements = append(e.Elements, v)
}

// elements is part of elementContainer
func (e *Enum) elements() []Visitee {
	return e.Elements
}

// takeLastComment is part of elementContainer
// removes and returns the last element of the list if it is a Comment.
func (e *Enum) takeLastComment() (last *Comment) {
	last, e.Elements = takeLastComment(e.Elements)
	return
}

func (e *Enum) parse(p *Parser) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != tIDENT {
		if !isKeyword(tok) {
			return p.unexpected(lit, "enum identifier", e)
		}
	}
	e.Name = lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != tLEFTCURLY {
		return p.unexpected(lit, "enum opening {", e)
	}
	for {
		tok, lit = p.scanIgnoreWhitespace()
		switch tok {
		case tCOMMENT:
			if com := mergeOrReturnComment(e.elements(), lit, p.s.line); com != nil { // not merged?
				e.Elements = append(e.Elements, com)
			}
		case tOPTION:
			v := new(Option)
			v.Comment = e.takeLastComment()
			err := v.parse(p)
			if err != nil {
				return err
			}
			e.Elements = append(e.Elements, v)
		case tRIGHTCURLY, tEOF:
			goto done
		case tSEMICOLON:
			maybeScanInlineComment(p, e)
		default:
			p.unscan()
			f := new(EnumField)
			f.Comment = e.takeLastComment()
			err := f.parse(p)
			if err != nil {
				return err
			}
			e.Elements = append(e.Elements, f)
		}
	}
done:
	if tok != tRIGHTCURLY {
		return p.unexpected(lit, "enum closing }", e)
	}
	return nil
}

// EnumField is part of the body of an Enum.
type EnumField struct {
	Comment       *Comment
	Name          string
	Integer       int
	ValueOption   *Option
	InlineComment *Comment
}

// Accept dispatches the call to the visitor.
func (f *EnumField) Accept(v Visitor) {
	v.VisitEnumField(f)
}

// inlineComment is part of commentInliner.
func (f *EnumField) inlineComment(c *Comment) {
	f.InlineComment = c
}

// Doc is part of Documented
func (f *EnumField) Doc() *Comment {
	return f.Comment
}

// columns returns printable source tokens
func (f EnumField) columns() (cols []aligned) {
	cols = append(cols, leftAligned(f.Name), alignedEquals, rightAligned(strconv.Itoa(f.Integer)))
	if f.ValueOption != nil {
		cols = append(cols, f.ValueOption.columns()...)
	}
	cols = append(cols, alignedSemicolon)
	if f.InlineComment != nil {
		cols = append(cols, f.InlineComment.alignedInlinePrefix(), notAligned(f.InlineComment.Message()))
	}
	return
}

func (f *EnumField) parse(p *Parser) error {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != tIDENT {
		if !isKeyword(tok) {
			return p.unexpected(lit, "enum field identifier", f)
		}
	}
	f.Name = lit
	tok, lit = p.scanIgnoreWhitespace()
	if tok != tEQUALS {
		return p.unexpected(lit, "enum field =", f)
	}
	i, err := p.s.scanInteger()
	if err != nil {
		return p.unexpected(lit, "enum field integer", f)
	}
	f.Integer = i
	tok, lit = p.scanIgnoreWhitespace()
	if tok == tLEFTSQUARE {
		o := new(Option)
		o.IsEmbedded = true
		err := o.parse(p)
		if err != nil {
			return err
		}
		f.ValueOption = o
		tok, lit = p.scanIgnoreWhitespace()
		if tok != tRIGHTSQUARE {
			return p.unexpected(lit, "option closing ]", f)
		}
	}
	if tSEMICOLON == tok {
		p.unscan() // put back this token for scanning inline comment
	}
	return nil
}
