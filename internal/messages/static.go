package messages

// Frequently used static messages

var (
	HelloMsg       string
	GetPeersMsg    string
	GetMempoolMsg  string
	GetChainTipMsg string
)

func init() {
	var err error

	agent := "marabobos"

	HelloMsg, err = MakeHelloMessage("0.10.0", &agent)
	if err != nil {
		panic(err)
	}
	GetPeersMsg, err = MakeGetPeersMessage()
	if err != nil {
		panic(err)
	}
	GetMempoolMsg, err = MakeGetMempoolMessage()
	if err != nil {
		panic(err)
	}
	GetChainTipMsg, err = MakeGetChainTipMessage()
	if err != nil {
		panic(err)
	}
}
