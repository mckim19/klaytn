package backend

import (
	"fmt"
	"testing"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	istanbulCore "github.com/klaytn/klaytn/consensus/istanbul/core"
	"github.com/klaytn/klaytn/rlp"
)

func TestMy(t *testing.T) {
	var err error
	hHash := common.HexToHash("0x2fbd84bc89545c655e998f0b670f155053ab497865bf9b9a0b409b09a0bfefd3")
	proposalSeal := istanbulCore.PrepareCommittedSeal(hHash)

	ed := common.FromHex("0xd683010804846b6c617986676f312e3138856c696e7578000000000000000000f89ed5942004b5d6385f501f072712b086f9c5758b9b20beb8415f1088e033532e8bf581f196381c277bfa5bd0d6bfcc634e538b0862e6c3f5b611ae6e55bd39ecf4ab08fbb11ec0c3f0b4377b0016aca10416383aa5a86b178d00f843b8416ed8814448b878e56da8723a4ee2442d11c6f8077a28c847704fb03f0e14031b0528d4df431388d30589e4c66cfac10969a71daa41e3c36173c7c7c3ffa38f3901")
	var istanbulExtra *types.IstanbulExtra
	err = rlp.DecodeBytes(ed[types.IstanbulExtraVanity:], &istanbulExtra)
	fmt.Println(err)

	extra := istanbulExtra
	for _, seal := range extra.CommittedSeal {
		addr, err := cacheSignatureAddresses(proposalSeal, seal)
		fmt.Println(addr.Hex(), err)
	}
}
