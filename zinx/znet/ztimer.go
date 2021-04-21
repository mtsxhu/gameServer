package znet

import (
	"mmorpg/zinx/ziface"
)


type TimePile struct {
	TimerList []*ziface.ZTimer
}

func NewTimePile()*TimePile {
	return &TimePile{
		TimerList: make([]*ziface.ZTimer,10,20),
	}
}
func (this *TimePile)AddTimer(timer *ziface.ZTimer){

}
func (this *TimePile)RemoveTimer(timer *ziface.ZTimer){
}
func (this *TimePile)GetTop()*ziface.ZTimer{
	return nil
}
func (this *TimePile)RemoveTop()*ziface.ZTimer{
	return nil
}
func (this *TimePile) Tick() {

}