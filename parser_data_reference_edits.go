package mayaascii

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type DataReferenceEditsParser struct {
	token *[]string
	errs  []string
	cur   int

	CurToken  *string
	PeekToken *string
}

func (p *DataReferenceEditsParser) parseReferenceEdit() (*ReferenceEdit, error) {
	re := &ReferenceEdit{
		ReferenceNode: strings.Trim(*p.CurToken, "\""),
	}
	if !p.PeekTokenIsNumber() {
		return nil, errors.New("parseReferenceEditError: commandNum is not int")
	}
	p.NextToken() // skip reference node name
	re.CommandNum, _ = strconv.Atoi(*p.CurToken)
	return re, nil
}

func (p *DataReferenceEditsParser) parseParent() (*RECmdParent, error) {
	p.NextToken() // skip 0
	prt := &RECmdParent{
		NodeA: strings.Trim(*p.CurToken, "\""),
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseParentError: not enough tokens")
	}
	p.NextToken() // skip nodeA
	prt.NodeB = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseParentError: not enough tokens")
	}
	p.NextToken() // skip nodeB
	prt.Arguments = (*p.CurToken)[1:len(*p.CurToken)-1]
	return prt, nil
}

func (p *DataReferenceEditsParser) parseAddAttr() (*RECmdAddAttr, error) {
	p.NextToken() // skip 1
	aa := &RECmdAddAttr{
		Node: *p.CurToken,
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseAddAttrError: not enough tokens")
	}
	p.NextToken() // skip node name
	aa.LongAttr = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseAddAttrError: not enough tokens")
	}
	p.NextToken() // skip long attr name
	aa.ShortAttr = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseAddAttrError: not enough tokens")
	}
	p.NextToken() // skip shot attr name
	aa.Arguments = (*p.CurToken)[1:len(*p.CurToken)-1]
	return aa, nil
}

func (p *DataReferenceEditsParser) parseSetAttr() (*RECmdSetAttr, error) {
	//fmt.Println("setAttr 2 ", *p.CurToken)
	p.NextToken() // skip 2
	//fmt.Println("setAttr Node ", *p.CurToken)
	sa := &RECmdSetAttr{
		Node: strings.Trim(*p.CurToken, "\""),
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseSetAttrError: not enough tokens")
	}
	p.NextToken() // skip node name
	//fmt.Println("setAttr AttrValue ", *p.CurToken)
	sa.Attr = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseSetAttrError: not enough tokens")
	}
	p.NextToken() // skip attr name
	//fmt.Println("setAttr Arguments ", *p.CurToken)
	sa.Arguments = (*p.CurToken)[1:len(*p.CurToken)-1]
	return sa, nil
}

func (p *DataReferenceEditsParser) parseDisconnectAttr() (*RECmdDisconnectAttr, error) {
	p.NextToken() // skip 3
	da := &RECmdDisconnectAttr{
		SourcePlug: strings.Trim(*p.CurToken, "\""),
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseDisconnectAttr: not enough tokens")
	}
	p.NextToken() // skip source plug
	da.DistPlug = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseDisconnectAttr: not enough tokens")
	}
	p.NextToken()
	da.Arguments = (*p.CurToken)[1:len(*p.CurToken)-1]
	return da, nil
}

func (p *DataReferenceEditsParser) parseDeleteAttr() (*RECmdDeleteAttr, error) {
	p.NextToken() // skip 4
	da := &RECmdDeleteAttr{
		Node: strings.Trim(*p.CurToken, "\""),
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseDeleteAttr: not enough tokens")
	}
	p.NextToken()
	da.Attr = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseDeleteAttr: not enough tokens")
	}
	p.NextToken()
	da.Arguments = (*p.CurToken)[1:len(*p.CurToken)-1]
	return da, nil
}

