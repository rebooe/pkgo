package queue

type Queue interface {
	// 推送任务到队列
	//	name: 队列名称
	//	arg: 执行任务的参数
	Push(name string, msg Msg) error

	// 获取队列中的任务,队列中无任务时该方法会阻塞
	//	name: 队列名称
	//	ack: 是否在队列中删除任务
	Pop(name string, opts ...OptionFunc) (*Msg, error)
}

type Msg struct {
	ContentType string
	Body        any
}

type OptionFunc func(opt *Option)

type Option struct {
	Ack  bool // 是否自动删除任务
	Wait bool // 是否阻塞
}

func defaultOption() *Option {
	return &Option{
		Ack:  true,
		Wait: true,
	}
}

func OptAck(ok bool) OptionFunc {
	return func(opt *Option) {
		opt.Ack = ok
	}
}

func OptWait(ok bool) OptionFunc {
	return func(opt *Option) {
		opt.Wait = ok
	}
}
