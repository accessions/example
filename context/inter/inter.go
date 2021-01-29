package inter

type Inter int

func (*Inter) Run ()  {

}

func (*Inter) Stop ()  {

}

type Ms interface {
	Run()
	Stop()
}