func (p *DataReferenceEditsParser) parseConnectAttr() (*RECmdConnectAttr, error) {
	p.NextToken() // skip 5
	magic, err := strconv.Atoi(*p.CurToken)
	if err != nil {
		return nil, err
	}
	ca := &RECmdConnectAttr{
		MagicNumber: magic,
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseConnectAttr: not enough tokens")
	}
	p.NextToken()
	ca.ReferenceNode = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseConnectAttr: not enough tokens")
	}
	p.NextToken()
	ca.SourcePlug = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseConnectAttr: not enough tokens")
	}
	p.NextToken()
	ca.DistPlug = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil {
		return nil, errors.New("parseConnectAttr: not enough tokens")
	}
	p.NextToken()
	if ca.MagicNumber == 0 {
		sourcePHL := strings.Trim(*p.CurToken, "\"")
		ca.SourcePHL = &sourcePHL
		if p.PeekToken == nil {
			return nil, errors.New("parseConnectAttr: not enough tokens")
		}
		p.NextToken()
		distPHL := strings.Trim(*p.CurToken, "\"")
		ca.DistPHL = &distPHL
		if p.PeekToken == nil {
			return nil, errors.New("parseConnectAttr: not enough tokens")
		}
		p.NextToken()
	}
	ca.Arguments = (*p.CurToken)[1:len(*p.CurToken)-1]
	return ca, nil
}

func (p *DataReferenceEditsParser) parseRelationship() (*RECmdRelationship, error) {
	p.NextToken() // skip 7
	rs := &RECmdRelationship{
		Type: strings.Trim(*p.CurToken, "\""),
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseRelationship: not enough tokens")
	}
	p.NextToken()
	rs.NodeName = strings.Trim(*p.CurToken, "\"")
	if p.PeekToken == nil && !p.PeekTokenIsNumber() {
		return nil, errors.New("parseRelationship: not enough tokens")
	}
	p.NextToken()
	var err error
	rs.CommandNum, err = strconv.Atoi(*p.CurToken)
	if err != nil {
		return nil, fmt.Errorf("rs.CommandNum strconv.Atoi(\"%s\"), %s",
			*p. CurToken, err)
	}
	if rs.CommandNum > 0 {
		rs.Commands = []string{}
	}
	for i := 0; i < rs.CommandNum; i++ {
		if p.PeekToken == nil && !p.PeekTokenIsNumber() {
			return nil, errors.New("parseRelationship: not enough tokens")
		}
		p.NextToken()
		rs.Commands = append(rs.Commands, strings.Trim(*p.CurToken, "\""))
	}
	if p.PeekToken == nil && !p.PeekTokenIsNumber() {
		return nil, errors.New("parseRelationship: not enough tokens")
	}
	p.NextToken() // skip last 0.
	return rs, nil
}

func (p *DataReferenceEditsParser) parseLock() (*RECmdLock, error) {
	p.NextToken() // skip 8
	ulk := &RECmdLock{
		Node: strings.Trim(*p.CurToken, "\""),
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseLock: not enough tokens")
	}
	p.NextToken() // skip node
	ulk.Attr = strings.Trim(*p.CurToken, "\"")
	return ulk, nil
}

func (p *DataReferenceEditsParser) parseUnlock() (*RECmdUnlock, error) {
	p.NextToken() // skip 9
	ulk := &RECmdUnlock{
		Node: strings.Trim(*p.CurToken, "\""),
	}
	if p.PeekToken == nil {
		return nil, errors.New("parseUnlock: not enough tokens")
	}
	p.NextToken() // skip node
	ulk.Attr = strings.Trim(*p.CurToken, "\"")
	return ulk, nil
}

