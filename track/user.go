package track

//玩家进入
func UserIn(roomKind int32, roomId uint64, userId uint64, credit int64) {
	t := GetTrack()
	userInEvent := newUserInEvent(
		t.gameKind,
		roomKind,
		roomId,
		userId,
		credit,
	)
	t.stream.CommitEvent(userInEvent)
}

//玩家离开
func UserOut(roomKind int32, roomId uint64, userId uint64, credit int64) {
	t := GetTrack()
	userInEvent := newUserOutEvent(
		t.gameKind,
		roomKind,
		roomId,
		userId,
		credit,
	)
	t.stream.CommitEvent(userInEvent)
}
