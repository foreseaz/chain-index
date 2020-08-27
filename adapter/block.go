package adapter

import (
	"time"

	"github.com/luci/go-render/render"
)

type RDbBlockCommittedCouncilNodeRow struct {
	BlockHeight        uint64    `json:"-"`
	ID                 uint64    `json:"id"`
	Name               string    `json:"name"`
	CouncilNodeAddress string    `json:"address"`
	Signature          string    `json:"signature"`
	IsProposer         bool      `json:"is_proposer"`
	CommitTime         time.Time `json:"timestamp"`
}

func (row *RDbBlockCommittedCouncilNodeRow) String() string {
	return render.Render(row)
}
