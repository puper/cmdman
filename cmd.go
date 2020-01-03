package cmdman

const (
	CMDTypeClose = "cmd.close"
)

type CMD struct {
	Owner  string
	Type   string
	Params interface{}

	Result chan interface{}
}

func NewCMD(owner string, typ string, params interface{}) *CMD {
	return &CMD{
		Owner:  owner,
		Type:   typ,
		Params: params,
		Result: make(chan interface{}, 1),
	}
}

func (this *CMD) SetResult(r interface{}) {
	this.Result <- r
}

func (this *CMD) GetResult() interface{} {
	return <-this.Result
}
