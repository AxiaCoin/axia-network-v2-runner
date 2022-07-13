// Copyright (C) 2019-2021, Axia Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package beacon

import (
	"github.com/axiacoin/axia-network-v2/ids"
	"github.com/axiacoin/axia-network-v2/utils"
)

// TODO: remove this in favor of an exported utility from axia

var _ Beacon = &beacon{}

type Beacon interface {
	ID() ids.ShortID
	IP() utils.IPDesc
}

type beacon struct {
	id ids.ShortID
	ip utils.IPDesc
}

func New(id ids.ShortID, ip utils.IPDesc) Beacon {
	return &beacon{
		id: id,
		ip: ip,
	}
}

func (b *beacon) ID() ids.ShortID  { return b.id }
func (b *beacon) IP() utils.IPDesc { return b.ip }
