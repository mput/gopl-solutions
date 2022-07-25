package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type stack []xml.StartElement

func (st *stack) Push(el xml.StartElement) {
	*st = append(*st, el)
}

func (st *stack) Pop() xml.StartElement {
	el := (*st)[len(*st)-1]
	*st = (*st)[:len(*st)-1]
	return el
}


type Selector []xml.StartElement

func (s *Selector) String() string {
	res := make([]string, 0, len(*s))
	for _,e := range *s {
		es := e.Name.Local
		for _, atr := range e.Attr {
			var val string
			if atr.Value != "" {
				val = fmt.Sprintf("=%s", atr.Value)
			}
			es += fmt.Sprintf("[%s%s]", atr.Name.Local, val)
		}
		res = append(res, es)

	}

	return strings.Join(res, " ")
}

func main() {
	selector, err := parseSelector(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	dec := xml.NewDecoder(os.Stdin)
	stack := make(stack,0)

	for {
		token, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)

		}

		switch token := token.(type) {
		case xml.StartElement:
			stack.Push(token)
		case xml.EndElement:
			stack.Pop()
		case xml.CharData:
			if isMatches(stack, selector) {
				s := Selector(stack)
				fmt.Printf(`Selector: ~%s~*
Tags: ~%+v~
Value: "%s"

`, &selector, &s,token)
			}
		default:
			continue
		}

	}


}

func parseSelector(args []string) (Selector, error) {
	res := make(Selector, 0, len(args))
	r := regexp.MustCompile("(\\w+)|(?:\\[(\\w+)=?(\\w+)?\\])")
	for _, arg := range args {
		var elm xml.StartElement
		matches := r.FindAllStringSubmatch(arg, -1)
		for _, match := range matches {
			if match[1] != "" {
				elm.Name.Local = match[1]
				continue
			}
			if match[2] != "" {
				atr := xml.Attr{Name: xml.Name{Local: match[2]}, Value: match[3]}
				elm.Attr = append(elm.Attr, atr)
				continue
			}
			return res, fmt.Errorf("wrong argument format %s", arg)

		}
		res = append(res, elm)
	}
	return res, nil
}


func attrMap (v xml.StartElement) map[string]string {
	res := make(map[string]string)
	for _,atr := range v.Attr {
		res[atr.Name.Local] = atr.Value
	}
	return res
}

func isMatcesElms(x , tmplt xml.StartElement) bool {
	if tmplt.Name.Local != "" && x.Name.Local != tmplt.Name.Local {
		return false
	}
	xatr, tatr := attrMap(x), attrMap(tmplt)
	for k, v := range tatr {
		xv, ok := xatr[k]
		if !ok {
			return false
		}
		if v != "" && v != xv {
			return false
		}
	}

	return true
}

func isMatches(x []xml.StartElement, tmplt Selector) bool {
	xlen := len(x)
	tmpllen := len(tmplt)
	for {
		if tmpllen == 0 {
			return true
		}
		if xlen == 0 {
			return false
		}
		if isMatcesElms(x[xlen-1], tmplt[tmpllen-1]) {
			xlen -= 1
			tmpllen -=1
		} else {
			return false
		}
	}
}