func (p *DataReferenceEditsParser) ParseToken() []*ReferenceEdit {
	var res []*ReferenceEdit
	commandNum := 0
	for p.CurToken != nil {
		var err error
		if commandNum == 0 {
			var re *ReferenceEdit
			re, err = p.parseReferenceEdit()
			if res == nil {
				res = []*ReferenceEdit{}
			}
			res = append(res, re)
			if re != nil {
				commandNum = (*re).CommandNum
			} else {
				fmt.Printf("commandNum == 0, referenceEdit is nil... %s %d %s %s %s %s\n",
					*p.CurToken, p.cur, (*p.token)[p.cur-2], (*p.token)[p.cur-1], (*p.token)[p.cur], (*p.token)[p.cur+1])
			}
		} else if p.CurTokenIs(string(RETypePArent)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var prt *RECmdParent
			prt, err = p.parseParent()
			if re.Parents == nil {
				re.Parents = []*RECmdParent{}
			}
			re.Parents = append(re.Parents, prt)
			commandNum--
		} else if p.CurTokenIs(string(RETypeAddAttr)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var add *RECmdAddAttr
			add, err = p.parseAddAttr()
			if re.AddAttrs == nil {
				re.AddAttrs = []*RECmdAddAttr{}
			}
			re.AddAttrs = append(re.AddAttrs, add)
			commandNum--
		} else if p.CurTokenIs(string(RETypeSetAttr)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var set *RECmdSetAttr
			set, err = p.parseSetAttr()
			if re.SetAttrs == nil {
				re.SetAttrs = []*RECmdSetAttr{}
			}
			re.SetAttrs = append(re.SetAttrs, set)
			commandNum--
		} else if p.CurTokenIs(string(RETypeDisconnectAttr)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var dis *RECmdDisconnectAttr
			dis, err = p.parseDisconnectAttr()
			if re.DisconnectAttrs == nil {
				re.DisconnectAttrs = []*RECmdDisconnectAttr{}
			}
			re.DisconnectAttrs = append(re.DisconnectAttrs, dis)
			commandNum--
		} else if p.CurTokenIs(string(RETypeDeleteAttr)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var del *RECmdDeleteAttr
			del, err = p.parseDeleteAttr()
			if re.DeleteAttrs == nil {
				re.DeleteAttrs = []*RECmdDeleteAttr{}
			}
			re.DeleteAttrs = append(re.DeleteAttrs, del)
			commandNum--
		} else if p.CurTokenIs(string(RETypeConnectAttr)) && p.PeekTokenIsNumber() {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var con *RECmdConnectAttr
			con, err = p.parseConnectAttr()
			if re.ConnectAttrs == nil {
				re.ConnectAttrs = []*RECmdConnectAttr{}
			}
			re.ConnectAttrs = append(re.ConnectAttrs, con)
			commandNum--
		} else if p.CurTokenIs(string(RETypeRelationship)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var rs *RECmdRelationship
			rs, err = p.parseRelationship()
			if re.Relationships == nil {
				re.Relationships = []*RECmdRelationship{}
			}
			re.Relationships = append(re.Relationships, rs)
			commandNum--
		} else if p.CurTokenIs(string(RETypeLock)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var lk *RECmdLock
			lk, err = p.parseLock()
			if re.Locks == nil {
				re.Locks = []*RECmdLock{}
			}
			re.Locks = append(re.Locks, lk)
			commandNum--
		} else if p.CurTokenIs(string(RETypeUnlock)) {
			if res == nil {
				p.errs = append(p.errs, "res is nil...", *p.CurToken)
				continue
			}
			re := res[len(res)-1] // last re
			var ulk *RECmdUnlock
			ulk, err = p.parseUnlock()
			if re.Unlocks == nil {
				re.Unlocks = []*RECmdUnlock{}
			}
			re.Unlocks = append(re.Unlocks, ulk)
			commandNum--
		}

		if err != nil {
			p.errs = append(p.errs, err.Error())
		}
		p.NextToken()
	}
	return res
}

func (p *DataReferenceEditsParser) ErrorCheck() {
	if p.errs != nil {
		for _, err := range p.errs {
			fmt.Println(err)
		}
	}
}

func (p *DataReferenceEditsParser) CurTokenIs(t string) bool {
	return *p.CurToken == t
}

func (p *DataReferenceEditsParser) PeekTokenIs(t string) bool {
	if p.PeekToken == nil {
		return false
	}
	return *p.PeekToken == t
}

func (p *DataReferenceEditsParser) PeekTokenIsNumber() bool {
	if _, err := strconv.Atoi(*p.PeekToken); err == nil {
		return true
	}
	return false
}

func (p *DataReferenceEditsParser) NextToken() {
	p.CurToken = p.PeekToken
	p.cur++
	if p.cur < len(*p.token) {
		p.PeekToken = &(*p.token)[p.cur]
	} else {
		p.PeekToken = nil
	}
}

func NewDataReferenceEditsParser(token *[]string, start int) *DataReferenceEditsParser {
	p := &DataReferenceEditsParser{
		token: token,
		cur:   start,
	}
	p.NextToken()
	p.NextToken()
	return p
}

func MakeDataReferenceEdits(token *[]string, start int) ([]AttrValue, int, error) {
	referenceNode := (*token)[start]
	re := AttrDataReferenceEdits{
		TopReferenceNode: strings.Trim(referenceNode, "\""),
		ReferenceEdits:   []*ReferenceEdit{},
	}
	p := NewDataReferenceEditsParser(token, start)
	p.ErrorCheck()
	re.ReferenceEdits = p.ParseToken()

	a := []AttrValue{&re}
	return a, len(*token) - start, nil
}
