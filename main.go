package main

import (
	"fmt"
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
	"os"
	"strconv"
	"time"
)

func createWindow(conn *xgb.Conn, window xproto.Window) {
	screen := xproto.Setup(conn).DefaultScreen(conn)
	xproto.CreateWindow(conn, xproto.WindowClassCopyFromParent, window, screen.Root,
		0, 0, 300, 200, 0,
		xproto.WindowClassInputOnly, xproto.WindowClassCopyFromParent, xproto.CwEventMask,
		[]uint32{xproto.EventMaskPointerMotion})
}

func mouseLocation(conn *xgb.Conn, window xproto.Window) (int16, int16) {
	mousePos, err := xproto.QueryPointer(conn, window).Reply()
	if err != nil {
		fmt.Println("panic: " + err.Error())
		os.Exit(1)
	}

	return mousePos.RootX, mousePos.RootY

}

func moveMouse(conn *xgb.Conn, window xproto.Window, xLimit int16) {
	mousePosX, mousePosY := mouseLocation(conn, window)
	moveAmount := -(mousePosX - xLimit)

	if mousePosX > xLimit {
		xproto.WarpPointer(conn, xproto.WindowNone, xproto.WindowNone,
			mousePosX, mousePosY, 0, 0, moveAmount, 0)
	}
}

func main() {
	xLimit, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Enter an integer")
		os.Exit(1)
	}
	fmt.Println("Putting mouse in jail, no second monitor for you!")

	conn, _ := xgb.NewConn()
	wid, _ := xproto.NewWindowId(conn)
	createWindow(conn, wid)

	xproto.MapWindow(conn, wid)
	xLimit16 := int16(xLimit)

	for {
		time.Sleep(time.Microsecond)
		moveMouse(conn, wid, xLimit16)
	}
}
