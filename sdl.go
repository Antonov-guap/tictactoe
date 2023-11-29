package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func runSDL(_ Game) {
	sdl.Main(func() {
		err := sdl.Init(sdl.INIT_EVERYTHING)
		fatalOn(err)

		for attr, val := range map[sdl.GLattr]int{
			sdl.GL_CONTEXT_PROFILE_MASK:  sdl.GL_CONTEXT_PROFILE_CORE,
			sdl.GL_CONTEXT_MAJOR_VERSION: 3,
			sdl.GL_CONTEXT_MINOR_VERSION: 3,
		} {
			if err := sdl.GLSetAttribute(attr, val); err != nil {
				log.Fatalln(err)
			}
		}

		if err := gl.Init(); err != nil {
			log.Fatalln(err)
		}

		window, err := sdl.CreateWindow(
			"TicTacToe", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_OPENGL,
		)
		if err != nil {
			log.Fatalln(err)
		}
		defer window.Destroy()

		if _, err = window.GLCreateContext(); err != nil {
			log.Fatalln(err)
		}

		version := gl.GoStr(gl.GetString(gl.VERSION))
		println("OpenGL Version:", version)

		var frames uint64
		go func() {
			for range time.Tick(1 * time.Second) {
				f := atomic.LoadUint64(&frames)
				fmt.Printf("%02d\n", f)
				atomic.StoreUint64(&frames, 0)
			}
		}()

		type shader struct {
			id    *uint32
			src   string
			sType uint32
		}

		ss := []shader{
			{
				id:    new(uint32),
				sType: gl.VERTEX_SHADER,
				src: `#version 330 core
				layout (location = 0) in vec3 aPos;
				void main()
				{
				   gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
    			}`,
			},
			{
				id:    new(uint32),
				sType: gl.FRAGMENT_SHADER,
				src: `#version 330 core
				out vec4 FragColor;
				void main()
				{
				    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
				}`,
			},
		}

		var srcLen int32
		progID := gl.CreateProgram()
		for i, s := range ss {
			src, free := gl.Strs(s.src)
			*s.id = gl.CreateShader(s.sType)
			srcLen = int32(len(s.src))
			gl.ShaderSource(*s.id, 1, src, &srcLen)
			gl.CompileShader(*s.id)
			free()
			var v int32
			gl.GetShaderiv(*s.id, gl.COMPILE_STATUS, &v)
			if v == gl.FALSE {
				var msg []uint8
				gl.GetShaderiv(*s.id, gl.INFO_LOG_LENGTH, &v)
				if v != 0 {
					msg = make([]uint8, int(v))
					gl.GetShaderInfoLog(*s.id, v, nil, (*uint8)(gl.Ptr(msg)))
				}
				log.Fatalf("opengl: shader %d compile failed: %s", i, string(msg))
			}
			gl.AttachShader(progID, *s.id)
		}
		gl.LinkProgram(progID)
		gl.UseProgram(progID)
		for _, s := range ss {
			gl.DeleteShader(*s.id)
		}

		for {
			sdl.PollEvent()
			drawgl()
			window.GLSwap()
			atomic.AddUint64(&frames, 1)
		}
	})
}

func drawgl() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.3, 0.2, 0.2, 1)
}

func fatalOn(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
